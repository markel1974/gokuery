package nodes

import (
	"errors"
	"fmt"
	"github.com/markel1974/kuery/src/config"
	"github.com/markel1974/kuery/src/context"
	"github.com/markel1974/kuery/src/objects"
	"strings"
)

type FunctionRange struct {
	field   INode
	opNode  INode
	argNode INode
}

func NewFunctionRange(field INode, opNode INode, argNode INode) INode {
	f := &FunctionRange{
		field:   field,
		opNode:  opNode,
		argNode: argNode,
	}
	return f
}

func (f *FunctionRange) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionRange) GetValue() interface{} {
	return nil
}

func (f *FunctionRange) SetValue(_ interface{}) {
}

func (f *FunctionRange) Clone() INode {
	return NewFunctionRange(f.field, f.opNode, f.argNode)
}

func (f *FunctionRange) Compile(indexPattern *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	if f.field == nil {
		return nil, errors.New("range: missing field")
	}

	fieldNameArg := f.field
	var op string
	var arg interface{}

	if f.opNode != nil {
		if v := f.opNode.GetValue(); v != nil {
			op = strings.ToLower(fmt.Sprintf("%v", v))
		}
	}

	switch op {
	case "gt", "lt", "gte", "lte", "format":
	default:
		return nil, errors.New("range: invalid op " + op)
	}

	if f.argNode != nil {
		if v := f.argNode.GetValue(); v != nil {
			arg = v
		}
	}

	var fields []*objects.Field
	if indexPattern != nil {
		fields = GetFields(fieldNameArg, indexPattern)
	}
	if len(fields) == 0 {
		v, _ := fieldNameArg.Compile(nil, nil, nil)
		name := fmt.Sprintf("%v", v)
		fields = append(fields, &objects.Field{Name: name, Scripted: false})
	}

	var queries []interface{}

	for _, field := range fields {
		wrapWithNestedQuery := func(query interface{}) interface{} {
			var nested *objects.Nested
			var nestedPath string
			if field.SubType != nil && field.SubType.Nested != nil {
				nested = field.SubType.Nested
				nestedPath = nested.Path
			}
			if !(fieldNameArg.GetType() == TypeWildcard) || nested == nil || ctx.Nested != nil {
				return query
			} else {
				return map[string]interface{}{
					"nested": map[string]interface{}{
						"path":       nestedPath,
						"query":      query,
						"score_mode": "none",
					},
				}
			}
		}

		if field.Scripted {
			q := map[string]interface{}{
				//TODO getRangeScript
				//"script": getRangeScript(field, queryParams),
			}
			queries = append(queries, q)
			continue
		}

		if field.Type == "date" {
			qRange := map[string]interface{}{
				op: arg,
			}
			if cfg != nil && cfg.HasTimeZone() {
				qRange["time_zone"] = cfg.GetTimeZone()
			}
			q := map[string]interface{}{
				"range": map[string]interface{}{
					field.Name: qRange,
				},
			}
			queries = append(queries, wrapWithNestedQuery(q))
			continue
		}

		q := map[string]interface{}{
			"range": map[string]interface{}{
				field.Name: map[string]interface{}{
					op: arg,
				},
			},
		}
		queries = append(queries, wrapWithNestedQuery(q))
	}

	return map[string]interface{}{
		"bool": map[string]interface{}{
			"should":               queries,
			"minimum_should_match": 1,
		},
	}, nil
}
