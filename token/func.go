package token

//Char returns *
func (t BoldToken) Char() rune { return '*' }

//Char returns /
func (t ItalicToken) Char() rune { return '/' }

//Char returns +
func (t HeaderToken) Char() rune { return '+' }

//Char returns \n
func (t NewLineToken) Char() rune { return '\n' }

//Char returns \t
func (t TabToken) Char() rune { return '\t' }

//Char returns -
func (t LessToken) Char() rune { return '-' }

//Char returns =
func (t EqualToken) Char() rune { return '=' }

//Char returns [
func (t SBracketOpenToken) Char() rune { return '[' }

//Char returns ]
func (t SBracketCloseToken) Char() rune { return ']' }
