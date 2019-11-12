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


func get(x interface{}, s *hclStrings) {


  obj, ok := x.(*ast.ObjectType)

  if ! ok {
    fmt.Println("Not object type")
  } else {
    for _, item := range obj.List.Items {
      get(item, s)
    }
  }

  objItem, ok := x.(*ast.ObjectItem)

  if ! ok {
    fmt.Println("Not object item")
  } else {
    value := objItem.Val.(*ast.LiteralType)
    token := value.Token
    s.AddString(token.Text)
  }
}



