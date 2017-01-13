package mparse

import (
	"testing"
)

func TestOnlyPointer(t *testing.T) {
	msg := ""

	type test struct{}

	v := test{}

	err := Parse(msg, v)
	if err == nil {
		t.Fail()
	}
}

func TestOnlyStruct(t *testing.T) {
	msg := ""

	v := 2

	err := Parse(msg, &v)
	if err == nil {
		t.Fail()
	}
}

func TestSimpleStruct(t *testing.T) {
	msg := "/main test"

	type test struct {
		Main string
	}

	v := &test{}

	err := Parse(msg, v)
	if err != nil {
		t.Error(err)
	}
	if v.Main != "test" {
		t.Errorf("Main field should be [%s], is [%s]", "test", v.Main)
	}
}

func TestUnexportedField(t *testing.T) {
	msg := "/main test"

	type test struct {
		main string
	}

	v := &test{}

	err := Parse(msg, v)
	if err != nil {
		t.Error(err)
	}
	if v.main != "" {
		t.Errorf("main field should be [%s], is [%s]", "", v.main)
	}
}

func TestFieldTag(t *testing.T) {
	msg := "/main test"

	type test struct {
		Main    string `mparse:"notMain"`
		NotMain string `mparse:"main"`
	}

	v := &test{}

	err := Parse(msg, v)
	if err != nil {
		t.Error(err)
	}
	if v.Main != "" {
		t.Errorf("Main field should be [%s], is [%s]", "", v.Main)
	}
	if v.NotMain != "test" {
		t.Errorf("NotMain field should be [%s], is [%s]", "test", v.NotMain)
	}
}

func TestDefault(t *testing.T) {
	msg := "/main test"

	type test struct {
		Main string `mparse:"test"`
		Abc string `mparse:"default"`
	}

	v := &test{}

	err := Parse(msg, v)
	if err != nil {
		t.Error(err)
	}
	if v.Main != "" {
		t.Errorf("Main field should be [%s], is [%s]", "", v.Main)
	}
	if v.Abc != msg {
		t.Errorf("Abc field should be [%s], is [%s]", msg, v.Abc)
	}
}

func TestNestedStruct(t *testing.T) {
	msg := `/a true
	/b true
	/c 12312412
	/d test1
	/d test2
	msg`
	d := "/d test1\n/d test2\nmsg"

	type inter struct {
		A string
		B bool
	}

	type test struct {
		In inter `mparse:"some text"`
		C  int64
		D  string `mparse:"default"`
		A  bool   `mparse:"abc"`
	}

	v := &test{}
	err := Parse(msg, v)
	if err != nil {
		t.Error(err)
	}
	if v.In.A != "true" {
		t.Errorf("In.A field should be [%s], is [%s]", "true", v.In.A)
	}
	if v.In.B != true {
		t.Errorf("In.B field should be [%b], is [%b]", true, v.In.B)
	}
	if v.C != 12312412 {
		t.Errorf("C field should be [%d], is [%d]",12312412, v.C)
	}
	if v.D != d {
		t.Errorf("D field should be [%s], is [%s]",d, v.D)
	}
}
