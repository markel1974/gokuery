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

type CursorNode struct {
	suggestionTypes []string
	nestedPath      string
	fieldName       string
	start           int
	end             int
	prefix          string
	suffix          string
	text            string
	value           interface{}
}

func createCursorNode(start int, end int, prefix string, suffix string, text string) *CursorNode {
	node := &CursorNode{
		suggestionTypes: nil,
		nestedPath:      "",
		start:           start,
		end:             end,
		prefix:          prefix,
		suffix:          suffix,
		text:            text,
	}
	return node
}

func NewCursorNode(start int, end int, prefix string, suffix string, text string) INode {
	return createCursorNode(start, end, prefix, suffix, text)
}

func (f *CursorNode) GetType() NodeType {
	return TypeCursor
}

func (f *CursorNode) Compile(_ *objects.IndexPattern, _ *config.Config, _ *context.Context) (interface{}, error) {
	return nil, nil
}

func (f *CursorNode) GetValue() interface{} {
	return f.value
}

func (f *CursorNode) SetValue(value interface{}) {
	f.value = value
}

func (f *CursorNode) Clone() INode {
	return createCursorNode(f.start, f.end, f.prefix, f.suffix, f.text)
}

func (f *CursorNode) Copy() *CursorNode {
	return createCursorNode(f.start, f.end, f.prefix, f.suffix, f.text)
}

func (f *CursorNode) GetSuggestionTypes() []string {
	return f.suggestionTypes
}

func (f *CursorNode) SetSuggestionTypes(suggestionTypes []string) {
	f.suggestionTypes = suggestionTypes
}

func (f *CursorNode) SetFieldName(fieldName string) {
	f.fieldName = fieldName
}

func (f *CursorNode) GetFieldName() string {
	return f.fieldName
}

func (f *CursorNode) SetNestedPath(nestedPath string) {
	f.nestedPath = nestedPath
}

func (f *CursorNode) GetNestedPath() string {
	return f.nestedPath
}

func (f *CursorNode) GetPrefix() string {
	return f.prefix
}

func (f *CursorNode) GetSuffix() string {
	return f.suffix
}
