package main


import (
  //"github.com/davecgh/go-spew/spew"
  "github.com/hashicorp/hcl"
  "github.com/hashicorp/hcl/hcl/ast"
	"io/ioutil"
  "fmt"
)


var (

  file = "test.hcl"
)


func main() {

  hs := &hclStrings{}

  b, err := ioutil.ReadFile(file)

  if err != nil {
    fmt.Println(err)
  }

  parsedAst, err := hcl.ParseBytes(b)
  
  if err != nil {
    fmt.Println(err)
  }

  var node *ast.ObjectList 

  node = parsedAst.Node.(*ast.ObjectList)

  get(node, hs)


  for _, items := range node.Items {

    ttt := items.Val

    get(ttt, hs)

  }

  fmt.Println(hs.Strings)

}

func get(x ast.Node, s *HCLTokens) {


  obj, ok := x.(*ast.ObjectType)

  if ! ok {
    fmt.Println("Not object type")
  } else {
    for _, item := range obj.List.Items {
      get(item, s)
    }
    return
  }

  list, ok := x.(*ast.ObjectList)

  if ! ok {
    fmt.Println("Not Object List")

    fmt.Println(spew.Sdump(list))
  } else {
    for _, item := range list.Items {
      get(item, s)
    }
    return
  }


  listType, ok := x.(*ast.ListType)

  if ! ok {
    fmt.Println("Not List type")

  } else {
    for _, item := range listType.List {
      get(item, s)
    }
    return
  }

  objItem, ok := x.(*ast.ObjectItem)

  if ! ok {
    fmt.Println("Not object item")
  } else {
    value, ok := objItem.Val.(*ast.LiteralType)

    if ! ok {
      fmt.Println("Unable to convert to LiteralType.")
      return
    }

    token := value.Token
    s.AddString(token.Text, token.Pos)
    return
  }

  return
}
 

