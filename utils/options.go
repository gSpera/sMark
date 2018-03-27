package eNote

import (
	"fmt"
	"html/template"
	"strconv"

	"github.com/pkg/errors"
)

//Options are the command line parameters
type Options struct {
	InputFile  *string
	OutputFile *string
	NewLine    *bool
	CustomCSS  *string
	InlineCSS  *string
	EnableFont *bool
	OnlyBody   *bool
	Title      *string
}

//OptionsTemplate is a copy of Options but members arte not pointer
type OptionsTemplate struct {
	InputFile  string
	OutputFile string
	NewLine    bool
	CustomCSS  template.URL
	InlineCSS  template.CSS
	EnableFont bool
	OnlyBody   bool
	Title      string
}

//ToTemplate returtns a copy of Options with out pointers
func (o Options) ToTemplate() OptionsTemplate {
	return OptionsTemplate{
		InputFile:  *o.InputFile,
		OutputFile: *o.OutputFile,
		NewLine:    *o.NewLine,
		CustomCSS:  template.URL(*o.CustomCSS),
		InlineCSS:  template.CSS(*o.InlineCSS),
		EnableFont: *o.EnableFont,
		OnlyBody:   *o.OnlyBody,
		Title:      *o.Title,
	}
}

//AddString updates the OptionsTemplate with the passed key and his value, returns and error if the key is not valid
func (o *OptionsTemplate) AddString(key, sValue string) error {
	switch key {
	case "NewLine":
		fmt.Println("NewLine")
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.NewLine = value
	case "Title":
		fmt.Println("Title")
		o.Title = sValue
	case "InlineCSS":
		fmt.Println("InlineCSS")
		o.InlineCSS = template.CSS(sValue)
	case "EnableFont":
		fmt.Println("EnableFont")
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.EnableFont = value
	case "OnlyBody":
		fmt.Println("OnlyBody")
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.OnlyBody = value
	default:
		return errors.New("Key is not valid")
	}

	return nil
}

//Update adds non specified value from the passed OptionsTemplate
func (o *Options) Update(oo OptionsTemplate) {
	if *o.Title == "" {
		o.Title = &oo.Title
	}
	if *o.InlineCSS == "" {
		str := string(oo.InlineCSS)
		o.InlineCSS = &str
	}
}
