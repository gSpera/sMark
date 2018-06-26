package token

//Token is an interface for simple and complex tokens
type Token interface {
	IsToken()

	//Type returns the Type of the token
	//It is used for recognizing the Token
	Type() Type
}

//TextToken is a complex token which contains text and some attribute.
type TextToken struct {
	Token

	Indentation int
	Text        string

	//Attributes
	Bold   bool
	Italic bool
}
