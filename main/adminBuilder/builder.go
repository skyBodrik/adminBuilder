package adminBuilder

import (
	"github.com/z7zmey/php-parser/php7"
	"bytes"
	"os"
	"./visitor"
	"./builders"
	"io/ioutil"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/charmap"
	"fmt"
	"flag"
)

func Run()  {
	var encode, output, input string
	flag.StringVar(&encode, "encode", "utf-8", "a string var")
	flag.StringVar(&output, "output", "output", "a string var")
	flag.StringVar(&input, "input", "input", "a string var")
	flag.Parse()
	fmt.Println("Source: " + input)
	fmt.Println("Destination: " + output)
	fdInput, _ := os.Open(input)//"/root/GolandProjects/adminBuilder/main/example/Tropic.php")
	fileInfo, _ := fdInput.Stat()
	var rawData []byte
	var err error
	switch encode {
		case "window-1251":
			tr := transform.NewReader(fdInput, charmap.Windows1251.NewDecoder())
			rawData, err = ioutil.ReadAll(tr)
		default:
			rawData, err = ioutil.ReadAll(fdInput)
	}
	if err != err {
		fmt.Printf("Reading error: %s", fileInfo.Name())
	}
	src := bytes.NewBuffer(rawData)
	nodes, comments, positions := php7.Parse(src, "Tropic.php")
	visitor := visitor.Ast{
		Indent:    "",
		Comments:  comments,
		Positions: positions,
	}
	nodes.Walk(visitor)
	builder := builders.Builder{
		SnippetsPath:  "main/adminBuilder/snippets/*/",
		OutputPath: output,
		Charset: "utf-8",
	}
	builder.Build(visitor)
	fmt.Println("Done")
	//fmt.Println(visitor.GetActionList())

}