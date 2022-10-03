package adapters

type Condition struct {
	Fields []string

	Equal map[string]interface{}
	Like  map[string]interface{}
	Inc   map[string]interface{}
	Lte   map[string]interface{}
	Gte   map[string]interface{}
	Not   map[string]interface{}
	Or    map[string]interface{}

	Offset int
	Limit  int
}
