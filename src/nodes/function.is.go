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
	"github.com/markel1974/gokuery/src/config"
	"github.com/markel1974/gokuery/src/context"
	"github.com/markel1974/gokuery/src/objects"
	"github.com/markel1974/gokuery/src/utils"
)

type FunctionIs struct {
	fieldNameArg INode
	valueArg     INode
	isPhraseArg  INode
	value        interface{}
}

func NewFunctionIs(fieldNameArg INode, valueArg INode, isPhraseArg INode) INode {
	f := &FunctionIs{
		fieldNameArg: fieldNameArg,
		valueArg:     valueArg,
		isPhraseArg:  isPhraseArg,
	}
	return f
}

func (f *FunctionIs) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionIs) GetValue() interface{} {
	return f.value
}

func (f *FunctionIs) SetValue(value interface{}) {
	f.value = value
}

func (f *FunctionIs) Clone() INode {
	out := &FunctionIs{
		fieldNameArg: f.fieldNameArg,
		valueArg:     f.valueArg,
		isPhraseArg:  f.isPhraseArg,
	}
	return out
}

func (f *FunctionIs) Compile(indexPattern *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	if f.fieldNameArg == nil {
		return nil, errors.New("missing field")
	}
	if f.valueArg == nil {
		return nil, errors.New("missing value")
	}

	var path string

	if ctx != nil && ctx.Nested != nil {
		path = ctx.Nested.Path
	}

	fullFieldNameArg, err := GetFullFieldNameNode(f.fieldNameArg, indexPattern, path)
	if err != nil {
		return nil, err
	}

	elField, err := fullFieldNameArg.Compile(nil, nil, nil)
	if err != nil {
		return nil, err
	}

	elValue, err := f.valueArg.Compile(nil, nil, nil)
	if err != nil {
		return nil, err
	}

	var kind string
	if f.isPhraseArg.GetValue() != nil {
		kind = "phrase"
	} else {
		kind = "best_fields"
	}

	if elField == nil {
		if f.valueArg.GetType() == TypeWildcard {
			wildcard, _ := f.valueArg.(*WildcardNode)
			q := map[string]interface{}{
				"query_string": map[string]interface{}{
					"query": wildcard.ToQueryStringQuery(cfg.EscapeQueryString),
				},
			}
			return q, nil
		}
		q := map[string]interface{}{
			"multi_match": map[string]interface{}{
				"type":    kind,
				"query":   elValue,
				"lenient": true,
			},
		}
		return q, nil
	}

	fieldName := fmt.Sprintf("%v", elField)
	var fields []*objects.Field

	if indexPattern != nil {
		fields = GetFields(fullFieldNameArg, indexPattern)
	}

	if len(fields) == 0 {
		fields = append(fields, &objects.Field{Name: fieldName, Scripted: false})
	}
	isExistsQuery := f.valueArg.GetType() == TypeWildcard && elValue == "*"
	isAllFieldsQuery := fullFieldNameArg.GetType() == TypeWildcard && fieldName == "*" || indexPattern != nil && len(fields) == indexPattern.FieldsLen()
	isMatchAllQuery := isExistsQuery && isAllFieldsQuery

	if isMatchAllQuery {
		q := map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
		return q, nil
	}

	var initial map[string]interface{}

	queries, err := utils.Reduce(fields, initial, func(accumulator map[string]interface{}, field *objects.Field, idx int) interface{} {
		wrapWithNestedQuery := func(query interface{}) interface{} {
			var nested *objects.Nested
			var nestedPath string
			if field.SubType != nil && field.SubType.Nested != nil {
				nested = field.SubType.Nested
				nestedPath = nested.Path
			}
			if !(fullFieldNameArg.GetType() == TypeWildcard) || nested == nil || ctx.Nested != nil {
				return query
			} else {
				return map[string]interface{}{"nested": map[string]interface{}{"path": nestedPath, "query": query, "score_mode": "none"}}
			}
		}

		if field.Scripted {
			if !isExistsQuery {
				script, err := GetPhraseScript(field, fmt.Sprintf("%v", elValue))
				if err != nil {
					return err
				}
				var x []interface{}
				if accumulator != nil {
					x = append(x, accumulator)
				}
				q := map[string]interface{}{
					"script": script,
				}
				x = append(x, q)
				return x
			}
			return nil
		}

		if isExistsQuery {
			var x []interface{}
			if accumulator != nil {
				x = append(x, accumulator)
			}
			q := map[string]interface{}{
				"exists": map[string]interface{}{
					"field": field.Name},
			}
			x = append(x, wrapWithNestedQuery(q))
			return x
		}

		if f.valueArg.GetType() == TypeWildcard {
			wildcard, _ := f.valueArg.(*WildcardNode)
			var x []interface{}
			if accumulator != nil {
				x = append(x, accumulator)
			}
			q := map[string]interface{}{
				"query_string": map[string]interface{}{
					"fields": []string{field.Name},
					"query":  wildcard.ToQueryStringQuery(cfg.EscapeQueryString),
				},
			}
			x = append(x, wrapWithNestedQuery(q))
			return x
		}

		if field.Type == "date" {
			var x []interface{}
			if accumulator != nil {
				x = append(x, accumulator)
			}
			qRange := map[string]interface{}{
				"gte": elValue,
				"lte": elValue,
			}
			if cfg != nil && cfg.HasTimeZone() {
				qRange["time_zone"] = cfg.GetTimeZone()
			}
			q := map[string]interface{}{
				"range": map[string]interface{}{
					field.Name: qRange,
				},
			}
			x = append(x, wrapWithNestedQuery(q))
			return x
		}

		var queryType string
		if kind == "phrase" {
			queryType = "match_phrase"
		} else {
			queryType = "match"
		}
		var x []interface{}
		if accumulator != nil {
			x = append(x, accumulator)
		}
		q := map[string]interface{}{
			queryType: map[string]interface{}{
				field.Name: elValue,
			},
		}
		x = append(x, wrapWithNestedQuery(q))
		return x
	})

	return map[string]interface{}{
		"bool": map[string]interface{}{
			"should":               queries,
			"minimum_should_match": 1,
		}}, nil
}
