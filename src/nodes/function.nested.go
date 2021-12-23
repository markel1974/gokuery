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
)

type FunctionNested struct {
	path  INode
	child INode
}

func NewFunctionNested(path INode, child INode) INode {
	f := &FunctionNested{
		path:  path,
		child: child,
	}
	return f
}

func (f *FunctionNested) GetType() NodeType {
	return TypeFunction
}

func (f *FunctionNested) GetValue() interface{} {
	return nil
}

func (f *FunctionNested) SetValue(_ interface{}) {
}

func (f *FunctionNested) Clone() INode {
	return NewFunctionNested(f.path, f.child)
}

func (f *FunctionNested) Compile(indexPattern *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	if f.path == nil {
		return nil, errors.New("nested: nil path")
	}
	if f.child == nil {
		return nil, errors.New("nested: nil child")
	}

	stringPath, err := f.path.Compile(indexPattern, cfg, ctx)
	if err != nil {
		return nil, err
	}
	var fullPath string
	if ctx.Nested != nil && len(ctx.Nested.Path) > 0 {
		fullPath = ctx.Nested.Path + "." + fmt.Sprintf("%v", stringPath)
	}

	z := ctx.Clone()
	if z.Nested == nil {
		z.Nested = context.NewNested()
	}
	z.Nested.Path = fullPath

	child, err := f.child.Compile(indexPattern, cfg, z)
	if err != nil {
		return nil, err
	}

	q := map[string]interface{}{
		"nested": map[string]interface{}{
			"path":       fullPath,
			"query":      child,
			"score_mode": "none",
		},
	}
	return q, nil
}
