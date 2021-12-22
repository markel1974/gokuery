package nodes

import (
	"errors"
	"github.com/markel1974/kuery/src/config"
	"github.com/markel1974/kuery/src/context"
	"github.com/markel1974/kuery/src/objects"
)

type FunctionOr struct {
	left  INode
	right INode
}

func NewFunctionOr(left INode, right INode) INode {
	f := &FunctionOr{
		left:  left,
		right: right,
	}
	return f
}

func (f *FunctionOr) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionOr) GetValue() interface{} {
	return nil
}

func (f *FunctionOr) SetValue(_ interface{}) {
}

func (f *FunctionOr) Clone() INode {
	return NewFunctionOr(f.left, f.right)
}

func (f *FunctionOr) Compile(indexPattern *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	if f.left == nil {
		return nil, errors.New("missing left value")
	}
	if f.right == nil {
		return nil, errors.New("missing right value")
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
			"should": []interface{}{ql, qr},
		},
	}
	return out, nil
}
