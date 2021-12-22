package nodes

import (
	"errors"
	"fmt"
	"github.com/markel1974/kuery/src/config"
	"github.com/markel1974/kuery/src/context"
	"github.com/markel1974/kuery/src/objects"
)

type FunctionExists struct {
	fieldName    string
	fieldNameArg INode
}

func NewFunctionExists(fieldName string) *FunctionExists {
	f := &FunctionExists{
		fieldName:    fieldName,
		fieldNameArg: NewLiteralNode(fieldName),
	}
	return f
}

func (f *FunctionExists) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionExists) GetValue() interface{} {
	return nil
}

func (f *FunctionExists) SetValue(_ interface{}) {
}

func (f *FunctionExists) Clone() INode {
	return NewFunctionExists(f.fieldName)
}

func (f *FunctionExists) Compile(indexPattern *objects.IndexPattern, _ *config.Config, ctx *context.Context) (interface{}, error) {
	if f.fieldNameArg == nil {
		return nil, errors.New("missing field")
	}
	fullFieldNameArg := f.fieldNameArg.Clone()

	if v := f.fieldNameArg.GetValue(); v != nil {
		field := fmt.Sprintf("%v", v)
		if ctx != nil && ctx.Nested != nil {
			fullFieldNameArg.SetValue(ctx.Nested.Path + "." + field)
		} else {
			fullFieldNameArg.SetValue(field)
		}
	}

	fieldName, err := fullFieldNameArg.Compile(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	if indexPattern != nil {
		fields := indexPattern.Find(fmt.Sprintf("%v", fieldName))
		if fields != nil && fields.Scripted {
			return nil, errors.New("exists query does not support scripted fields")
		}
	}
	out := map[string]interface{}{
		"exists": map[string]interface{}{
			"field": fieldName,
		},
	}
	return out, nil
}
