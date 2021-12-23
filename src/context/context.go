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

package context

import (
	"encoding/json"
)

type Context struct {
	Nested *Nested `json:"nested"`
}

func New() *Context {
	return &Context{}
}

func Unmarshal(data []byte) (*Context, error) {
	var ctx Context
	err := json.Unmarshal(data, &ctx)
	if err != nil {
		return nil, err
	}
	return &ctx, nil
}

func (ctx *Context) Clone() *Context {
	z := New()
	if ctx.Nested != nil {
		z.Nested = ctx.Nested.Clone()
	}
	return z
}

func (ctx *Context) UpdatePath(path string) string {
	if len(path) > 0 && ctx.Nested != nil && len(ctx.Nested.Path) > 0 {
		return ctx.Nested.Path + "." + path
	}
	return path
}
