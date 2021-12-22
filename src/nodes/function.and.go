package nodes

import (
	"errors"
	"github.com/markel1974/kuery/src/config"
	"github.com/markel1974/kuery/src/context"
	"github.com/markel1974/kuery/src/objects"
)

type FunctionAnd struct {
	left  INode
	right INode
}

func NewFunctionAnd(left INode, right INode) INode {
	f := &FunctionAnd{
		left:  left,
		right: right,
	}
	return f
}

func (f *FunctionAnd) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionAnd) GetValue() interface{} {
	return nil
}

func (f *FunctionAnd) SetValue(_ interface{}) {
}

func (f *FunctionAnd) Clone() INode {
	return NewFunctionAnd(f.left, f.right)
}

func (f *FunctionAnd) Compile(indexPattern *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	if f.left == nil {
		return nil, errors.New("and: missing left value")
	}
	if f.right == nil {
		return nil, errors.New("and: missing right value")
	}
	ql, err := f.left.Compile(indexPattern, cfg, ctx)
	if err != nil {
		return nil, err
	}
	qr, err := f.right.Compile(indexPattern, cfg, ctx)
	if err != nil {
		return nil, err
	}
	out := map[string]interface{}{
		"bool": map[string]interface{}{
			"filter": []interface{}{ql, qr},
		},
	}
	return out, nil
}
