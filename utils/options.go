package sMark

import "strings"

//Options are the parameters to the compilation.
//They can be obtained throught command line parameters or HeaderParagraphs
type Options struct {
	String  map[string]string
	Bool    map[string]bool
	Generic map[string]interface{}
}

//NewOptions create a Options struct and initialize all maps
func NewOptions() Options {
	o := Options{}
	o.String = map[string]string{}
	o.Bool = map[string]bool{}
	o.Generic = map[string]interface{}{}
	return o
}

//Update adds non specified value from the passed Options
func (o *Options) Update(oo Options) {
	for k, v := range oo.String {
		o.String[k] = v
	}
	for k, v := range oo.Bool {
		o.Bool[k] = v
	}
	for k, v := range oo.Generic {
		o.Generic[k] = v
	}
}

//Insert inserts a single value in the Option parsing it
//T, TRUE, 1, ON rapresent a Bool True
//F, FALSE, 0, OFF rapresent a Bool False
//All the other value are inserted as String
//If you need to insert a Generic Option you will need to insert it manually
func (o *Options) Insert(key, value string) {
	switch strings.ToUpper(value) {
	case "T", "TRUE", "1", "ON":
		o.Bool[key] = true
	case "F", "FALSE", "0", "OFF":
		o.Bool[key] = false
	default:
		o.String[key] = value
	}
}

//ParseOption parses an Option String returning the Key and Value
func ParseOption(line string) (string, string) {
	key := ""
	buffer := []rune{}

	for i, ch := range line {
		switch ch {
		case '=':
			key = string(buffer)
			buffer = []rune{}
		case ';': //Comment
			if line[i-1] == ' ' {
				return key, string(buffer)
			}
			buffer = append(buffer, ch)
		case ' ':
			//Ignore if last character
			if len(line) <= i+1 {
				continue
			}

			if line[i-1] == '=' || line[i+1] == '=' { //Space are allowed around equal sign
				continue
			}

			if len(line) > i && line[i+1] == ';' { //Comment incoming
				continue
			}

			//Else insert it
			fallthrough
		default:
			buffer = append(buffer, ch)
		}
	}

	return key, string(buffer)
}
