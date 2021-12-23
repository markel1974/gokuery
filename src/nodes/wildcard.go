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
	"regexp"
	"strings"
)

const WildcardSymbol = "@kuery-wildcard@"

var _escapeRgx *regexp.Regexp
var _escapeQueryStringRgx *regexp.Regexp

func init() {
	_escapeRgx = regexp.MustCompile(`[.*+?^${}()|[\]\\]`)
	_escapeQueryStringRgx = regexp.MustCompile(`[+\-=&|><!(){}[\]^"~*?:\\/]`)
}

func EscapeRegExp(src string) string {
	res := _escapeRgx.ReplaceAllStringFunc(src, func(in string) string {
		return "\\" + in
	})
	return res
}

func EscapeQueryString(src string) string {
	res := _escapeQueryStringRgx.ReplaceAllStringFunc(src, func(in string) string {
		return "\\" + in
	})
	return res
}

type WildcardNode struct {
	value string
}

func NewWildcardNode(value string) INode {
	w := &WildcardNode{
		value: value,
	}
	return w
}

func (w *WildcardNode) GetType() NodeType {
	return TypeWildcard
}

func (w *WildcardNode) GetValue() interface{} {
	return w.value
}

func (w *WildcardNode) SetValue(value interface{}) {
	if v, ok := value.(string); ok {
		w.value = v
	}
}

func (w *WildcardNode) Clone() INode {
	return NewWildcardNode(w.value)
}

func (w *WildcardNode) Compile(_ *objects.IndexPattern, _ *config.Config, _ *context.Context) (interface{}, error) {
	out := strings.Replace(w.value, WildcardSymbol, "*", -1)
	return out, nil
}

func (w *WildcardNode) ToQueryStringQuery(escape bool) string {
	values := strings.Split(w.value, WildcardSymbol)
	var res []string
	for _, v := range values {
		if escape {
			res = append(res, EscapeQueryString(v))
		} else {
			res = append(res, v)
		}
	}
	out := strings.Join(res, "*")
	return out
}

func (w *WildcardNode) HasLeadingWildcard() bool {
	out := strings.HasPrefix(w.value, WildcardSymbol) && len(strings.Replace(w.value, WildcardSymbol, "", -1)) > 0
	return out
}

func (w *WildcardNode) Test(src string) bool {
	var p []string
	for _, v := range strings.Split(w.value, WildcardSymbol) {
		if r := EscapeRegExp(v); len(r) > 0 {
			p = append(p, r)
		}
	}
	regex := strings.Join(p, "[\\s\\S]*")
	rgx, err := regexp.Compile("^" + regex + "$")
	if err != nil {
		return false
	}
	return rgx.MatchString(src)
}
