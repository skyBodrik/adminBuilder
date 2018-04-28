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
	"regexp"
)

// Служебные константы
const TEMP_FOLDER = "temp"

const DEFAULT_SECTION_NAME = "Без группировки"
const CUSTOM_SECTION_NAME = "Ссылки"

// Сниппеты
const MAIN_SNIPPET = "index"

// Имена комманд
const CMD_FIELD = "field"
const CMD_FIELDS_FOR_DISPLAY = "show"
const CMD_PAGE_CONFIG = "page"
const CMD_ACTION = "action"
const CMD_REQUEST = "request"
const CMD_INPUT = "input"
const CMD_PATHS = "paths"
const CMD_SNIPPETS = "snippets"
const CMD_MENU = "menu"

// Имена свойств
const PROP_CAPTION = "caption"
const PROP_TITLE = "title"
const PROP_NAME = "name"
const PROP_CHARSET = "charset"
const PROP_URL = "url"
const PROP_QUERY = "query"
const PROP_SNIPPETS = "snippets"
const PROP_OUTPUT = "output"
const PROP_SECTION = "section"
const PROP_SET = "set"

// Имена типов
const TYPE_ARRAY = "array"

type Builder struct {
	SnippetsPath     string
	SnippetsSet      string
	RequestUrl       string
	QueryString      string
	OutputPath       string
	PageCaption      string
	PageName         string
	Charset          string
	TemplatesList    map[string]string
	ObjectIdCounters map[string]int
	MenuSections     map[string]map[string]string
}

/**
 * Строим админку по AST исходного файла
 */
