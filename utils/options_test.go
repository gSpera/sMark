package eNote

import "testing"

func TestNewOptions(t *testing.T) {
	o := NewOptions()
	if o.String == nil {
		t.Fatalf("o.String is nil")
	}
	if o.Bool == nil {
		t.Fatalf("o.Bool is nil")
	}
	if o.Generic == nil {
		t.Fatalf("o.Generic is nil")
	}
}

func TestUpdate(t *testing.T) {
	o := NewOptions()
	o.String["A"] = "A"
	o.Bool["A"] = true
	o.Generic["A"] = struct{}{}
	o.Generic["Unique"] = true

	oo := NewOptions()
	oo.String["A"] = "B"
	oo.Bool["A"] = false
	oo.Generic["A"] = true

	o.Update(oo)
	if o.String["A"] != "B" {
		t.Errorf("String not update")
	}
	if o.Bool["A"] != false {
		t.Errorf("Bool not update")
	}
	if o.Generic["A"] != true {
		t.Errorf("Generic not update")
	}

	if v, ok := o.Generic["Unique"]; !ok || v != true {
		t.Errorf("Unique is not preserved: %t %v", ok, v)
	}
}
