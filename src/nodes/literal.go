package nodes

import (
	"markel/home/kuery/src/config"
	"markel/home/kuery/src/context"
	"markel/home/kuery/src/objects"
)

type LiteralNode struct {
	value   interface{}
}

func NewLiteralNode(value interface{}) INode {
	l := &LiteralNode{
		value: value,
	}
	return l
}

func (l * LiteralNode) GetType() NodeType {
	return TypeLiteral
}

func (l * LiteralNode) GetValue() interface{} {
	return l.value
}

func (l * LiteralNode) SetValue(value interface{}) {
	l.value = value
}

func (l * LiteralNode) Clone() INode {
	return NewLiteralNode(l.value)
}

func (l * LiteralNode) Compile(_ * objects.IndexPattern, _ * config.Config, _ * context.Context) (interface{}, error) {
	return l.value, nil
}