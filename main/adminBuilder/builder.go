package adminBuilder

import (
	"github.com/z7zmey/php-parser/php7"
	"bytes"
	"os"
	"./visitor"
	"./builders"
)

func Run()  {
	fdInput, _ := os.Open("/root/GolandProjects/adminBuilder/main/example/Tropic.php")
	fileInfo, _ := fdInput.Stat()
	rawData := make([]byte, fileInfo.Size())
	fdInput.Read(rawData)
	src := bytes.NewBuffer(rawData)
	nodes, comments, positions := php7.Parse(src, "Tropic.php")
	visitor := visitor.Ast{
		Indent:    "",
		Comments:  comments,
		Positions: positions,
	}
	nodes.Walk(visitor)
	builder := builders.Builder{
		SnippetsPath:  "/root/GolandProjects/adminBuilder/main/adminBuilder/snippets/stories/",
	}
	builder.Build(visitor)
	//fmt.Println(visitor.GetActionList())

}