func (b *Builder)Build(v visitor.Ast) {
	// Инициализация
	b.TemplatesList = make(map[string]string, 1)
	b.ObjectIdCounters = make(map[string]int, 1)
	b.MenuSections = make(map[string]map[string]string, 1)
	actionList := v.GetActionList()
	bodyContent := bytes.Buffer{}
	pageContent := bytes.Buffer{}
	outputFileName := b.OutputPath //+ "/" + b.PageName + ".html"

	// Функции доступные в шаблоне
	funcMap := template.FuncMap{
		"URL": func(s string) string {
			return s
		},
		"CHOOSE": func(args ...interface{}) string {
			for _, value := range args {
				if value != nil && value != "" {
					return value.(string)
				}
			}
			return ""
		},
		"NOT": func(s bool) bool {
			return !s
		},
		"ifExists": func(s interface{}) string {
			if s == nil {
				return ""
			}
			return s.(string)
		},
		"JSON": func(data interface{}) string {
			strJson, _ := json.Marshal(data)
			return string(strJson)
		},
		"JSEscape": func(s string) string {
			return template.JSEscapeString(s)
		},
		"checkType": func(obj interface{}, typeName string) bool {
			return reflect.TypeOf(obj).String() == typeName
		},
	}

	// Собираем боди контент, рендерим экшены
	tpl := template.Must(template.New("array").Funcs(funcMap).ParseGlob(b.SnippetsPath + "*.html"))
	bodyContent.WriteString("{{define \"bodyContent\"}}")
	for _, action := range actionList {
		b.renderAction(action, tpl, &bodyContent)
	}
	bodyContent.WriteString("{{end}}")

	// Теперь собираем всю страницу
	tpl = template.Must(tpl.Parse(bodyContent.String()))
	params := map[string]interface{}{
		"pageCaption": b.PageCaption,
		"pageCharset": b.Charset,
		"templatesList": b.TemplatesList,
		"menuSections": b.MenuSections,
	}
	pageContent.WriteString("{{define \"pageContent\"}}")
	fmt.Println(tpl.ExecuteTemplate(&pageContent, b.SnippetsSet + "/" + MAIN_SNIPPET, params))
	pageContent.WriteString("{{end}}")

	tempFileName := TEMP_FOLDER + "/" + b.PageName + ".html"
	// Теперь собираем всю страницу, сохраняем во временный файл
	outfileTemp, _ := os.Create(tempFileName)
	tpl = template.Must(tpl.Parse(pageContent.String()))
	fmt.Println(tpl.ExecuteTemplate(outfileTemp, "core/main", params))
	//fmt.Println(tempFileName)
	defer outfileTemp.Close()

	// Кодируем как нам надо и сохраняем в целевой файл
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

/**
 * Рендерим экшены
 */
func (b *Builder)renderAction(action visitor.Action, t *template.Template, writer io.Writer) bool {
	actionTitle, actionSection, actionNeedRender := "", DEFAULT_SECTION_NAME, false
	// Извлекаем базовую инфу из экшена
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
		case CMD_SNIPPETS:
			var p0 interface{} = cmd.Params
			for propKey, propValue := range p0.(map[string]interface{}) {
				switch propKey {
				case PROP_SET:
					b.SnippetsSet = propValue.(string)
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
		case CMD_FIELD:
			var p0 interface{} = cmd.Params
			var fullFieldName = action.ActionName + "_" + p0.(map[string]interface{})["_fieldName"].(string)
			buffer := bytes.Buffer{}
			p0.(map[string]interface{})["_initUrl"] = b.RequestUrl + action.ActionName + b.QueryString
			p0.(map[string]interface{})["_withLoadDataScript"] = false
			//p0.(map[string]interface{})["_isAction"] = false
			p0.(map[string]interface{})["_fullFieldName"] = fullFieldName
			p0.(map[string]interface{})["_actionName"] = action.ActionName
			p0.(map[string]interface{})["_value"] = "{{ ._value_" + fullFieldName + " }}"
			p0.(map[string]interface{})["_id"] = "{{ ._id_" + fullFieldName + " }}"
			t.ExecuteTemplate(&buffer, b.SnippetsSet + "/" + p0.(map[string]interface{})["snippet"].(string), cmd.Params)
			b.TemplatesList[(p0.(map[string]interface{})["_fullFieldName"]).(string)] = buffer.String()
			p0.(map[string]interface{})["_withLoadDataScript"] = true
			buffer = bytes.Buffer{}
			t.ExecuteTemplate(&buffer, b.SnippetsSet + "/" + p0.(map[string]interface{})["snippet"].(string), cmd.Params)
			b.TemplatesList["withLoadDataScript:" + fullFieldName] = buffer.String()
			action.FieldsDescription[p0.(map[string]interface{})["_fieldName"].(string)] = p0.(map[string]interface{})
			//fmt.Println(buffer.String())
		case CMD_ACTION:
			var p0 interface{} = cmd.Params
			for propKey, propValue := range p0.(map[string]interface{}) {
				switch propKey {
				case PROP_TITLE:
					actionTitle = propValue.(string)
				case PROP_SECTION:
					actionSection = propValue.(string)
				}
			}
			sectionData, _ := b.MenuSections[actionSection]
			if len(sectionData) == 0 {
				b.MenuSections[actionSection] = make(map[string]string, 1)
			}
			b.MenuSections[actionSection]["#" + action.ActionName] = actionTitle
		case CMD_MENU:
			var p0 interface{} = cmd.Params
			caption, url, section := "", "#", CUSTOM_SECTION_NAME
			for sectionCaption, sectionContent := range p0.(map[string]interface{}) {
				r0, _ := regexp.Compile(`(?m:^_+)`)
				if r0.Match([]byte(sectionCaption)) {
					continue
				}
				if reflect.TypeOf(sectionContent).String() == "[]interface {}" {
					section = sectionCaption
					sectionData, _ := b.MenuSections[section]
					if len(sectionData) == 0 {
						b.MenuSections[section] = make(map[string]string, 1)
					}
					for _, itemContent := range sectionContent.([]interface{}) {
						for propKey, propValue := range itemContent.(map[string]interface{}) {
							switch propKey {
							case PROP_CAPTION:
								caption = propValue.(string)
							case PROP_URL:
								url = propValue.(string)
							}
						}
						b.MenuSections[section][url] = caption
					}
				}
			}
		case CMD_FIELDS_FOR_DISPLAY:
			actionNeedRender = true
		case CMD_INPUT:
			//var p0 interface{} = cmd.Params
			//for _, fieldId := range p0.([]string) {
			//
			//}
		default:
		}
	}
	if !actionNeedRender {
		return false
	}
	// Компилируем предварительный HTML код экшена
	fmt.Println(t.ExecuteTemplate(writer, "core/action", map[string]interface{}{
		"_actionName": action.ActionName,
		"_actionTitle": actionTitle,
		"_beforeRender": true,
	}))
	// Формируем инфу об отображаемых в админке полях
	for _, cmd := range action.Cmds {
		if cmd.Name == CMD_FIELDS_FOR_DISPLAY {
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
				action.Cmds[hashKey].Params.(map[string]interface{})["_id"] = action.Cmds[hashKey].Params.(map[string]interface{})["_fullFieldName"]
				fields, ok := p1.(map[string]interface{})["fields"]
				if ok && fields != nil {
					columnTitles := map[string]string{}
					columnTitlesArray := []string{}
					//fields = append(fields.([]interface{}), fieldName)
					var p2 interface{} = fields
					for _, field := range p2.([]interface{}) {
						if reflect.TypeOf(field).String() == "[]interface {}" {
							var titles []string
							for _, subField := range field.([]interface{}) {
								h := md5.New()
								h.Write([]byte(CMD_FIELD + ":" + subField.(string)))
								hashKey2 := hex.EncodeToString(h.Sum(nil))
								var p3 interface{} = action.Cmds[hashKey2].Params
								columnTitles[hashKey2] = p3.(map[string]interface{})[PROP_CAPTION].(string)
								titles = append(titles, columnTitles[hashKey2])
								p3.(map[string]interface{})["_withLoadDataScript"] = action.Cmds[hashKey].Params.(map[string]interface{})["type"] != TYPE_ARRAY
							}
							columnTitlesArray = append(columnTitlesArray, strings.Join(titles, " / "))
						} else {
							h := md5.New()
							h.Write([]byte(CMD_FIELD + ":" + field.(string)))
							hashKey2 := hex.EncodeToString(h.Sum(nil))
							var p3 interface{} = action.Cmds[hashKey2].Params
							columnTitles[hashKey2] = p3.(map[string]interface{})[PROP_CAPTION].(string)
							columnTitlesArray = append(columnTitlesArray, columnTitles[hashKey2])
							p3.(map[string]interface{})["_withLoadDataScript"] = action.Cmds[hashKey].Params.(map[string]interface{})["type"] != TYPE_ARRAY
						}
					}
					action.Cmds[hashKey].Params.(map[string]interface{})["_columnTitles"] = columnTitlesArray
					action.Cmds[hashKey].Params.(map[string]interface{})["_actionName"] = action.ActionName
					action.Cmds[hashKey].Params.(map[string]interface{})["_id"] = action.Cmds[hashKey].Params.(map[string]interface{})["_fullFieldName"]
					action.Cmds[hashKey].Params.(map[string]interface{})["_withLoadDataScript"] = true
				}
				action.Cmds[hashKey].Params.(map[string]interface{})["_initUrl"] = b.RequestUrl + action.ActionName + b.QueryString

				// *** Компоновка параметров, клонирование params
				paramsBuffer := map[string]interface{}{}
				for key, value := range action.Cmds[hashKey].Params.(map[string]interface{}) {
					paramsBuffer[key] = value
				}
				paramsBuffer["_fieldsDescription"] = action.FieldsDescription
				// ***

				//action.Cmds[hashKey].Params.(map[string]interface{})["_isAction"] = true
				writer.Write([]byte("<div class=\"form-group\">"))
				fmt.Println(t.ExecuteTemplate(writer, b.SnippetsSet + "/" + p1.(map[string]interface{})["snippet"].(string), paramsBuffer))
				writer.Write([]byte("</div>"))
				//fmt.Println(action.Cmds[hashKey].Params)
			}
		}
	}
	// Рендерим окончательный HTML код для экшена
	fmt.Println(t.ExecuteTemplate(writer, "core/action", map[string]interface{}{
		"_actionName": action.ActionName,
		"_beforeRender": false,
	}))
	return true
}