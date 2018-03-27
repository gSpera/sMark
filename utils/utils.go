package eNote

import "html/template"

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
