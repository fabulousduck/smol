package losp

type node struct {
	t    string
	body string
}

type parser struct {
	ast []node
}

func NewParser() *parser {
	return new(parser)
}

// func (p *parser) parse(tokens []token) {
// 	for i := 0; i < len(tokens); i++ {
// 		node := node{}

// 	}
// }
