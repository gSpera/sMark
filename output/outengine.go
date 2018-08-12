package output

import (
	"github.com/gSpera/sMark/token"
	sMark "github.com/gSpera/sMark/utils"
)

//Engine is an interface implemented by outputs engine
//Output takes as parameters a slice of paragraphs and the final options
//Output generates a slice of bytes containing the generated data
//This could be a html document for example
type Engine func([]token.ParagraphToken, sMark.Options) ([]byte, error)
