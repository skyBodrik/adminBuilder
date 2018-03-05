package builders

import (
	"../visitor"
	"fmt"
	"os"
	"crypto/md5"
	"encoding/hex"
	"io"
	"text/template"
	"log"
	"io/ioutil"
	"github.com/s-yata/go-iconv"
	"encoding/json"
	"bytes"
	"reflect"
	"strings"
	"strconv"
)

const TEMP_FOLDER = "temp/"

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
	UsedSnippets map[string]string
	ObjectIdCounters map[string]int
}

func (b *Builder)Build(v visitor.Ast) {
	b.UsedSnippets = make(map[string]string, 1)
	b.ObjectIdCounters = make(map[string]int, 1)
	funcMap := template.FuncMap{
		"URL": func(s string) string {
			return s
		},
		"JSON": func(data interface{}) string {
			json, _ := json.Marshal(data)
			fmt.Println(string(json))
			return string(json)
		},
		"JSEscape": func(s string) string {
			return template.JSEscapeString(s)
		},
	}
	t := template.Must(template.New("array").Funcs(funcMap).ParseGlob(b.SnippetsPath + "*.html"))
	actionList := v.GetActionList()
	buffFile, err := os.Create(TEMP_FOLDER + "bodyContent.tpl")
	buffFile.WriteString("{{define \"bodyContent\"}}")
	fmt.Println(err)
	for _, action := range actionList {
		b.renderAction(action, t, buffFile)
	}
	buffFile.WriteString("{{end}}")
	buffFile.Close()
	tpl, _ := template.New("array").Funcs(funcMap).ParseGlob(b.SnippetsPath + "*.html")
	t2 := template.Must(tpl.ParseFiles(TEMP_FOLDER + "bodyContent.tpl"))
	outputFileName := b.OutputPath + "/" + b.PageName + ".html"
	tempFileName := TEMP_FOLDER + "/" + b.PageName + ".html"
	outfileTemp, _ := os.Create(tempFileName)
	fmt.Println(b.UsedSnippets["textbox"])
	fmt.Println(t2.ExecuteTemplate(outfileTemp, MAIN_SNIPPET, map[string]interface{}{
		"pageCaption": b.PageCaption,
		"pageCharset": b.Charset,
		"usedSnippets": b.UsedSnippets,
	}))
	defer outfileTemp.Close()
	h, err := iconv.Open(b.Charset, "UTF-8")
	if err != nil {
		log.Fatalf("iconv.Open failed: %v", err)
	}
	defer h.Close()
	outfile, _ := os.Create(outputFileName)
	data, _ := ioutil.ReadFile(tempFileName)
	newData, _ := h.Conv(data)
	outfile.Write(newData)
	defer outfile.Close()
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
					var fullName = action.ActionName + "_" + p1.(map[string]interface{})["_fieldName"].(string)
					_, err := b.ObjectIdCounters[fullName]
					if err {
						b.ObjectIdCounters[fullName] = 0
					} else {
						b.ObjectIdCounters[fullName]++
					}
					action.Cmds[hashKey].Params.(map[string]interface{})["_fullFieldName"] = fullName + "-" + strconv.Itoa(b.ObjectIdCounters[fullName]) + ""
						fields, ok := p1.(map[string]interface{})["fields"]
					if ok && fields != nil {
						columnTitles := map[string]string{}
						columnTitlesArray := []string{}
						fieldsDescription := map[string]interface{}{}
						//fields = append(fields.([]interface{}), fieldName)
						var p2 interface{} = fields
						fmt.Println(fields)
						//fields = append(fields.([]string), "test")
						for _, field := range p2.([]interface{}) {
							if reflect.TypeOf(field).String() == "[]interface {}" {
								var titles []string
								for _, subField := range field.([]interface{}) {
									h := md5.New()
									h.Write([]byte(CMD_FIELD + ":" + subField.(string)))
									hashKey := hex.EncodeToString(h.Sum(nil))
									var p3 interface{} = action.Cmds[hashKey].Params
									columnTitles[hashKey] = p3.(map[string]interface{})[PROP_CAPTION].(string)
									titles = append(titles, columnTitles[hashKey])
									fieldsDescription[subField.(string)] = p3.(map[string]interface{})
								}
								columnTitlesArray = append(columnTitlesArray, strings.Join(titles, "/"))
							} else {
								h := md5.New()
								h.Write([]byte(CMD_FIELD + ":" + field.(string)))
								hashKey := hex.EncodeToString(h.Sum(nil))
								var p3 interface{} = action.Cmds[hashKey].Params
								columnTitles[hashKey] = p3.(map[string]interface{})[PROP_CAPTION].(string)
								columnTitlesArray = append(columnTitlesArray, columnTitles[hashKey])
								fieldsDescription[field.(string)] = p3.(map[string]interface{})
							}
						}
						// *** Получение параметров для родительского поля
						currentFieldDescription := map[string]interface{}{}
						for key, value := range p1.(map[string]interface{}) {
							currentFieldDescription[key] = value
						}
						fieldsDescription[fieldName] = currentFieldDescription
						// ***
						action.Cmds[hashKey].Params.(map[string]interface{})["_columnTitles"] = columnTitlesArray
						action.Cmds[hashKey].Params.(map[string]interface{})["_fieldsDescription"] = fieldsDescription
						action.Cmds[hashKey].Params.(map[string]interface{})["_actionName"] = action.ActionName
						action.Cmds[hashKey].Params.(map[string]interface{})["_id"] = action.Cmds[hashKey].Params.(map[string]interface{})["_fullFieldName"]
					}
					action.Cmds[hashKey].Params.(map[string]interface{})["_initUrl"] = b.RequestUrl + action.ActionName + b.QueryString
					action.Cmds[hashKey].Params.(map[string]interface{})["_isAction"] = true
					fmt.Println(t.ExecuteTemplate(writer, p1.(map[string]interface{})["snippet"].(string), action.Cmds[hashKey].Params))
					//fmt.Println(action.Cmds[hashKey].Params)
				}
			case CMD_FIELD:
				var p0 interface{} = cmd.Params
				buffer := bytes.Buffer{}
				p0.(map[string]interface{})["_isAction"] = false
				p0.(map[string]interface{})["_fullFieldName"] = action.ActionName + "_" + p0.(map[string]interface{})["_fieldName"].(string)
				p0.(map[string]interface{})["_value"] = "{{ .value }}"
				p0.(map[string]interface{})["_id"] = "{{ .id }}"
				t.ExecuteTemplate(&buffer, p0.(map[string]interface{})["snippet"].(string), cmd.Params)
				b.UsedSnippets[(p0.(map[string]interface{})["snippet"]).(string)] = buffer.String()
				//fmt.Println(buffer.String())
		case CMD_INPUT:
				//var p0 interface{} = cmd.Params
				//for _, fieldId := range p0.([]string) {
				//
				//}
			default:

		}
	}
}