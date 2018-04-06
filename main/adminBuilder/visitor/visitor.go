package visitor

import (
	"reflect"

	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/position"

	"github.com/z7zmey/php-parser/comment"
	"github.com/z7zmey/php-parser/walker"
	"github.com/z7zmey/php-parser/node/stmt"
	"../parser"
	"regexp"
)

var ActionList = []Action{}

func isWalkerImplementsNodeInterface(w walker.Walkable) bool {
	switch w.(type) {
	case node.Node:
		return true
	default:
		return false
	}
}


func actionHandler(d Ast, m *stmt.ClassMethod, w walker.Walkable) bool {
	//fmt.Println(m.MethodName)
	r0, _ := regexp.Compile(`(Action$)`)
	methodName := (m.MethodName.Attributes()["Value"]).(string)
	methodName = r0.ReplaceAllString(methodName, "")
	cmdList := parser.RunPhpDocParser(m.PhpDocComment)
	ActionList = append(ActionList, Action{methodName, cmdList, map[string]interface{}{}})
	//for i := range ActionList {
	//	fmt.Println(i)
	//}
	return false
}

func classHandler(d Ast, m *stmt.Class, w walker.Walkable) bool {
	//className := (m.ClassName.Attributes()["Value"]).(string)
	ActionList = append(ActionList, Action{"INIT", parser.RunPhpDocParser(m.PhpDocComment), map[string]interface{}{}})
	return false
}

type Action struct {
	ActionName string
	Cmds map[string]parser.Cmd
	FieldsDescription map[string]interface{}
}


// Dumper prints ast hierarchy to stdout
// Also prints comments and positions attached to nodes
type Ast struct {
	Indent    string
	Comments  comment.Comments
	Positions position.Positions
}

// EnterNode is invoked at every node in heirerchy
func (d Ast) EnterNode(w walker.Walkable) bool {
	if !isWalkerImplementsNodeInterface(w) {
		return false
	}

	n := w.(node.Node)

	if reflect.TypeOf(n) == reflect.TypeOf(&stmt.ClassMethod{}) {
		var classMethodStmt interface{} = n
		var stmt = classMethodStmt.(*stmt.ClassMethod)
		actionHandler(d, stmt, w)
	} else if reflect.TypeOf(n) == reflect.TypeOf(&stmt.Class{}) {
		var classStmt interface{} = n
		var stmt = classStmt.(*stmt.Class)
		classHandler(d, stmt, w)
	}
	//fmt.Printf("%v%v", d.Indent, reflect.TypeOf(n))
	//if p := d.Positions[n]; p != nil {
	//	fmt.Printf(" %v", *p)
	//}
	//if a := n.Attributes(); len(a) > 0 {
	//	fmt.Printf(" %v", a)
	//}
	//fmt.Println()
	//
	//if c := d.Comments[n]; len(c) > 0 {
	//	fmt.Printf("%vComments:\n", d.Indent+"  ")
	//	for _, cc := range c {
	//		fmt.Printf("%v%q\n", d.Indent+"    ", cc)
	//	}
	//}

	return true
}

// GetChildrenVisitor is invoked at every node parameter that contains children nodes
func (d Ast) GetChildrenVisitor(key string) walker.Visitor {
	//fmt.Printf("%v%q:\n", d.Indent+"  ", key)
	return Ast{d.Indent + "    ", d.Comments, d.Positions}
}

// LeaveNode is invoked after node process
func (d Ast) LeaveNode(n walker.Walkable) {
	// do nothing
}

func (d Ast) GetActionList() []Action {
	return ActionList
}