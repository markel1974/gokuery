package nodes

import (
	"errors"
	"github.com/markel1974/gokuery/src/config"
	"github.com/markel1974/gokuery/src/context"
	"github.com/markel1974/gokuery/src/objects"
	"strings"
)

func ParseKueryString(kql string, ip *objects.IndexPattern, cfg *config.Config, ctx *context.Context) (interface{}, error) {
	got, err := ParseReader("", strings.NewReader(kql), GlobalStore("config", cfg))
	if err != nil {
		return nil, err
	}
	res, _ := got.(INode)
	if res == nil {
		return nil, errors.New("can't cast to inode")
	}
	out, err := res.Compile(ip, cfg, ctx)
	if err != nil {
		return nil, err
	}
	return out, nil
}
