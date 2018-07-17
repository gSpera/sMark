package eNote

import (
	"html/template"
	"log"
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

//Options are the command line parameters
type Options struct {
	InputFile    *string
	OutputFile   *string
	NewLine      *bool
	CustomCSS    *string
	InlineCSS    *string
	EnableFont   *bool
	OnlyBody     *bool
	Title        *string
	TabWidth     *uint
	Verbose      *bool
	HTMLOut      *bool
	TelegraphOut *bool
	Prettify     *string
	Watch        *bool
}

//OptionsTemplate is a copy of Options but members are not pointer
//You may use it only within a template
type OptionsTemplate struct {
	InputFile    string
	OutputFile   string
	NewLine      bool
	CustomCSS    template.URL
	InlineCSS    template.CSS
	EnableFont   bool
	OnlyBody     bool
	Title        string
	TabWidth     uint
	Verbose      bool
	HTMLOut      bool
	TelegraphOut bool
	Prettify     string
	Watch        bool
}

//ToTemplate returtns a copy of Options with out pointers
func (o Options) ToTemplate() OptionsTemplate {
	return OptionsTemplate{
		InputFile:    *o.InputFile,
		OutputFile:   *o.OutputFile,
		NewLine:      *o.NewLine,
		CustomCSS:    template.URL(*o.CustomCSS),
		InlineCSS:    template.CSS(*o.InlineCSS),
		EnableFont:   *o.EnableFont,
		OnlyBody:     *o.OnlyBody,
		Title:        *o.Title,
		TabWidth:     *o.TabWidth,
		Verbose:      *o.Verbose,
		HTMLOut:      *o.HTMLOut,
		TelegraphOut: *o.TelegraphOut,
		Prettify:     *o.Prettify,
		Watch:        *o.Watch,
	}
}

//AddString updates the Options with the passed key and his value, returns and error if the key is not valid
func (o *Options) AddString(key, sValue string) error {
	switch key {
	case "NewLine":
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.NewLine = &value
	case "Title":
		o.Title = &sValue
	case "InlineCSS":
		o.InlineCSS = &sValue
	case "EnableFont":
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.EnableFont = &value
	case "OnlyBody":
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.OnlyBody = &value
	default:
		return errors.New("Key is not valid")
	}

	return nil
}

//Update adds non specified value from the passed Options
func (o *Options) Update(oo Options) {
	oR := reflect.ValueOf(o).Elem()
	ooR := reflect.ValueOf(&oo).Elem()
	structType := reflect.TypeOf(*o)

	for i := 0; i < structType.NumField(); i++ {
		oF := oR.Field(i)
		ooF := ooR.Field(i)

		if ooF.IsNil() {
			continue
		}

		if !oF.CanAddr() {
			log.Fatal("Cannot take address of Options")
		}

		if !ooF.CanAddr() {
			log.Fatal("Cannot take address of OOptions")
		}

		log.Printf("%v = %v", structType.Field(i).Name, ooF.Elem())
		oF.Set(ooF)
	}
}
