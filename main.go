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

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/markel1974/gokuery/src/config"
	"github.com/markel1974/gokuery/src/context"
	"github.com/markel1974/gokuery/src/nodes"
	"github.com/markel1974/gokuery/src/objects"
	"github.com/markel1974/gokuery/src/version"
)

func main() {
	var help bool
	var showVersion bool
	var kql string
	var indexPattern string
	var escapeQueryString bool

	flag.StringVar(&indexPattern, "i", "", "index pattern (json)")
	flag.StringVar(&kql, "k", "", "kql statement")
	flag.BoolVar(&escapeQueryString, "e", true, "escape querystring")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&help, "h", false, "show this help")
	flag.Parse()

	if help {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println(version.AppName + " " + version.AppVersion)
		return
	}

	cfg := config.New()
	cfg.EscapeQueryString = escapeQueryString

	ctx := context.New()

	var ip *objects.IndexPattern
	if len(indexPattern) > 0 {
		var err error
		ip, err = objects.UnmarshalIndexPattern([]byte(indexPattern))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	//const kRange = `account_number >= 100 and items_sold <= 200`
	//const kRange = "(account_number >= 100 and items_sold <= 200) or k:10"
	//const kNot = `NOT a:true* OR (k:11)`
	//const t = `\+`
	out, err := nodes.ParseKueryString(kql, ip, cfg, ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	j, _ := json.Marshal(out)
	fmt.Println(string(j))
}
