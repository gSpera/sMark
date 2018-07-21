package token

//LineContainerFromString generate a LineContainer from a given text,
//it parses the tab and add them.
func LineContainerFromString(txt string) LineContainer {
	tokens := []Token{}
	indentation := 0

	for _, ch := range txt {
		if ch == '\t' {
			indentation++
		} else {
			break
		}
	}

	//Remove starting tabs
	tokens = append(tokens, TextToken{Text: txt[indentation:]})
	return LineContainer{
		Tokens:      tokens,
		Indentation: indentation,
	}
}
