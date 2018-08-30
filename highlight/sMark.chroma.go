
package lexers

import (
    . "github.com/alecthomas/chroma" // nolint
)

// Smark lexer.
var Smark = Register(MustNewLexer(
    &Config{
        Name:      "sMark",
        Aliases:   []string{ "sm", "smark",  },
        Filenames: []string{ "*.sm",  },
        MimeTypes: []string{ "text/x-sMark",  },
    },
    Rules{
        "root": {
            { ` .*\n`, Text, nil },
            { `^\s*.+\n\s*=+\n`, GenericHeading, nil },
            { `^\s*.+\n\s*-+\n`, GenericSubheading, nil },
            { `^\|.*\n`, Literal, nil },
            { `^-+ .*\n`, Literal, nil },
            { `\[.\]`, GenericStrong, nil },
            { `\[.{2,}\]\n`, GenericOutput, Push("code") },
            { `\n-{2,}\n`, GenericEmph, nil },
            { `".+"@".+"`, LiteralStringSingle, nil },
            { `\*.*\*`, GenericStrong, nil },
            { `/.*/`, GenericEmph, nil },
            { `_.*_`, GenericEmph, nil },
            { `\-.*\-`, GenericEmph, nil },
            { `\++\n(.*\n)+\++\n`, GenericOutput, nil },
            { `.*\n`, Text, nil },
        },
        "code": {
            { `.*\n`, GenericOutput, nil },
            { `\[end\]\n`, GenericOutput, Pop(1) },
        },
    },
))

