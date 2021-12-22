package nodes

import (
	"markel/home/kuery/src/config"
	"markel/home/kuery/src/context"
	"markel/home/kuery/src/objects"
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

func (w * WildcardNode) GetType() NodeType {
	return TypeWildcard
}

func (w * WildcardNode) GetValue() interface{} {
	return w.value
}

func (w * WildcardNode) SetValue(value interface{}) {
	if v, ok := value.(string); ok {
		w.value = v
	}
}

func (w * WildcardNode) Clone() INode {
	return NewWildcardNode(w.value)
}

func (w * WildcardNode) Compile(_ * objects.IndexPattern, _ * config.Config, _ * context.Context) (interface{}, error) {
	out := strings.Replace(w.value, WildcardSymbol, "*", -1)
	return out, nil
}

func (w * WildcardNode) ToQueryStringQuery(escape bool) string {
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

func (w * WildcardNode) HasLeadingWildcard() bool {
	out := strings.HasPrefix(w.value, WildcardSymbol) && len(strings.Replace(w.value, WildcardSymbol, "", -1)) > 0
	return out
}

func (w * WildcardNode) Test(src string) bool {
	var p []string
	for _, v := range strings.Split(w.value, WildcardSymbol) {
		if r := EscapeRegExp(v); len(r) > 0 {
			p = append(p, r)
		}
	}
	regex := strings.Join(p,"[\\s\\S]*")
	rgx, err := regexp.Compile("^$" + regex + "$")
	if err != nil {
		return false
	}
	return rgx.MatchString(src)
}