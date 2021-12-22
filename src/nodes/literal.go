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

type LiteralNode struct {
	value interface{}
}

func NewLiteralNode(value interface{}) INode {
	l := &LiteralNode{
		value: value,
	}
	return l
}

func (l *LiteralNode) GetType() NodeType {
	return TypeLiteral
}

func (l *LiteralNode) GetValue() interface{} {
	return l.value
}

func (l *LiteralNode) SetValue(value interface{}) {
	l.value = value
}

func (l *LiteralNode) Clone() INode {
	return NewLiteralNode(l.value)
}

func (l *LiteralNode) Compile(_ *objects.IndexPattern, _ *config.Config, _ *context.Context) (interface{}, error) {
	return l.value, nil
}
