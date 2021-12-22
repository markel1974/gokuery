/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nodes

import (
	"errors"
	"fmt"
	"github.com/markel1974/gokuery/src/objects"
	"reflect"
	"strings"
)

func Concat(src interface{}) string {
	var seq []byte
	data, _ := src.([]interface{})
	for _, x := range data {
		switch d := x.(type) {
		case byte:
			seq = append(seq, d)
		case []byte:
			seq = append(seq, d...)
		case string:
			seq = append(seq, []byte(d)...)
		}
	}
	return string(seq)
}

func Trim(src string) string {
	return strings.Trim(src, " \r\n\t")
}

func Reduce(source, initialValue, reducer interface{}) (interface{}, error) {
	srcV := reflect.ValueOf(source)
	kind := srcV.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, errors.New("source value is not an array")
	}
	if reducer == nil {
		return nil, errors.New("reducer function cannot be nil")
	}
	rv := reflect.ValueOf(reducer)
	if rv.Kind() != reflect.Func {
		return nil, errors.New("reducer argument must be a function")
	}
	accumulator := initialValue
	accV := reflect.ValueOf(accumulator)
	for i := 0; i < srcV.Len(); i++ {
		entry := srcV.Index(i)
		reduceResults := rv.Call([]reflect.Value{
			accV,
			entry,
			reflect.ValueOf(i),
		})
		accV = reduceResults[0]
	}
	return accV.Interface(), nil
}

func GetFields(node INode, indexPattern *objects.IndexPattern) []*objects.Field {
	if node.GetType() == TypeLiteral {
		q, err := node.Compile(nil, nil, nil)
		if err != nil {
			return nil
		}
		fieldName, ok := q.(string)
		if !ok {
			return nil
		}
		field := indexPattern.Find(fieldName)
		if field == nil {
			return nil
		}
		return []*objects.Field{field}
	}
	if node.GetType() == TypeWildcard {
		if wn, ok := node.(*WildcardNode); ok && wn != nil {
			fields := indexPattern.Filter(func(field string) bool { return wn.Test(field) })
			return fields
		}
	}
	return nil
}

func GetFullFieldNameNode(rootNameNode INode, indexPattern *objects.IndexPattern, nestedPath string) (INode, error) {
	fullFieldNameNode := rootNameNode.Clone()

	if v := rootNameNode.GetValue(); v != nil {
		path := fmt.Sprintf("%v", v)
		if len(nestedPath) > 0 {
			fullFieldNameNode.SetValue(nestedPath + "." + path)
		} else {
			fullFieldNameNode.SetValue(path)
		}
	}

	if indexPattern == nil || fullFieldNameNode.GetType() == TypeWildcard && len(nestedPath) == 0 {
		return fullFieldNameNode, nil
	}
	fields := GetFields(fullFieldNameNode, indexPattern)
	var errs []string
	for _, field := range fields {
		var nestedPathFromField string
		if field.SubType != nil && field.SubType.Nested != nil && len(field.SubType.Nested.Path) > 0 {
			nestedPathFromField = field.SubType.Nested.Path
		}
		if len(nestedPath) > 0 && len(nestedPathFromField) == 0 {
			errs = append(errs, field.Name+" is not a nested field but is in nested group"+nestedPath+"in the KQL expression")
			continue
		}
		if len(nestedPathFromField) > 0 && len(nestedPath) == 0 {
			errs = append(errs, field.Name+" is a nested field, but is not in a nested group in the KQL expression")
			continue
		}
		if nestedPathFromField != nestedPath {
			errs = append(errs, field.Name+" is being queried with the incorrect nested path. The correct path is "+nestedPathFromField)
			continue
		}
	}
	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return fullFieldNameNode, nil
}

func GetConvertedValueForField(field *objects.Field, value interface{}) (interface{}, error) {
	if _, isBoolean := value.(bool); !isBoolean && field.Type == "boolean" {
		switch x := value.(type) {
		case bool:
			return x, nil
		case int:
			if x == 0 {
				return false, nil
			} else {
				return true, nil
			}
		case string:
			if strings.ToLower(x) == "false" {
				return false, nil
			} else {
				return true, nil
			}
		default:
			return nil, errors.New(fmt.Sprintf("%v is not a valid boolean value for boolean field %s", value, field.Name))
		}
	}
	return value, nil
}

func BuildInlineScriptForPhraseFilter(scriptedField *objects.Field) string {
	if scriptedField.Lang == "painless" {
		return `boolean compare(Supplier s, def v) {return s.get() == v;} compare(() -> { ${scriptedField.script} }, params.value);`
	}
	return "(${" + scriptedField.Script + "}) == value"
}

func GetPhraseScript(field *objects.Field, value string) (interface{}, error) {
	convertedValue, err := GetConvertedValueForField(field, value)
	if err != nil {
		return nil, err
	}
	script := BuildInlineScriptForPhraseFilter(field)
	q := map[string]interface{}{
		"script": map[string]interface{}{
			"source": script,
			"lang":   field.Lang,
			"params": map[string]interface{}{
				value: convertedValue,
			},
		},
	}
	return q, nil
}
