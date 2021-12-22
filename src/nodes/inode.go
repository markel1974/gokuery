package nodes

import (
	"markel/home/kuery/src/config"
	"markel/home/kuery/src/context"
	"markel/home/kuery/src/objects"
)

type NodeType int

const (
	TypeFunction NodeType = iota
	TypeLiteral  NodeType = iota
	TypeWildcard NodeType = iota
	TypeCursor   NodeType = iota
)

type INode interface {
	GetType() NodeType
	Compile(indexPattern * objects.IndexPattern, cfg * config.Config, ctx * context.Context) (interface{}, error)
	GetValue() interface{}
	SetValue(value interface{})
	Clone() INode
}