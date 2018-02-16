package builders

import (
	"../visitor"
	"fmt"
	"reflect"
	"os"
	"bufio"
	"strings"
)

type Builder struct {
	SnippetsPath string
	RequestUrl   string
	OutputPath   string
}

func (b *Builder)Build(v visitor.Ast) {
	actionList := v.GetActionList()
	for _, action := range actionList {
		for _, cmd := range action.Cmds {
			fmt.Println(reflect.TypeOf(map[string]interface{}{}));
			fmt.Println(reflect.TypeOf(cmd.Params));
			if reflect.TypeOf(cmd.Params) == reflect.TypeOf(map[string]interface{}{}) {
				var p interface{} = cmd.Params
				sPath := b.SnippetsPath + p.(map[string]interface{})["snippet"].(string) + ".html"
				file, err := os.Open(sPath)
				if err != nil {
					fmt.Printf("ERROR: Read snippet(%s) has been fail", sPath)
					return
				}
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					str := scanner.Text()
					for key, value := range p.(map[string]interface{}) {
						str = strings.Replace(str, "$" + key, value.(string), -1)
					}
					fmt.Println(str)
				}
				defer file.Close()
			}
		}
	}
	fmt.Println(actionList);
}