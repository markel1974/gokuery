package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"markel/home/kuery/src/config"
	"markel/home/kuery/src/nodes"
	"markel/home/kuery/src/version"
	"strings"
)

func main() {
	var help bool
	var showVersion bool
	var kql string

	flag.StringVar(&kql, "k", "", "kql statement")
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
	cfg := config.NewConfig()
	cfg.EscapeQueryString = true
	out, err := res.Compile(nil, cfg, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	j, _ := json.Marshal(out)
	fmt.Println(string(j))
}
