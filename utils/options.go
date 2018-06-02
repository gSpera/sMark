package eNote

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
	TabWidth   *uint
}

//OptionsTemplate is a copy of Options but members arte not pointer
//You may use it only within a template
type OptionsTemplate struct {
	InputFile  string
	OutputFile string
	NewLine    bool
	CustomCSS  template.URL
	InlineCSS  template.CSS
	EnableFont bool
	OnlyBody   bool
	Title      string
	TabWidth   uint
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
		TabWidth:   *o.TabWidth,
	}
}

//AddString updates the Options with the passed key and his value, returns and error if the key is not valid
func (o *Options) AddString(key, sValue string) error {
	switch key {
	case "NewLine":
		fmt.Println("NewLine")
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.NewLine = &value
	case "Title":
		fmt.Println("Title")
		o.Title = &sValue
	case "InlineCSS":
		fmt.Println("InlineCSS")
		o.InlineCSS = &sValue
	case "EnableFont":
		fmt.Println("EnableFont")
		value, err := strconv.ParseBool(key)
		if err != nil {
			return errors.Wrap(err, "Could not parse value to bool")
		}
		o.EnableFont = &value
	case "OnlyBody":
		fmt.Println("OnlyBody")
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
	fmt.Println("Update")
	spew.Dump(oo)
	oR := reflect.ValueOf(o).Elem()
	ooR := reflect.ValueOf(&oo).Elem()
	structType := reflect.TypeOf(*o)

	fmt.Println(reflect.TypeOf(oR).NumField())
	fmt.Println(oR.Type(), ooR.Type())
	for i := 0; i < structType.NumField(); i++ {
		oF := oR.Field(i)
		ooF := ooR.Field(i)
		fmt.Println("Field", i, structType.Field(i).Name, oF.Elem(), ooF.Elem())

		if ooF.IsNil() {
			fmt.Println("ooF is nil")
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

	fmt.Println("Update Done")
	spew.Dump(o)
}

//EscapeHTML escapes the html string passed as parameter,
//it sanitaze more html entities than html.EscapeString
func EscapeHTML(data string) string {
	data = html.EscapeString(data)
	data = strings.NewReplacer(
		"<", " &lt;",
		">", " &gt;",
		"&", "&amp;",
		"\"", "&quot;",
		"'", "&apos;",
		"¢", "&cent;",
		"£", "&pound;",
		"¥", "&yen;",
		"€", "&euro;",
		"©", "&copy;",
		"®", "&reg;",
	).Replace(data)

	return data
}
