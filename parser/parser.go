package parser

type Parser struct{}

type ParserInterface interface {
	DecodeCLIXML(clixml string) (string, error)
}

func NewParser() *Parser {
	return &Parser{}
}
