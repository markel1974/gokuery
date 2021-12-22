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
	"github.com/markel1974/gokuery/src/nodes"
	"github.com/markel1974/gokuery/src/objects"
	"github.com/markel1974/gokuery/src/version"
	"strings"
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

	cfg := config.NewConfig()
	cfg.EscapeQueryString = escapeQueryString

	var ip *objects.IndexPattern
	if len(indexPattern) > 0 {
		ip = objects.NewIndexPattern()
		err := json.Unmarshal([]byte(indexPattern), &ip)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	//const kRange = `account_number >= 100 and items_sold <= 200`
	//const kRange = "(account_number >= 100 and items_sold <= 200) or k:10"
	//const kNot = `NOT a:true* OR (k:11)`
	//const t = `\+`
	got, err := nodes.ParseReader("", strings.NewReader(kql))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	res, _ := got.(nodes.INode)
	if res == nil {
		fmt.Println("can't cast to inode")
		return
	}
	out, err := res.Compile(ip, cfg, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	j, _ := json.Marshal(out)
	fmt.Println(string(j))
}
