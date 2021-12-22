package nodes

/*
import (
	"fmt"
	"markel/home/kuery/src/config"
	"markel/home/kuery/src/context"
	"markel/home/kuery/src/objects"
)


func buildNode(value interface{}) INode {
	var argumentNode INode
	if v, ok := value.(INode); ok {
		if v.GetType() == TypeLiteral {
			argumentNode = v
		}
	}
	if argumentNode == nil {
		argumentNode = NewLiteralNode(value)
	}
	return argumentNode
}


type NamedArgNode struct {
	name  string
	value interface{}
}

func NewNamedArgNode(name string, value interface{}) INode {
	n := &NamedArgNode{
		name: name,
		value: value,
	}
	return n
}

func (n * NamedArgNode) GetType() NodeType {
	return TypeNamedArg //n.node.GetType()
}

func (n * NamedArgNode) GetValue() interface{} {
	return fmt.Sprintf("%v", n.value)
}

func (n * NamedArgNode) SetValue(value interface{}) {
	n.value = value
}

func (n * NamedArgNode) Clone() INode {
	return NewNamedArgNode(n.name, n.value)
}

func (n * NamedArgNode) Compile(indexPattern * objects.IndexPattern, cfg * config.Config, ctx * context.Context) (interface{}, error) {
	node := buildNode(n.value)
	return node.Compile(indexPattern, cfg, ctx)
	//return n.node.Compile(indexPattern, cfg, ctx)
}

/*
import _ from 'lodash';
import * as ast from '../ast';
import { nodeTypes } from '../node_types';

export function buildNode(name, value) {
  const argumentNode =
      _.get(value, 'type') === 'literal' ? value : nodeTypes.literal.buildNode(value);
  return {
    type: 'namedArg',
    name,
    value: argumentNode,
  };
}

export function toElasticsearchQuery(node) {
  return ast.toElasticsearchQuery(node.value);
}

 */