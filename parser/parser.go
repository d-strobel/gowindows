package parser

type Parser struct {
	DecodeCLIXML func(xmlErr string) (string, error)
}

func New() Parser {
	return Parser{
		DecodeCLIXML: DecodeCLIXML,
	}
}
