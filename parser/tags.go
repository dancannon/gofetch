package parser

var tagActionMap map[string]interface{} = map[string]interface{}{

}

type TagAction interface {
	Start(name string)
	End(name string)
}
