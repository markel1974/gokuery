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
	"github.com/markel1974/gokuery/src/config"
	"github.com/markel1974/gokuery/src/context"
	"github.com/markel1974/gokuery/src/objects"
)

type FunctionNot struct {
	node INode
}

func NewFunctionNot(node INode) INode {
	f := &FunctionNot{
		node: node,
	}
	return f
}

func (f *FunctionNot) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionNot) GetValue() interface{} {
	return nil
}

func (f *FunctionNot) SetValue(_ interface{}) {
}

func (f *FunctionNot) Clone() INode {
	return NewFunctionNot(f.node)
}

func (f *FunctionNot) Compile(indexPattern *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	var val interface{}
	if f.node != nil {
		var err error
		if val, err = f.node.Compile(indexPattern, cfg, ctx); err != nil {
			return nil, err
		}
	}
	q := map[string]interface{}{
		"bool": map[string]interface{}{
			"must_not": val,
		},
	}
	return q, nil
}
