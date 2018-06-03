package token

//LineContainerFromString generate a LineContainer from a given text,
//it parses the tab and add them.
func LineContainerFromString(txt string) LineContainer {
	tokens := []Token{}

	for _, ch := range txt {
		if ch == '\t' {
			tokens = append(tokens, TabToken{})
		} else {
			break
		}
	}

	//Remove starting tabs
	tokens = append(tokens, TextToken{Text: txt[len(tokens):]})
	return LineContainer{
		Tokens: tokens,
	}
}
