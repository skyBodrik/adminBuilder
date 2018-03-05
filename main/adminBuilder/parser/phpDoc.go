package parser

import (
	"regexp"
	"fmt"
	"strings"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
)

const END_OF_STR_CHARACTERS = "<strend>"

type Cmd struct {
	Name string
	Params interface{}
}

/**
 * Парсим PHPDoc
 */
func RunPhpDocParser(strPhpDoc string) map[string]Cmd {
	r0, _ := regexp.Compile("\n")
	strPhpDoc = r0.ReplaceAllString(strPhpDoc, END_OF_STR_CHARACTERS)
	r1, _ := regexp.Compile(`(?m:[*\s]+)`)
	strPhpDoc = r1.ReplaceAllString(strPhpDoc, " ")
	r2, _ := regexp.Compile(`(?m:@adminBuilder(.*?);` + END_OF_STR_CHARACTERS + `)`)
	listCmdGroups := r2.FindAllString(strPhpDoc, -1)

	r3, _ := regexp.Compile(`(?m:\s+[^\s]*\s*(.*?)[,;]` + END_OF_STR_CHARACTERS + `)`)
	r4, _ := regexp.Compile(`(?m:\s)`)
	r5, _ := regexp.Compile(`(?m:[,;]` + END_OF_STR_CHARACTERS + `)`)
	r6, _ := regexp.Compile(`(?m:\{)`)
	r7, _ := regexp.Compile(`(?m:(\w*)\:\s*([\"\d\[\{]))`)
	r8, _ := regexp.Compile(`(?m:[,\s]+)`)
	var list = map[string]Cmd{}
	for _, group := range listCmdGroups {
		//fmt.Println(group)
		listCmd := r3.FindAllString(group, -1)
		for _, cmd := range listCmd {
			cmd = r5.ReplaceAllString(cmd, "")
			cmd = strings.Trim(cmd, " ")
			parts := r4.Split(cmd, 2)
			if len(parts) >= 2 {
				partsOfParams := r6.Split(parts[1], 2)
				h := md5.New()
				//fmt.Println(parts[1])
				if len(partsOfParams) >= 2 {
					h.Write([]byte(parts[0] + ":" + partsOfParams[0]))
					hashKey := h.Sum(nil)
					partsOfParams[1] = r7.ReplaceAllString(partsOfParams[1], "\"$1\": $2")
					var params interface{}
					json.Unmarshal([]byte("{\"_fieldName\": \"" + partsOfParams[0] + "\", " + partsOfParams[1]), &params)
					list[hex.EncodeToString(hashKey)] = Cmd{
						Name:   parts[0],
						Params: params,
					}
				} else {
					h.Write([]byte(parts[0]))
					hashKey := h.Sum(nil)
					list[hex.EncodeToString(hashKey)] = Cmd{
						Name:   parts[0],
						Params: r8.Split(partsOfParams[0], -1),
					}
				}
				fmt.Println(list)
			}
		}
	}
	return list
}