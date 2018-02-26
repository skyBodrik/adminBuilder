package builders

import (
	"../visitor"
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"io"
	"text/template"
)

const MAIN_SNIPPET = "index"

const CMD_FIELD = "field"
const CMD_FIELDS_FOR_DISPLAY = "show"
const CMD_PAGE_CONFIG = "page"
const CMD_REQUEST = "request"
const CMD_INPUT = "input"
const CMD_PATHS = "paths"

const PROP_CAPTION = "caption"
const PROP_NAME = "name"
const PROP_CHARSET = "charset"
const PROP_URL = "url"
const PROP_QUERY = "query"
const PROP_SNIPPETS = "snippets"
const PROP_OUTPUT = "output"

type Builder struct {
	SnippetsPath string
	RequestUrl   string
	QueryString   string
	OutputPath   string
	PageCaption  string
	PageName     string
	Charset      string
}

func (b *Builder)Build(v visitor.Ast) {
	funcMap := template.FuncMap{"URL": func(s string) string {
		return s
	}}
	t := template.Must(template.New("array").Funcs(funcMap).ParseGlob(b.SnippetsPath + "*.html"))
	actionList := v.GetActionList()
	buffFile, err := os.Create("temp/bodyContent.tpl")
	buffFile.WriteString("{{define \"bodyContent\"}}")
	fmt.Println(err)
	for _, action := range actionList {
		b.renderAction(action, t, buffFile)
	}
	buffFile.WriteString("{{end}}")
	buffFile.Close()
	tpl, _ := template.New("array").Funcs(funcMap).ParseGlob(b.SnippetsPath + "*.html")
	t2 := template.Must(tpl.ParseFiles("temp/bodyContent.tpl"))
	outfile, err := os.Create(b.OutputPath + "/" + b.PageName + ".html")
	fmt.Println(t2.ExecuteTemplate(outfile, MAIN_SNIPPET, map[string]interface{}{
		"pageCaption": b.PageCaption,
		"pageCharset": b.Charset,
	}))
	outfile.Close()
}

func (b *Builder)renderAction(action visitor.Action, t *template.Template, writer io.Writer) {
	for _, cmd := range action.Cmds {
		switch cmd.Name {
			case CMD_PATHS:
				var p0 interface{} = cmd.Params
				for propKey, propValue := range p0.(map[string]interface{}) {
					switch propKey {
					case PROP_SNIPPETS:
						b.SnippetsPath = propValue.(string)
					case PROP_OUTPUT:
						b.OutputPath = propValue.(string)
					}
				}
			case CMD_REQUEST:
				var p0 interface{} = cmd.Params
				for propKey, propValue := range p0.(map[string]interface{}) {
					switch propKey {
						case PROP_QUERY:
							b.QueryString = propValue.(string)
						case PROP_URL:
							b.RequestUrl = propValue.(string)
					}
				}
			case CMD_PAGE_CONFIG:
				var p0 interface{} = cmd.Params
				for propKey, propValue := range p0.(map[string]interface{}) {
					switch propKey {
						case PROP_CAPTION:
							b.PageCaption = propValue.(string)
						case PROP_NAME:
							b.PageName = propValue.(string)
						case PROP_CHARSET:
							b.Charset = propValue.(string)
					}
				}
			case CMD_FIELDS_FOR_DISPLAY:
				var p0 interface{} = cmd.Params
				for _, fieldName := range p0.([]string) {
					h := md5.New()
					h.Write([]byte(CMD_FIELD + ":" + fieldName))
					hashKey := hex.EncodeToString(h.Sum(nil))
					var p1 interface{} = action.Cmds[hashKey].Params
					action.Cmds[hashKey].Params.(map[string]interface{})["fullFieldName"] = action.ActionName + "_" + p1.(map[string]interface{})["fieldName"].(string)
					fields, ok := p1.(map[string]interface{})["fields"]
					if ok && fields != nil {
						columnTitles := map[string]string{}
						columnTitlesArray := []string{}
						var p1 interface{} = fields
						for _, field := range p1.([]interface{}) {
							h := md5.New()
							h.Write([]byte(CMD_FIELD + ":" + field.(string)))
							hashKey := hex.EncodeToString(h.Sum(nil))
							var p2 interface{} = action.Cmds[hashKey].Params
							columnTitles[hashKey] = p2.(map[string]interface{})[PROP_CAPTION].(string)
							columnTitlesArray = append(columnTitlesArray, columnTitles[hashKey])
						}
						action.Cmds[hashKey].Params.(map[string]interface{})["columnTitles"] = columnTitlesArray
					}
					action.Cmds[hashKey].Params.(map[string]interface{})["initUrl"] = b.RequestUrl + action.ActionName + b.QueryString
					fmt.Println(t.ExecuteTemplate(writer, p1.(map[string]interface{})["snippet"].(string), action.Cmds[hashKey].Params))
				}
			case CMD_FIELD:
			case CMD_INPUT:
				//var p0 interface{} = cmd.Params
				//for _, fieldId := range p0.([]string) {
				//
				//}
			default:

		}
	}
}