package eNote

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
