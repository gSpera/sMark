package token

//Token is an interface for simple and complex tokens
type Token interface {
	IsToken()

	//String returns a string rapresentation of the token.
	//String is deprecated
	String() string

	//Type returns the Type of the token
	//Type is deprecated
	Type() Type
}

//TextToken is a complex token which contains text and some attribute.
type TextToken struct {
	Token

	Text string

	//Attributes
	Bold   bool
	Italic bool
}
