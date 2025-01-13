package validator

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/fr"
	"github.com/go-playground/locales/nl"
	ut "github.com/go-playground/universal-translator"
)

// NOTES:
// - Run "go test" to run tests
// - Run "gocov test | gocov report" to report on test coverage by file
// - Run "gocov test | gocov annotate -" to report on all code and functions, those, marked with "MISS" were never called
//
// or
//
// -- may be a good idea to change to output path to somewherelike /tmp
// go test -coverprofile cover.out && go tool cover -html=cover.out -o cover.html
//
//
// go test -cpuprofile cpu.out
// ./validator.test -test.bench=. -test.cpuprofile=cpu.prof
// go tool pprof validator.test cpu.prof
//
//
// go test -memprofile mem.out

type I interface {
	Foo() string
}

type Impl struct {
	F string `validate:"len=3"`
}

func (i *Impl) Foo() string {
	return i.F
}

type SubTest struct {
	Test string `validate:"required"`
}

type TestInterface struct {
	Iface I
}

type TestString struct {
	BlankTag  string `validate:""`
	Required  string `validate:"required"`
	Len       string `validate:"len=10"`
	Min       string `validate:"min=1"`
	Max       string `validate:"max=10"`
	MinMax    string `validate:"min=1,max=10"`
	Lt        string `validate:"lt=10"`
	Lte       string `validate:"lte=10"`
	Gt        string `validate:"gt=10"`
	Gte       string `validate:"gte=10"`
	OmitEmpty string `validate:"omitempty,min=1,max=10"`
	Boolean   string `validate:"boolean"`
	Sub       *SubTest
	SubIgnore *SubTest `validate:"-"`
	Anonymous struct {
		A string `validate:"required"`
	}
	Iface I
}

type TestUint64 struct {
	Required  uint64 `validate:"required"`
	Len       uint64 `validate:"len=10"`
	Min       uint64 `validate:"min=1"`
	Max       uint64 `validate:"max=10"`
	MinMax    uint64 `validate:"min=1,max=10"`
	OmitEmpty uint64 `validate:"omitempty,min=1,max=10"`
}

type TestFloat64 struct {
	Required  float64 `validate:"required"`
	Len       float64 `validate:"len=10"`
	Min       float64 `validate:"min=1"`
	Max       float64 `validate:"max=10"`
	MinMax    float64 `validate:"min=1,max=10"`
	Lte       float64 `validate:"lte=10"`
	OmitEmpty float64 `validate:"omitempty,min=1,max=10"`
}

type TestSlice struct {
	Required  []int `validate:"required"`
	Len       []int `validate:"len=10"`
	Min       []int `validate:"min=1"`
	Max       []int `validate:"max=10"`
	MinMax    []int `validate:"min=1,max=10"`
	OmitEmpty []int `validate:"omitempty,min=1,max=10"`
}

func AssertError(t *testing.T, err error, nsKey, structNsKey, field, structField, expectedTag string) {
	errs := err.(ValidationErrors)

	found := false
	var fe FieldError

	for i := 0; i < len(errs); i++ {
		if errs[i].Namespace() == nsKey && errs[i].StructNamespace() == structNsKey {
			found = true
			fe = errs[i]
			break
		}
	}

	EqualSkip(t, 2, found, true)
	NotEqualSkip(t, 2, fe, nil)
	EqualSkip(t, 2, fe.Field(), field)
	EqualSkip(t, 2, fe.StructField(), structField)
	EqualSkip(t, 2, fe.Tag(), expectedTag)
}

func AssertDeepError(t *testing.T, err error, nsKey, structNsKey, field, structField, expectedTag, actualTag string) {
	errs := err.(ValidationErrors)

	found := false
	var fe FieldError

	for i := 0; i < len(errs); i++ {
		if errs[i].Namespace() == nsKey && errs[i].StructNamespace() == structNsKey && errs[i].Tag() == expectedTag && errs[i].ActualTag() == actualTag {
			found = true
			fe = errs[i]
			break
		}
	}

	EqualSkip(t, 2, found, true)
	NotEqualSkip(t, 2, fe, nil)
	EqualSkip(t, 2, fe.Field(), field)
	EqualSkip(t, 2, fe.StructField(), structField)
}

func getError(err error, nsKey, structNsKey string) FieldError {
	errs := err.(ValidationErrors)

	var fe FieldError

	for i := 0; i < len(errs); i++ {
		if errs[i].Namespace() == nsKey && errs[i].StructNamespace() == structNsKey {
			fe = errs[i]
			break
		}
	}

	return fe
}

type valuer struct {
	Name string
}

func (v valuer) Value() (driver.Value, error) {
	if v.Name == "errorme" {
		panic("SQL Driver Valuer error: some kind of error")
		// return nil, errors.New("some kind of error")
	}

	if len(v.Name) == 0 {
		return nil, nil
	}

	return v.Name, nil
}

type MadeUpCustomType struct {
	FirstName string
	LastName  string
}

func ValidateCustomType(field reflect.Value) interface{} {
	if cust, ok := field.Interface().(MadeUpCustomType); ok {

		if len(cust.FirstName) == 0 || len(cust.LastName) == 0 {
			return ""
		}

		return cust.FirstName + " " + cust.LastName
	}

	return ""
}

func OverrideIntTypeForSomeReason(field reflect.Value) interface{} {
	if i, ok := field.Interface().(int); ok {
		if i == 1 {
			return "1"
		}

		if i == 2 {
			return "12"
		}
	}

	return ""
}

type CustomMadeUpStruct struct {
	MadeUp        MadeUpCustomType `validate:"required"`
	OverriddenInt int              `validate:"gt=1"`
}

func ValidateValuerType(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()
		if err != nil {
			// handle the error how you want
			return nil
		}

		return val
	}

	return nil
}

type TestPartial struct {
	NoTag     string
	BlankTag  string     `validate:""`
	Required  string     `validate:"required"`
	SubSlice  []*SubTest `validate:"required,dive"`
	Sub       *SubTest
	SubIgnore *SubTest `validate:"-"`
	Anonymous struct {
		A         string     `validate:"required"`
		ASubSlice []*SubTest `validate:"required,dive"`

		SubAnonStruct []struct {
			Test      string `validate:"required"`
			OtherTest string `validate:"required"`
		} `validate:"required,dive"`
	}
}

type TestStruct struct {
	String string `validate:"required" json:"StringVal"`
}

func StructValidationTestStructSuccess(sl StructLevel) {
	st := sl.Current().Interface().(TestStruct)

	if st.String != "good value" {
		sl.ReportError(st.String, "StringVal", "String", "badvalueteststruct", "good value")
	}
}

func StructValidationTestStruct(sl StructLevel) {
	st := sl.Current().Interface().(TestStruct)

	if st.String != "bad value" {
		sl.ReportError(st.String, "StringVal", "String", "badvalueteststruct", "bad value")
	}
}

func StructValidationNoTestStructCustomName(sl StructLevel) {
	st := sl.Current().Interface().(TestStruct)

	if st.String != "bad value" {
		sl.ReportError(st.String, "String", "", "badvalueteststruct", "bad value")
	}
}

func StructValidationTestStructInvalid(sl StructLevel) {
	st := sl.Current().Interface().(TestStruct)

	if st.String != "bad value" {
		sl.ReportError(nil, "StringVal", "String", "badvalueteststruct", "bad value")
	}
}

func StructValidationTestStructReturnValidationErrors(sl StructLevel) {
	s := sl.Current().Interface().(TestStructReturnValidationErrors)

	errs := sl.Validator().Struct(s.Inner1.Inner2)
	if errs == nil {
		return
	}

	sl.ReportValidationErrors("Inner1.", "Inner1.", errs.(ValidationErrors))
}

func StructValidationTestStructReturnValidationErrors2(sl StructLevel) {
	s := sl.Current().Interface().(TestStructReturnValidationErrors)

	errs := sl.Validator().Struct(s.Inner1.Inner2)
	if errs == nil {
		return
	}

	sl.ReportValidationErrors("Inner1JSON.", "Inner1.", errs.(ValidationErrors))
}

type TestStructReturnValidationErrorsInner2 struct {
	String string `validate:"required" json:"JSONString"`
}

type TestStructReturnValidationErrorsInner1 struct {
	Inner2 *TestStructReturnValidationErrorsInner2
}

type TestStructReturnValidationErrors struct {
	Inner1 *TestStructReturnValidationErrorsInner1 `json:"Inner1JSON"`
}

type StructLevelInvalidErr struct {
	Value string
}

func StructLevelInvalidError(sl StructLevel) {
	top := sl.Top().Interface().(StructLevelInvalidErr)
	s := sl.Current().Interface().(StructLevelInvalidErr)

	if top.Value == s.Value {
		sl.ReportError(nil, "Value", "Value", "required", "")
	}
}

func stringPtr(v string) *string {
	return &v
}

func intPtr(v int) *int {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}

func TestStructLevelInvalidError(t *testing.T) {
	validate := New()
	validate.RegisterStructValidation(StructLevelInvalidError, StructLevelInvalidErr{})

	var test StructLevelInvalidErr

	err := validate.Struct(test)
	NotEqual(t, err, nil)

	errs, ok := err.(ValidationErrors)
	Equal(t, ok, true)

	fe := errs[0]
	Equal(t, fe.Field(), "Value")
	Equal(t, fe.StructField(), "Value")
	Equal(t, fe.Namespace(), "StructLevelInvalidErr.Value")
	Equal(t, fe.StructNamespace(), "StructLevelInvalidErr.Value")
	Equal(t, fe.Tag(), "required")
	Equal(t, fe.ActualTag(), "required")
	Equal(t, fe.Kind(), reflect.Invalid)
	Equal(t, fe.Type(), reflect.TypeOf(nil))
}

func TestNameNamespace(t *testing.T) {
	type Inner2Namespace struct {
		String []string `validate:"dive,required" json:"JSONString"`
	}

	type Inner1Namespace struct {
		Inner2 *Inner2Namespace `json:"Inner2JSON"`
	}

	type Namespace struct {
		Inner1 *Inner1Namespace `json:"Inner1JSON"`
	}

	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	i2 := &Inner2Namespace{String: []string{"ok", "ok", "ok"}}
	i1 := &Inner1Namespace{Inner2: i2}
	ns := &Namespace{Inner1: i1}

	errs := validate.Struct(ns)
	Equal(t, errs, nil)

	i2.String[1] = ""

	errs = validate.Struct(ns)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "Namespace.Inner1JSON.Inner2JSON.JSONString[1]", "Namespace.Inner1.Inner2.String[1]", "JSONString[1]", "String[1]", "required")

	fe := getError(ve, "Namespace.Inner1JSON.Inner2JSON.JSONString[1]", "Namespace.Inner1.Inner2.String[1]")
	NotEqual(t, fe, nil)
	Equal(t, fe.Field(), "JSONString[1]")
	Equal(t, fe.StructField(), "String[1]")
	Equal(t, fe.Namespace(), "Namespace.Inner1JSON.Inner2JSON.JSONString[1]")
	Equal(t, fe.StructNamespace(), "Namespace.Inner1.Inner2.String[1]")
}

func TestAnonymous(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type Test struct {
		Anonymous struct {
			A string `validate:"required" json:"EH"`
		}
		AnonymousB struct {
			B string `validate:"required" json:"BEE"`
		}
		anonymousC struct {
			c string `validate:"required"`
		}
	}

	tst := &Test{
		Anonymous: struct {
			A string `validate:"required" json:"EH"`
		}{
			A: "1",
		},
		AnonymousB: struct {
			B string `validate:"required" json:"BEE"`
		}{
			B: "",
		},
		anonymousC: struct {
			c string `validate:"required"`
		}{
			c: "",
		},
	}

	Equal(t, tst.anonymousC.c, "")

	err := validate.Struct(tst)
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)

	Equal(t, len(errs), 1)
	AssertError(t, errs, "Test.AnonymousB.BEE", "Test.AnonymousB.B", "BEE", "B", "required")

	fe := getError(errs, "Test.AnonymousB.BEE", "Test.AnonymousB.B")
	NotEqual(t, fe, nil)
	Equal(t, fe.Field(), "BEE")
	Equal(t, fe.StructField(), "B")

	s := struct {
		c string `validate:"required"`
	}{
		c: "",
	}

	err = validate.Struct(s)
	Equal(t, err, nil)
}

func TestAnonymousSameStructDifferentTags(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type Test struct {
		A interface{}
	}

	tst := &Test{
		A: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
	}

	err := validate.Struct(tst)
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)

	Equal(t, len(errs), 1)
	AssertError(t, errs, "Test.A.A", "Test.A.A", "A", "A", "required")

	tst = &Test{
		A: struct {
			A string `validate:"omitempty,required"`
		}{
			A: "",
		},
	}

	err = validate.Struct(tst)
	Equal(t, err, nil)
}

func TestStructLevelReturnValidationErrors(t *testing.T) {
	validate := New()
	validate.RegisterStructValidation(StructValidationTestStructReturnValidationErrors, TestStructReturnValidationErrors{})

	inner2 := &TestStructReturnValidationErrorsInner2{
		String: "I'm HERE",
	}

	inner1 := &TestStructReturnValidationErrorsInner1{
		Inner2: inner2,
	}

	val := &TestStructReturnValidationErrors{
		Inner1: inner1,
	}

	errs := validate.Struct(val)
	Equal(t, errs, nil)

	inner2.String = ""

	errs = validate.Struct(val)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 2)
	AssertError(t, errs, "TestStructReturnValidationErrors.Inner1.Inner2.String", "TestStructReturnValidationErrors.Inner1.Inner2.String", "String", "String", "required")
	// this is an extra error reported from struct validation
	AssertError(t, errs, "TestStructReturnValidationErrors.Inner1.TestStructReturnValidationErrorsInner2.String", "TestStructReturnValidationErrors.Inner1.TestStructReturnValidationErrorsInner2.String", "String", "String", "required")
}

func TestStructLevelReturnValidationErrorsWithJSON(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	validate.RegisterStructValidation(StructValidationTestStructReturnValidationErrors2, TestStructReturnValidationErrors{})

	inner2 := &TestStructReturnValidationErrorsInner2{
		String: "I'm HERE",
	}

	inner1 := &TestStructReturnValidationErrorsInner1{
		Inner2: inner2,
	}

	val := &TestStructReturnValidationErrors{
		Inner1: inner1,
	}

	errs := validate.Struct(val)
	Equal(t, errs, nil)

	inner2.String = ""

	errs = validate.Struct(val)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 2)
	AssertError(t, errs, "TestStructReturnValidationErrors.Inner1JSON.Inner2.JSONString", "TestStructReturnValidationErrors.Inner1.Inner2.String", "JSONString", "String", "required")
	// this is an extra error reported from struct validation, it's a badly formatted one, but on purpose
	AssertError(t, errs, "TestStructReturnValidationErrors.Inner1JSON.TestStructReturnValidationErrorsInner2.JSONString", "TestStructReturnValidationErrors.Inner1.TestStructReturnValidationErrorsInner2.String", "JSONString", "String", "required")

	fe := getError(errs, "TestStructReturnValidationErrors.Inner1JSON.Inner2.JSONString", "TestStructReturnValidationErrors.Inner1.Inner2.String")
	NotEqual(t, fe, nil)

	// check for proper JSON namespace
	Equal(t, fe.Field(), "JSONString")
	Equal(t, fe.StructField(), "String")
	Equal(t, fe.Namespace(), "TestStructReturnValidationErrors.Inner1JSON.Inner2.JSONString")
	Equal(t, fe.StructNamespace(), "TestStructReturnValidationErrors.Inner1.Inner2.String")

	fe = getError(errs, "TestStructReturnValidationErrors.Inner1JSON.TestStructReturnValidationErrorsInner2.JSONString", "TestStructReturnValidationErrors.Inner1.TestStructReturnValidationErrorsInner2.String")
	NotEqual(t, fe, nil)

	// check for proper JSON namespace
	Equal(t, fe.Field(), "JSONString")
	Equal(t, fe.StructField(), "String")
	Equal(t, fe.Namespace(), "TestStructReturnValidationErrors.Inner1JSON.TestStructReturnValidationErrorsInner2.JSONString")
	Equal(t, fe.StructNamespace(), "TestStructReturnValidationErrors.Inner1.TestStructReturnValidationErrorsInner2.String")
}

func TestStructLevelValidations(t *testing.T) {
	v1 := New()
	v1.RegisterStructValidation(StructValidationTestStruct, TestStruct{})

	tst := &TestStruct{
		String: "good value",
	}

	errs := v1.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestStruct.StringVal", "TestStruct.String", "StringVal", "String", "badvalueteststruct")

	v2 := New()
	v2.RegisterStructValidation(StructValidationNoTestStructCustomName, TestStruct{})

	errs = v2.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestStruct.String", "TestStruct.String", "String", "String", "badvalueteststruct")

	v3 := New()
	v3.RegisterStructValidation(StructValidationTestStructInvalid, TestStruct{})

	errs = v3.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestStruct.StringVal", "TestStruct.String", "StringVal", "String", "badvalueteststruct")

	v4 := New()
	v4.RegisterStructValidation(StructValidationTestStructSuccess, TestStruct{})

	errs = v4.Struct(tst)
	Equal(t, errs, nil)
}

func TestAliasTags(t *testing.T) {
	validate := New()
	validate.RegisterAlias("iscoloralias", "hexcolor|rgb|rgba|hsl|hsla")

	s := "rgb(255,255,255)"
	errs := validate.Var(s, "iscoloralias")
	Equal(t, errs, nil)

	s = ""
	errs = validate.Var(s, "omitempty,iscoloralias")
	Equal(t, errs, nil)

	s = "rgb(255,255,0)"
	errs = validate.Var(s, "iscoloralias,len=5")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "len")

	type Test struct {
		Color string `validate:"iscoloralias"`
	}

	tst := &Test{
		Color: "#000",
	}

	errs = validate.Struct(tst)
	Equal(t, errs, nil)

	tst.Color = "cfvre"
	errs = validate.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Color", "Test.Color", "Color", "Color", "iscoloralias")

	fe := getError(errs, "Test.Color", "Test.Color")
	NotEqual(t, fe, nil)
	Equal(t, fe.ActualTag(), "hexcolor|rgb|rgba|hsl|hsla")

	validate.RegisterAlias("req", "required,dive,iscoloralias")
	arr := []string{"val1", "#fff", "#000"}

	errs = validate.Var(arr, "req")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "[0]", "[0]", "[0]", "[0]", "iscoloralias")

	PanicMatches(t, func() { validate.RegisterAlias("exists!", "gt=5,lt=10") }, "Alias 'exists!' either contains restricted characters or is the same as a restricted tag needed for normal operation")
}

func TestNilValidator(t *testing.T) {
	type TestStruct struct {
		Test string `validate:"required"`
	}

	ts := TestStruct{}

	var val *Validate

	fn := func(fl FieldLevel) bool {
		return fl.Parent().String() == fl.Field().String()
	}

	PanicMatches(t, func() { val.RegisterCustomTypeFunc(ValidateCustomType, MadeUpCustomType{}) }, "runtime error: invalid memory address or nil pointer dereference")
	PanicMatches(t, func() { _ = val.RegisterValidation("something", fn) }, "runtime error: invalid memory address or nil pointer dereference")
	PanicMatches(t, func() { _ = val.Var(ts.Test, "required") }, "runtime error: invalid memory address or nil pointer dereference")
	PanicMatches(t, func() { _ = val.VarWithValue("test", ts.Test, "required") }, "runtime error: invalid memory address or nil pointer dereference")
	PanicMatches(t, func() { _ = val.Struct(ts) }, "runtime error: invalid memory address or nil pointer dereference")
	PanicMatches(t, func() { _ = val.StructExcept(ts, "Test") }, "runtime error: invalid memory address or nil pointer dereference")
	PanicMatches(t, func() { _ = val.StructPartial(ts, "Test") }, "runtime error: invalid memory address or nil pointer dereference")
}

func TestStructPartial(t *testing.T) {
	p1 := []string{
		"NoTag",
		"Required",
	}

	p2 := []string{
		"SubSlice[0].Test",
		"Sub",
		"SubIgnore",
		"Anonymous.A",
	}

	p3 := []string{
		"SubTest.Test",
	}

	p4 := []string{
		"A",
	}

	tPartial := &TestPartial{
		NoTag:    "NoTag",
		Required: "Required",

		SubSlice: []*SubTest{
			{

				Test: "Required",
			},
			{

				Test: "Required",
			},
		},

		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A             string     `validate:"required"`
			ASubSlice     []*SubTest `validate:"required,dive"`
			SubAnonStruct []struct {
				Test      string `validate:"required"`
				OtherTest string `validate:"required"`
			} `validate:"required,dive"`
		}{
			A: "1",
			ASubSlice: []*SubTest{
				{
					Test: "Required",
				},
				{
					Test: "Required",
				},
			},

			SubAnonStruct: []struct {
				Test      string `validate:"required"`
				OtherTest string `validate:"required"`
			}{
				{"Required", "RequiredOther"},
				{"Required", "RequiredOther"},
			},
		},
	}

	validate := New()

	// the following should all return no errors as everything is valid in
	// the default state
	errs := validate.StructPartialCtx(context.Background(), tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructPartial(tPartial, p2...)
	Equal(t, errs, nil)

	// this isn't really a robust test, but is meant to illustrate the ANON CASE below
	errs = validate.StructPartial(tPartial.SubSlice[0], p3...)
	Equal(t, errs, nil)

	errs = validate.StructExceptCtx(context.Background(), tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructExcept(tPartial, p2...)
	Equal(t, errs, nil)

	// mod tPartial for required field and re-test making sure invalid fields are NOT required:
	tPartial.Required = ""

	errs = validate.StructExcept(tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructPartial(tPartial, p2...)
	Equal(t, errs, nil)

	// inversion and retesting Partial to generate failures:
	errs = validate.StructPartial(tPartial, p1...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Required", "TestPartial.Required", "Required", "Required", "required")

	errs = validate.StructExcept(tPartial, p2...)
	AssertError(t, errs, "TestPartial.Required", "TestPartial.Required", "Required", "Required", "required")

	// reset Required field, and set nested struct
	tPartial.Required = "Required"
	tPartial.Anonymous.A = ""

	// will pass as unset fields is not going to be tested
	errs = validate.StructPartial(tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructExcept(tPartial, p2...)
	Equal(t, errs, nil)

	// ANON CASE the response here is strange, it clearly does what it is being told to
	errs = validate.StructExcept(tPartial.Anonymous, p4...)
	Equal(t, errs, nil)

	// will fail as unset field is tested
	errs = validate.StructPartial(tPartial, p2...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Anonymous.A", "TestPartial.Anonymous.A", "A", "A", "required")

	errs = validate.StructExcept(tPartial, p1...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Anonymous.A", "TestPartial.Anonymous.A", "A", "A", "required")

	// reset nested struct and unset struct in slice
	tPartial.Anonymous.A = "Required"
	tPartial.SubSlice[0].Test = ""

	// these will pass as unset item is NOT tested
	errs = validate.StructPartial(tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructExcept(tPartial, p2...)
	Equal(t, errs, nil)

	// these will fail as unset item IS tested
	errs = validate.StructExcept(tPartial, p1...)
	AssertError(t, errs, "TestPartial.SubSlice[0].Test", "TestPartial.SubSlice[0].Test", "Test", "Test", "required")
	Equal(t, len(errs.(ValidationErrors)), 1)

	errs = validate.StructPartial(tPartial, p2...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.SubSlice[0].Test", "TestPartial.SubSlice[0].Test", "Test", "Test", "required")
	Equal(t, len(errs.(ValidationErrors)), 1)

	// Unset second slice member concurrently to test dive behavior:
	tPartial.SubSlice[1].Test = ""

	errs = validate.StructPartial(tPartial, p1...)
	Equal(t, errs, nil)

	// NOTE: When specifying nested items, it is still the users responsibility
	// to specify the dive tag, the library does not override this.
	errs = validate.StructExcept(tPartial, p2...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.SubSlice[1].Test", "TestPartial.SubSlice[1].Test", "Test", "Test", "required")

	errs = validate.StructExcept(tPartial, p1...)
	Equal(t, len(errs.(ValidationErrors)), 2)
	AssertError(t, errs, "TestPartial.SubSlice[0].Test", "TestPartial.SubSlice[0].Test", "Test", "Test", "required")
	AssertError(t, errs, "TestPartial.SubSlice[1].Test", "TestPartial.SubSlice[1].Test", "Test", "Test", "required")

	errs = validate.StructPartial(tPartial, p2...)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestPartial.SubSlice[0].Test", "TestPartial.SubSlice[0].Test", "Test", "Test", "required")

	// reset struct in slice, and unset struct in slice in unset position
	tPartial.SubSlice[0].Test = "Required"

	// these will pass as the unset item is NOT tested
	errs = validate.StructPartial(tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructPartial(tPartial, p2...)
	Equal(t, errs, nil)

	// testing for missing item by exception, yes it dives and fails
	errs = validate.StructExcept(tPartial, p1...)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestPartial.SubSlice[1].Test", "TestPartial.SubSlice[1].Test", "Test", "Test", "required")

	errs = validate.StructExcept(tPartial, p2...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.SubSlice[1].Test", "TestPartial.SubSlice[1].Test", "Test", "Test", "required")

	tPartial.SubSlice[1].Test = "Required"

	tPartial.Anonymous.SubAnonStruct[0].Test = ""
	// these will pass as the unset item is NOT tested
	errs = validate.StructPartial(tPartial, p1...)
	Equal(t, errs, nil)

	errs = validate.StructPartial(tPartial, p2...)
	Equal(t, errs, nil)

	errs = validate.StructExcept(tPartial, p1...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Anonymous.SubAnonStruct[0].Test", "TestPartial.Anonymous.SubAnonStruct[0].Test", "Test", "Test", "required")

	errs = validate.StructExcept(tPartial, p2...)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Anonymous.SubAnonStruct[0].Test", "TestPartial.Anonymous.SubAnonStruct[0].Test", "Test", "Test", "required")

	// Test for unnamed struct
	testStruct := &TestStruct{
		String: "test",
	}
	unnamedStruct := struct {
		String string `validate:"required" json:"StringVal"`
	}{String: "test"}
	composedUnnamedStruct := struct{ *TestStruct }{&TestStruct{String: "test"}}

	errs = validate.StructPartial(testStruct, "String")
	Equal(t, errs, nil)

	errs = validate.StructPartial(unnamedStruct, "String")
	Equal(t, errs, nil)

	errs = validate.StructPartial(composedUnnamedStruct, "TestStruct.String")
	Equal(t, errs, nil)

	testStruct.String = ""
	errs = validate.StructPartial(testStruct, "String")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestStruct.String", "TestStruct.String", "String", "String", "required")

	unnamedStruct.String = ""
	errs = validate.StructPartial(unnamedStruct, "String")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "String", "String", "String", "String", "required")

	composedUnnamedStruct.String = ""
	errs = validate.StructPartial(composedUnnamedStruct, "TestStruct.String")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestStruct.String", "TestStruct.String", "String", "String", "required")
}

func TestCrossStructLteFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	type Inner struct {
		CreatedAt *time.Time
		String    string
		Int       int
		Uint      uint
		Float     float64
		Array     []string
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"ltecsfield=Inner.CreatedAt"`
		String    string     `validate:"ltecsfield=Inner.String"`
		Int       int        `validate:"ltecsfield=Inner.Int"`
		Uint      uint       `validate:"ltecsfield=Inner.Uint"`
		Float     float64    `validate:"ltecsfield=Inner.Float"`
		Array     []string   `validate:"ltecsfield=Inner.Array"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	inner := &Inner{
		CreatedAt: &then,
		String:    "abcd",
		Int:       13,
		Uint:      13,
		Float:     1.13,
		Array:     []string{"val1", "val2"},
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
		String:    "abc",
		Int:       12,
		Uint:      12,
		Float:     1.12,
		Array:     []string{"val1"},
	}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	test.CreatedAt = &then
	test.String = "abcd"
	test.Int = 13
	test.Uint = 13
	test.Float = 1.13
	test.Array = []string{"val1", "val2"}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	after := now.Add(time.Hour * 10)

	test.CreatedAt = &after
	test.String = "abce"
	test.Int = 14
	test.Uint = 14
	test.Float = 1.14
	test.Array = []string{"val1", "val2", "val3"}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "ltecsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "ltecsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "ltecsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "ltecsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "ltecsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "ltecsfield")

	errs = validate.VarWithValueCtx(context.Background(), 1, "", "ltecsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltecsfield")

	// this test is for the WARNING about unforeseen validation issues.
	errs = validate.VarWithValue(test, now, "ltecsfield")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 6)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "ltecsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "ltecsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "ltecsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "ltecsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "ltecsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "ltecsfield")

	type Other struct {
		Value string
	}

	type Test2 struct {
		Value Other
		Time  time.Time `validate:"ltecsfield=Value"`
	}

	tst := Test2{
		Value: Other{Value: "StringVal"},
		Time:  then,
	}

	errs = validate.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test2.Time", "Test2.Time", "Time", "Time", "ltecsfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "ltecsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "ltecsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "ltecsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltecsfield")

	errs = validate.VarWithValue(time.Duration(0), -time.Minute, "omitempty,ltecsfield")
	Equal(t, errs, nil)

	// -- Validations for a struct and an inner struct with time.Duration type fields.

	type TimeDurationInner struct {
		Duration time.Duration
	}
	var timeDurationInner *TimeDurationInner

	type TimeDurationTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"ltecsfield=Inner.Duration"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationInner = &TimeDurationInner{time.Hour + time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour - time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "ltecsfield")

	type TimeDurationOmitemptyTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"omitempty,ltecsfield=Inner.Duration"`
	}
	var timeDurationOmitemptyTest *TimeDurationOmitemptyTest

	timeDurationInner = &TimeDurationInner{-time.Minute}
	timeDurationOmitemptyTest = &TimeDurationOmitemptyTest{timeDurationInner, time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestCrossStructLtFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	type Inner struct {
		CreatedAt *time.Time
		String    string
		Int       int
		Uint      uint
		Float     float64
		Array     []string
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"ltcsfield=Inner.CreatedAt"`
		String    string     `validate:"ltcsfield=Inner.String"`
		Int       int        `validate:"ltcsfield=Inner.Int"`
		Uint      uint       `validate:"ltcsfield=Inner.Uint"`
		Float     float64    `validate:"ltcsfield=Inner.Float"`
		Array     []string   `validate:"ltcsfield=Inner.Array"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	inner := &Inner{
		CreatedAt: &then,
		String:    "abcd",
		Int:       13,
		Uint:      13,
		Float:     1.13,
		Array:     []string{"val1", "val2"},
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
		String:    "abc",
		Int:       12,
		Uint:      12,
		Float:     1.12,
		Array:     []string{"val1"},
	}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	test.CreatedAt = &then
	test.String = "abcd"
	test.Int = 13
	test.Uint = 13
	test.Float = 1.13
	test.Array = []string{"val1", "val2"}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "ltcsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "ltcsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "ltcsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "ltcsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "ltcsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "ltcsfield")

	errs = validate.VarWithValue(1, "", "ltcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltcsfield")

	// this test is for the WARNING about unforeseen validation issues.
	errs = validate.VarWithValue(test, now, "ltcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "ltcsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "ltcsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "ltcsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "ltcsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "ltcsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "ltcsfield")

	type Other struct {
		Value string
	}

	type Test2 struct {
		Value Other
		Time  time.Time `validate:"ltcsfield=Value"`
	}

	tst := Test2{
		Value: Other{Value: "StringVal"},
		Time:  then,
	}

	errs = validate.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test2.Time", "Test2.Time", "Time", "Time", "ltcsfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "ltcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "ltcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltcsfield")

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "ltcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltcsfield")

	errs = validate.VarWithValue(time.Duration(0), -time.Minute, "omitempty,ltcsfield")
	Equal(t, errs, nil)

	// -- Validations for a struct and an inner struct with time.Duration type fields.

	type TimeDurationInner struct {
		Duration time.Duration
	}
	var timeDurationInner *TimeDurationInner

	type TimeDurationTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"ltcsfield=Inner.Duration"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationInner = &TimeDurationInner{time.Hour + time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "ltcsfield")

	timeDurationInner = &TimeDurationInner{time.Hour - time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "ltcsfield")

	type TimeDurationOmitemptyTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"omitempty,ltcsfield=Inner.Duration"`
	}
	var timeDurationOmitemptyTest *TimeDurationOmitemptyTest

	timeDurationInner = &TimeDurationInner{-time.Minute}
	timeDurationOmitemptyTest = &TimeDurationOmitemptyTest{timeDurationInner, time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestCrossStructGteFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	type Inner struct {
		CreatedAt *time.Time
		String    string
		Int       int
		Uint      uint
		Float     float64
		Array     []string
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"gtecsfield=Inner.CreatedAt"`
		String    string     `validate:"gtecsfield=Inner.String"`
		Int       int        `validate:"gtecsfield=Inner.Int"`
		Uint      uint       `validate:"gtecsfield=Inner.Uint"`
		Float     float64    `validate:"gtecsfield=Inner.Float"`
		Array     []string   `validate:"gtecsfield=Inner.Array"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)

	inner := &Inner{
		CreatedAt: &then,
		String:    "abcd",
		Int:       13,
		Uint:      13,
		Float:     1.13,
		Array:     []string{"val1", "val2"},
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
		String:    "abcde",
		Int:       14,
		Uint:      14,
		Float:     1.14,
		Array:     []string{"val1", "val2", "val3"},
	}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	test.CreatedAt = &then
	test.String = "abcd"
	test.Int = 13
	test.Uint = 13
	test.Float = 1.13
	test.Array = []string{"val1", "val2"}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	before := now.Add(time.Hour * -10)

	test.CreatedAt = &before
	test.String = "abc"
	test.Int = 12
	test.Uint = 12
	test.Float = 1.12
	test.Array = []string{"val1"}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "gtecsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "gtecsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "gtecsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "gtecsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "gtecsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "gtecsfield")

	errs = validate.VarWithValue(1, "", "gtecsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtecsfield")

	// this test is for the WARNING about unforeseen validation issues.
	errs = validate.VarWithValue(test, now, "gtecsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "gtecsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "gtecsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "gtecsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "gtecsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "gtecsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "gtecsfield")

	type Other struct {
		Value string
	}

	type Test2 struct {
		Value Other
		Time  time.Time `validate:"gtecsfield=Value"`
	}

	tst := Test2{
		Value: Other{Value: "StringVal"},
		Time:  then,
	}

	errs = validate.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test2.Time", "Test2.Time", "Time", "Time", "gtecsfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "gtecsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "gtecsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "gtecsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtecsfield")

	errs = validate.VarWithValue(time.Duration(0), time.Hour, "omitempty,gtecsfield")
	Equal(t, errs, nil)

	// -- Validations for a struct and an inner struct with time.Duration type fields.

	type TimeDurationInner struct {
		Duration time.Duration
	}
	var timeDurationInner *TimeDurationInner

	type TimeDurationTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"gtecsfield=Inner.Duration"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationInner = &TimeDurationInner{time.Hour - time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour + time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "gtecsfield")

	type TimeDurationOmitemptyTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"omitempty,gtecsfield=Inner.Duration"`
	}
	var timeDurationOmitemptyTest *TimeDurationOmitemptyTest

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationOmitemptyTest = &TimeDurationOmitemptyTest{timeDurationInner, time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestCrossStructGtFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	type Inner struct {
		CreatedAt *time.Time
		String    string
		Int       int
		Uint      uint
		Float     float64
		Array     []string
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"gtcsfield=Inner.CreatedAt"`
		String    string     `validate:"gtcsfield=Inner.String"`
		Int       int        `validate:"gtcsfield=Inner.Int"`
		Uint      uint       `validate:"gtcsfield=Inner.Uint"`
		Float     float64    `validate:"gtcsfield=Inner.Float"`
		Array     []string   `validate:"gtcsfield=Inner.Array"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * -5)

	inner := &Inner{
		CreatedAt: &then,
		String:    "abcd",
		Int:       13,
		Uint:      13,
		Float:     1.13,
		Array:     []string{"val1", "val2"},
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
		String:    "abcde",
		Int:       14,
		Uint:      14,
		Float:     1.14,
		Array:     []string{"val1", "val2", "val3"},
	}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	test.CreatedAt = &then
	test.String = "abcd"
	test.Int = 13
	test.Uint = 13
	test.Float = 1.13
	test.Array = []string{"val1", "val2"}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "gtcsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "gtcsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "gtcsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "gtcsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "gtcsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "gtcsfield")

	errs = validate.VarWithValue(1, "", "gtcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtcsfield")

	// this test is for the WARNING about unforeseen validation issues.
	errs = validate.VarWithValue(test, now, "gtcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "gtcsfield")
	AssertError(t, errs, "Test.String", "Test.String", "String", "String", "gtcsfield")
	AssertError(t, errs, "Test.Int", "Test.Int", "Int", "Int", "gtcsfield")
	AssertError(t, errs, "Test.Uint", "Test.Uint", "Uint", "Uint", "gtcsfield")
	AssertError(t, errs, "Test.Float", "Test.Float", "Float", "Float", "gtcsfield")
	AssertError(t, errs, "Test.Array", "Test.Array", "Array", "Array", "gtcsfield")

	type Other struct {
		Value string
	}

	type Test2 struct {
		Value Other
		Time  time.Time `validate:"gtcsfield=Value"`
	}

	tst := Test2{
		Value: Other{Value: "StringVal"},
		Time:  then,
	}

	errs = validate.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test2.Time", "Test2.Time", "Time", "Time", "gtcsfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "gtcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "gtcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtcsfield")

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "gtcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtcsfield")

	errs = validate.VarWithValue(time.Duration(0), time.Hour, "omitempty,gtcsfield")
	Equal(t, errs, nil)

	// -- Validations for a struct and an inner struct with time.Duration type fields.

	type TimeDurationInner struct {
		Duration time.Duration
	}
	var timeDurationInner *TimeDurationInner

	type TimeDurationTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"gtcsfield=Inner.Duration"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationInner = &TimeDurationInner{time.Hour - time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "gtcsfield")

	timeDurationInner = &TimeDurationInner{time.Hour + time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "gtcsfield")

	type TimeDurationOmitemptyTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"omitempty,gtcsfield=Inner.Duration"`
	}
	var timeDurationOmitemptyTest *TimeDurationOmitemptyTest

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationOmitemptyTest = &TimeDurationOmitemptyTest{timeDurationInner, time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestCrossStructNeFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	type Inner struct {
		CreatedAt *time.Time
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"necsfield=Inner.CreatedAt"`
	}

	now := time.Now().UTC()
	then := now.Add(time.Hour * 5)

	inner := &Inner{
		CreatedAt: &then,
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
	}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	test.CreatedAt = &then

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "necsfield")

	var j uint64
	var k float64
	var j2 uint64
	var k2 float64
	s := "abcd"
	i := 1
	j = 1
	k = 1.543
	b := true
	arr := []string{"test"}

	s2 := "abcd"
	i2 := 1
	j2 = 1
	k2 = 1.543
	b2 := true
	arr2 := []string{"test"}
	arr3 := []string{"test", "test2"}
	now2 := now

	errs = validate.VarWithValue(s, s2, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(i2, i, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(j2, j, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(k2, k, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(b2, b, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(arr2, arr, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(now2, now, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(arr3, arr, "necsfield")
	Equal(t, errs, nil)

	type SInner struct {
		Name string
	}

	type TStruct struct {
		Inner     *SInner
		CreatedAt *time.Time `validate:"necsfield=Inner"`
	}

	sinner := &SInner{
		Name: "NAME",
	}

	test2 := &TStruct{
		Inner:     sinner,
		CreatedAt: &now,
	}

	errs = validate.Struct(test2)
	Equal(t, errs, nil)

	test2.Inner = nil
	errs = validate.Struct(test2)
	Equal(t, errs, nil)

	errs = validate.VarWithValue(nil, 1, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "necsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "necsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "necsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "necsfield")

	errs = validate.VarWithValue(time.Duration(0), time.Duration(0), "omitempty,necsfield")
	Equal(t, errs, nil)

	// -- Validations for a struct and an inner struct with time.Duration type fields.

	type TimeDurationInner struct {
		Duration time.Duration
	}
	var timeDurationInner *TimeDurationInner

	type TimeDurationTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"necsfield=Inner.Duration"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationInner = &TimeDurationInner{time.Hour - time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour + time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "necsfield")

	type TimeDurationOmitemptyTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"omitempty,necsfield=Inner.Duration"`
	}
	var timeDurationOmitemptyTest *TimeDurationOmitemptyTest

	timeDurationInner = &TimeDurationInner{time.Duration(0)}
	timeDurationOmitemptyTest = &TimeDurationOmitemptyTest{timeDurationInner, time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestCrossStructEqFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	type Inner struct {
		CreatedAt *time.Time
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"eqcsfield=Inner.CreatedAt"`
	}

	now := time.Now().UTC()

	inner := &Inner{
		CreatedAt: &now,
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
	}

	errs = validate.Struct(test)
	Equal(t, errs, nil)

	newTime := time.Now().Add(time.Hour).UTC()
	test.CreatedAt = &newTime

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.CreatedAt", "Test.CreatedAt", "CreatedAt", "CreatedAt", "eqcsfield")

	var j uint64
	var k float64
	s := "abcd"
	i := 1
	j = 1
	k = 1.543
	b := true
	arr := []string{"test"}

	var j2 uint64
	var k2 float64
	s2 := "abcd"
	i2 := 1
	j2 = 1
	k2 = 1.543
	b2 := true
	arr2 := []string{"test"}
	arr3 := []string{"test", "test2"}
	now2 := now

	errs = validate.VarWithValue(s, s2, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(i2, i, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(j2, j, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(k2, k, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(b2, b, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(arr2, arr, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(now2, now, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(arr3, arr, "eqcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqcsfield")

	type SInner struct {
		Name string
	}

	type TStruct struct {
		Inner     *SInner
		CreatedAt *time.Time `validate:"eqcsfield=Inner"`
	}

	sinner := &SInner{
		Name: "NAME",
	}

	test2 := &TStruct{
		Inner:     sinner,
		CreatedAt: &now,
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TStruct.CreatedAt", "TStruct.CreatedAt", "CreatedAt", "CreatedAt", "eqcsfield")

	test2.Inner = nil
	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TStruct.CreatedAt", "TStruct.CreatedAt", "CreatedAt", "CreatedAt", "eqcsfield")

	errs = validate.VarWithValue(nil, 1, "eqcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqcsfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour, "eqcsfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "eqcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqcsfield")

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "eqcsfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqcsfield")

	errs = validate.VarWithValue(time.Duration(0), time.Hour, "omitempty,eqcsfield")
	Equal(t, errs, nil)

	// -- Validations for a struct and an inner struct with time.Duration type fields.

	type TimeDurationInner struct {
		Duration time.Duration
	}
	var timeDurationInner *TimeDurationInner

	type TimeDurationTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"eqcsfield=Inner.Duration"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationInner = &TimeDurationInner{time.Hour - time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "eqcsfield")

	timeDurationInner = &TimeDurationInner{time.Hour + time.Minute}
	timeDurationTest = &TimeDurationTest{timeDurationInner, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "eqcsfield")

	type TimeDurationOmitemptyTest struct {
		Inner    *TimeDurationInner
		Duration time.Duration `validate:"omitempty,eqcsfield=Inner.Duration"`
	}
	var timeDurationOmitemptyTest *TimeDurationOmitemptyTest

	timeDurationInner = &TimeDurationInner{time.Hour}
	timeDurationOmitemptyTest = &TimeDurationOmitemptyTest{timeDurationInner, time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestCrossNamespaceFieldValidation(t *testing.T) {
	type SliceStruct struct {
		Name string
	}

	type Inner struct {
		CreatedAt        *time.Time
		Slice            []string
		SliceStructs     []*SliceStruct
		SliceSlice       [][]string
		SliceSliceStruct [][]*SliceStruct
		SliceMap         []map[string]string
		Map              map[string]string
		MapMap           map[string]map[string]string
		MapStructs       map[string]*SliceStruct
		MapMapStruct     map[string]map[string]*SliceStruct
		MapSlice         map[string][]string
		MapInt           map[int]string
		MapInt8          map[int8]string
		MapInt16         map[int16]string
		MapInt32         map[int32]string
		MapInt64         map[int64]string
		MapUint          map[uint]string
		MapUint8         map[uint8]string
		MapUint16        map[uint16]string
		MapUint32        map[uint32]string
		MapUint64        map[uint64]string
		MapFloat32       map[float32]string
		MapFloat64       map[float64]string
		MapBool          map[bool]string
	}

	type Test struct {
		Inner     *Inner
		CreatedAt *time.Time
	}

	now := time.Now()

	inner := &Inner{
		CreatedAt:        &now,
		Slice:            []string{"val1", "val2", "val3"},
		SliceStructs:     []*SliceStruct{{Name: "name1"}, {Name: "name2"}, {Name: "name3"}},
		SliceSlice:       [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
		SliceSliceStruct: [][]*SliceStruct{{{Name: "name1"}, {Name: "name2"}, {Name: "name3"}}, {{Name: "name4"}, {Name: "name5"}, {Name: "name6"}}, {{Name: "name7"}, {Name: "name8"}, {Name: "name9"}}},
		SliceMap:         []map[string]string{{"key1": "val1", "key2": "val2", "key3": "val3"}, {"key4": "val4", "key5": "val5", "key6": "val6"}},
		Map:              map[string]string{"key1": "val1", "key2": "val2", "key3": "val3"},
		MapStructs:       map[string]*SliceStruct{"key1": {Name: "name1"}, "key2": {Name: "name2"}, "key3": {Name: "name3"}},
		MapMap:           map[string]map[string]string{"key1": {"key1-1": "val1"}, "key2": {"key2-1": "val2"}, "key3": {"key3-1": "val3"}},
		MapMapStruct:     map[string]map[string]*SliceStruct{"key1": {"key1-1": {Name: "name1"}}, "key2": {"key2-1": {Name: "name2"}}, "key3": {"key3-1": {Name: "name3"}}},
		MapSlice:         map[string][]string{"key1": {"1", "2", "3"}, "key2": {"4", "5", "6"}, "key3": {"7", "8", "9"}},
		MapInt:           map[int]string{1: "val1", 2: "val2", 3: "val3"},
		MapInt8:          map[int8]string{1: "val1", 2: "val2", 3: "val3"},
		MapInt16:         map[int16]string{1: "val1", 2: "val2", 3: "val3"},
		MapInt32:         map[int32]string{1: "val1", 2: "val2", 3: "val3"},
		MapInt64:         map[int64]string{1: "val1", 2: "val2", 3: "val3"},
		MapUint:          map[uint]string{1: "val1", 2: "val2", 3: "val3"},
		MapUint8:         map[uint8]string{1: "val1", 2: "val2", 3: "val3"},
		MapUint16:        map[uint16]string{1: "val1", 2: "val2", 3: "val3"},
		MapUint32:        map[uint32]string{1: "val1", 2: "val2", 3: "val3"},
		MapUint64:        map[uint64]string{1: "val1", 2: "val2", 3: "val3"},
		MapFloat32:       map[float32]string{1.01: "val1", 2.02: "val2", 3.03: "val3"},
		MapFloat64:       map[float64]string{1.01: "val1", 2.02: "val2", 3.03: "val3"},
		MapBool:          map[bool]string{true: "val1", false: "val2"},
	}

	test := &Test{
		Inner:     inner,
		CreatedAt: &now,
	}

	val := reflect.ValueOf(test)

	vd := New()
	v := &validate{
		v: vd,
	}

	current, kind, _, ok := v.getStructFieldOKInternal(val, "Inner.CreatedAt")
	Equal(t, ok, true)
	Equal(t, kind, reflect.Struct)
	tm, ok := current.Interface().(time.Time)
	Equal(t, ok, true)
	Equal(t, tm, now)

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.Slice[1]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	_, _, _, ok = v.getStructFieldOKInternal(val, "Inner.CrazyNonExistantField")
	Equal(t, ok, false)

	_, _, _, ok = v.getStructFieldOKInternal(val, "Inner.Slice[101]")
	Equal(t, ok, false)

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.Map[key3]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val3")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapMap[key2][key2-1]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapStructs[key2].Name")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "name2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapMapStruct[key3][key3-1].Name")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "name3")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.SliceSlice[2][0]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "7")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.SliceSliceStruct[2][1].Name")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "name8")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.SliceMap[1][key5]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val5")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapSlice[key3][2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "9")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapInt[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapInt8[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapInt16[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapInt32[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapInt64[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapUint[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapUint8[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapUint16[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapUint32[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapUint64[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapFloat32[3.03]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val3")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapFloat64[2.02]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val2")

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.MapBool[true]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.String)
	Equal(t, current.String(), "val1")

	inner = &Inner{
		CreatedAt:        &now,
		Slice:            []string{"val1", "val2", "val3"},
		SliceStructs:     []*SliceStruct{{Name: "name1"}, {Name: "name2"}, nil},
		SliceSlice:       [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
		SliceSliceStruct: [][]*SliceStruct{{{Name: "name1"}, {Name: "name2"}, {Name: "name3"}}, {{Name: "name4"}, {Name: "name5"}, {Name: "name6"}}, {{Name: "name7"}, {Name: "name8"}, {Name: "name9"}}},
		SliceMap:         []map[string]string{{"key1": "val1", "key2": "val2", "key3": "val3"}, {"key4": "val4", "key5": "val5", "key6": "val6"}},
		Map:              map[string]string{"key1": "val1", "key2": "val2", "key3": "val3"},
		MapStructs:       map[string]*SliceStruct{"key1": {Name: "name1"}, "key2": {Name: "name2"}, "key3": {Name: "name3"}},
		MapMap:           map[string]map[string]string{"key1": {"key1-1": "val1"}, "key2": {"key2-1": "val2"}, "key3": {"key3-1": "val3"}},
		MapMapStruct:     map[string]map[string]*SliceStruct{"key1": {"key1-1": {Name: "name1"}}, "key2": {"key2-1": {Name: "name2"}}, "key3": {"key3-1": {Name: "name3"}}},
		MapSlice:         map[string][]string{"key1": {"1", "2", "3"}, "key2": {"4", "5", "6"}, "key3": {"7", "8", "9"}},
	}

	test = &Test{
		Inner:     inner,
		CreatedAt: nil,
	}

	val = reflect.ValueOf(test)

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.SliceStructs[2]")
	Equal(t, ok, true)
	Equal(t, kind, reflect.Ptr)
	Equal(t, current.String(), "<*validator.SliceStruct Value>")
	Equal(t, current.IsNil(), true)

	current, kind, _, ok = v.getStructFieldOKInternal(val, "Inner.SliceStructs[2].Name")
	Equal(t, ok, false)
	Equal(t, kind, reflect.Ptr)
	Equal(t, current.String(), "<*validator.SliceStruct Value>")
	Equal(t, current.IsNil(), true)

	PanicMatches(t, func() { v.getStructFieldOKInternal(reflect.ValueOf(1), "crazyinput") }, "Invalid field namespace")
}

func TestExistsValidation(t *testing.T) {
	jsonText := "{ \"truthiness2\": true }"

	type Thing struct {
		Truthiness *bool `json:"truthiness" validate:"required"`
	}

	var ting Thing

	err := json.Unmarshal([]byte(jsonText), &ting)
	Equal(t, err, nil)
	NotEqual(t, ting, nil)
	Equal(t, ting.Truthiness, nil)

	validate := New()
	errs := validate.Struct(ting)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Thing.Truthiness", "Thing.Truthiness", "Truthiness", "Truthiness", "required")

	jsonText = "{ \"truthiness\": true }"

	err = json.Unmarshal([]byte(jsonText), &ting)
	Equal(t, err, nil)
	NotEqual(t, ting, nil)
	Equal(t, ting.Truthiness, true)

	errs = validate.Struct(ting)
	Equal(t, errs, nil)
}

func TestSQLValue2Validation(t *testing.T) {
	validate := New()
	validate.RegisterCustomTypeFunc(ValidateValuerType, valuer{}, (*driver.Valuer)(nil), sql.NullString{}, sql.NullInt64{}, sql.NullBool{}, sql.NullFloat64{})
	validate.RegisterCustomTypeFunc(ValidateCustomType, MadeUpCustomType{})
	validate.RegisterCustomTypeFunc(OverrideIntTypeForSomeReason, 1)

	val := valuer{
		Name: "",
	}

	errs := validate.Var(val, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	val.Name = "Valid Name"
	errs = validate.VarCtx(context.Background(), val, "required")
	Equal(t, errs, nil)

	val.Name = "errorme"

	PanicMatches(t, func() { _ = validate.Var(val, "required") }, "SQL Driver Valuer error: some kind of error")

	myVal := valuer{
		Name: "",
	}

	errs = validate.Var(myVal, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	cust := MadeUpCustomType{
		FirstName: "Joey",
		LastName:  "Bloggs",
	}

	c := CustomMadeUpStruct{MadeUp: cust, OverriddenInt: 2}

	errs = validate.Struct(c)
	Equal(t, errs, nil)

	c.MadeUp.FirstName = ""
	c.OverriddenInt = 1

	errs = validate.Struct(c)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 2)
	AssertError(t, errs, "CustomMadeUpStruct.MadeUp", "CustomMadeUpStruct.MadeUp", "MadeUp", "MadeUp", "required")
	AssertError(t, errs, "CustomMadeUpStruct.OverriddenInt", "CustomMadeUpStruct.OverriddenInt", "OverriddenInt", "OverriddenInt", "gt")
}

func TestSQLValueValidation(t *testing.T) {
	validate := New()
	validate.RegisterCustomTypeFunc(ValidateValuerType, (*driver.Valuer)(nil), valuer{})
	validate.RegisterCustomTypeFunc(ValidateCustomType, MadeUpCustomType{})
	validate.RegisterCustomTypeFunc(OverrideIntTypeForSomeReason, 1)

	val := valuer{
		Name: "",
	}

	errs := validate.Var(val, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	val.Name = "Valid Name"
	errs = validate.Var(val, "required")
	Equal(t, errs, nil)

	val.Name = "errorme"

	PanicMatches(t, func() { errs = validate.Var(val, "required") }, "SQL Driver Valuer error: some kind of error")

	myVal := valuer{
		Name: "",
	}

	errs = validate.Var(myVal, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	cust := MadeUpCustomType{
		FirstName: "Joey",
		LastName:  "Bloggs",
	}

	c := CustomMadeUpStruct{MadeUp: cust, OverriddenInt: 2}

	errs = validate.Struct(c)
	Equal(t, errs, nil)

	c.MadeUp.FirstName = ""
	c.OverriddenInt = 1

	errs = validate.Struct(c)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 2)
	AssertError(t, errs, "CustomMadeUpStruct.MadeUp", "CustomMadeUpStruct.MadeUp", "MadeUp", "MadeUp", "required")
	AssertError(t, errs, "CustomMadeUpStruct.OverriddenInt", "CustomMadeUpStruct.OverriddenInt", "OverriddenInt", "OverriddenInt", "gt")

	// Test for empty field on structs without tags
	type InvalidValuePanicSafetyTest struct {
		V valuer
	}
	errs = validate.Struct(InvalidValuePanicSafetyTest{})
	Equal(t, errs, nil)
}

func TestMACValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"3D:F2:C9:A6:B3:4F", true},
		{"3D-F2-C9-A6-B3:4F", false},
		{"123", false},
		{"", false},
		{"abacaba", false},
		{"00:25:96:FF:FE:12:34:56", true},
		{"0025:96FF:FE12:3456", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "mac")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d mac failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d mac failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "mac" {
					t.Fatalf("Index: %d mac failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestIPValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.0.1", true},
		{"192.168.255.254", true},
		{"192.168.255.256", false},
		{"172.16.255.254", true},
		{"172.16.256.255", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", true},
		{"2001:cdba:0:0:0:0:3257:9652", true},
		{"2001:cdba::3257:9652", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ip")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ip" {
					t.Fatalf("Index: %d ip failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestIPv6Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"10.0.0.1", false},
		{"172.16.0.1", false},
		{"192.168.0.1", false},
		{"192.168.255.254", false},
		{"192.168.255.256", false},
		{"172.16.255.254", false},
		{"172.16.256.255", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", true},
		{"2001:cdba:0:0:0:0:3257:9652", true},
		{"2001:cdba::3257:9652", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ipv6")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ipv6 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ipv6 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ipv6" {
					t.Fatalf("Index: %d ipv6 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestIPv4Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.0.1", true},
		{"192.168.255.254", true},
		{"192.168.255.256", false},
		{"172.16.255.254", true},
		{"172.16.256.255", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", false},
		{"2001:cdba:0:0:0:0:3257:9652", false},
		{"2001:cdba::3257:9652", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ipv4")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ipv4 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ipv4 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ipv4" {
					t.Fatalf("Index: %d ipv4 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestCIDRValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"10.0.0.0/0", true},
		{"10.0.0.1/8", true},
		{"172.16.0.1/16", true},
		{"192.168.0.1/24", true},
		{"192.168.255.254/24", true},
		{"192.168.255.254/48", false},
		{"192.168.255.256/24", false},
		{"172.16.255.254/16", true},
		{"172.16.256.255/16", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652/64", true},
		{"2001:cdba:0000:0000:0000:0000:3257:9652/256", false},
		{"2001:cdba:0:0:0:0:3257:9652/32", true},
		{"2001:cdba::3257:9652/16", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "cidr")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d cidr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d cidr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "cidr" {
					t.Fatalf("Index: %d cidr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestCIDRv6Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"10.0.0.0/0", false},
		{"10.0.0.1/8", false},
		{"172.16.0.1/16", false},
		{"192.168.0.1/24", false},
		{"192.168.255.254/24", false},
		{"192.168.255.254/48", false},
		{"192.168.255.256/24", false},
		{"172.16.255.254/16", false},
		{"172.16.256.255/16", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652/64", true},
		{"2001:cdba:0000:0000:0000:0000:3257:9652/256", false},
		{"2001:cdba:0:0:0:0:3257:9652/32", true},
		{"2001:cdba::3257:9652/16", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "cidrv6")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d cidrv6 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d cidrv6 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "cidrv6" {
					t.Fatalf("Index: %d cidrv6 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestCIDRv4Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"0.0.0.0/0", true},
		{"10.0.0.0/0", false},
		{"10.0.0.0/8", true},
		{"10.0.0.1/8", false},
		{"172.16.0.0/16", true},
		{"172.16.0.1/16", false},
		{"192.168.0.0/24", true},
		{"192.168.0.1/24", false},
		{"192.168.255.0/24", true},
		{"192.168.255.254/24", false},
		{"192.168.255.254/48", false},
		{"192.168.255.256/24", false},
		{"172.16.0.0/16", true},
		{"172.16.255.254/16", false},
		{"172.16.256.255/16", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652/64", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652/256", false},
		{"2001:cdba:0:0:0:0:3257:9652/32", false},
		{"2001:cdba::3257:9652/16", false},
		{"172.56.1.0/16", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "cidrv4")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d cidrv4 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d cidrv4 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "cidrv4" {
					t.Fatalf("Index: %d cidrv4 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestTCPAddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{":80", false},
		{"127.0.0.1:80", true},
		{"[::1]:80", true},
		{"256.0.0.0:1", false},
		{"[::1]", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "tcp_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d tcp_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d tcp_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "tcp_addr" {
					t.Fatalf("Index: %d tcp_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestTCP6AddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{":80", false},
		{"127.0.0.1:80", false},
		{"[::1]:80", true},
		{"256.0.0.0:1", false},
		{"[::1]", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "tcp6_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d tcp6_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d tcp6_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "tcp6_addr" {
					t.Fatalf("Index: %d tcp6_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestTCP4AddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{":80", false},
		{"127.0.0.1:80", true},
		{"[::1]:80", false}, // https://github.com/golang/go/issues/14037
		{"256.0.0.0:1", false},
		{"[::1]", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "tcp4_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d tcp4_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Log(test.param, IsEqual(errs, nil))
				t.Fatalf("Index: %d tcp4_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "tcp4_addr" {
					t.Fatalf("Index: %d tcp4_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUDPAddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{":80", false},
		{"127.0.0.1:80", true},
		{"[::1]:80", true},
		{"256.0.0.0:1", false},
		{"[::1]", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "udp_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d udp_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d udp_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "udp_addr" {
					t.Fatalf("Index: %d udp_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUDP6AddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{":80", false},
		{"127.0.0.1:80", false},
		{"[::1]:80", true},
		{"256.0.0.0:1", false},
		{"[::1]", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "udp6_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d udp6_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d udp6_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "udp6_addr" {
					t.Fatalf("Index: %d udp6_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUDP4AddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{":80", false},
		{"127.0.0.1:80", true},
		{"[::1]:80", false}, // https://github.com/golang/go/issues/14037
		{"256.0.0.0:1", false},
		{"[::1]", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "udp4_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d udp4_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Log(test.param, IsEqual(errs, nil))
				t.Fatalf("Index: %d udp4_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "udp4_addr" {
					t.Fatalf("Index: %d udp4_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestIPAddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"127.0.0.1:80", false},
		{"::1", true},
		{"256.0.0.0", false},
		{"localhost", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "ip_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ip_addr" {
					t.Fatalf("Index: %d ip_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestIP6AddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"127.0.0.1", false}, // https://github.com/golang/go/issues/14037
		{"127.0.0.1:80", false},
		{"::1", true},
		{"0:0:0:0:0:0:0:1", true},
		{"256.0.0.0", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "ip6_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip6_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip6_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ip6_addr" {
					t.Fatalf("Index: %d ip6_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestIP4AddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"127.0.0.1", true},
		{"127.0.0.1:80", false},
		{"::1", false}, // https://github.com/golang/go/issues/14037
		{"256.0.0.0", false},
		{"localhost", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "ip4_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ip4_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Log(test.param, IsEqual(errs, nil))
				t.Fatalf("Index: %d ip4_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ip4_addr" {
					t.Fatalf("Index: %d ip4_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUnixAddrValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"v.sock", true},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "unix_addr")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d unix_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Log(test.param, IsEqual(errs, nil))
				t.Fatalf("Index: %d unix_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "unix_addr" {
					t.Fatalf("Index: %d unix_addr failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestSliceMapArrayChanFuncPtrInterfaceRequiredValidation(t *testing.T) {
	validate := New()

	var m map[string]string

	errs := validate.Var(m, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	m = map[string]string{}
	errs = validate.Var(m, "required")
	Equal(t, errs, nil)

	var arr [5]string
	errs = validate.Var(arr, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	arr[0] = "ok"
	errs = validate.Var(arr, "required")
	Equal(t, errs, nil)

	var s []string
	errs = validate.Var(s, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	s = []string{}
	errs = validate.Var(s, "required")
	Equal(t, errs, nil)

	var c chan string
	errs = validate.Var(c, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	c = make(chan string)
	errs = validate.Var(c, "required")
	Equal(t, errs, nil)

	var tst *int
	errs = validate.Var(tst, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	one := 1
	tst = &one
	errs = validate.Var(tst, "required")
	Equal(t, errs, nil)

	var iface interface{}

	errs = validate.Var(iface, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	errs = validate.Var(iface, "omitempty,required")
	Equal(t, errs, nil)

	errs = validate.Var(iface, "")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(nil, iface, "")
	Equal(t, errs, nil)

	var f func(string)

	errs = validate.Var(f, "required")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "required")

	f = func(name string) {}

	errs = validate.Var(f, "required")
	Equal(t, errs, nil)
}

func TestDatePtrValidationIssueValidation(t *testing.T) {
	type Test struct {
		LastViewed *time.Time
		Reminder   *time.Time
	}

	test := &Test{}

	validate := New()
	errs := validate.Struct(test)
	Equal(t, errs, nil)
}

func TestCommaAndPipeObfuscationValidation(t *testing.T) {
	s := "My Name Is, |joeybloggs|"

	validate := New()

	errs := validate.Var(s, "excludesall=0x2C")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "excludesall")

	errs = validate.Var(s, "excludesall=0x7C")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "excludesall")
}

func TestBadKeyValidation(t *testing.T) {
	type Test struct {
		Name string `validate:"required, "`
	}

	tst := &Test{
		Name: "test",
	}

	validate := New()

	PanicMatches(t, func() { _ = validate.Struct(tst) }, "Undefined validation function ' ' on field 'Name'")

	type Test2 struct {
		Name string `validate:"required,,len=2"`
	}

	tst2 := &Test2{
		Name: "test",
	}

	PanicMatches(t, func() { _ = validate.Struct(tst2) }, "Invalid validation tag on field 'Name'")
}

func TestInterfaceErrValidation(t *testing.T) {
	var v2 interface{} = 1
	var v1 interface{} = v2

	validate := New()
	errs := validate.Var(v1, "len=1")
	Equal(t, errs, nil)

	errs = validate.Var(v2, "len=1")
	Equal(t, errs, nil)

	type ExternalCMD struct {
		Userid string      `json:"userid"`
		Action uint32      `json:"action"`
		Data   interface{} `json:"data,omitempty" validate:"required"`
	}

	s := &ExternalCMD{
		Userid: "123456",
		Action: 10000,
		// Data:   1,
	}

	errs = validate.Struct(s)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "ExternalCMD.Data", "ExternalCMD.Data", "Data", "Data", "required")

	type ExternalCMD2 struct {
		Userid string      `json:"userid"`
		Action uint32      `json:"action"`
		Data   interface{} `json:"data,omitempty" validate:"len=1"`
	}

	s2 := &ExternalCMD2{
		Userid: "123456",
		Action: 10000,
		// Data:   1,
	}

	errs = validate.Struct(s2)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "ExternalCMD2.Data", "ExternalCMD2.Data", "Data", "Data", "len")

	s3 := &ExternalCMD2{
		Userid: "123456",
		Action: 10000,
		Data:   2,
	}

	errs = validate.Struct(s3)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "ExternalCMD2.Data", "ExternalCMD2.Data", "Data", "Data", "len")

	type Inner struct {
		Name string `validate:"required"`
	}

	inner := &Inner{
		Name: "",
	}

	s4 := &ExternalCMD{
		Userid: "123456",
		Action: 10000,
		Data:   inner,
	}

	errs = validate.Struct(s4)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "ExternalCMD.Data.Name", "ExternalCMD.Data.Name", "Name", "Name", "required")

	type TestMapStructPtr struct {
		Errs map[int]interface{} `validate:"gt=0,dive,required"`
	}

	mip := map[int]interface{}{0: &Inner{"ok"}, 3: nil, 4: &Inner{"ok"}}

	msp := &TestMapStructPtr{
		Errs: mip,
	}

	errs = validate.Struct(msp)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestMapStructPtr.Errs[3]", "TestMapStructPtr.Errs[3]", "Errs[3]", "Errs[3]", "required")

	type TestMultiDimensionalStructs struct {
		Errs [][]interface{} `validate:"gt=0,dive,dive"`
	}

	var errStructArray [][]interface{}

	errStructArray = append(errStructArray, []interface{}{&Inner{"ok"}, &Inner{""}, &Inner{""}})
	errStructArray = append(errStructArray, []interface{}{&Inner{"ok"}, &Inner{""}, &Inner{""}})

	tms := &TestMultiDimensionalStructs{
		Errs: errStructArray,
	}

	errs = validate.Struct(tms)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 4)
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[0][1].Name", "TestMultiDimensionalStructs.Errs[0][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[0][2].Name", "TestMultiDimensionalStructs.Errs[0][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[1][1].Name", "TestMultiDimensionalStructs.Errs[1][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[1][2].Name", "TestMultiDimensionalStructs.Errs[1][2].Name", "Name", "Name", "required")

	type TestMultiDimensionalStructsPtr2 struct {
		Errs [][]*Inner `validate:"gt=0,dive,dive,required"`
	}

	var errStructPtr2Array [][]*Inner

	errStructPtr2Array = append(errStructPtr2Array, []*Inner{{"ok"}, {""}, {""}})
	errStructPtr2Array = append(errStructPtr2Array, []*Inner{{"ok"}, {""}, {""}})
	errStructPtr2Array = append(errStructPtr2Array, []*Inner{{"ok"}, {""}, nil})

	tmsp2 := &TestMultiDimensionalStructsPtr2{
		Errs: errStructPtr2Array,
	}

	errs = validate.Struct(tmsp2)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 6)
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[0][1].Name", "TestMultiDimensionalStructsPtr2.Errs[0][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[0][2].Name", "TestMultiDimensionalStructsPtr2.Errs[0][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[1][1].Name", "TestMultiDimensionalStructsPtr2.Errs[1][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[1][2].Name", "TestMultiDimensionalStructsPtr2.Errs[1][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[2][1].Name", "TestMultiDimensionalStructsPtr2.Errs[2][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[2][2]", "TestMultiDimensionalStructsPtr2.Errs[2][2]", "Errs[2][2]", "Errs[2][2]", "required")

	m := map[int]interface{}{0: "ok", 3: "", 4: "ok"}

	errs = validate.Var(m, "len=3,dive,len=2")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "[3]", "[3]", "[3]", "[3]", "len")

	errs = validate.Var(m, "len=2,dive,required")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "", "", "", "", "len")

	arr := []interface{}{"ok", "", "ok"}

	errs = validate.Var(arr, "len=3,dive,len=2")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "[1]", "[1]", "[1]", "[1]", "len")

	errs = validate.Var(arr, "len=2,dive,required")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "", "", "", "", "len")

	type MyStruct struct {
		A, B string
		C    interface{}
	}

	var a MyStruct

	a.A = "value"
	a.C = "nu"

	errs = validate.Struct(a)
	Equal(t, errs, nil)
}

func TestMapDiveValidation(t *testing.T) {
	validate := New()

	n := map[int]interface{}{0: nil}
	errs := validate.Var(n, "omitempty,required")
	Equal(t, errs, nil)

	m := map[int]string{0: "ok", 3: "", 4: "ok"}

	errs = validate.Var(m, "len=3,dive,required")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "[3]", "[3]", "[3]", "[3]", "required")

	errs = validate.Var(m, "len=2,dive,required")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "", "", "", "", "len")

	type Inner struct {
		Name string `validate:"required"`
	}

	type TestMapStruct struct {
		Errs map[int]Inner `validate:"gt=0,dive"`
	}

	mi := map[int]Inner{0: {"ok"}, 3: {""}, 4: {"ok"}}

	ms := &TestMapStruct{
		Errs: mi,
	}

	errs = validate.Struct(ms)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestMapStruct.Errs[3].Name", "TestMapStruct.Errs[3].Name", "Name", "Name", "required")

	// for full test coverage
	s := fmt.Sprint(errs.Error())
	NotEqual(t, s, "")

	type TestMapInterface struct {
		Errs map[int]interface{} `validate:"dive"`
	}

	mit := map[int]interface{}{0: Inner{"ok"}, 1: Inner{""}, 3: nil, 5: "string", 6: 33}

	msi := &TestMapInterface{
		Errs: mit,
	}

	errs = validate.Struct(msi)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestMapInterface.Errs[1].Name", "TestMapInterface.Errs[1].Name", "Name", "Name", "required")

	type TestMapTimeStruct struct {
		Errs map[int]*time.Time `validate:"gt=0,dive,required"`
	}

	t1 := time.Now().UTC()

	mta := map[int]*time.Time{0: &t1, 3: nil, 4: nil}

	mt := &TestMapTimeStruct{
		Errs: mta,
	}

	errs = validate.Struct(mt)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 2)
	AssertError(t, errs, "TestMapTimeStruct.Errs[3]", "TestMapTimeStruct.Errs[3]", "Errs[3]", "Errs[3]", "required")
	AssertError(t, errs, "TestMapTimeStruct.Errs[4]", "TestMapTimeStruct.Errs[4]", "Errs[4]", "Errs[4]", "required")

	type TestMapStructPtr struct {
		Errs map[int]*Inner `validate:"gt=0,dive,required"`
	}

	mip := map[int]*Inner{0: {"ok"}, 3: nil, 4: {"ok"}}

	msp := &TestMapStructPtr{
		Errs: mip,
	}

	errs = validate.Struct(msp)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestMapStructPtr.Errs[3]", "TestMapStructPtr.Errs[3]", "Errs[3]", "Errs[3]", "required")

	type TestMapStructPtr2 struct {
		Errs map[int]*Inner `validate:"gt=0,dive,omitempty,required"`
	}

	mip2 := map[int]*Inner{0: {"ok"}, 3: nil, 4: {"ok"}}

	msp2 := &TestMapStructPtr2{
		Errs: mip2,
	}

	errs = validate.Struct(msp2)
	Equal(t, errs, nil)

	v2 := New()
	v2.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type MapDiveJSONTest struct {
		Map map[string]string `validate:"required,gte=1,dive,gte=1" json:"MyName"`
	}

	mdjt := &MapDiveJSONTest{
		Map: map[string]string{
			"Key1": "Value1",
			"Key2": "",
		},
	}

	err := v2.Struct(mdjt)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)
	fe := getError(errs, "MapDiveJSONTest.MyName[Key2]", "MapDiveJSONTest.Map[Key2]")
	NotEqual(t, fe, nil)
	Equal(t, fe.Tag(), "gte")
	Equal(t, fe.ActualTag(), "gte")
	Equal(t, fe.Field(), "MyName[Key2]")
	Equal(t, fe.StructField(), "Map[Key2]")
}

func TestArrayDiveValidation(t *testing.T) {
	validate := New()

	arr := []string{"ok", "", "ok"}

	errs := validate.Var(arr, "len=3,dive,required")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "[1]", "[1]", "[1]", "[1]", "required")

	errs = validate.Var(arr, "len=2,dive,required")
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "", "", "", "", "len")

	type BadDive struct {
		Name string `validate:"dive"`
	}

	bd := &BadDive{
		Name: "TEST",
	}

	PanicMatches(t, func() { _ = validate.Struct(bd) }, "dive error! can't dive on a non slice or map")

	type Test struct {
		Errs []string `validate:"gt=0,dive,required"`
	}

	test := &Test{
		Errs: []string{"ok", "", "ok"},
	}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "Test.Errs[1]", "Test.Errs[1]", "Errs[1]", "Errs[1]", "required")

	test = &Test{
		Errs: []string{"ok", "ok", ""},
	}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "Test.Errs[2]", "Test.Errs[2]", "Errs[2]", "Errs[2]", "required")

	type TestMultiDimensional struct {
		Errs [][]string `validate:"gt=0,dive,dive,required"`
	}

	var errArray [][]string

	errArray = append(errArray, []string{"ok", "", ""})
	errArray = append(errArray, []string{"ok", "", ""})

	tm := &TestMultiDimensional{
		Errs: errArray,
	}

	errs = validate.Struct(tm)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 4)
	AssertError(t, errs, "TestMultiDimensional.Errs[0][1]", "TestMultiDimensional.Errs[0][1]", "Errs[0][1]", "Errs[0][1]", "required")
	AssertError(t, errs, "TestMultiDimensional.Errs[0][2]", "TestMultiDimensional.Errs[0][2]", "Errs[0][2]", "Errs[0][2]", "required")
	AssertError(t, errs, "TestMultiDimensional.Errs[1][1]", "TestMultiDimensional.Errs[1][1]", "Errs[1][1]", "Errs[1][1]", "required")
	AssertError(t, errs, "TestMultiDimensional.Errs[1][2]", "TestMultiDimensional.Errs[1][2]", "Errs[1][2]", "Errs[1][2]", "required")

	type Inner struct {
		Name string `validate:"required"`
	}

	type TestMultiDimensionalStructs struct {
		Errs [][]Inner `validate:"gt=0,dive,dive"`
	}

	var errStructArray [][]Inner

	errStructArray = append(errStructArray, []Inner{{"ok"}, {""}, {""}})
	errStructArray = append(errStructArray, []Inner{{"ok"}, {""}, {""}})

	tms := &TestMultiDimensionalStructs{
		Errs: errStructArray,
	}

	errs = validate.Struct(tms)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 4)
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[0][1].Name", "TestMultiDimensionalStructs.Errs[0][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[0][2].Name", "TestMultiDimensionalStructs.Errs[0][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[1][1].Name", "TestMultiDimensionalStructs.Errs[1][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructs.Errs[1][2].Name", "TestMultiDimensionalStructs.Errs[1][2].Name", "Name", "Name", "required")

	type TestMultiDimensionalStructsPtr struct {
		Errs [][]*Inner `validate:"gt=0,dive,dive"`
	}

	var errStructPtrArray [][]*Inner

	errStructPtrArray = append(errStructPtrArray, []*Inner{{"ok"}, {""}, {""}})
	errStructPtrArray = append(errStructPtrArray, []*Inner{{"ok"}, {""}, {""}})
	errStructPtrArray = append(errStructPtrArray, []*Inner{{"ok"}, {""}, nil})

	tmsp := &TestMultiDimensionalStructsPtr{
		Errs: errStructPtrArray,
	}

	errs = validate.Struct(tmsp)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 5)
	AssertError(t, errs, "TestMultiDimensionalStructsPtr.Errs[0][1].Name", "TestMultiDimensionalStructsPtr.Errs[0][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr.Errs[0][2].Name", "TestMultiDimensionalStructsPtr.Errs[0][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr.Errs[1][1].Name", "TestMultiDimensionalStructsPtr.Errs[1][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr.Errs[1][2].Name", "TestMultiDimensionalStructsPtr.Errs[1][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr.Errs[2][1].Name", "TestMultiDimensionalStructsPtr.Errs[2][1].Name", "Name", "Name", "required")

	// for full test coverage
	s := fmt.Sprint(errs.Error())
	NotEqual(t, s, "")

	type TestMultiDimensionalStructsPtr2 struct {
		Errs [][]*Inner `validate:"gt=0,dive,dive,required"`
	}

	var errStructPtr2Array [][]*Inner

	errStructPtr2Array = append(errStructPtr2Array, []*Inner{{"ok"}, {""}, {""}})
	errStructPtr2Array = append(errStructPtr2Array, []*Inner{{"ok"}, {""}, {""}})
	errStructPtr2Array = append(errStructPtr2Array, []*Inner{{"ok"}, {""}, nil})

	tmsp2 := &TestMultiDimensionalStructsPtr2{
		Errs: errStructPtr2Array,
	}

	errs = validate.Struct(tmsp2)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 6)
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[0][1].Name", "TestMultiDimensionalStructsPtr2.Errs[0][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[0][2].Name", "TestMultiDimensionalStructsPtr2.Errs[0][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[1][1].Name", "TestMultiDimensionalStructsPtr2.Errs[1][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[1][2].Name", "TestMultiDimensionalStructsPtr2.Errs[1][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[2][1].Name", "TestMultiDimensionalStructsPtr2.Errs[2][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr2.Errs[2][2]", "TestMultiDimensionalStructsPtr2.Errs[2][2]", "Errs[2][2]", "Errs[2][2]", "required")

	type TestMultiDimensionalStructsPtr3 struct {
		Errs [][]*Inner `validate:"gt=0,dive,dive,omitempty"`
	}

	var errStructPtr3Array [][]*Inner

	errStructPtr3Array = append(errStructPtr3Array, []*Inner{{"ok"}, {""}, {""}})
	errStructPtr3Array = append(errStructPtr3Array, []*Inner{{"ok"}, {""}, {""}})
	errStructPtr3Array = append(errStructPtr3Array, []*Inner{{"ok"}, {""}, nil})

	tmsp3 := &TestMultiDimensionalStructsPtr3{
		Errs: errStructPtr3Array,
	}

	errs = validate.Struct(tmsp3)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 5)
	AssertError(t, errs, "TestMultiDimensionalStructsPtr3.Errs[0][1].Name", "TestMultiDimensionalStructsPtr3.Errs[0][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr3.Errs[0][2].Name", "TestMultiDimensionalStructsPtr3.Errs[0][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr3.Errs[1][1].Name", "TestMultiDimensionalStructsPtr3.Errs[1][1].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr3.Errs[1][2].Name", "TestMultiDimensionalStructsPtr3.Errs[1][2].Name", "Name", "Name", "required")
	AssertError(t, errs, "TestMultiDimensionalStructsPtr3.Errs[2][1].Name", "TestMultiDimensionalStructsPtr3.Errs[2][1].Name", "Name", "Name", "required")

	type TestMultiDimensionalTimeTime struct {
		Errs [][]*time.Time `validate:"gt=0,dive,dive,required"`
	}

	var errTimePtr3Array [][]*time.Time

	t1 := time.Now().UTC()
	t2 := time.Now().UTC()
	t3 := time.Now().UTC().Add(time.Hour * 24)

	errTimePtr3Array = append(errTimePtr3Array, []*time.Time{&t1, &t2, &t3})
	errTimePtr3Array = append(errTimePtr3Array, []*time.Time{&t1, &t2, nil})
	errTimePtr3Array = append(errTimePtr3Array, []*time.Time{&t1, nil, nil})

	tmtp3 := &TestMultiDimensionalTimeTime{
		Errs: errTimePtr3Array,
	}

	errs = validate.Struct(tmtp3)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 3)
	AssertError(t, errs, "TestMultiDimensionalTimeTime.Errs[1][2]", "TestMultiDimensionalTimeTime.Errs[1][2]", "Errs[1][2]", "Errs[1][2]", "required")
	AssertError(t, errs, "TestMultiDimensionalTimeTime.Errs[2][1]", "TestMultiDimensionalTimeTime.Errs[2][1]", "Errs[2][1]", "Errs[2][1]", "required")
	AssertError(t, errs, "TestMultiDimensionalTimeTime.Errs[2][2]", "TestMultiDimensionalTimeTime.Errs[2][2]", "Errs[2][2]", "Errs[2][2]", "required")

	type TestMultiDimensionalTimeTime2 struct {
		Errs [][]*time.Time `validate:"gt=0,dive,dive,required"`
	}

	var errTimeArray [][]*time.Time

	t1 = time.Now().UTC()
	t2 = time.Now().UTC()
	t3 = time.Now().UTC().Add(time.Hour * 24)

	errTimeArray = append(errTimeArray, []*time.Time{&t1, &t2, &t3})
	errTimeArray = append(errTimeArray, []*time.Time{&t1, &t2, nil})
	errTimeArray = append(errTimeArray, []*time.Time{&t1, nil, nil})

	tmtp := &TestMultiDimensionalTimeTime2{
		Errs: errTimeArray,
	}

	errs = validate.Struct(tmtp)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 3)
	AssertError(t, errs, "TestMultiDimensionalTimeTime2.Errs[1][2]", "TestMultiDimensionalTimeTime2.Errs[1][2]", "Errs[1][2]", "Errs[1][2]", "required")
	AssertError(t, errs, "TestMultiDimensionalTimeTime2.Errs[2][1]", "TestMultiDimensionalTimeTime2.Errs[2][1]", "Errs[2][1]", "Errs[2][1]", "required")
	AssertError(t, errs, "TestMultiDimensionalTimeTime2.Errs[2][2]", "TestMultiDimensionalTimeTime2.Errs[2][2]", "Errs[2][2]", "Errs[2][2]", "required")
}

func TestNilStructPointerValidation(t *testing.T) {
	type Inner struct {
		Data string
	}

	type Outer struct {
		Inner *Inner `validate:"omitempty"`
	}

	inner := &Inner{
		Data: "test",
	}

	outer := &Outer{
		Inner: inner,
	}

	validate := New()
	errs := validate.Struct(outer)
	Equal(t, errs, nil)

	outer = &Outer{
		Inner: nil,
	}

	errs = validate.Struct(outer)
	Equal(t, errs, nil)

	type Inner2 struct {
		Data string
	}

	type Outer2 struct {
		Inner2 *Inner2 `validate:"required"`
	}

	inner2 := &Inner2{
		Data: "test",
	}

	outer2 := &Outer2{
		Inner2: inner2,
	}

	errs = validate.Struct(outer2)
	Equal(t, errs, nil)

	outer2 = &Outer2{
		Inner2: nil,
	}

	errs = validate.Struct(outer2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Outer2.Inner2", "Outer2.Inner2", "Inner2", "Inner2", "required")

	type Inner3 struct {
		Data string
	}

	type Outer3 struct {
		Inner3 *Inner3
	}

	inner3 := &Inner3{
		Data: "test",
	}

	outer3 := &Outer3{
		Inner3: inner3,
	}

	errs = validate.Struct(outer3)
	Equal(t, errs, nil)

	type Inner4 struct {
		Data string
	}

	type Outer4 struct {
		Inner4 *Inner4 `validate:"-"`
	}

	inner4 := &Inner4{
		Data: "test",
	}

	outer4 := &Outer4{
		Inner4: inner4,
	}

	errs = validate.Struct(outer4)
	Equal(t, errs, nil)
}

func TestSSNValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"00-90-8787", false},
		{"66690-76", false},
		{"191 60 2869", true},
		{"191-60-2869", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ssn")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d SSN failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d SSN failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ssn" {
					t.Fatalf("Index: %d Latitude failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestLongitudeValidation(t *testing.T) {
	tests := []struct {
		param    interface{}
		expected bool
	}{
		{"", false},
		{"-180.000", true},
		{"180.1", false},
		{"+73.234", true},
		{"+382.3811", false},
		{"23.11111111", true},
		{uint(180), true},
		{float32(-180.0), true},
		{-180, true},
		{180.1, false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "longitude")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d Longitude failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d Longitude failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "longitude" {
					t.Fatalf("Index: %d Longitude failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() { _ = validate.Var(true, "longitude") }, "Bad field type bool")
}

func TestLatitudeValidation(t *testing.T) {
	tests := []struct {
		param    interface{}
		expected bool
	}{
		{"", false},
		{"-90.000", true},
		{"+90", true},
		{"47.1231231", true},
		{"+99.9", false},
		{"108", false},
		{uint(90), true},
		{float32(-90.0), true},
		{-90, true},
		{90.1, false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "latitude")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d Latitude failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d Latitude failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "latitude" {
					t.Fatalf("Index: %d Latitude failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() { _ = validate.Var(true, "latitude") }, "Bad field type bool")
}

func TestDataURIValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"data:image/png;base64,TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4=", true},
		{"data:text/plain;base64,Vml2YW11cyBmZXJtZW50dW0gc2VtcGVyIHBvcnRhLg==", true},
		{"image/gif;base64,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", false},
		{
			"data:image/gif;base64,MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuMPNS1Ufof9EW/M98FNw" +
				"UAKrwflsqVxaxQjBQnHQmiI7Vac40t8x7pIb8gLGV6wL7sBTJiPovJ0V7y7oc0Ye" +
				"rhKh0Rm4skP2z/jHwwZICgGzBvA0rH8xlhUiTvcwDCJ0kc+fh35hNt8srZQM4619" +
				"FTgB66Xmp4EtVyhpQV+t02g6NzK72oZI0vnAvqhpkxLeLiMCyrI416wHm5Tkukhx" +
				"QmcL2a6hNOyu0ixX/x2kSFXApEnVrJ+/IxGyfyw8kf4N2IZpW5nEP847lpfj0SZZ" +
				"Fwrd1mnfnDbYohX2zRptLy2ZUn06Qo9pkG5ntvFEPo9bfZeULtjYzIl6K8gJ2uGZ" + "HQIDAQAB", true,
		},
		{"data:image/png;base64,12345", false},
		{"", false},
		{"data:text,:;base85,U3VzcGVuZGlzc2UgbGVjdHVzIGxlbw==", false},
		{"data:image/jpeg;key=value;base64,UEsDBBQAAAAI", true},
		{"data:image/jpeg;key=value,UEsDBBQAAAAI", true},
		{"data:;base64;sdfgsdfgsdfasdfa=s,UEsDBBQAAAAI", true},
		{"data:,UEsDBBQAAAAI", true},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "datauri")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d DataURI failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d DataURI failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "datauri" {
					t.Fatalf("Index: %d DataURI failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestMultibyteValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"abc", false},
		{"123", false},
		{"<>@;.-=", false},
		{"ひらがな・カタカナ、．漢字", true},
		{"あいうえお foobar", true},
		{"test＠example.com", true},
		{"test＠example.com", true},
		{"1234abcDEｘｙｚ", true},
		{"ｶﾀｶﾅ", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "multibyte")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d Multibyte failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d Multibyte failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "multibyte" {
					t.Fatalf("Index: %d Multibyte failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestPrintableASCIIValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"ｆｏｏbar", false},
		{"ｘｙｚ０９８", false},
		{"１２３456", false},
		{"ｶﾀｶﾅ", false},
		{"foobar", true},
		{"0987654321", true},
		{"test@example.com", true},
		{"1234abcDEF", true},
		{"newline\n", false},
		{"\x19test\x7F", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "printascii")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d Printable ASCII failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d Printable ASCII failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "printascii" {
					t.Fatalf("Index: %d Printable ASCII failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestASCIIValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", true},
		{"ｆｏｏbar", false},
		{"ｘｙｚ０９８", false},
		{"１２３456", false},
		{"ｶﾀｶﾅ", false},
		{"foobar", true},
		{"0987654321", true},
		{"test@example.com", true},
		{"1234abcDEF", true},
		{"", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ascii")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ASCII failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ASCII failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ascii" {
					t.Fatalf("Index: %d ASCII failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUUID5Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{

		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"9c858901-8a57-4791-81fe-4c455b099bc9", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"987fbc97-4bed-5078-af07-9141ba07c9f3", true},
		{"987fbc97-4bed-5078-9f07-9141ba07c9f3", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid5")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID5 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID5 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid5" {
					t.Fatalf("Index: %d UUID5 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUUID4Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-5078-af07-9141ba07c9f3", false},
		{"934859", false},
		{"57b73598-8764-4ad0-a76a-679bb6640eb1", true},
		{"625e63f3-58f5-40b7-83a1-a72ad31acffb", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid4")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID4 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID4 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid4" {
					t.Fatalf("Index: %d UUID4 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUUID3Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"412452646", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-4078-8f07-9141ba07c9f3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid3")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID3 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID3 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid3" {
					t.Fatalf("Index: %d UUID3 failed Error: %s", i, errs)
				}
			}
		}
	}
}

type uuidTestType struct {
	val string
}

func (u uuidTestType) String() string {
	return u.val
}

type uuidAlias string

func (u uuidAlias) String() string {
	return "This is a UUID " + string(u)
}

var _ fmt.Stringer = uuidTestType{}
var _ fmt.Stringer = uuidAlias("")

func TestUUIDValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3xxx", false},
		{"a987fbc94bed3078cf079141ba07c9f3", false},
		{"934859", false},
		{"987fbc9-4bed-3078-cf07a-9141ba07c9f3", false},
		{"aaaaaaaa-1111-1111-aaag-111111111111", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9f3", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid" {
					t.Fatalf("Index: %d UUID failed Error: %s", i, errs)
				}
			}
		}
	}

	// Test UUID validation on uuid structs type that implements Stringer interface.
	structWithValidUUID := struct {
		UUID uuidTestType `validate:"uuid"`
	}{
		UUID: uuidTestType{val: "a987fbc9-4bed-3078-cf07-9141ba07c9f3"},
	}
	structWithInvalidUUID := struct {
		UUID uuidTestType `validate:"uuid"`
	}{
		UUID: uuidTestType{val: "934859"},
	}

	if err := validate.Struct(structWithValidUUID); err != nil {
		t.Fatalf("UUID failed Error: %s", err)
	}
	if err := validate.Struct(structWithInvalidUUID); err == nil {
		t.Fatal("UUID failed Error expected but received nil")
	}

	// Test on Alias type with Stringer interface.
	alias := uuidAlias("a987fbc9-4bed-3078-cf07-9141ba07c9f3")
	if err := validate.Var(alias, "uuid"); err != nil {
		t.Fatalf("UUID failed Error: %s", err)
	}
}

func TestUUID5RFC4122Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{

		{"", false},
		{"xxxa987Fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"9c858901-8a57-4791-81Fe-4c455b099bc9", false},
		{"a987Fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"987Fbc97-4bed-5078-af07-9141ba07c9f3", true},
		{"987Fbc97-4bed-5078-9f07-9141ba07c9f3", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid5_rfc4122")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID5RFC4122 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID5RFC4122 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid5_rfc4122" {
					t.Fatalf("Index: %d UUID5RFC4122 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUUID4RFC4122Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9F3", false},
		{"a987fbc9-4bed-5078-af07-9141ba07c9F3", false},
		{"934859", false},
		{"57b73598-8764-4ad0-a76A-679bb6640eb1", true},
		{"625e63f3-58f5-40b7-83a1-a72ad31acFfb", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid4_rfc4122")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID4RFC4122 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID4RFC4122 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid4_rfc4122" {
					t.Fatalf("Index: %d UUID4RFC4122 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUUID3RFC4122Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"412452646", false},
		{"xxxa987fbc9-4bed-3078-cf07-9141ba07c9F3", false},
		{"a987fbc9-4bed-4078-8f07-9141ba07c9F3", false},
		{"a987fbc9-4bed-3078-cf07-9141ba07c9F3", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid3_rfc4122")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID3RFC4122 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUID3RFC4122 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid3_rfc4122" {
					t.Fatalf("Index: %d UUID3RFC4122 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestUUIDRFC4122Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"xxxa987Fbc9-4bed-3078-cf07-9141ba07c9f3", false},
		{"a987Fbc9-4bed-3078-cf07-9141ba07c9f3xxx", false},
		{"a987Fbc94bed3078cf079141ba07c9f3", false},
		{"934859", false},
		{"987fbc9-4bed-3078-cf07a-9141ba07c9F3", false},
		{"aaaaaaaa-1111-1111-aaaG-111111111111", false},
		{"a987Fbc9-4bed-3078-cf07-9141ba07c9f3", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uuid_rfc4122")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUIDRFC4122 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d UUIDRFC4122 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uuid_rfc4122" {
					t.Fatalf("Index: %d UUIDRFC4122 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestULIDValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"01BX5ZZKBKACT-V9WEVGEMMVRZ", false},
		{"01bx5zzkbkactav9wevgemmvrz", true},
		{"a987Fbc9-4bed-3078-cf07-9141ba07c9f3xxx", false},
		{"01BX5ZZKBKACTAV9WEVGEMMVRZABC", false},
		{"01BX5ZZKBKACTAV9WEVGEMMVRZABC", false},
		{"0IBX5ZZKBKACTAV9WEVGEMMVRZ", false},
		{"O1BX5ZZKBKACTAV9WEVGEMMVRZ", false},
		{"01BX5ZZKBKACTAVLWEVGEMMVRZ", false},
		{"01BX5ZZKBKACTAV9WEVGEMMVRZ", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ulid")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ULID failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ULID failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ulid" {
					t.Fatalf("Index: %d ULID failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestMD4Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "md4")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d MD4 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d MD4 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "md4" {
					t.Fatalf("Index: %d MD4 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestMD5Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "md5")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d MD5 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d MD5 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "md5" {
					t.Fatalf("Index: %d MD5 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestSHA256Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "sha256")
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d SHA256 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d SHA256 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "sha256" {
					t.Fatalf("Index: %d SHA256 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestSHA384Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "sha384")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d SHA384 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d SHA384 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "sha384" {
					t.Fatalf("Index: %d SHA384 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestSHA512Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc46f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "sha512")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d SHA512 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d SHA512 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "sha512" {
					t.Fatalf("Index: %d SHA512 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestRIPEMD128Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ripemd128")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d RIPEMD128 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d RIPEMD128 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ripemd128" {
					t.Fatalf("Index: %d RIPEMD128 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestRIPEMD160Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "ripemd160")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d RIPEMD160 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d RIPEMD160 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "ripemd160" {
					t.Fatalf("Index: %d RIPEMD160 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestTIGER128Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "tiger128")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d TIGER128 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d TIGER128 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "tiger128" {
					t.Fatalf("Index: %d TIGER128 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestTIGER160Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "tiger160")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d TIGER160 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d TIGER160 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "tiger160" {
					t.Fatalf("Index: %d TIGER160 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestTIGER192Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"6f5902ac237024bd6f5902ac237024bdd0c176cb93063dc4", true},
		{"6f5902ac237024bd6f5902ac237024bdd0c176cb93063dc-", false},
		{"6f5902ac237024bd6f5902ac237024bdd0c176cb93063dc41", false},
		{"6f5902ac237024bd6f5902ac237024bdd0c176cb93063dcC", false},
		{"6f5902ac237024bd6f5902ac237024bdd0c176cb93063dc", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "tiger192")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d TIGER192 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d TIGER192 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "tiger192" {
					t.Fatalf("Index: %d TIGER192 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestISBNValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "isbn")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISBN failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISBN failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "isbn" {
					t.Fatalf("Index: %d ISBN failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestISBN13Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"3-8362-2119-5", false},
		{"01234567890ab", false},
		{"978 3 8362 2119 0", false},
		{"9784873113685", true},
		{"978-4-87311-368-5", true},
		{"978 3401013190", true},
		{"978-3-8362-2119-1", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "isbn13")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISBN13 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISBN13 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "isbn13" {
					t.Fatalf("Index: %d ISBN13 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestISBN10Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"3423214121", false},
		{"978-3836221191", false},
		{"3-423-21412-1", false},
		{"3 423 21412 1", false},
		{"3836221195", true},
		{"1-61729-085-8", true},
		{"3 423 21412 0", true},
		{"3 401 01319 X", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "isbn10")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISBN10 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISBN10 failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "isbn10" {
					t.Fatalf("Index: %d ISBN10 failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestISSNValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"foo", false},
		{"20519990", false},
		{"2051-9991", false},
		{"2051-999X", false},
		{"1050-124X", true},
		{"0317-8471", true},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.param, "issn")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISSN failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d ISSN failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "issn" {
					t.Fatalf("Index: %d ISSN failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestExcludesRuneValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"excludesrune=☻"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "a☺b☻c☹d", Tag: "excludesrune=☻", ExpectedNil: false},
		{Value: "abcd", Tag: "excludesrune=☻", ExpectedNil: true},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestExcludesAllValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"excludesall=@!{}[]"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "abcd@!jfk", Tag: "excludesall=@!{}[]", ExpectedNil: false},
		{Value: "abcdefg", Tag: "excludesall=@!{}[]", ExpectedNil: true},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}

	username := "joeybloggs "

	errs := validate.Var(username, "excludesall=@ ")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "excludesall")

	excluded := ","

	errs = validate.Var(excluded, "excludesall=!@#$%^&*()_+.0x2C?")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "excludesall")

	excluded = "="

	errs = validate.Var(excluded, "excludesall=!@#$%^&*()_+.0x2C=?")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "excludesall")
}

func TestExcludesValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"excludes=@"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "abcd@!jfk", Tag: "excludes=@", ExpectedNil: false},
		{Value: "abcdq!jfk", Tag: "excludes=@", ExpectedNil: true},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestContainsRuneValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"containsrune=☻"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "a☺b☻c☹d", Tag: "containsrune=☻", ExpectedNil: true},
		{Value: "abcd", Tag: "containsrune=☻", ExpectedNil: false},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestContainsAnyValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"containsany=@!{}[]"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "abcd@!jfk", Tag: "containsany=@!{}[]", ExpectedNil: true},
		{Value: "abcdefg", Tag: "containsany=@!{}[]", ExpectedNil: false},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestContainsValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"contains=@"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "abcd@!jfk", Tag: "contains=@", ExpectedNil: true},
		{Value: "abcdq!jfk", Tag: "contains=@", ExpectedNil: false},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestIsNeFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	var j uint64
	var k float64
	s := "abcd"
	i := 1
	j = 1
	k = 1.543
	b := true
	arr := []string{"test"}
	now := time.Now().UTC()

	var j2 uint64
	var k2 float64
	s2 := "abcdef"
	i2 := 3
	j2 = 2
	k2 = 1.5434456
	b2 := false
	arr2 := []string{"test", "test2"}
	arr3 := []string{"test"}
	now2 := now

	errs = validate.VarWithValue(s, s2, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(i2, i, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(j2, j, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(k2, k, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(b2, b, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(arr2, arr, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(now2, now, "nefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "nefield")

	errs = validate.VarWithValue(arr3, arr, "nefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "nefield")

	type Test struct {
		Start *time.Time `validate:"nefield=End"`
		End   *time.Time
	}

	sv := &Test{
		Start: &now,
		End:   &now,
	}

	errs = validate.Struct(sv)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Start", "Test.Start", "Start", "Start", "nefield")

	now3 := time.Now().Add(time.Hour).UTC()

	sv = &Test{
		Start: &now,
		End:   &now3,
	}

	errs = validate.Struct(sv)
	Equal(t, errs, nil)

	errs = validate.VarWithValue(nil, 1, "nefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "nefield")

	errs = validate.VarWithValue(sv, now, "nefield")
	Equal(t, errs, nil)

	type Test2 struct {
		Start *time.Time `validate:"nefield=NonExistantField"`
		End   *time.Time
	}

	sv2 := &Test2{
		Start: &now,
		End:   &now,
	}

	errs = validate.Struct(sv2)
	Equal(t, errs, nil)

	type Other struct {
		Value string
	}

	type Test3 struct {
		Value Other
		Time  time.Time `validate:"nefield=Value"`
	}

	tst := Test3{
		Value: Other{Value: "StringVal"},
		Time:  now,
	}

	errs = validate.Struct(tst)
	Equal(t, errs, nil)

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "nefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "nefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "nefield")

	errs = validate.VarWithValue(time.Duration(0), time.Duration(0), "omitempty,nefield")
	Equal(t, errs, nil)

	// -- Validations for a struct with time.Duration type fields.

	type TimeDurationTest struct {
		First  time.Duration `validate:"nefield=Second"`
		Second time.Duration
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "nefield")

	type TimeDurationOmitemptyTest struct {
		First  time.Duration `validate:"omitempty,nefield=Second"`
		Second time.Duration
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsNeValidation(t *testing.T) {
	var errs error
	validate := New()

	var j uint64
	var k float64
	s := "abcdef"
	i := 3
	j = 2
	k = 1.5434
	arr := []string{"test"}
	now := time.Now().UTC()

	errs = validate.Var(s, "ne=abcd")
	Equal(t, errs, nil)

	errs = validate.Var(i, "ne=1")
	Equal(t, errs, nil)

	errs = validate.Var(j, "ne=1")
	Equal(t, errs, nil)

	errs = validate.Var(k, "ne=1.543")
	Equal(t, errs, nil)

	errs = validate.Var(arr, "ne=2")
	Equal(t, errs, nil)

	errs = validate.Var(arr, "ne=1")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ne")

	PanicMatches(t, func() { _ = validate.Var(now, "ne=now") }, "Bad field type time.Time")

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour-time.Minute, "ne=1h")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+time.Minute, "ne=1h")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour, "ne=1h")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ne")

	errs = validate.Var(time.Duration(0), "omitempty,ne=0")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"ne=1h"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "ne")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,ne=0"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsNeIgnoreCaseValidation(t *testing.T) {
	var errs error
	validate := New()
	s := "abcd"
	now := time.Now()

	errs = validate.Var(s, "ne_ignore_case=efgh")
	Equal(t, errs, nil)

	errs = validate.Var(s, "ne_ignore_case=AbCd")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ne_ignore_case")

	PanicMatches(
		t, func() { _ = validate.Var(now, "ne_ignore_case=abcd") }, "Bad field type time.Time",
	)
}

func TestIsEqFieldValidation(t *testing.T) {
	var errs error
	validate := New()

	var j uint64
	var k float64
	s := "abcd"
	i := 1
	j = 1
	k = 1.543
	b := true
	arr := []string{"test"}
	now := time.Now().UTC()

	var j2 uint64
	var k2 float64
	s2 := "abcd"
	i2 := 1
	j2 = 1
	k2 = 1.543
	b2 := true
	arr2 := []string{"test"}
	arr3 := []string{"test", "test2"}
	now2 := now

	errs = validate.VarWithValue(s, s2, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(i2, i, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(j2, j, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(k2, k, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(b2, b, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(arr2, arr, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(now2, now, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(arr3, arr, "eqfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqfield")

	type Test struct {
		Start *time.Time `validate:"eqfield=End"`
		End   *time.Time
	}

	sv := &Test{
		Start: &now,
		End:   &now,
	}

	errs = validate.Struct(sv)
	Equal(t, errs, nil)

	now3 := time.Now().Add(time.Hour).UTC()

	sv = &Test{
		Start: &now,
		End:   &now3,
	}

	errs = validate.Struct(sv)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Start", "Test.Start", "Start", "Start", "eqfield")

	errs = validate.VarWithValue(nil, 1, "eqfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqfield")

	channel := make(chan string)
	errs = validate.VarWithValue(5, channel, "eqfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqfield")

	errs = validate.VarWithValue(5, now, "eqfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqfield")

	type Test2 struct {
		Start *time.Time `validate:"eqfield=NonExistantField"`
		End   *time.Time
	}

	sv2 := &Test2{
		Start: &now,
		End:   &now,
	}

	errs = validate.Struct(sv2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test2.Start", "Test2.Start", "Start", "Start", "eqfield")

	type Inner struct {
		Name string
	}

	type TStruct struct {
		Inner     *Inner
		CreatedAt *time.Time `validate:"eqfield=Inner"`
	}

	inner := &Inner{
		Name: "NAME",
	}

	test := &TStruct{
		Inner:     inner,
		CreatedAt: &now,
	}

	errs = validate.Struct(test)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TStruct.CreatedAt", "TStruct.CreatedAt", "CreatedAt", "CreatedAt", "eqfield")

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour, "eqfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "eqfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqfield")

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "eqfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eqfield")

	errs = validate.VarWithValue(time.Duration(0), time.Hour, "omitempty,eqfield")
	Equal(t, errs, nil)

	// -- Validations for a struct with time.Duration type fields.

	type TimeDurationTest struct {
		First  time.Duration `validate:"eqfield=Second"`
		Second time.Duration
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "eqfield")

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "eqfield")

	type TimeDurationOmitemptyTest struct {
		First  time.Duration `validate:"omitempty,eqfield=Second"`
		Second time.Duration
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), time.Hour}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsEqFieldValidationWithAliasTime(t *testing.T) {
	var errs error
	validate := New()

	type CustomTime time.Time

	type Test struct {
		Start CustomTime `validate:"eqfield=End"`
		End   *time.Time
	}

	now := time.Now().UTC()

	sv := &Test{
		Start: CustomTime(now),
		End:   &now,
	}

	errs = validate.Struct(sv)
	Equal(t, errs, nil)
}

func TestIsEqValidation(t *testing.T) {
	var errs error
	validate := New()

	var j uint64
	var k float64
	s := "abcd"
	i := 1
	j = 1
	k = 1.543
	arr := []string{"test"}
	now := time.Now().UTC()

	errs = validate.Var(s, "eq=abcd")
	Equal(t, errs, nil)

	errs = validate.Var(i, "eq=1")
	Equal(t, errs, nil)

	errs = validate.Var(j, "eq=1")
	Equal(t, errs, nil)

	errs = validate.Var(k, "eq=1.543")
	Equal(t, errs, nil)

	errs = validate.Var(arr, "eq=1")
	Equal(t, errs, nil)

	errs = validate.Var(arr, "eq=2")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eq")

	PanicMatches(t, func() { _ = validate.Var(now, "eq=now") }, "Bad field type time.Time")

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "eq=1h")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-time.Minute, "eq=1h")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eq")

	errs = validate.Var(time.Hour+time.Minute, "eq=1h")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "eq")

	errs = validate.Var(time.Duration(0), "omitempty,eq=1h")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"eq=1h"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "eq")

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "eq")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,eq=1h"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsEqIgnoreCaseValidation(t *testing.T) {
	var errs error
	validate := New()
	s := "abcd"
	now := time.Now()

	errs = validate.Var(s, "eq_ignore_case=abcd")
	Equal(t, errs, nil)

	errs = validate.Var(s, "eq_ignore_case=AbCd")
	Equal(t, errs, nil)

	PanicMatches(
		t, func() { _ = validate.Var(now, "eq_ignore_case=abcd") }, "Bad field type time.Time",
	)
}

func TestOneOfValidation(t *testing.T) {
	validate := New()

	passSpecs := []struct {
		f interface{}
		t string
	}{
		{f: "red", t: "oneof=red green"},
		{f: "green", t: "oneof=red green"},
		{f: "red green", t: "oneof='red green' blue"},
		{f: "blue", t: "oneof='red green' blue"},
		{f: 5, t: "oneof=5 6"},
		{f: 6, t: "oneof=5 6"},
		{f: int8(6), t: "oneof=5 6"},
		{f: int16(6), t: "oneof=5 6"},
		{f: int32(6), t: "oneof=5 6"},
		{f: int64(6), t: "oneof=5 6"},
		{f: uint(6), t: "oneof=5 6"},
		{f: uint8(6), t: "oneof=5 6"},
		{f: uint16(6), t: "oneof=5 6"},
		{f: uint32(6), t: "oneof=5 6"},
		{f: uint64(6), t: "oneof=5 6"},
	}

	for _, spec := range passSpecs {
		t.Logf("%#v", spec)
		errs := validate.Var(spec.f, spec.t)
		Equal(t, errs, nil)
	}

	failSpecs := []struct {
		f interface{}
		t string
	}{
		{f: "", t: "oneof=red green"},
		{f: "yellow", t: "oneof=red green"},
		{f: "green", t: "oneof='red green' blue"},
		{f: 5, t: "oneof=red green"},
		{f: 6, t: "oneof=red green"},
		{f: 6, t: "oneof=7"},
		{f: uint(6), t: "oneof=7"},
		{f: int8(5), t: "oneof=red green"},
		{f: int16(5), t: "oneof=red green"},
		{f: int32(5), t: "oneof=red green"},
		{f: int64(5), t: "oneof=red green"},
		{f: uint(5), t: "oneof=red green"},
		{f: uint8(5), t: "oneof=red green"},
		{f: uint16(5), t: "oneof=red green"},
		{f: uint32(5), t: "oneof=red green"},
		{f: uint64(5), t: "oneof=red green"},
	}

	for _, spec := range failSpecs {
		t.Logf("%#v", spec)
		errs := validate.Var(spec.f, spec.t)
		AssertError(t, errs, "", "", "", "", "oneof")
	}

	PanicMatches(t, func() {
		_ = validate.Var(3.14, "oneof=red green")
	}, "Bad field type float64")
}

func TestOneOfCIValidation(t *testing.T) {
	validate := New()

	passSpecs := []struct {
		f interface{}
		t string
	}{
		{f: "red", t: "oneofci=RED GREEN"},
		{f: "RED", t: "oneofci=red green"},
		{f: "red", t: "oneofci=red green"},
		{f: "RED", t: "oneofci=RED GREEN"},
		{f: "green", t: "oneofci=red green"},
		{f: "red green", t: "oneofci='red green' blue"},
		{f: "blue", t: "oneofci='red green' blue"},
		{f: "GREEN", t: "oneofci=Red Green"},
		{f: "ReD", t: "oneofci=RED GREEN"},
		{f: "gReEn", t: "oneofci=rEd GrEeN"},
		{f: "RED GREEN", t: "oneofci='red green' blue"},
		{f: "red Green", t: "oneofci='RED GREEN' Blue"},
		{f: "Red green", t: "oneofci='Red Green' BLUE"},
		{f: "rEd GrEeN", t: "oneofci='ReD gReEn' BlUe"},
		{f: "BLUE", t: "oneofci='Red Green' BLUE"},
		{f: "BlUe", t: "oneofci='RED GREEN' Blue"},
		{f: "bLuE", t: "oneofci='red green' BLUE"},
	}

	for _, spec := range passSpecs {
		t.Logf("%#v", spec)
		errs := validate.Var(spec.f, spec.t)
		Equal(t, errs, nil)
	}

	failSpecs := []struct {
		f interface{}
		t string
	}{
		{f: "", t: "oneofci=red green"},
		{f: "yellow", t: "oneofci=red green"},
		{f: "green", t: "oneofci='red green' blue"},
	}

	for _, spec := range failSpecs {
		t.Logf("%#v", spec)
		errs := validate.Var(spec.f, spec.t)
		AssertError(t, errs, "", "", "", "", "oneofci")
	}

	panicSpecs := []struct {
		f interface{}
		t string
	}{
		{f: 3.14, t: "oneofci=red green"},
		{f: 5, t: "oneofci=red green"},
		{f: uint(6), t: "oneofci=7"},
		{f: int8(5), t: "oneofci=red green"},
		{f: int16(5), t: "oneofci=red green"},
		{f: int32(5), t: "oneofci=red green"},
		{f: int64(5), t: "oneofci=red green"},
		{f: uint(5), t: "oneofci=red green"},
		{f: uint8(5), t: "oneofci=red green"},
		{f: uint16(5), t: "oneofci=red green"},
		{f: uint32(5), t: "oneofci=red green"},
		{f: uint64(5), t: "oneofci=red green"},
	}

	panicCount := 0
	for _, spec := range panicSpecs {
		t.Logf("%#v", spec)
		PanicMatches(t, func() {
			_ = validate.Var(spec.f, spec.t)
		}, fmt.Sprintf("Bad field type %T", spec.f))
		panicCount++
	}
	Equal(t, panicCount, len(panicSpecs))
}

func TestBase32Validation(t *testing.T) {
	validate := New()

	s := "ABCD2345"
	errs := validate.Var(s, "base32")
	Equal(t, errs, nil)

	s = "AB======"
	errs = validate.Var(s, "base32")
	Equal(t, errs, nil)

	s = "ABCD2==="
	errs = validate.Var(s, "base32")
	Equal(t, errs, nil)

	s = "ABCD===="
	errs = validate.Var(s, "base32")
	Equal(t, errs, nil)

	s = "ABCD234="
	errs = validate.Var(s, "base32")
	Equal(t, errs, nil)

	s = ""
	errs = validate.Var(s, "base32")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "base32")

	s = "ABCabc1890== foo bar"
	errs = validate.Var(s, "base32")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "base32")
}

func TestBase64Validation(t *testing.T) {
	validate := New()

	s := "dW5pY29ybg=="

	errs := validate.Var(s, "base64")
	Equal(t, errs, nil)

	s = "dGhpIGlzIGEgdGVzdCBiYXNlNjQ="
	errs = validate.Var(s, "base64")
	Equal(t, errs, nil)

	s = ""
	errs = validate.Var(s, "base64")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "base64")

	s = "dW5pY29ybg== foo bar"
	errs = validate.Var(s, "base64")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "base64")
}

func TestBase64URLValidation(t *testing.T) {
	validate := New()

	testCases := []struct {
		decoded, encoded string
		success          bool
	}{
		// empty string, although a valid base64 string, should fail
		{"", "", false},
		// invalid length
		{"", "a", false},
		// base64 with padding
		{"f", "Zg==", true},
		{"fo", "Zm8=", true},
		// base64 without padding
		{"foo", "Zm9v", true},
		{"", "Zg", false},
		{"", "Zm8", false},
		// base64 URL safe encoding with invalid, special characters '+' and '/'
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l+", false},
		{"\x14\xfb\x9c\x03\xf9\x73", "FPucA/lz", false},
		// base64 URL safe encoding with valid, special characters '-' and '_'
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l-", true},
		{"\x14\xfb\x9c\x03\xf9\x73", "FPucA_lz", true},
		// non base64 characters
		{"", "@mc=", false},
		{"", "Zm 9", false},
	}
	for _, tc := range testCases {
		err := validate.Var(tc.encoded, "base64url")
		if tc.success {
			Equal(t, err, nil)
			// make sure encoded value is decoded back to the expected value
			d, innerErr := base64.URLEncoding.DecodeString(tc.encoded)
			Equal(t, innerErr, nil)
			Equal(t, tc.decoded, string(d))
		} else {
			NotEqual(t, err, nil)
			if len(tc.encoded) > 0 {
				// make sure that indeed the encoded value was faulty
				_, err := base64.URLEncoding.DecodeString(tc.encoded)
				NotEqual(t, err, nil)
			}
		}
	}
}

func TestBase64RawURLValidation(t *testing.T) {
	validate := New()

	testCases := []struct {
		decoded, encoded string
		success          bool
	}{
		// empty string, although a valid base64 string, should fail
		{"", "", false},
		// invalid length
		{"", "a", false},
		// base64 with padding should fail
		{"f", "Zg==", false},
		{"fo", "Zm8=", false},
		// base64 without padding
		{"foo", "Zm9v", true},
		{"hello", "aGVsbG8", true},
		{"", "aGVsb", false},
		// // base64 URL safe encoding with invalid, special characters '+' and '/'
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l+", false},
		{"\x14\xfb\x9c\x03\xf9\x73", "FPucA/lz", false},
		// // base64 URL safe encoding with valid, special characters '-' and '_'
		{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l-", true},
		{"\x14\xfb\x9c\x03\xf9\x73", "FPucA_lz", true},
		// non base64 characters
		{"", "@mc=", false},
		{"", "Zm 9", false},
	}
	for _, tc := range testCases {
		err := validate.Var(tc.encoded, "base64rawurl")
		if tc.success {
			Equal(t, err, nil)
			// make sure encoded value is decoded back to the expected value
			d, innerErr := base64.RawURLEncoding.DecodeString(tc.encoded)
			Equal(t, innerErr, nil)
			Equal(t, tc.decoded, string(d))
		} else {
			NotEqual(t, err, nil)
			if len(tc.encoded) > 0 {
				// make sure that indeed the encoded value was faulty
				_, err := base64.RawURLEncoding.DecodeString(tc.encoded)
				NotEqual(t, err, nil)
			}
		}
	}
}

func TestFileValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		title    string
		param    string
		expected bool
	}{
		{"empty path", "", false},
		{"regular file", filepath.Join("testdata", "a.go"), true},
		{"missing file", filepath.Join("testdata", "no.go"), false},
		{"directory, not a file", "testdata", false},
	}

	for _, test := range tests {
		errs := validate.Var(test.param, "file")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(6, "file")
	}, "Bad field type int")
}

func TestImageValidation(t *testing.T) {
	validate := New()

	paths := map[string]string{
		"empty":     "",
		"directory": "testdata",
		"missing":   filepath.Join("testdata", "none.png"),
		"png":       filepath.Join("testdata", "image.png"),
		"jpeg":      filepath.Join("testdata", "image.jpg"),
		"mp3":       filepath.Join("testdata", "music.mp3"),
	}

	tests := []struct {
		title        string
		param        string
		expected     bool
		createImage  func()
		destroyImage func()
	}{
		{
			"empty path",
			paths["empty"], false,
			func() {},
			func() {},
		},
		{
			"directory, not a file",
			paths["directory"],
			false,
			func() {},
			func() {},
		},
		{
			"missing file",
			paths["missing"],
			false,
			func() {},
			func() {},
		},
		{
			"valid png",
			paths["png"],
			true,
			func() {
				img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{10, 10}})
				f, _ := os.Create(paths["png"])
				err := png.Encode(f, img)
				if err != nil {
					panic(fmt.Sprintf("Could not encode file in PNG. Error: %s", err))
				}
			},
			func() {
				os.Remove(paths["png"])
			},
		},
		{
			"valid jpeg",
			paths["jpeg"],
			true,
			func() {
				var opt jpeg.Options
				img := image.NewGray(image.Rect(0, 0, 10, 10))
				f, _ := os.Create(paths["jpeg"])
				err := jpeg.Encode(f, img, &opt)
				if err != nil {
					panic(fmt.Sprintf("Could not encode file in JPEG. Error: %s", err))
				}
			},
			func() {
				os.Remove(paths["jpeg"])
			},
		},
		{
			"valid mp3",
			paths["mp3"],
			false,
			func() {},
			func() {},
		},
	}

	for _, test := range tests {
		test.createImage()
		errs := validate.Var(test.param, "image")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		}
		test.destroyImage()
	}

	PanicMatches(t, func() {
		_ = validate.Var(6, "file")
	}, "Bad field type int")
}

func TestFilePathValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		title    string
		param    string
		expected bool
	}{
		{"empty filepath", "", false},
		{"valid filepath", filepath.Join("testdata", "a.go"), true},
		{"invalid filepath", filepath.Join("testdata", "no\000.go"), false},
		{"directory, not a filepath", "testdata" + string(os.PathSeparator), false},
		{"directory", "testdata", false},
	}

	for _, test := range tests {
		errs := validate.Var(test.param, "filepath")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		}

	}

	PanicMatches(t, func() {
		_ = validate.Var(6, "filepath")
	}, "Bad field type int")
}

func TestEthereumAddressValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		param    string
		expected bool
	}{
		// All caps.
		{"0x52908400098527886E0F7030069857D2E4169EE7", true},
		{"0x8617E340B3D01FA5F11F306F4090FD50E238070D", true},

		// All lower.
		{"0xde709f2102306220921060314715629080e2fb77", true},
		{"0x27b1fdb04752bbc536007a920d24acb045561c26", true},
		{"0x123f681646d4a755815f9cb19e1acc8565a0c2ac", true},

		// Mixed case: runs checksum validation.
		{"0x02F9AE5f22EA3fA88F05780B30385bECFacbf130", true},
		{"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed", true},
		{"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359", true},
		{"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB", true},
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", true},
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDB", true}, // Invalid checksum, but valid address.

		// Other.
		{"", false},
		{"D1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", false},    // Missing "0x" prefix.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDbc", false}, // More than 40 hex digits.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aD", false},   // Less than 40 hex digits.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDw", false},  // Invalid hex digit "w".
	}

	for i, test := range tests {

		errs := validate.Var(test.param, "eth_addr")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d eth_addr failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d eth_addr failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "eth_addr" {
					t.Fatalf("Index: %d Latitude failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestEthereumAddressChecksumValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		param    string
		expected bool
	}{
		// All caps.
		{"0x52908400098527886E0F7030069857D2E4169EE7", true},
		{"0x8617E340B3D01FA5F11F306F4090FD50E238070D", true},

		// All lower.
		{"0x27b1fdb04752bbc536007a920d24acb045561c26", true},
		{"0x123f681646d4a755815f9cb19e1acc8565a0c2ac", false},

		// Mixed case: runs checksum validation.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", true},
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDB", false}, // Invalid checksum.
		{"0x000000000000000000000000000000000000dead", false}, // Invalid checksum.
		{"0x000000000000000000000000000000000000dEaD", true},  // Valid checksum.

		// Other.
		{"", false},
		{"D1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb", false},    // Missing "0x" prefix.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDbc", false}, // More than 40 hex digits.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aD", false},   // Less than 40 hex digits.
		{"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDw", false},  // Invalid hex digit "w".
	}

	for i, test := range tests {

		errs := validate.Var(test.param, "eth_addr_checksum")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d eth_addr_checksum failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d eth_addr_checksum failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "eth_addr_checksum" {
					t.Fatalf("Index: %d Latitude failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestBitcoinAddressValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"x", false},
		{"0x02F9AE5f22EA3fA88F05780B30385bEC", false},
		{"1A1zP1ePQGefi2DMPTifTL5SLmv7DivfNa", false},
		{"1P9RQEr2XeE3PEb44ZE35sfZRRW1JH8Uqx", false},
		{"3P14159I73E4gFr7JterCCQh9QjiTjiZrG", false},
		{"3P141597f3E4gFr7JterCCQh9QjiTjiZrG", false},
		{"37qgekLpCCHrQuSjvX3fs496FWTGsHFHizjJAs6NPcR47aefnnCWECAhHV6E3g4YN7u7Yuwod5Y", false},
		{"dzb7VV1Ui55BARxv7ATxAtCUeJsANKovDGWFVgpTbhq9gvPqP3yv", false},
		{"MuNu7ZAEDFiHthiunm7dPjwKqrVNCM3mAz6rP9zFveQu14YA8CxExSJTHcVP9DErn6u84E6Ej7S", false},
		{"rPpQpYknyNQ5AEHuY6H8ijJJrYc2nDKKk9jjmKEXsWzyAQcFGpDLU2Zvsmoi8JLR7hAwoy3RQWf", false},
		{"4Uc3FmN6NQ6zLBK5QQBXRBUREaaHwCZYsGCueHauuDmJpZKn6jkEskMB2Zi2CNgtb5r6epWEFfUJq", false},
		{"7aQgR5DFQ25vyXmqZAWmnVCjL3PkBcdVkBUpjrjMTcghHx3E8wb", false},
		{"17QpPprjeg69fW1DV8DcYYCKvWjYhXvWkov6MJ1iTTvMFj6weAqW7wybZeH57WTNxXVCRH4veVs", false},
		{"KxuACDviz8Xvpn1xAh9MfopySZNuyajYMZWz16Dv2mHHryznWUp3", false},
		{"7nK3GSmqdXJQtdohvGfJ7KsSmn3TmGqExug49583bDAL91pVSGq5xS9SHoAYL3Wv3ijKTit65th", false},
		{"cTivdBmq7bay3RFGEBBuNfMh2P1pDCgRYN2Wbxmgwr4ki3jNUL2va", false},
		{"gjMV4vjNjyMrna4fsAr8bWxAbwtmMUBXJS3zL4NJt5qjozpbQLmAfK1uA3CquSqsZQMpoD1g2nk", false},
		{"emXm1naBMoVzPjbk7xpeTVMFy4oDEe25UmoyGgKEB1gGWsK8kRGs", false},
		{"7VThQnNRj1o3Zyvc7XHPRrjDf8j2oivPTeDXnRPYWeYGE4pXeRJDZgf28ppti5hsHWXS2GSobdqyo", false},
		{"1G9u6oCVCPh2o8m3t55ACiYvG1y5BHewUkDSdiQarDcYXXhFHYdzMdYfUAhfxn5vNZBwpgUNpso", false},
		{"31QQ7ZMLkScDiB4VyZjuptr7AEc9j1SjstF7pRoLhHTGkW4Q2y9XELobQmhhWxeRvqcukGd1XCq", false},
		{"DHqKSnpxa8ZdQyH8keAhvLTrfkyBMQxqngcQA5N8LQ9KVt25kmGN", false},
		{"2LUHcJPbwLCy9GLH1qXmfmAwvadWw4bp4PCpDfduLqV17s6iDcy1imUwhQJhAoNoN1XNmweiJP4i", false},
		{"7USRzBXAnmck8fX9HmW7RAb4qt92VFX6soCnts9s74wxm4gguVhtG5of8fZGbNPJA83irHVY6bCos", false},
		{"1DGezo7BfVebZxAbNT3XGujdeHyNNBF3vnficYoTSp4PfK2QaML9bHzAMxke3wdKdHYWmsMTJVu", false},
		{"2D12DqDZKwCxxkzs1ZATJWvgJGhQ4cFi3WrizQ5zLAyhN5HxuAJ1yMYaJp8GuYsTLLxTAz6otCfb", false},
		{"8AFJzuTujXjw1Z6M3fWhQ1ujDW7zsV4ePeVjVo7D1egERqSW9nZ", false},
		{"163Q17qLbTCue8YY3AvjpUhotuaodLm2uqMhpYirsKjVqnxJRWTEoywMVY3NbBAHuhAJ2cF9GAZ", false},
		{"2MnmgiRH4eGLyLc9eAqStzk7dFgBjFtUCtu", false},
		{"461QQ2sYWxU7H2PV4oBwJGNch8XVTYYbZxU", false},
		{"2UCtv53VttmQYkVU4VMtXB31REvQg4ABzs41AEKZ8UcB7DAfVzdkV9JDErwGwyj5AUHLkmgZeobs", false},
		{"cSNjAsnhgtiFMi6MtfvgscMB2Cbhn2v1FUYfviJ1CdjfidvmeW6mn", false},
		{"gmsow2Y6EWAFDFE1CE4Hd3Tpu2BvfmBfG1SXsuRARbnt1WjkZnFh1qGTiptWWbjsq2Q6qvpgJVj", false},
		{"nksUKSkzS76v8EsSgozXGMoQFiCoCHzCVajFKAXqzK5on9ZJYVHMD5CKwgmX3S3c7M1U3xabUny", false},
		{"L3favK1UzFGgdzYBF2oBT5tbayCo4vtVBLJhg2iYuMeePxWG8SQc", false},
		{"7VxLxGGtYT6N99GdEfi6xz56xdQ8nP2dG1CavuXx7Rf2PrvNMTBNevjkfgs9JmkcGm6EXpj8ipyPZ ", false},
		{"2mbZwFXF6cxShaCo2czTRB62WTx9LxhTtpP", false},
		{"dB7cwYdcPSgiyAwKWL3JwCVwSk6epU2txw", false},
		{"HPhFUhUAh8ZQQisH8QQWafAxtQYju3SFTX", false},
		{"4ctAH6AkHzq5ioiM1m9T3E2hiYEev5mTsB", false},
		{"31uEbMgunupShBVTewXjtqbBv5MndwfXhb", false},
		{"175tWpb8K1S7NmH4Zx6rewF9WQrcZv245W", false},
		{"Hn1uFi4dNexWrqARpjMqgT6cX1UsNPuV3cHdGg9ExyXw8HTKadbktRDtdeVmY3M1BxJStiL4vjJ", false},
		{"Sq3fDbvutABmnAHHExJDgPLQn44KnNC7UsXuT7KZecpaYDMU9Txs", false},
		{"6TqWyrqdgUEYDQU1aChMuFMMEimHX44qHFzCUgGfqxGgZNMUVWJ", false},
		{"giqJo7oWqFxNKWyrgcBxAVHXnjJ1t6cGoEffce5Y1y7u649Noj5wJ4mmiUAKEVVrYAGg2KPB3Y4", false},
		{"cNzHY5e8vcmM3QVJUcjCyiKMYfeYvyueq5qCMV3kqcySoLyGLYUK", false},
		{"37uTe568EYc9WLoHEd9jXEvUiWbq5LFLscNyqvAzLU5vBArUJA6eydkLmnMwJDjkL5kXc2VK7ig", false},
		{"EsYbG4tWWWY45G31nox838qNdzksbPySWc", false},
		{"nbuzhfwMoNzA3PaFnyLcRxE9bTJPDkjZ6Rf6Y6o2ckXZfzZzXBT", false},
		{"cQN9PoxZeCWK1x56xnz6QYAsvR11XAce3Ehp3gMUdfSQ53Y2mPzx", false},
		{"1Gm3N3rkef6iMbx4voBzaxtXcmmiMTqZPhcuAepRzYUJQW4qRpEnHvMojzof42hjFRf8PE2jPde", false},
		{"2TAq2tuN6x6m233bpT7yqdYQPELdTDJn1eU", false},
		{"ntEtnnGhqPii4joABvBtSEJG6BxjT2tUZqE8PcVYgk3RHpgxgHDCQxNbLJf7ardf1dDk2oCQ7Cf", false},
		{"Ky1YjoZNgQ196HJV3HpdkecfhRBmRZdMJk89Hi5KGfpfPwS2bUbfd", false},
		{"2A1q1YsMZowabbvta7kTy2Fd6qN4r5ZCeG3qLpvZBMzCixMUdkN2Y4dHB1wPsZAeVXUGD83MfRED", false},
		{"1AGNa15ZQXAZUgFiqJ2i7Z2DPU2J6hW62i", true},
		{"1Ax4gZtb7gAit2TivwejZHYtNNLT18PUXJ", true},
		{"1C5bSj1iEGUgSTbziymG7Cn18ENQuT36vv", true},
		{"1Gqk4Tv79P91Cc1STQtU3s1W6277M2CVWu", true},
		{"1JwMWBVLtiqtscbaRHai4pqHokhFCbtoB4", true},
		{"19dcawoKcZdQz365WpXWMhX6QCUpR9SY4r", true},
		{"13p1ijLwsnrcuyqcTvJXkq2ASdXqcnEBLE", true},
		{"1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2", true},
		{"3P14159f73E4gFr7JterCCQh9QjiTjiZrG", true},
		{"3CMNFxN1oHBc4R1EpboAL5yzHGgE611Xou", true},
		{"3QjYXhTkvuj8qPaXHTTWb5wjXhdsLAAWVy", true},
		{"3AnNxabYGoTxYiTEZwFEnerUoeFXK2Zoks", true},
		{"33vt8ViH5jsr115AGkW6cEmEz9MpvJSwDk", true},
		{"3QCzvfL4ZRvmJFiWWBVwxfdaNBT8EtxB5y", true},
		{"37Sp6Rv3y4kVd1nQ1JV5pfqXccHNyZm1x3", true},
		{"3ALJH9Y951VCGcVZYAdpA3KchoP9McEj1G", true},
		{"12KYrjTdVGjFMtaxERSk3gphreJ5US8aUP", true},
		{"12QeMLzSrB8XH8FvEzPMVoRxVAzTr5XM2y", true},
		{"1oNLrsHnBcR6dpaBpwz3LSwutbUNkNSjs", true},
		{"1SQHtwR5oJRKLfiWQ2APsAd9miUc4k2ez", true},
		{"116CGDLddrZhMrTwhCVJXtXQpxygTT1kHd", true},
		{"3NJZLcZEEYBpxYEUGewU4knsQRn1WM5Fkt", true},
	}

	for i, test := range tests {

		errs := validate.Var(test.param, "btc_addr")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d btc_addr failed with Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d btc_addr failed with Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "btc_addr" {
					t.Fatalf("Index: %d Latitude failed with Error: %s", i, errs)
				}
			}
		}
	}
}

func TestBitcoinBech32AddressValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"bc1rw5uspcuh", false},
		{"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5", false},
		{"BC13W508D6QEJXTDG4Y5R3ZARVARY0C5XW7KN40WF2", false},
		{"qw508d6qejxtdg4y5r3zarvary0c5xw7kg3g4ty", false},
		{"bc1rw5uspcuh", false},
		{"bc10w508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7kw5rljs90", false},
		{"BC1QW508d6QEJxTDG4y5R3ZArVARY0C5XW7KV8F3T4", false},
		{"BC1QR508D6QEJXTDG4Y5R3ZARVARYV98GJ9P", false},
		{"bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t5", false},
		{"bc10w508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7kw5rljs90", false},
		{"bc1pw508d6qejxtdg4y5r3zarqfsj6c3", false},
		{"bc1zw508d6qejxtdg4y5r3zarvaryvqyzf3du", false},
		{"bc1gmk9yu", false},
		{"bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3pjxtptv", false},
		{"BC1QW508D6QEJXTDG4Y5R3ZARVARY0C5XW7KV8F3T4", true},
		{"bc1pw508d6qejxtdg4y5r3zarvary0c5xw7kw508d6qejxtdg4y5r3zarvary0c5xw7k7grplx", true},
		{"bc1qrp33g0q5c5txsp9arysrx4k6zdkfs4nce4xj0gdcccefvpysxf3qccfmv3", true},
		{"BC1SW50QA3JX3S", true},
		{"bc1zw508d6qejxtdg4y5r3zarvaryvg6kdaj", true},
	}

	for i, test := range tests {

		errs := validate.Var(test.param, "btc_addr_bech32")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d btc_addr_bech32 failed with Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d btc_addr_bech32 failed with Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "btc_addr_bech32" {
					t.Fatalf("Index: %d Latitude failed with Error: %s", i, errs)
				}
			}
		}
	}
}

func TestNoStructLevelValidation(t *testing.T) {
	type Inner struct {
		Test string `validate:"len=5"`
	}

	type Outer struct {
		InnerStruct    Inner  `validate:"required,nostructlevel"`
		InnerStructPtr *Inner `validate:"required,nostructlevel"`
	}

	outer := &Outer{
		InnerStructPtr: nil,
		InnerStruct:    Inner{},
	}

	// test with struct required failing on
	validate := New(WithRequiredStructEnabled())

	errs := validate.Struct(outer)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Outer.InnerStruct", "Outer.InnerStruct", "InnerStruct", "InnerStruct", "required")
	AssertError(t, errs, "Outer.InnerStructPtr", "Outer.InnerStructPtr", "InnerStructPtr", "InnerStructPtr", "required")

	inner := Inner{
		Test: "1234",
	}

	outer = &Outer{
		InnerStruct:    inner,
		InnerStructPtr: &inner,
	}

	errs = validate.Struct(outer)
	Equal(t, errs, nil)

	// test with struct required failing off

	outer = &Outer{
		InnerStructPtr: nil,
		InnerStruct:    Inner{},
	}
	validate = New()

	errs = validate.Struct(outer)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Outer.InnerStructPtr", "Outer.InnerStructPtr", "InnerStructPtr", "InnerStructPtr", "required")

	inner = Inner{
		Test: "1234",
	}

	outer = &Outer{
		InnerStruct:    inner,
		InnerStructPtr: &inner,
	}

	errs = validate.Struct(outer)
	Equal(t, errs, nil)
}

func TestStructOnlyValidation(t *testing.T) {
	type Inner struct {
		Test string `validate:"len=5"`
	}

	type Outer struct {
		InnerStruct    Inner  `validate:"required,structonly"`
		InnerStructPtr *Inner `validate:"required,structonly"`
	}

	outer := &Outer{
		InnerStruct:    Inner{},
		InnerStructPtr: nil,
	}

	// without required struct on
	validate := New()

	errs := validate.Struct(outer)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Outer.InnerStructPtr", "Outer.InnerStructPtr", "InnerStructPtr", "InnerStructPtr", "required")

	// with required struct on
	validate.requiredStructEnabled = true

	errs = validate.Struct(outer)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Outer.InnerStruct", "Outer.InnerStruct", "InnerStruct", "InnerStruct", "required")
	AssertError(t, errs, "Outer.InnerStructPtr", "Outer.InnerStructPtr", "InnerStructPtr", "InnerStructPtr", "required")

	inner := Inner{
		Test: "1234",
	}

	outer = &Outer{
		InnerStruct:    inner,
		InnerStructPtr: &inner,
	}

	errs = validate.Struct(outer)
	Equal(t, errs, nil)

	// Address houses a users address information
	type Address struct {
		Street string `validate:"required"`
		City   string `validate:"required"`
		Planet string `validate:"required"`
		Phone  string `validate:"required"`
	}

	type User struct {
		FirstName      string     `json:"fname"`
		LastName       string     `json:"lname"`
		Age            uint8      `validate:"gte=0,lte=130"`
		Number         string     `validate:"required,e164"`
		Email          string     `validate:"required,email"`
		FavouriteColor string     `validate:"hexcolor|rgb|rgba"`
		Addresses      []*Address `validate:"required"`   // a person can have a home and cottage...
		Address        Address    `validate:"structonly"` // a person can have a home and cottage...
	}

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
		City:   "Unknown",
	}

	user := &User{
		FirstName:      "",
		LastName:       "",
		Age:            45,
		Number:         "+1123456789",
		Email:          "Badger.Smith@gmail.com",
		FavouriteColor: "#000",
		Addresses:      []*Address{address},
		Address: Address{
			// Street: "Eavesdown Docks",
			Planet: "Persphone",
			Phone:  "none",
			City:   "Unknown",
		},
	}

	errs = validate.Struct(user)
	Equal(t, errs, nil)
}

func TestGtField(t *testing.T) {
	var errs error
	validate := New()

	type TimeTest struct {
		Start *time.Time `validate:"required,gt"`
		End   *time.Time `validate:"required,gt,gtfield=Start"`
	}

	now := time.Now()
	start := now.Add(time.Hour * 24)
	end := start.Add(time.Hour * 24)

	timeTest := &TimeTest{
		Start: &start,
		End:   &end,
	}

	errs = validate.Struct(timeTest)
	Equal(t, errs, nil)

	timeTest = &TimeTest{
		Start: &end,
		End:   &start,
	}

	errs = validate.Struct(timeTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest.End", "TimeTest.End", "End", "End", "gtfield")

	errs = validate.VarWithValue(&end, &start, "gtfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(&start, &end, "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	errs = validate.VarWithValue(&end, &start, "gtfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(&timeTest, &end, "gtfield")
	NotEqual(t, errs, nil)

	errs = validate.VarWithValue("test bigger", "test", "gtfield")
	Equal(t, errs, nil)

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "gtfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	errs = validate.VarWithValue(time.Duration(0), time.Hour, "omitempty,gtfield")
	Equal(t, errs, nil)

	// -- Validations for a struct with time.Duration type fields.

	type TimeDurationTest struct {
		First  time.Duration `validate:"gtfield=Second"`
		Second time.Duration
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "gtfield")

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "gtfield")

	type TimeDurationOmitemptyTest struct {
		First  time.Duration `validate:"omitempty,gtfield=Second"`
		Second time.Duration
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), time.Hour}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)

	// Tests for Ints types.

	type IntTest struct {
		Val1 int `validate:"required"`
		Val2 int `validate:"required,gtfield=Val1"`
	}

	intTest := &IntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(intTest)
	Equal(t, errs, nil)

	intTest = &IntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(intTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "IntTest.Val2", "IntTest.Val2", "Val2", "Val2", "gtfield")

	errs = validate.VarWithValue(int(5), int(1), "gtfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(int(1), int(5), "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	type UIntTest struct {
		Val1 uint `validate:"required"`
		Val2 uint `validate:"required,gtfield=Val1"`
	}

	uIntTest := &UIntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(uIntTest)
	Equal(t, errs, nil)

	uIntTest = &UIntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(uIntTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "UIntTest.Val2", "UIntTest.Val2", "Val2", "Val2", "gtfield")

	errs = validate.VarWithValue(uint(5), uint(1), "gtfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(uint(1), uint(5), "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	type FloatTest struct {
		Val1 float64 `validate:"required"`
		Val2 float64 `validate:"required,gtfield=Val1"`
	}

	floatTest := &FloatTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(floatTest)
	Equal(t, errs, nil)

	floatTest = &FloatTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(floatTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "FloatTest.Val2", "FloatTest.Val2", "Val2", "Val2", "gtfield")

	errs = validate.VarWithValue(float32(5), float32(1), "gtfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(float32(1), float32(5), "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	errs = validate.VarWithValue(nil, 1, "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	errs = validate.VarWithValue(5, "T", "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	errs = validate.VarWithValue(5, start, "gtfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtfield")

	type TimeTest2 struct {
		Start *time.Time `validate:"required"`
		End   *time.Time `validate:"required,gtfield=NonExistantField"`
	}

	timeTest2 := &TimeTest2{
		Start: &start,
		End:   &end,
	}

	errs = validate.Struct(timeTest2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest2.End", "TimeTest2.End", "End", "End", "gtfield")

	type Other struct {
		Value string
	}

	type Test struct {
		Value Other
		Time  time.Time `validate:"gtfield=Value"`
	}

	tst := Test{
		Value: Other{Value: "StringVal"},
		Time:  end,
	}

	errs = validate.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Time", "Test.Time", "Time", "Time", "gtfield")
}

func TestLtField(t *testing.T) {
	var errs error
	validate := New()

	type TimeTest struct {
		Start *time.Time `validate:"required,lt,ltfield=End"`
		End   *time.Time `validate:"required,lt"`
	}

	now := time.Now()
	start := now.Add(time.Hour * 24 * -1 * 2)
	end := start.Add(time.Hour * 24)

	timeTest := &TimeTest{
		Start: &start,
		End:   &end,
	}

	errs = validate.Struct(timeTest)
	Equal(t, errs, nil)

	timeTest = &TimeTest{
		Start: &end,
		End:   &start,
	}

	errs = validate.Struct(timeTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest.Start", "TimeTest.Start", "Start", "Start", "ltfield")

	errs = validate.VarWithValue(&start, &end, "ltfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(&end, &start, "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue(&end, timeTest, "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue("tes", "test", "ltfield")
	Equal(t, errs, nil)

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "ltfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue(time.Duration(0), -time.Minute, "omitempty,ltfield")
	Equal(t, errs, nil)

	// -- Validations for a struct with time.Duration type fields.

	type TimeDurationTest struct {
		First  time.Duration `validate:"ltfield=Second"`
		Second time.Duration
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "ltfield")

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "ltfield")

	type TimeDurationOmitemptyTest struct {
		First  time.Duration `validate:"omitempty,ltfield=Second"`
		Second time.Duration
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), -time.Minute}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)

	// Tests for Ints types.

	type IntTest struct {
		Val1 int `validate:"required"`
		Val2 int `validate:"required,ltfield=Val1"`
	}

	intTest := &IntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(intTest)
	Equal(t, errs, nil)

	intTest = &IntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(intTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "IntTest.Val2", "IntTest.Val2", "Val2", "Val2", "ltfield")

	errs = validate.VarWithValue(int(1), int(5), "ltfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(int(5), int(1), "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	type UIntTest struct {
		Val1 uint `validate:"required"`
		Val2 uint `validate:"required,ltfield=Val1"`
	}

	uIntTest := &UIntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(uIntTest)
	Equal(t, errs, nil)

	uIntTest = &UIntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(uIntTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "UIntTest.Val2", "UIntTest.Val2", "Val2", "Val2", "ltfield")

	errs = validate.VarWithValue(uint(1), uint(5), "ltfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(uint(5), uint(1), "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	type FloatTest struct {
		Val1 float64 `validate:"required"`
		Val2 float64 `validate:"required,ltfield=Val1"`
	}

	floatTest := &FloatTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(floatTest)
	Equal(t, errs, nil)

	floatTest = &FloatTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(floatTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "FloatTest.Val2", "FloatTest.Val2", "Val2", "Val2", "ltfield")

	errs = validate.VarWithValue(float32(1), float32(5), "ltfield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(float32(5), float32(1), "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue(nil, 5, "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue(1, "T", "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	errs = validate.VarWithValue(1, end, "ltfield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltfield")

	type TimeTest2 struct {
		Start *time.Time `validate:"required"`
		End   *time.Time `validate:"required,ltfield=NonExistantField"`
	}

	timeTest2 := &TimeTest2{
		Start: &end,
		End:   &start,
	}

	errs = validate.Struct(timeTest2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest2.End", "TimeTest2.End", "End", "End", "ltfield")
}

func TestFieldContains(t *testing.T) {
	validate := New()

	type StringTest struct {
		Foo string `validate:"fieldcontains=Bar"`
		Bar string
	}

	stringTest := &StringTest{
		Foo: "foobar",
		Bar: "bar",
	}

	errs := validate.Struct(stringTest)
	Equal(t, errs, nil)

	stringTest = &StringTest{
		Foo: "foo",
		Bar: "bar",
	}

	errs = validate.Struct(stringTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "StringTest.Foo", "StringTest.Foo", "Foo", "Foo", "fieldcontains")

	errs = validate.VarWithValue("foo", "bar", "fieldcontains")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "fieldcontains")

	errs = validate.VarWithValue("bar", "foobarfoo", "fieldcontains")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "fieldcontains")

	errs = validate.VarWithValue("foobarfoo", "bar", "fieldcontains")
	Equal(t, errs, nil)

	type StringTestMissingField struct {
		Foo string `validate:"fieldcontains=Bar"`
	}

	stringTestMissingField := &StringTestMissingField{
		Foo: "foo",
	}

	errs = validate.Struct(stringTestMissingField)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "StringTestMissingField.Foo", "StringTestMissingField.Foo", "Foo", "Foo", "fieldcontains")
}

func TestFieldExcludes(t *testing.T) {
	validate := New()

	type StringTest struct {
		Foo string `validate:"fieldexcludes=Bar"`
		Bar string
	}

	stringTest := &StringTest{
		Foo: "foobar",
		Bar: "bar",
	}

	errs := validate.Struct(stringTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "StringTest.Foo", "StringTest.Foo", "Foo", "Foo", "fieldexcludes")

	stringTest = &StringTest{
		Foo: "foo",
		Bar: "bar",
	}

	errs = validate.Struct(stringTest)
	Equal(t, errs, nil)

	errs = validate.VarWithValue("foo", "bar", "fieldexcludes")
	Equal(t, errs, nil)

	errs = validate.VarWithValue("bar", "foobarfoo", "fieldexcludes")
	Equal(t, errs, nil)

	errs = validate.VarWithValue("foobarfoo", "bar", "fieldexcludes")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "fieldexcludes")

	type StringTestMissingField struct {
		Foo string `validate:"fieldexcludes=Bar"`
	}

	stringTestMissingField := &StringTestMissingField{
		Foo: "foo",
	}

	errs = validate.Struct(stringTestMissingField)
	Equal(t, errs, nil)
}

func TestContainsAndExcludes(t *testing.T) {
	validate := New()

	type ImpossibleStringTest struct {
		Foo string `validate:"fieldcontains=Bar"`
		Bar string `validate:"fieldexcludes=Foo"`
	}

	impossibleStringTest := &ImpossibleStringTest{
		Foo: "foo",
		Bar: "bar",
	}

	errs := validate.Struct(impossibleStringTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "ImpossibleStringTest.Foo", "ImpossibleStringTest.Foo", "Foo", "Foo", "fieldcontains")

	impossibleStringTest = &ImpossibleStringTest{
		Foo: "bar",
		Bar: "foo",
	}

	errs = validate.Struct(impossibleStringTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "ImpossibleStringTest.Foo", "ImpossibleStringTest.Foo", "Foo", "Foo", "fieldcontains")
}

func TestLteField(t *testing.T) {
	var errs error
	validate := New()

	type TimeTest struct {
		Start *time.Time `validate:"required,lte,ltefield=End"`
		End   *time.Time `validate:"required,lte"`
	}

	now := time.Now()
	start := now.Add(time.Hour * 24 * -1 * 2)
	end := start.Add(time.Hour * 24)

	timeTest := &TimeTest{
		Start: &start,
		End:   &end,
	}

	errs = validate.Struct(timeTest)
	Equal(t, errs, nil)

	timeTest = &TimeTest{
		Start: &end,
		End:   &start,
	}

	errs = validate.Struct(timeTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest.Start", "TimeTest.Start", "Start", "Start", "ltefield")

	errs = validate.VarWithValue(&start, &end, "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(&end, &start, "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	errs = validate.VarWithValue(&end, timeTest, "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	errs = validate.VarWithValue("tes", "test", "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue("test", "test", "ltefield")
	Equal(t, errs, nil)

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	errs = validate.VarWithValue(time.Duration(0), -time.Minute, "omitempty,ltefield")
	Equal(t, errs, nil)

	// -- Validations for a struct with time.Duration type fields.

	type TimeDurationTest struct {
		First  time.Duration `validate:"ltefield=Second"`
		Second time.Duration
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "ltefield")

	type TimeDurationOmitemptyTest struct {
		First  time.Duration `validate:"omitempty,ltefield=Second"`
		Second time.Duration
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), -time.Minute}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)

	// Tests for Ints types.

	type IntTest struct {
		Val1 int `validate:"required"`
		Val2 int `validate:"required,ltefield=Val1"`
	}

	intTest := &IntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(intTest)
	Equal(t, errs, nil)

	intTest = &IntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(intTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "IntTest.Val2", "IntTest.Val2", "Val2", "Val2", "ltefield")

	errs = validate.VarWithValue(int(1), int(5), "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(int(5), int(1), "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	type UIntTest struct {
		Val1 uint `validate:"required"`
		Val2 uint `validate:"required,ltefield=Val1"`
	}

	uIntTest := &UIntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(uIntTest)
	Equal(t, errs, nil)

	uIntTest = &UIntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(uIntTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "UIntTest.Val2", "UIntTest.Val2", "Val2", "Val2", "ltefield")

	errs = validate.VarWithValue(uint(1), uint(5), "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(uint(5), uint(1), "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	type FloatTest struct {
		Val1 float64 `validate:"required"`
		Val2 float64 `validate:"required,ltefield=Val1"`
	}

	floatTest := &FloatTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(floatTest)
	Equal(t, errs, nil)

	floatTest = &FloatTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(floatTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "FloatTest.Val2", "FloatTest.Val2", "Val2", "Val2", "ltefield")

	errs = validate.VarWithValue(float32(1), float32(5), "ltefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(float32(5), float32(1), "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	errs = validate.VarWithValue(nil, 5, "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	errs = validate.VarWithValue(1, "T", "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	errs = validate.VarWithValue(1, end, "ltefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "ltefield")

	type TimeTest2 struct {
		Start *time.Time `validate:"required"`
		End   *time.Time `validate:"required,ltefield=NonExistantField"`
	}

	timeTest2 := &TimeTest2{
		Start: &end,
		End:   &start,
	}

	errs = validate.Struct(timeTest2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest2.End", "TimeTest2.End", "End", "End", "ltefield")
}

func TestGteField(t *testing.T) {
	var errs error
	validate := New()

	type TimeTest struct {
		Start *time.Time `validate:"required,gte"`
		End   *time.Time `validate:"required,gte,gtefield=Start"`
	}

	now := time.Now()
	start := now.Add(time.Hour * 24)
	end := start.Add(time.Hour * 24)

	timeTest := &TimeTest{
		Start: &start,
		End:   &end,
	}

	errs = validate.Struct(timeTest)
	Equal(t, errs, nil)

	timeTest = &TimeTest{
		Start: &end,
		End:   &start,
	}

	errs = validate.Struct(timeTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest.End", "TimeTest.End", "End", "End", "gtefield")

	errs = validate.VarWithValue(&end, &start, "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(&start, &end, "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	errs = validate.VarWithValue(&start, timeTest, "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	errs = validate.VarWithValue("test", "test", "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue("test bigger", "test", "gtefield")
	Equal(t, errs, nil)

	// Tests for time.Duration type.

	// -- Validations for variables of time.Duration type.

	errs = validate.VarWithValue(time.Hour, time.Hour-time.Minute, "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour, "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(time.Hour, time.Hour+time.Minute, "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	errs = validate.VarWithValue(time.Duration(0), time.Hour, "omitempty,gtefield")
	Equal(t, errs, nil)

	// -- Validations for a struct with time.Duration type fields.

	type TimeDurationTest struct {
		First  time.Duration `validate:"gtefield=Second"`
		Second time.Duration
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour, time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.First", "TimeDurationTest.First", "First", "First", "gtefield")

	type TimeDurationOmitemptyTest struct {
		First  time.Duration `validate:"omitempty,gtefield=Second"`
		Second time.Duration
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0), time.Hour}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)

	// Tests for Ints types.

	type IntTest struct {
		Val1 int `validate:"required"`
		Val2 int `validate:"required,gtefield=Val1"`
	}

	intTest := &IntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(intTest)
	Equal(t, errs, nil)

	intTest = &IntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(intTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "IntTest.Val2", "IntTest.Val2", "Val2", "Val2", "gtefield")

	errs = validate.VarWithValue(int(5), int(1), "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(int(1), int(5), "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	type UIntTest struct {
		Val1 uint `validate:"required"`
		Val2 uint `validate:"required,gtefield=Val1"`
	}

	uIntTest := &UIntTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(uIntTest)
	Equal(t, errs, nil)

	uIntTest = &UIntTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(uIntTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "UIntTest.Val2", "UIntTest.Val2", "Val2", "Val2", "gtefield")

	errs = validate.VarWithValue(uint(5), uint(1), "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(uint(1), uint(5), "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	type FloatTest struct {
		Val1 float64 `validate:"required"`
		Val2 float64 `validate:"required,gtefield=Val1"`
	}

	floatTest := &FloatTest{
		Val1: 1,
		Val2: 5,
	}

	errs = validate.Struct(floatTest)
	Equal(t, errs, nil)

	floatTest = &FloatTest{
		Val1: 5,
		Val2: 1,
	}

	errs = validate.Struct(floatTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "FloatTest.Val2", "FloatTest.Val2", "Val2", "Val2", "gtefield")

	errs = validate.VarWithValue(float32(5), float32(1), "gtefield")
	Equal(t, errs, nil)

	errs = validate.VarWithValue(float32(1), float32(5), "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	errs = validate.VarWithValue(nil, 1, "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	errs = validate.VarWithValue(5, "T", "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	errs = validate.VarWithValue(5, start, "gtefield")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gtefield")

	type TimeTest2 struct {
		Start *time.Time `validate:"required"`
		End   *time.Time `validate:"required,gtefield=NonExistantField"`
	}

	timeTest2 := &TimeTest2{
		Start: &start,
		End:   &end,
	}

	errs = validate.Struct(timeTest2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeTest2.End", "TimeTest2.End", "End", "End", "gtefield")
}

func TestValidateByTagAndValue(t *testing.T) {
	validate := New()

	val := "test"
	field := "test"
	errs := validate.VarWithValue(val, field, "required")
	Equal(t, errs, nil)

	fn := func(fl FieldLevel) bool {
		return fl.Parent().String() == fl.Field().String()
	}

	errs = validate.RegisterValidation("isequaltestfunc", fn)
	Equal(t, errs, nil)

	errs = validate.VarWithValue(val, field, "isequaltestfunc")
	Equal(t, errs, nil)

	val = "unequal"

	errs = validate.VarWithValue(val, field, "isequaltestfunc")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "isequaltestfunc")
}

func TestAddFunctions(t *testing.T) {
	fn := func(fl FieldLevel) bool {
		return true
	}

	fnCtx := func(ctx context.Context, fl FieldLevel) bool {
		return true
	}

	validate := New()

	errs := validate.RegisterValidation("new", fn)
	Equal(t, errs, nil)

	errs = validate.RegisterValidation("", fn)
	NotEqual(t, errs, nil)

	errs = validate.RegisterValidation("new", nil)
	NotEqual(t, errs, nil)

	errs = validate.RegisterValidation("new", fn)
	Equal(t, errs, nil)

	errs = validate.RegisterValidationCtx("new", fnCtx)
	Equal(t, errs, nil)

	PanicMatches(t, func() { _ = validate.RegisterValidation("dive", fn) }, "Tag 'dive' either contains restricted characters or is the same as a restricted tag needed for normal operation")
}

func TestChangeTag(t *testing.T) {
	validate := New()
	validate.SetTagName("val")

	type Test struct {
		Name string `val:"len=4"`
	}
	s := &Test{
		Name: "TEST",
	}

	errs := validate.Struct(s)
	Equal(t, errs, nil)

	s.Name = ""

	errs = validate.Struct(s)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Name", "Test.Name", "Name", "Name", "len")
}

func TestUnexposedStruct(t *testing.T) {
	validate := New()

	type Test struct {
		Name      string
		unexposed struct {
			A string `validate:"required"`
		}
	}

	s := &Test{
		Name: "TEST",
	}
	Equal(t, s.unexposed.A, "")

	errs := validate.Struct(s)
	Equal(t, errs, nil)
}

func TestBadParams(t *testing.T) {
	validate := New()
	i := 1
	errs := validate.Var(i, "-")
	Equal(t, errs, nil)

	PanicMatches(t, func() { _ = validate.Var(i, "len=a") }, "strconv.ParseInt: parsing \"a\": invalid syntax")
	PanicMatches(t, func() { _ = validate.Var(i, "len=a") }, "strconv.ParseInt: parsing \"a\": invalid syntax")

	var ui uint = 1
	PanicMatches(t, func() { _ = validate.Var(ui, "len=a") }, "strconv.ParseUint: parsing \"a\": invalid syntax")

	f := 1.23
	PanicMatches(t, func() { _ = validate.Var(f, "len=a") }, "strconv.ParseFloat: parsing \"a\": invalid syntax")
}

func TestLength(t *testing.T) {
	validate := New()
	i := true
	PanicMatches(t, func() { _ = validate.Var(i, "len") }, "Bad field type bool")
}

func TestIsGt(t *testing.T) {
	var errs error
	validate := New()

	myMap := map[string]string{}
	errs = validate.Var(myMap, "gt=0")
	NotEqual(t, errs, nil)

	f := 1.23
	errs = validate.Var(f, "gt=5")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gt")

	var ui uint = 5
	errs = validate.Var(ui, "gt=10")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gt")

	i := true
	PanicMatches(t, func() { _ = validate.Var(i, "gt") }, "Bad field type bool")

	tm := time.Now().UTC()
	tm = tm.Add(time.Hour * 24)

	errs = validate.Var(tm, "gt")
	Equal(t, errs, nil)

	t2 := time.Now().UTC().Add(-time.Hour)

	errs = validate.Var(t2, "gt")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gt")

	type Test struct {
		Now *time.Time `validate:"gt"`
	}
	s := &Test{
		Now: &tm,
	}

	errs = validate.Struct(s)
	Equal(t, errs, nil)

	s = &Test{
		Now: &t2,
	}

	errs = validate.Struct(s)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Now", "Test.Now", "Now", "Now", "gt")

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "gt=59m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-time.Minute, "gt=59m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gt")

	errs = validate.Var(time.Hour-2*time.Minute, "gt=59m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gt")

	errs = validate.Var(time.Duration(0), "omitempty,gt=59m")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"gt=59m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "gt")

	timeDurationTest = &TimeDurationTest{time.Hour - 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "gt")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,gt=59m"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsGte(t *testing.T) {
	var errs error
	validate := New()

	i := true
	PanicMatches(t, func() { _ = validate.Var(i, "gte") }, "Bad field type bool")

	t1 := time.Now().UTC()
	t1 = t1.Add(time.Hour * 24)

	errs = validate.Var(t1, "gte")
	Equal(t, errs, nil)

	t2 := time.Now().UTC().Add(-time.Hour)

	errs = validate.Var(t2, "gte")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gte")

	type Test struct {
		Now *time.Time `validate:"gte"`
	}
	s := &Test{
		Now: &t1,
	}

	errs = validate.Struct(s)
	Equal(t, errs, nil)

	s = &Test{
		Now: &t2,
	}

	errs = validate.Struct(s)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Now", "Test.Now", "Now", "Now", "gte")

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "gte=59m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-time.Minute, "gte=59m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-2*time.Minute, "gte=59m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "gte")

	errs = validate.Var(time.Duration(0), "omitempty,gte=59m")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"gte=59m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "gte")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,gte=59m"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestMinValidation(t *testing.T) {
	var errs error
	validate := New()

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "min=59m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-time.Minute, "min=59m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-2*time.Minute, "min=59m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "min")

	errs = validate.Var(time.Duration(0), "omitempty,min=59m")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"min=59m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "min")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,min=59m"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestMaxValidation(t *testing.T) {
	var errs error
	validate := New()

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "max=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+time.Minute, "max=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+2*time.Minute, "max=1h1m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "max")

	errs = validate.Var(time.Duration(0), "omitempty,max=-1s")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"max=1h1m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "max")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,max=-1s"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestMinMaxValidation(t *testing.T) {
	var errs error
	validate := New()

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "min=59m,max=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-time.Minute, "min=59m,max=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+time.Minute, "min=59m,max=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-2*time.Minute, "min=59m,max=1h1m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "min")

	errs = validate.Var(time.Hour+2*time.Minute, "min=59m,max=1h1m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "max")

	errs = validate.Var(time.Duration(0), "omitempty,min=59m,max=1h1m")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"min=59m,max=1h1m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "min")

	timeDurationTest = &TimeDurationTest{time.Hour + 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "max")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,min=59m,max=1h1m"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestLenValidation(t *testing.T) {
	var errs error
	validate := New()

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "len=1h")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour-time.Minute, "len=1h")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "len")

	errs = validate.Var(time.Hour+time.Minute, "len=1h")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "len")

	errs = validate.Var(time.Duration(0), "omitempty,len=1h")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"len=1h"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour - time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "len")

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "len")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,len=1h"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsLt(t *testing.T) {
	var errs error
	validate := New()

	myMap := map[string]string{}
	errs = validate.Var(myMap, "lt=0")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lt")

	f := 1.23
	errs = validate.Var(f, "lt=0")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lt")

	var ui uint = 5
	errs = validate.Var(ui, "lt=0")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lt")

	i := true
	PanicMatches(t, func() { _ = validate.Var(i, "lt") }, "Bad field type bool")

	t1 := time.Now().UTC().Add(-time.Hour)

	errs = validate.Var(t1, "lt")
	Equal(t, errs, nil)

	t2 := time.Now().UTC()
	t2 = t2.Add(time.Hour * 24)

	errs = validate.Var(t2, "lt")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lt")

	type Test struct {
		Now *time.Time `validate:"lt"`
	}

	s := &Test{
		Now: &t1,
	}

	errs = validate.Struct(s)
	Equal(t, errs, nil)

	s = &Test{
		Now: &t2,
	}

	errs = validate.Struct(s)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.Now", "Test.Now", "Now", "Now", "lt")

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "lt=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+time.Minute, "lt=1h1m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lt")

	errs = validate.Var(time.Hour+2*time.Minute, "lt=1h1m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lt")

	errs = validate.Var(time.Duration(0), "omitempty,lt=0")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"lt=1h1m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "lt")

	timeDurationTest = &TimeDurationTest{time.Hour + 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "lt")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,lt=0"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestIsLte(t *testing.T) {
	var errs error
	validate := New()

	i := true
	PanicMatches(t, func() { _ = validate.Var(i, "lte") }, "Bad field type bool")

	t1 := time.Now().UTC().Add(-time.Hour)

	errs = validate.Var(t1, "lte")
	Equal(t, errs, nil)

	t2 := time.Now().UTC()
	t2 = t2.Add(time.Hour * 24)

	errs = validate.Var(t2, "lte")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lte")

	type Test struct {
		Now *time.Time `validate:"lte"`
	}

	s := &Test{
		Now: &t1,
	}

	errs = validate.Struct(s)
	Equal(t, errs, nil)

	s = &Test{
		Now: &t2,
	}

	errs = validate.Struct(s)
	NotEqual(t, errs, nil)

	// Tests for time.Duration type.

	// -- Validations for a variable of time.Duration type.

	errs = validate.Var(time.Hour, "lte=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+time.Minute, "lte=1h1m")
	Equal(t, errs, nil)

	errs = validate.Var(time.Hour+2*time.Minute, "lte=1h1m")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "lte")

	errs = validate.Var(time.Duration(0), "omitempty,lte=-1s")
	Equal(t, errs, nil)

	// -- Validations for a struct with a time.Duration type field.

	type TimeDurationTest struct {
		Duration time.Duration `validate:"lte=1h1m"`
	}
	var timeDurationTest *TimeDurationTest

	timeDurationTest = &TimeDurationTest{time.Hour}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + time.Minute}
	errs = validate.Struct(timeDurationTest)
	Equal(t, errs, nil)

	timeDurationTest = &TimeDurationTest{time.Hour + 2*time.Minute}
	errs = validate.Struct(timeDurationTest)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TimeDurationTest.Duration", "TimeDurationTest.Duration", "Duration", "Duration", "lte")

	type TimeDurationOmitemptyTest struct {
		Duration time.Duration `validate:"omitempty,lte=-1s"`
	}

	timeDurationOmitemptyTest := &TimeDurationOmitemptyTest{time.Duration(0)}
	errs = validate.Struct(timeDurationOmitemptyTest)
	Equal(t, errs, nil)
}

func TestUrnRFC2141(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"urn:a:b", true},
		{"urn:a::", true},
		{"urn:a:-", true},
		{"URN:simple:simple", true},
		{"urn:urna:simple", true},
		{"urn:burnout:nss", true},
		{"urn:burn:nss", true},
		{"urn:urnurnurn:x", true},
		{"urn:abcdefghilmnopqrstuvzabcdefghilm:x", true},
		{"URN:123:x", true},
		{"URN:abcd-:x", true},
		{"URN:abcd-abcd:x", true},
		{"urn:urnx:urn", true},
		{"urn:ciao:a:b:c", true},
		{"urn:aaa:x:y:", true},
		{"urn:ciao:-", true},
		{"urn:colon:::::nss", true},
		{"urn:ciao:@!=%2C(xyz)+a,b.*@g=$_'", true},
		{"URN:hexes:%25", true},
		{"URN:x:abc%1Dz%2F%3az", true},
		{"URN:foo:a123,456", true},
		{"urn:foo:a123,456", true},
		{"urn:FOO:a123,456", true},
		{"urn:foo:A123,456", true},
		{"urn:foo:a123%2C456", true},
		{"URN:FOO:a123%2c456", true},
		{"URN:FOO:ABC%FFabc123%2c456", true},
		{"URN:FOO:ABC%FFabc123%2C456%9A", true},
		{"urn:ietf:params:scim:schemas:core:2.0:User", true},
		{"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:meta.lastModified", true},
		{"URN:-xxx:x", false},
		{"urn::colon:nss", false},
		{"urn:abcdefghilmnopqrstuvzabcdefghilmn:specificstring", false},
		{"URN:a!?:x", false},
		{"URN:#,:x", false},
		{"urn:urn:NSS", false},
		{"urn:URN:NSS", false},
		{"urn:white space:NSS", false},
		{"urn:concat:no spaces", false},
		{"urn:a:%", false},
		{"urn:", false},
	}

	tag := "urn_rfc2141"

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d URN failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d URN failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != tag {
					t.Fatalf("Index: %d URN failed Error: %s", i, errs)
				}
			}
		}
	}

	i := 1
	PanicMatches(t, func() { _ = validate.Var(i, tag) }, "Bad field type int")
}

func TestUrl(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"", false},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"irc://#channel@network", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
		{"irc:", false},
		{"http://", false},
		{"file://path/to/file.txt", true},
		{"file:///c:/Windows/file.txt", true},
		{"file://localhost/path/to/file.txt", true},
		{"file://localhost/c:/WINDOWS/file.txt", true},
		{"file://", true},
		{"file:////remotehost/path/file.txt", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "url")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d URL failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d URL failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "url" {
					t.Fatalf("Index: %d URL failed Error: %s", i, errs)
				}
			}
		}
	}

	i := 1
	PanicMatches(t, func() { _ = validate.Var(i, "url") }, "Bad field type int")
}

func TestHttpUrl(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"HTTP://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", false},
		{"file:///etc/passwd", false},
		{"file://C:/windows/win.ini", false},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"", false},
		{"a://b", false},
		{"xyz://foobar.com", false},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", false},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", false},
		{"irc://irc.server.org/channel", false},
		{"irc://#channel@network", false},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
		{"http:", false},
		{"http://", false},
		{"http://#invalid", false},
		{"https://1.1.1.1", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "http_url")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d HTTP URL failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d HTTP URL failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "http_url" {
					t.Fatalf("Index: %d HTTP URL failed Error: %s", i, errs)
				}
			}
		}
	}

	i := 1
	PanicMatches(t, func() { _ = validate.Var(i, "http_url") }, "Bad field type int")
}

func TestUri(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"", false},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"irc://#channel@network", true},
		{"/abs/test/dir", true},
		{"./rel/test/dir", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uri")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d URI failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d URI failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uri" {
					t.Fatalf("Index: %d URI failed Error: %s", i, errs)
				}
			}
		}
	}

	i := 1
	PanicMatches(t, func() { _ = validate.Var(i, "uri") }, "Bad field type int")
}

func TestOrTag(t *testing.T) {
	validate := New()

	s := "rgba(0,31,255,0.5)"
	errs := validate.Var(s, "rgb|rgba")
	Equal(t, errs, nil)

	s = "rgba(0,31,255,0.5)"
	errs = validate.Var(s, "rgb|rgba|len=18")
	Equal(t, errs, nil)

	s = "this ain't right"
	errs = validate.Var(s, "rgb|rgba")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb|rgba")

	s = "this ain't right"
	errs = validate.Var(s, "rgb|rgba|len=10")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb|rgba|len=10")

	s = "this is right"
	errs = validate.Var(s, "rgb|rgba|len=13")
	Equal(t, errs, nil)

	s = ""
	errs = validate.Var(s, "omitempty,rgb|rgba")
	Equal(t, errs, nil)

	s = "green"
	errs = validate.Var(s, "eq=|eq=blue,rgb|rgba") // should fail on first validation block
	NotEqual(t, errs, nil)
	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	Equal(t, ve[0].Tag(), "eq=|eq=blue")

	s = "this is right, but a blank or isn't"

	PanicMatches(t, func() { _ = validate.Var(s, "rgb||len=13") }, "Invalid validation tag on field ''")
	PanicMatches(t, func() { _ = validate.Var(s, "rgb|rgbaa|len=13") }, "Undefined validation function 'rgbaa' on field ''")

	v2 := New()
	v2.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type Colors struct {
		Fav string `validate:"rgb|rgba" json:"fc"`
	}

	c := Colors{Fav: "this ain't right"}

	err := v2.Struct(c)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)
	fe := getError(errs, "Colors.fc", "Colors.Fav")
	NotEqual(t, fe, nil)
}

func TestHsla(t *testing.T) {
	validate := New()

	s := "hsla(360,100%,100%,1)"
	errs := validate.Var(s, "hsla")
	Equal(t, errs, nil)

	s = "hsla(360,100%,100%,0.5)"
	errs = validate.Var(s, "hsla")
	Equal(t, errs, nil)

	s = "hsla(0,0%,0%, 0)"
	errs = validate.Var(s, "hsla")
	Equal(t, errs, nil)

	s = "hsl(361,100%,50%,1)"
	errs = validate.Var(s, "hsla")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsla")

	s = "hsl(361,100%,50%)"
	errs = validate.Var(s, "hsla")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsla")

	s = "hsla(361,100%,50%)"
	errs = validate.Var(s, "hsla")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsla")

	s = "hsla(360,101%,50%)"
	errs = validate.Var(s, "hsla")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsla")

	s = "hsla(360,100%,101%)"
	errs = validate.Var(s, "hsla")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsla")

	i := 1
	errs = validate.Var(i, "hsla")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsla")
}

func TestHsl(t *testing.T) {
	validate := New()

	s := "hsl(360,100%,50%)"
	errs := validate.Var(s, "hsl")
	Equal(t, errs, nil)

	s = "hsl(0,0%,0%)"
	errs = validate.Var(s, "hsl")
	Equal(t, errs, nil)

	s = "hsl(361,100%,50%)"
	errs = validate.Var(s, "hsl")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsl")

	s = "hsl(361,101%,50%)"
	errs = validate.Var(s, "hsl")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsl")

	s = "hsl(361,100%,101%)"
	errs = validate.Var(s, "hsl")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsl")

	s = "hsl(-10,100%,100%)"
	errs = validate.Var(s, "hsl")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsl")

	i := 1
	errs = validate.Var(i, "hsl")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hsl")
}

func TestRgba(t *testing.T) {
	validate := New()

	s := "rgba(0,31,255,0.5)"
	errs := validate.Var(s, "rgba")
	Equal(t, errs, nil)

	s = "rgba(0,31,255,0.12)"
	errs = validate.Var(s, "rgba")
	Equal(t, errs, nil)

	s = "rgba(12%,55%,100%,0.12)"
	errs = validate.Var(s, "rgba")
	Equal(t, errs, nil)

	s = "rgba( 0,  31, 255, 0.5)"
	errs = validate.Var(s, "rgba")
	Equal(t, errs, nil)

	s = "rgba(12%,55,100%,0.12)"
	errs = validate.Var(s, "rgba")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgba")

	s = "rgb(0,  31, 255)"
	errs = validate.Var(s, "rgba")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgba")

	s = "rgb(1,349,275,0.5)"
	errs = validate.Var(s, "rgba")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgba")

	s = "rgb(01,31,255,0.5)"
	errs = validate.Var(s, "rgba")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgba")

	i := 1
	errs = validate.Var(i, "rgba")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgba")
}

func TestRgb(t *testing.T) {
	validate := New()

	s := "rgb(0,31,255)"
	errs := validate.Var(s, "rgb")
	Equal(t, errs, nil)

	s = "rgb(0,  31, 255)"
	errs = validate.Var(s, "rgb")
	Equal(t, errs, nil)

	s = "rgb(10%,  50%, 100%)"
	errs = validate.Var(s, "rgb")
	Equal(t, errs, nil)

	s = "rgb(10%,  50%, 55)"
	errs = validate.Var(s, "rgb")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb")

	s = "rgb(1,349,275)"
	errs = validate.Var(s, "rgb")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb")

	s = "rgb(01,31,255)"
	errs = validate.Var(s, "rgb")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb")

	s = "rgba(0,31,255)"
	errs = validate.Var(s, "rgb")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb")

	i := 1
	errs = validate.Var(i, "rgb")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "rgb")
}

func TestEmail(t *testing.T) {
	validate := New()

	s := "test@mail.com"
	errs := validate.Var(s, "email")
	Equal(t, errs, nil)

	s = "Dörte@Sörensen.example.com"
	errs = validate.Var(s, "email")
	Equal(t, errs, nil)

	s = "θσερ@εχαμπλε.ψομ"
	errs = validate.Var(s, "email")
	Equal(t, errs, nil)

	s = "юзер@екзампл.ком"
	errs = validate.Var(s, "email")
	Equal(t, errs, nil)

	s = "उपयोगकर्ता@उदाहरण.कॉम"
	errs = validate.Var(s, "email")
	Equal(t, errs, nil)

	s = "用户@例子.广告"
	errs = validate.Var(s, "email")
	Equal(t, errs, nil)

	s = "mail@domain_with_underscores.org"
	errs = validate.Var(s, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")

	s = ""
	errs = validate.Var(s, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")

	s = "test@email"
	errs = validate.Var(s, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")

	s = "test@email."
	errs = validate.Var(s, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")

	s = "@email.com"
	errs = validate.Var(s, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")

	s = `"test test"@email.com`
	errs = validate.Var(s, "email")
	Equal(t, errs, nil)

	s = `"@email.com`
	errs = validate.Var(s, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")

	i := true
	errs = validate.Var(i, "email")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "email")
}

func TestHexColor(t *testing.T) {
	validate := New()

	s := "#fff"
	errs := validate.Var(s, "hexcolor")
	Equal(t, errs, nil)

	s = "#c2c2c2"
	errs = validate.Var(s, "hexcolor")
	Equal(t, errs, nil)

	s = "fff"
	errs = validate.Var(s, "hexcolor")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hexcolor")

	s = "fffFF"
	errs = validate.Var(s, "hexcolor")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hexcolor")

	i := true
	errs = validate.Var(i, "hexcolor")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hexcolor")
}

func TestHexadecimal(t *testing.T) {
	validate := New()

	s := "ff0044"
	errs := validate.Var(s, "hexadecimal")
	Equal(t, errs, nil)

	s = "0xff0044"
	errs = validate.Var(s, "hexadecimal")
	Equal(t, errs, nil)

	s = "0Xff0044"
	errs = validate.Var(s, "hexadecimal")
	Equal(t, errs, nil)

	s = "abcdefg"
	errs = validate.Var(s, "hexadecimal")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hexadecimal")

	i := true
	errs = validate.Var(i, "hexadecimal")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "hexadecimal")
}

func TestNumber(t *testing.T) {
	validate := New()

	s := "1"
	errs := validate.Var(s, "number")
	Equal(t, errs, nil)

	s = "+1"
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	s = "-1"
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	s = "1.12"
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	s = "+1.12"
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	s = "-1.12"
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	s = "1."
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	s = "1.o"
	errs = validate.Var(s, "number")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "number")

	i := 1
	errs = validate.Var(i, "number")
	Equal(t, errs, nil)
}

func TestNumeric(t *testing.T) {
	validate := New()

	s := "1"
	errs := validate.Var(s, "numeric")
	Equal(t, errs, nil)

	s = "+1"
	errs = validate.Var(s, "numeric")
	Equal(t, errs, nil)

	s = "-1"
	errs = validate.Var(s, "numeric")
	Equal(t, errs, nil)

	s = "1.12"
	errs = validate.Var(s, "numeric")
	Equal(t, errs, nil)

	s = "+1.12"
	errs = validate.Var(s, "numeric")
	Equal(t, errs, nil)

	s = "-1.12"
	errs = validate.Var(s, "numeric")
	Equal(t, errs, nil)

	s = "1."
	errs = validate.Var(s, "numeric")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "numeric")

	s = "1.o"
	errs = validate.Var(s, "numeric")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "numeric")

	i := 1
	errs = validate.Var(i, "numeric")
	Equal(t, errs, nil)
}
func TestBoolean(t *testing.T) {
	validate := New()

	b := true
	errs := validate.Var(b, "boolean")
	Equal(t, errs, nil)

	b = false
	errs = validate.Var(b, "boolean")
	Equal(t, errs, nil)

	s := "true"
	errs = validate.Var(s, "boolean")
	Equal(t, errs, nil)

	s = "false"
	errs = validate.Var(s, "boolean")
	Equal(t, errs, nil)

	s = "0"
	errs = validate.Var(s, "boolean")
	Equal(t, errs, nil)

	s = "1"
	errs = validate.Var(s, "boolean")
	Equal(t, errs, nil)

	s = "xyz"
	errs = validate.Var(s, "boolean")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "boolean")

	s = "1."
	errs = validate.Var(s, "boolean")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "boolean")
}

func TestAlphaNumeric(t *testing.T) {
	validate := New()

	s := "abcd123"
	errs := validate.Var(s, "alphanum")
	Equal(t, errs, nil)

	s = "abc!23"
	errs = validate.Var(s, "alphanum")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alphanum")

	errs = validate.Var(1, "alphanum")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alphanum")
}

func TestAlpha(t *testing.T) {
	validate := New()

	s := "abcd"
	errs := validate.Var(s, "alpha")
	Equal(t, errs, nil)

	s = "abc®"
	errs = validate.Var(s, "alpha")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alpha")

	s = "abc÷"
	errs = validate.Var(s, "alpha")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alpha")

	s = "abc1"
	errs = validate.Var(s, "alpha")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alpha")

	s = "this is a test string"
	errs = validate.Var(s, "alpha")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alpha")

	errs = validate.Var(1, "alpha")
	NotEqual(t, errs, nil)
	AssertError(t, errs, "", "", "", "", "alpha")
}

func TestStructStringValidation(t *testing.T) {
	validate := New()

	tSuccess := &TestString{
		Required:  "Required",
		Len:       "length==10",
		Min:       "min=1",
		Max:       "1234567890",
		MinMax:    "12345",
		Lt:        "012345678",
		Lte:       "0123456789",
		Gt:        "01234567890",
		Gte:       "0123456789",
		Boolean:   "true",
		OmitEmpty: "",
		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "1",
		},
		Iface: &Impl{
			F: "123",
		},
	}

	errs := validate.Struct(tSuccess)
	Equal(t, errs, nil)

	tFail := &TestString{
		Required:  "",
		Len:       "",
		Min:       "",
		Max:       "12345678901",
		MinMax:    "",
		Lt:        "0123456789",
		Lte:       "01234567890",
		Gt:        "1",
		Gte:       "1",
		OmitEmpty: "12345678901",
		Boolean:   "nope",
		Sub: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A string `validate:"required"`
		}{
			A: "",
		},
		Iface: &Impl{
			F: "12",
		},
	}

	errs = validate.Struct(tFail)

	// Assert Top Level
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 14)

	// Assert Fields
	AssertError(t, errs, "TestString.Required", "TestString.Required", "Required", "Required", "required")
	AssertError(t, errs, "TestString.Len", "TestString.Len", "Len", "Len", "len")
	AssertError(t, errs, "TestString.Min", "TestString.Min", "Min", "Min", "min")
	AssertError(t, errs, "TestString.Max", "TestString.Max", "Max", "Max", "max")
	AssertError(t, errs, "TestString.MinMax", "TestString.MinMax", "MinMax", "MinMax", "min")
	AssertError(t, errs, "TestString.Lt", "TestString.Lt", "Lt", "Lt", "lt")
	AssertError(t, errs, "TestString.Lte", "TestString.Lte", "Lte", "Lte", "lte")
	AssertError(t, errs, "TestString.Gt", "TestString.Gt", "Gt", "Gt", "gt")
	AssertError(t, errs, "TestString.Gte", "TestString.Gte", "Gte", "Gte", "gte")
	AssertError(t, errs, "TestString.OmitEmpty", "TestString.OmitEmpty", "OmitEmpty", "OmitEmpty", "max")
	AssertError(t, errs, "TestString.Boolean", "TestString.Boolean", "Boolean", "Boolean", "boolean")

	// Nested Struct Field Errs
	AssertError(t, errs, "TestString.Anonymous.A", "TestString.Anonymous.A", "A", "A", "required")
	AssertError(t, errs, "TestString.Sub.Test", "TestString.Sub.Test", "Test", "Test", "required")
	AssertError(t, errs, "TestString.Iface.F", "TestString.Iface.F", "F", "F", "len")
}

func TestStructInt32Validation(t *testing.T) {
	type TestInt32 struct {
		Required  int `validate:"required"`
		Len       int `validate:"len=10"`
		Min       int `validate:"min=1"`
		Max       int `validate:"max=10"`
		MinMax    int `validate:"min=1,max=10"`
		Lt        int `validate:"lt=10"`
		Lte       int `validate:"lte=10"`
		Gt        int `validate:"gt=10"`
		Gte       int `validate:"gte=10"`
		OmitEmpty int `validate:"omitempty,min=1,max=10"`
	}

	tSuccess := &TestInt32{
		Required:  1,
		Len:       10,
		Min:       1,
		Max:       10,
		MinMax:    5,
		Lt:        9,
		Lte:       10,
		Gt:        11,
		Gte:       10,
		OmitEmpty: 0,
	}

	validate := New()
	errs := validate.Struct(tSuccess)
	Equal(t, errs, nil)

	tFail := &TestInt32{
		Required:  0,
		Len:       11,
		Min:       -1,
		Max:       11,
		MinMax:    -1,
		Lt:        10,
		Lte:       11,
		Gt:        10,
		Gte:       9,
		OmitEmpty: 11,
	}

	errs = validate.Struct(tFail)

	// Assert Top Level
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 10)

	// Assert Fields
	AssertError(t, errs, "TestInt32.Required", "TestInt32.Required", "Required", "Required", "required")
	AssertError(t, errs, "TestInt32.Len", "TestInt32.Len", "Len", "Len", "len")
	AssertError(t, errs, "TestInt32.Min", "TestInt32.Min", "Min", "Min", "min")
	AssertError(t, errs, "TestInt32.Max", "TestInt32.Max", "Max", "Max", "max")
	AssertError(t, errs, "TestInt32.MinMax", "TestInt32.MinMax", "MinMax", "MinMax", "min")
	AssertError(t, errs, "TestInt32.Lt", "TestInt32.Lt", "Lt", "Lt", "lt")
	AssertError(t, errs, "TestInt32.Lte", "TestInt32.Lte", "Lte", "Lte", "lte")
	AssertError(t, errs, "TestInt32.Gt", "TestInt32.Gt", "Gt", "Gt", "gt")
	AssertError(t, errs, "TestInt32.Gte", "TestInt32.Gte", "Gte", "Gte", "gte")
	AssertError(t, errs, "TestInt32.OmitEmpty", "TestInt32.OmitEmpty", "OmitEmpty", "OmitEmpty", "max")
}

func TestStructUint64Validation(t *testing.T) {
	validate := New()

	tSuccess := &TestUint64{
		Required:  1,
		Len:       10,
		Min:       1,
		Max:       10,
		MinMax:    5,
		OmitEmpty: 0,
	}

	errs := validate.Struct(tSuccess)
	Equal(t, errs, nil)

	tFail := &TestUint64{
		Required:  0,
		Len:       11,
		Min:       0,
		Max:       11,
		MinMax:    0,
		OmitEmpty: 11,
	}

	errs = validate.Struct(tFail)

	// Assert Top Level
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 6)

	// Assert Fields
	AssertError(t, errs, "TestUint64.Required", "TestUint64.Required", "Required", "Required", "required")
	AssertError(t, errs, "TestUint64.Len", "TestUint64.Len", "Len", "Len", "len")
	AssertError(t, errs, "TestUint64.Min", "TestUint64.Min", "Min", "Min", "min")
	AssertError(t, errs, "TestUint64.Max", "TestUint64.Max", "Max", "Max", "max")
	AssertError(t, errs, "TestUint64.MinMax", "TestUint64.MinMax", "MinMax", "MinMax", "min")
	AssertError(t, errs, "TestUint64.OmitEmpty", "TestUint64.OmitEmpty", "OmitEmpty", "OmitEmpty", "max")
}

func TestStructFloat64Validation(t *testing.T) {
	validate := New()

	tSuccess := &TestFloat64{
		Required:  1,
		Len:       10,
		Min:       1,
		Max:       10,
		MinMax:    5,
		OmitEmpty: 0,
	}

	errs := validate.Struct(tSuccess)
	Equal(t, errs, nil)

	tFail := &TestFloat64{
		Required:  0,
		Len:       11,
		Min:       0,
		Max:       11,
		MinMax:    0,
		OmitEmpty: 11,
	}

	errs = validate.Struct(tFail)

	// Assert Top Level
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 6)

	// Assert Fields
	AssertError(t, errs, "TestFloat64.Required", "TestFloat64.Required", "Required", "Required", "required")
	AssertError(t, errs, "TestFloat64.Len", "TestFloat64.Len", "Len", "Len", "len")
	AssertError(t, errs, "TestFloat64.Min", "TestFloat64.Min", "Min", "Min", "min")
	AssertError(t, errs, "TestFloat64.Max", "TestFloat64.Max", "Max", "Max", "max")
	AssertError(t, errs, "TestFloat64.MinMax", "TestFloat64.MinMax", "MinMax", "MinMax", "min")
	AssertError(t, errs, "TestFloat64.OmitEmpty", "TestFloat64.OmitEmpty", "OmitEmpty", "OmitEmpty", "max")
}

func TestStructSliceValidation(t *testing.T) {
	validate := New()

	tSuccess := &TestSlice{
		Required:  []int{1},
		Len:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		Min:       []int{1, 2},
		Max:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		MinMax:    []int{1, 2, 3, 4, 5},
		OmitEmpty: nil,
	}

	errs := validate.Struct(tSuccess)
	Equal(t, errs, nil)

	tFail := &TestSlice{
		Required:  nil,
		Len:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1},
		Min:       []int{},
		Max:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1},
		MinMax:    []int{},
		OmitEmpty: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1},
	}

	errs = validate.Struct(tFail)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 6)

	// Assert Field Errors
	AssertError(t, errs, "TestSlice.Required", "TestSlice.Required", "Required", "Required", "required")
	AssertError(t, errs, "TestSlice.Len", "TestSlice.Len", "Len", "Len", "len")
	AssertError(t, errs, "TestSlice.Min", "TestSlice.Min", "Min", "Min", "min")
	AssertError(t, errs, "TestSlice.Max", "TestSlice.Max", "Max", "Max", "max")
	AssertError(t, errs, "TestSlice.MinMax", "TestSlice.MinMax", "MinMax", "MinMax", "min")
	AssertError(t, errs, "TestSlice.OmitEmpty", "TestSlice.OmitEmpty", "OmitEmpty", "OmitEmpty", "max")

	fe := getError(errs, "TestSlice.Len", "TestSlice.Len")
	NotEqual(t, fe, nil)
	Equal(t, fe.Field(), "Len")
	Equal(t, fe.StructField(), "Len")
	Equal(t, fe.Namespace(), "TestSlice.Len")
	Equal(t, fe.StructNamespace(), "TestSlice.Len")
	Equal(t, fe.Tag(), "len")
	Equal(t, fe.ActualTag(), "len")
	Equal(t, fe.Param(), "10")
	Equal(t, fe.Kind(), reflect.Slice)
	Equal(t, fe.Type(), reflect.TypeOf([]int{}))

	_, ok := fe.Value().([]int)
	Equal(t, ok, true)
}

func TestInvalidStruct(t *testing.T) {
	validate := New()

	s := &SubTest{
		Test: "1",
	}

	err := validate.Struct(s.Test)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "validator: (nil string)")

	err = validate.Struct(nil)
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "validator: (nil)")

	err = validate.StructPartial(nil, "SubTest.Test")
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "validator: (nil)")

	err = validate.StructExcept(nil, "SubTest.Test")
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "validator: (nil)")
}

func TestInvalidValidatorFunction(t *testing.T) {
	validate := New()

	s := &SubTest{
		Test: "1",
	}

	PanicMatches(t, func() { _ = validate.Var(s.Test, "zzxxBadFunction") }, "Undefined validation function 'zzxxBadFunction' on field ''")
}

func TestCustomFieldName(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("schema"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type A struct {
		B string `schema:"b" validate:"required"`
		C string `schema:"c" validate:"required"`
		D []bool `schema:"d" validate:"required"`
		E string `schema:"-" validate:"required"`
	}

	a := &A{}

	err := validate.Struct(a)
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)
	Equal(t, len(errs), 4)
	Equal(t, getError(errs, "A.b", "A.B").Field(), "b")
	Equal(t, getError(errs, "A.c", "A.C").Field(), "c")
	Equal(t, getError(errs, "A.d", "A.D").Field(), "d")
	Equal(t, getError(errs, "A.E", "A.E").Field(), "E")

	v2 := New()
	err = v2.Struct(a)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)
	Equal(t, len(errs), 4)
	Equal(t, getError(errs, "A.B", "A.B").Field(), "B")
	Equal(t, getError(errs, "A.C", "A.C").Field(), "C")
	Equal(t, getError(errs, "A.D", "A.D").Field(), "D")
	Equal(t, getError(errs, "A.E", "A.E").Field(), "E")
}

func TestMultipleRecursiveExtractStructCache(t *testing.T) {
	validate := New()

	type Recursive struct {
		Field *string `validate:"required,len=5,ne=string"`
	}

	var test Recursive

	current := reflect.ValueOf(test)
	name := "Recursive"
	proceed := make(chan struct{})

	sc := validate.extractStructCache(current, name)
	ptr := fmt.Sprintf("%p", sc)

	for i := 0; i < 100; i++ {
		go func() {
			<-proceed
			sc := validate.extractStructCache(current, name)
			Equal(t, ptr, fmt.Sprintf("%p", sc))
		}()
	}

	close(proceed)
}

// Thanks @robbrockbank, see https://github.com/go-playground/validator/issues/249
func TestPointerAndOmitEmpty(t *testing.T) {
	validate := New()

	type Test struct {
		MyInt *int `validate:"omitempty,gte=2,lte=255"`
	}

	val1 := 0
	val2 := 256

	t1 := Test{MyInt: &val1} // This should fail validation on gte because value is 0
	t2 := Test{MyInt: &val2} // This should fail validate on lte because value is 256
	t3 := Test{MyInt: nil}   // This should succeed validation because pointer is nil

	errs := validate.Struct(t1)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.MyInt", "Test.MyInt", "MyInt", "MyInt", "gte")

	errs = validate.Struct(t2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "Test.MyInt", "Test.MyInt", "MyInt", "MyInt", "lte")

	errs = validate.Struct(t3)
	Equal(t, errs, nil)

	type TestIface struct {
		MyInt interface{} `validate:"omitempty,gte=2,lte=255"`
	}

	ti1 := TestIface{MyInt: &val1} // This should fail validation on gte because value is 0
	ti2 := TestIface{MyInt: &val2} // This should fail validate on lte because value is 256
	ti3 := TestIface{MyInt: nil}   // This should succeed validation because pointer is nil

	errs = validate.Struct(ti1)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestIface.MyInt", "TestIface.MyInt", "MyInt", "MyInt", "gte")

	errs = validate.Struct(ti2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestIface.MyInt", "TestIface.MyInt", "MyInt", "MyInt", "lte")

	errs = validate.Struct(ti3)
	Equal(t, errs, nil)
}

func TestRequired(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type Test struct {
		Value interface{} `validate:"required"`
	}

	var test Test

	err := validate.Struct(test)
	NotEqual(t, err, nil)
	AssertError(t, err.(ValidationErrors), "Test.Value", "Test.Value", "Value", "Value", "required")
}

func TestBoolEqual(t *testing.T) {
	validate := New()

	type Test struct {
		Value bool `validate:"eq=true"`
	}

	var test Test

	err := validate.Struct(test)
	NotEqual(t, err, nil)
	AssertError(t, err.(ValidationErrors), "Test.Value", "Test.Value", "Value", "Value", "eq")

	test.Value = true
	err = validate.Struct(test)
	Equal(t, err, nil)
}

func TestTranslations(t *testing.T) {
	en := en.New()
	uni := ut.New(en, en, fr.New())

	trans, _ := uni.GetTranslator("en")
	fr, _ := uni.GetTranslator("fr")

	validate := New()
	err := validate.RegisterTranslation("required", trans,
		func(ut ut.Translator) (err error) {
			// using this stype because multiple translation may have to be added for the full translation
			if err = ut.Add("required", "{0} is a required field", false); err != nil {
				return
			}

			return
		}, func(ut ut.Translator, fe FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %#v", fe.(*fieldError))
				return fe.(*fieldError).Error()
			}

			return t
		})
	Equal(t, err, nil)

	err = validate.RegisterTranslation("required", fr,
		func(ut ut.Translator) (err error) {
			// using this stype because multiple translation may have to be added for the full translation
			if err = ut.Add("required", "{0} est un champ obligatoire", false); err != nil {
				return
			}

			return
		}, func(ut ut.Translator, fe FieldError) string {
			t, transErr := ut.T(fe.Tag(), fe.Field())
			if transErr != nil {
				fmt.Printf("warning: error translating FieldError: %#v", fe.(*fieldError))
				return fe.(*fieldError).Error()
			}

			return t
		})

	Equal(t, err, nil)

	type Test struct {
		Value interface{} `validate:"required"`
	}

	var test Test

	err = validate.Struct(test)
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)
	Equal(t, len(errs), 1)

	fe := errs[0]
	Equal(t, fe.Tag(), "required")
	Equal(t, fe.Namespace(), "Test.Value")
	Equal(t, fe.Translate(trans), fmt.Sprintf("%s is a required field", fe.Field()))
	Equal(t, fe.Translate(fr), fmt.Sprintf("%s est un champ obligatoire", fe.Field()))

	nl := nl.New()
	uni2 := ut.New(nl, nl)
	trans2, _ := uni2.GetTranslator("nl")
	Equal(t, fe.Translate(trans2), "Key: 'Test.Value' Error:Field validation for 'Value' failed on the 'required' tag")

	terrs := errs.Translate(trans)
	Equal(t, len(terrs), 1)

	v, ok := terrs["Test.Value"]
	Equal(t, ok, true)
	Equal(t, v, fmt.Sprintf("%s is a required field", fe.Field()))

	terrs = errs.Translate(fr)
	Equal(t, len(terrs), 1)

	v, ok = terrs["Test.Value"]
	Equal(t, ok, true)
	Equal(t, v, fmt.Sprintf("%s est un champ obligatoire", fe.Field()))

	type Test2 struct {
		Value string `validate:"gt=1"`
	}

	var t2 Test2

	err = validate.Struct(t2)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)
	Equal(t, len(errs), 1)

	fe = errs[0]
	Equal(t, fe.Tag(), "gt")
	Equal(t, fe.Namespace(), "Test2.Value")
	Equal(t, fe.Translate(trans), "Key: 'Test2.Value' Error:Field validation for 'Value' failed on the 'gt' tag")
}

func TestTranslationErrors(t *testing.T) {
	en := en.New()
	uni := ut.New(en, en, fr.New())

	trans, _ := uni.GetTranslator("en")
	err := trans.Add("required", "{0} is a required field", false) // using translator outside of validator also
	Equal(t, err, nil)

	validate := New()
	err = validate.RegisterTranslation("required", trans,
		func(ut ut.Translator) (err error) {
			// using this stype because multiple translation may have to be added for the full translation
			if err = ut.Add("required", "{0} is a required field", false); err != nil {
				return
			}

			return
		}, func(ut ut.Translator, fe FieldError) string {
			t, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				fmt.Printf("warning: error translating FieldError: %#v", fe.(*fieldError))
				return fe.(*fieldError).Error()
			}

			return t
		})

	NotEqual(t, err, nil)
	Equal(t, err.Error(), "error: conflicting key 'required' rule 'Unknown' with text '{0} is a required field' for locale 'en', value being ignored")
}

func TestStructFiltered(t *testing.T) {
	p1 := func(ns []byte) bool {
		if bytes.HasSuffix(ns, []byte("NoTag")) || bytes.HasSuffix(ns, []byte("Required")) {
			return false
		}

		return true
	}

	p2 := func(ns []byte) bool {
		if bytes.HasSuffix(ns, []byte("SubSlice[0].Test")) ||
			bytes.HasSuffix(ns, []byte("SubSlice[0]")) ||
			bytes.HasSuffix(ns, []byte("SubSlice")) ||
			bytes.HasSuffix(ns, []byte("Sub")) ||
			bytes.HasSuffix(ns, []byte("SubIgnore")) ||
			bytes.HasSuffix(ns, []byte("Anonymous")) ||
			bytes.HasSuffix(ns, []byte("Anonymous.A")) {
			return false
		}

		return true
	}

	p3 := func(ns []byte) bool {
		return !bytes.HasSuffix(ns, []byte("SubTest.Test"))
	}

	// p4 := []string{
	// 	"A",
	// }

	tPartial := &TestPartial{
		NoTag:    "NoTag",
		Required: "Required",

		SubSlice: []*SubTest{
			{

				Test: "Required",
			},
			{

				Test: "Required",
			},
		},

		Sub: &SubTest{
			Test: "1",
		},
		SubIgnore: &SubTest{
			Test: "",
		},
		Anonymous: struct {
			A             string     `validate:"required"`
			ASubSlice     []*SubTest `validate:"required,dive"`
			SubAnonStruct []struct {
				Test      string `validate:"required"`
				OtherTest string `validate:"required"`
			} `validate:"required,dive"`
		}{
			A: "1",
			ASubSlice: []*SubTest{
				{
					Test: "Required",
				},
				{
					Test: "Required",
				},
			},

			SubAnonStruct: []struct {
				Test      string `validate:"required"`
				OtherTest string `validate:"required"`
			}{
				{"Required", "RequiredOther"},
				{"Required", "RequiredOther"},
			},
		},
	}

	validate := New()

	// the following should all return no errors as everything is valid in
	// the default state
	errs := validate.StructFilteredCtx(context.Background(), tPartial, p1)
	Equal(t, errs, nil)

	errs = validate.StructFiltered(tPartial, p2)
	Equal(t, errs, nil)

	// this isn't really a robust test, but is meant to illustrate the ANON CASE below
	errs = validate.StructFiltered(tPartial.SubSlice[0], p3)
	Equal(t, errs, nil)

	// mod tPartial for required field and re-test making sure invalid fields are NOT required:
	tPartial.Required = ""

	// inversion and retesting Partial to generate failures:
	errs = validate.StructFiltered(tPartial, p1)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Required", "TestPartial.Required", "Required", "Required", "required")

	// reset Required field, and set nested struct
	tPartial.Required = "Required"
	tPartial.Anonymous.A = ""

	// will pass as unset fields is not going to be tested
	errs = validate.StructFiltered(tPartial, p1)
	Equal(t, errs, nil)

	// will fail as unset field is tested
	errs = validate.StructFiltered(tPartial, p2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.Anonymous.A", "TestPartial.Anonymous.A", "A", "A", "required")

	// reset nested struct and unset struct in slice
	tPartial.Anonymous.A = "Required"
	tPartial.SubSlice[0].Test = ""

	// these will pass as unset item is NOT tested
	errs = validate.StructFiltered(tPartial, p1)
	Equal(t, errs, nil)

	errs = validate.StructFiltered(tPartial, p2)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestPartial.SubSlice[0].Test", "TestPartial.SubSlice[0].Test", "Test", "Test", "required")
	Equal(t, len(errs.(ValidationErrors)), 1)

	// Unset second slice member concurrently to test dive behavior:
	tPartial.SubSlice[1].Test = ""

	errs = validate.StructFiltered(tPartial, p1)
	Equal(t, errs, nil)

	errs = validate.StructFiltered(tPartial, p2)
	NotEqual(t, errs, nil)
	Equal(t, len(errs.(ValidationErrors)), 1)
	AssertError(t, errs, "TestPartial.SubSlice[0].Test", "TestPartial.SubSlice[0].Test", "Test", "Test", "required")

	// reset struct in slice, and unset struct in slice in unset position
	tPartial.SubSlice[0].Test = "Required"

	// these will pass as the unset item is NOT tested
	errs = validate.StructFiltered(tPartial, p1)
	Equal(t, errs, nil)

	errs = validate.StructFiltered(tPartial, p2)
	Equal(t, errs, nil)

	tPartial.SubSlice[1].Test = "Required"
	tPartial.Anonymous.SubAnonStruct[0].Test = ""

	// these will pass as the unset item is NOT tested
	errs = validate.StructFiltered(tPartial, p1)
	Equal(t, errs, nil)

	errs = validate.StructFiltered(tPartial, p2)
	Equal(t, errs, nil)

	dt := time.Now()
	err := validate.StructFiltered(&dt, func(ns []byte) bool { return true })
	NotEqual(t, err, nil)
	Equal(t, err.Error(), "validator: (nil *time.Time)")
}

func TestRequiredPtr(t *testing.T) {
	type Test struct {
		Bool *bool `validate:"required"`
	}

	validate := New()

	f := false

	test := Test{
		Bool: &f,
	}

	err := validate.Struct(test)
	Equal(t, err, nil)

	tr := true

	test.Bool = &tr

	err = validate.Struct(test)
	Equal(t, err, nil)

	test.Bool = nil

	err = validate.Struct(test)
	NotEqual(t, err, nil)

	errs, ok := err.(ValidationErrors)
	Equal(t, ok, true)
	Equal(t, len(errs), 1)
	AssertError(t, errs, "Test.Bool", "Test.Bool", "Bool", "Bool", "required")

	type Test2 struct {
		Bool bool `validate:"required"`
	}

	var test2 Test2

	err = validate.Struct(test2)
	NotEqual(t, err, nil)

	errs, ok = err.(ValidationErrors)
	Equal(t, ok, true)
	Equal(t, len(errs), 1)
	AssertError(t, errs, "Test2.Bool", "Test2.Bool", "Bool", "Bool", "required")

	test2.Bool = true

	err = validate.Struct(test2)
	Equal(t, err, nil)

	type Test3 struct {
		Arr []string `validate:"required"`
	}

	var test3 Test3

	err = validate.Struct(test3)
	NotEqual(t, err, nil)

	errs, ok = err.(ValidationErrors)
	Equal(t, ok, true)
	Equal(t, len(errs), 1)
	AssertError(t, errs, "Test3.Arr", "Test3.Arr", "Arr", "Arr", "required")

	test3.Arr = make([]string, 0)

	err = validate.Struct(test3)
	Equal(t, err, nil)

	type Test4 struct {
		Arr *[]string `validate:"required"` // I know I know pointer to array, just making sure validation works as expected...
	}

	var test4 Test4

	err = validate.Struct(test4)
	NotEqual(t, err, nil)

	errs, ok = err.(ValidationErrors)
	Equal(t, ok, true)
	Equal(t, len(errs), 1)
	AssertError(t, errs, "Test4.Arr", "Test4.Arr", "Arr", "Arr", "required")

	arr := make([]string, 0)
	test4.Arr = &arr

	err = validate.Struct(test4)
	Equal(t, err, nil)
}

func TestAlphaUnicodeValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"this is a test string", false},
		{"这是一个测试字符串", true},
		{"123", false},
		{"<>@;.-=", false},
		{"ひらがな・カタカナ、．漢字", false},
		{"あいうえおfoobar", true},
		{"test＠example.com", false},
		{"1234abcDE", false},
		{"ｶﾀｶﾅ", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "alphaunicode")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d Alpha Unicode failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d Alpha Unicode failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "alphaunicode" {
					t.Fatalf("Index: %d Alpha Unicode failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestAlphanumericUnicodeValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"this is a test string", false},
		{"这是一个测试字符串", true},
		{"\u0031\u0032\u0033", true}, // unicode 5
		{"123", true},
		{"<>@;.-=", false},
		{"ひらがな・カタカナ、．漢字", false},
		{"あいうえおfoobar", true},
		{"test＠example.com", false},
		{"1234abcDE", true},
		{"ｶﾀｶﾅ", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "alphanumunicode")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d Alphanum Unicode failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d Alphanum Unicode failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "alphanumunicode" {
					t.Fatalf("Index: %d Alphanum Unicode failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestArrayStructNamespace(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type child struct {
		Name string `json:"name" validate:"required"`
	}
	var input struct {
		Children []child `json:"children" validate:"required,gt=0,dive"`
	}
	input.Children = []child{{"ok"}, {""}}

	errs := validate.Struct(input)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "children[1].name", "Children[1].Name", "name", "Name", "required")
}

func TestMapStructNamespace(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type child struct {
		Name string `json:"name" validate:"required"`
	}
	var input struct {
		Children map[int]child `json:"children" validate:"required,gt=0,dive"`
	}
	input.Children = map[int]child{
		0: {Name: "ok"},
		1: {Name: ""},
	}

	errs := validate.Struct(input)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "children[1].name", "Children[1].Name", "name", "Name", "required")
}

func TestFieldLevelName(t *testing.T) {
	type Test struct {
		String string            `validate:"custom1"      json:"json1"`
		Array  []string          `validate:"dive,custom2" json:"json2"`
		Map    map[string]string `validate:"dive,custom3" json:"json3"`
		Array2 []string          `validate:"custom4"      json:"json4"`
		Map2   map[string]string `validate:"custom5"      json:"json5"`
	}

	var res1, res2, res3, res4, res5, alt1, alt2, alt3, alt4, alt5 string
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	err := validate.RegisterValidation("custom1", func(fl FieldLevel) bool {
		res1 = fl.FieldName()
		alt1 = fl.StructFieldName()
		return true
	})
	Equal(t, err, nil)

	err = validate.RegisterValidation("custom2", func(fl FieldLevel) bool {
		res2 = fl.FieldName()
		alt2 = fl.StructFieldName()
		return true
	})
	Equal(t, err, nil)

	err = validate.RegisterValidation("custom3", func(fl FieldLevel) bool {
		res3 = fl.FieldName()
		alt3 = fl.StructFieldName()
		return true
	})
	Equal(t, err, nil)

	err = validate.RegisterValidation("custom4", func(fl FieldLevel) bool {
		res4 = fl.FieldName()
		alt4 = fl.StructFieldName()
		return true
	})
	Equal(t, err, nil)

	err = validate.RegisterValidation("custom5", func(fl FieldLevel) bool {
		res5 = fl.FieldName()
		alt5 = fl.StructFieldName()
		return true
	})
	Equal(t, err, nil)

	test := Test{
		String: "test",
		Array:  []string{"1"},
		Map:    map[string]string{"test": "test"},
	}

	errs := validate.Struct(test)
	Equal(t, errs, nil)
	Equal(t, res1, "json1")
	Equal(t, alt1, "String")
	Equal(t, res2, "json2[0]")
	Equal(t, alt2, "Array[0]")
	Equal(t, res3, "json3[test]")
	Equal(t, alt3, "Map[test]")
	Equal(t, res4, "json4")
	Equal(t, alt4, "Array2")
	Equal(t, res5, "json5")
	Equal(t, alt5, "Map2")
}

func TestValidateStructRegisterCtx(t *testing.T) {
	var ctxVal string

	fnCtx := func(ctx context.Context, fl FieldLevel) bool {
		ctxVal = ctx.Value(&ctxVal).(string)
		return true
	}

	var ctxSlVal string
	slFn := func(ctx context.Context, sl StructLevel) {
		ctxSlVal = ctx.Value(&ctxSlVal).(string)
	}

	type Test struct {
		Field string `validate:"val"`
	}

	var tst Test

	validate := New()
	err := validate.RegisterValidationCtx("val", fnCtx)
	Equal(t, err, nil)

	validate.RegisterStructValidationCtx(slFn, Test{})

	ctx := context.WithValue(context.Background(), &ctxVal, "testval")
	ctx = context.WithValue(ctx, &ctxSlVal, "slVal")
	errs := validate.StructCtx(ctx, tst)
	Equal(t, errs, nil)
	Equal(t, ctxVal, "testval")
	Equal(t, ctxSlVal, "slVal")
}

func TestHostnameRFC952Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"test.example.com", true},
		{"example.com", true},
		{"example24.com", true},
		{"test.example24.com", true},
		{"test24.example24.com", true},
		{"example", true},
		{"EXAMPLE", true},
		{"1.foo.com", false},
		{"test.example.com.", false},
		{"example.com.", false},
		{"example24.com.", false},
		{"test.example24.com.", false},
		{"test24.example24.com.", false},
		{"example.", false},
		{"192.168.0.1", false},
		{"email@example.com", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", false},
		{"2001:cdba:0:0:0:0:3257:9652", false},
		{"2001:cdba::3257:9652", false},
		{"example..........com", false},
		{"1234", false},
		{"abc1234", true},
		{"example. com", false},
		{"ex ample.com", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "hostname")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d hostname failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d hostname failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "hostname" {
					t.Fatalf("Index: %d hostname failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestHostnameRFC1123Validation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"test.example.com", true},
		{"example.com", true},
		{"example24.com", true},
		{"test.example24.com", true},
		{"test24.example24.com", true},
		{"example", true},
		{"1.foo.com", true},
		{"test.example.com.", false},
		{"example.com.", false},
		{"example24.com.", false},
		{"test.example24.com.", false},
		{"test24.example24.com.", false},
		{"example.", false},
		{"test_example", false},
		{"192.168.0.1", true},
		{"email@example.com", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", false},
		{"2001:cdba:0:0:0:0:3257:9652", false},
		{"2001:cdba::3257:9652", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "hostname_rfc1123")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Hostname: %v failed Error: %v", test, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Hostname: %v failed Error: %v", test, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "hostname_rfc1123" {
					t.Fatalf("Hostname: %v failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestHostnameRFC1123AliasValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"test.example.com", true},
		{"example.com", true},
		{"example24.com", true},
		{"test.example24.com", true},
		{"test24.example24.com", true},
		{"example", true},
		{"1.foo.com", true},
		{"test.example.com.", false},
		{"example.com.", false},
		{"example24.com.", false},
		{"test.example24.com.", false},
		{"test24.example24.com.", false},
		{"example.", false},
		{"test_example", false},
		{"192.168.0.1", true},
		{"email@example.com", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", false},
		{"2001:cdba:0:0:0:0:3257:9652", false},
		{"2001:cdba::3257:9652", false},
	}

	validate := New()
	validate.RegisterAlias("hostname", "hostname_rfc1123")

	for i, test := range tests {

		errs := validate.Var(test.param, "hostname")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d hostname failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d hostname failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "hostname" {
					t.Fatalf("Index: %d hostname failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestFQDNValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"test.example.com", true},
		{"example.com", true},
		{"example24.com", true},
		{"test.example24.com", true},
		{"test24.example24.com", true},
		{"test.example.com.", true},
		{"example.com.", true},
		{"example24.com.", true},
		{"test.example24.com.", true},
		{"test24.example24.com.", true},
		{"24.example24.com", true},
		{"test.24.example.com", true},
		{"test24.example24.com..", false},
		{"example", false},
		{"192.168.0.1", false},
		{"email@example.com", false},
		{"2001:cdba:0000:0000:0000:0000:3257:9652", false},
		{"2001:cdba:0:0:0:0:3257:9652", false},
		{"2001:cdba::3257:9652", false},
		{"", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "fqdn")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d fqdn failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d fqdn failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "fqdn" {
					t.Fatalf("Index: %d fqdn failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestIsDefault(t *testing.T) {
	validate := New()

	type Inner struct {
		String string `validate:"isdefault"`
	}
	type Test struct {
		String string `validate:"isdefault"`
		Inner  *Inner `validate:"isdefault"`
	}

	var tt Test

	errs := validate.Struct(tt)
	Equal(t, errs, nil)

	tt.Inner = &Inner{String: ""}
	errs = validate.Struct(tt)
	NotEqual(t, errs, nil)

	fe := errs.(ValidationErrors)[0]
	Equal(t, fe.Field(), "Inner")
	Equal(t, fe.Namespace(), "Test.Inner")
	Equal(t, fe.Tag(), "isdefault")

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	type Inner2 struct {
		String string `validate:"isdefault"`
	}

	type Test2 struct {
		Inner Inner2 `validate:"isdefault" json:"inner"`
	}

	var t2 Test2
	errs = validate.Struct(t2)
	Equal(t, errs, nil)

	t2.Inner.String = "Changed"
	errs = validate.Struct(t2)
	NotEqual(t, errs, nil)

	fe = errs.(ValidationErrors)[0]
	Equal(t, fe.Field(), "inner")
	Equal(t, fe.Namespace(), "Test2.inner")
	Equal(t, fe.Tag(), "isdefault")
}

func TestUniqueValidation(t *testing.T) {
	tests := []struct {
		param    interface{}
		expected bool
	}{
		// Arrays
		{[2]string{"a", "b"}, true},
		{[2]int{1, 2}, true},
		{[2]float64{1, 2}, true},
		{[2]interface{}{"a", "b"}, true},
		{[2]interface{}{"a", 1}, true},
		{[2]float64{1, 1}, false},
		{[2]int{1, 1}, false},
		{[2]string{"a", "a"}, false},
		{[2]interface{}{"a", "a"}, false},
		{[4]interface{}{"a", 1, "b", 1}, false},
		{[2]*string{stringPtr("a"), stringPtr("b")}, true},
		{[2]*int{intPtr(1), intPtr(2)}, true},
		{[2]*float64{float64Ptr(1), float64Ptr(2)}, true},
		{[2]*string{stringPtr("a"), stringPtr("a")}, false},
		{[2]*float64{float64Ptr(1), float64Ptr(1)}, false},
		{[2]*int{intPtr(1), intPtr(1)}, false},
		// Slices
		{[]string{"a", "b"}, true},
		{[]int{1, 2}, true},
		{[]float64{1, 2}, true},
		{[]interface{}{"a", "b"}, true},
		{[]interface{}{"a", 1}, true},
		{[]float64{1, 1}, false},
		{[]int{1, 1}, false},
		{[]string{"a", "a"}, false},
		{[]interface{}{"a", "a"}, false},
		{[]interface{}{"a", 1, "b", 1}, false},
		{[]*string{stringPtr("a"), stringPtr("b")}, true},
		{[]*int{intPtr(1), intPtr(2)}, true},
		{[]*float64{float64Ptr(1), float64Ptr(2)}, true},
		{[]*string{stringPtr("a"), stringPtr("a")}, false},
		{[]*float64{float64Ptr(1), float64Ptr(1)}, false},
		{[]*int{intPtr(1), intPtr(1)}, false},
		// Maps
		{map[string]string{"one": "a", "two": "b"}, true},
		{map[string]int{"one": 1, "two": 2}, true},
		{map[string]float64{"one": 1, "two": 2}, true},
		{map[string]interface{}{"one": "a", "two": "b"}, true},
		{map[string]interface{}{"one": "a", "two": 1}, true},
		{map[string]float64{"one": 1, "two": 1}, false},
		{map[string]int{"one": 1, "two": 1}, false},
		{map[string]string{"one": "a", "two": "a"}, false},
		{map[string]interface{}{"one": "a", "two": "a"}, false},
		{map[string]interface{}{"one": "a", "two": 1, "three": "b", "four": 1}, false},
		{map[string]*string{"one": stringPtr("a"), "two": stringPtr("a")}, false},
		{map[string]*string{"one": stringPtr("a"), "two": stringPtr("b")}, true},
		{map[string]*int{"one": intPtr(1), "two": intPtr(1)}, false},
		{map[string]*int{"one": intPtr(1), "two": intPtr(2)}, true},
		{map[string]*float64{"one": float64Ptr(1.1), "two": float64Ptr(1.1)}, false},
		{map[string]*float64{"one": float64Ptr(1.1), "two": float64Ptr(1.2)}, true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "unique")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d unique failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d unique failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "unique" {
					t.Fatalf("Index: %d unique failed Error: %v", i, errs)
				}
			}
		}
	}
	PanicMatches(t, func() { _ = validate.Var(1.0, "unique") }, "Bad field type float64")

	t.Run("struct", func(t *testing.T) {
		tests := []struct {
			param    interface{}
			expected bool
		}{
			{struct {
				A string `validate:"unique=B"`
				B string
			}{A: "abc", B: "bcd"}, true},
			{struct {
				A string `validate:"unique=B"`
				B string
			}{A: "abc", B: "abc"}, false},
		}
		validate := New()

		for i, test := range tests {
			errs := validate.Struct(test.param)
			if test.expected {
				if !IsEqual(errs, nil) {
					t.Fatalf("Index: %d unique failed Error: %v", i, errs)
				}
			} else {
				if IsEqual(errs, nil) {
					t.Fatalf("Index: %d unique failed Error: %v", i, errs)
				} else {
					val := getError(errs, "A", "A")
					if val.Tag() != "unique" {
						t.Fatalf("Index: %d unique failed Error: %v", i, errs)
					}
				}
			}
		}
	})
}

func TestUniqueValidationStructSlice(t *testing.T) {
	testStructs := []struct {
		A string
		B string
	}{
		{A: "one", B: "two"},
		{A: "one", B: "three"},
	}

	tests := []struct {
		target   interface{}
		param    string
		expected bool
	}{
		{testStructs, "unique", true},
		{testStructs, "unique=A", false},
		{testStructs, "unique=B", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.target, test.param)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d unique failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d unique failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "unique" {
					t.Fatalf("Index: %d unique failed Error: %v", i, errs)
				}
			}
		}
	}
	PanicMatches(t, func() { _ = validate.Var(testStructs, "unique=C") }, "Bad field name C")
}

func TestUniqueValidationStructPtrSlice(t *testing.T) {
	testStructs := []*struct {
		A *string
		B *string
	}{
		{A: stringPtr("one"), B: stringPtr("two")},
		{A: stringPtr("one"), B: stringPtr("three")},
		{},
	}

	tests := []struct {
		target   interface{}
		param    string
		expected bool
	}{
		{testStructs, "unique", true},
		{testStructs, "unique=A", false},
		{testStructs, "unique=B", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.target, test.param)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d unique failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d unique failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "unique" {
					t.Fatalf("Index: %d unique failed Error: %v", i, errs)
				}
			}
		}
	}
	PanicMatches(t, func() { _ = validate.Var(testStructs, "unique=C") }, "Bad field name C")
}

func TestHTMLValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"<html>", true},
		{"<script>", true},
		{"<stillworks>", true},
		{"</html", false},
		{"</script>", true},
		{"<//script>", false},
		{"<123nonsense>", false},
		{"test", false},
		{"&example", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "html")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d html failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d html failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "html" {
					t.Fatalf("Index: %d html failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestHTMLEncodedValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"&#x3c;", true},
		{"&#xaf;", true},
		{"&#x00;", true},
		{"&#xf0;", true},
		{"&#x3c", true},
		{"&#xaf", true},
		{"&#x00", true},
		{"&#xf0", true},
		{"&#ab", true},
		{"&lt;", true},
		{"&gt;", true},
		{"&quot;", true},
		{"&amp;", true},
		{"#x0a", false},
		{"&x00", false},
		{"&#x1z", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "html_encoded")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d html_encoded failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d html_encoded failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "html_encoded" {
					t.Fatalf("Index: %d html_encoded failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestURLEncodedValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"%20", true},
		{"%af", true},
		{"%ff", true},
		{"<%az", false},
		{"%test%", false},
		{"a%b", false},
		{"1%2", false},
		{"%%a%%", false},
		{"hello", true},
		{"", true},
		{"+", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "url_encoded")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d url_encoded failed Error: %v", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d url_encoded failed Error: %v", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "url_encoded" {
					t.Fatalf("Index: %d url_encoded failed Error: %v", i, errs)
				}
			}
		}
	}
}

func TestKeys(t *testing.T) {
	type Test struct {
		Test1 map[string]string `validate:"gt=0,dive,keys,eq=testkey,endkeys,eq=testval" json:"test1"`
		Test2 map[int]int       `validate:"gt=0,dive,keys,eq=3,endkeys,eq=4"             json:"test2"`
		Test3 map[int]int       `validate:"gt=0,dive,keys,eq=3,endkeys"                  json:"test3"`
	}

	var tst Test

	validate := New()
	err := validate.Struct(tst)
	NotEqual(t, err, nil)
	Equal(t, len(err.(ValidationErrors)), 3)
	AssertError(t, err.(ValidationErrors), "Test.Test1", "Test.Test1", "Test1", "Test1", "gt")
	AssertError(t, err.(ValidationErrors), "Test.Test2", "Test.Test2", "Test2", "Test2", "gt")
	AssertError(t, err.(ValidationErrors), "Test.Test3", "Test.Test3", "Test3", "Test3", "gt")

	tst.Test1 = map[string]string{
		"testkey": "testval",
	}

	tst.Test2 = map[int]int{
		3: 4,
	}

	tst.Test3 = map[int]int{
		3: 4,
	}

	err = validate.Struct(tst)
	Equal(t, err, nil)

	tst.Test1["badtestkey"] = "badtestvalue"
	tst.Test2[10] = 11

	err = validate.Struct(tst)
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)

	Equal(t, len(errs), 4)

	AssertDeepError(t, errs, "Test.Test1[badtestkey]", "Test.Test1[badtestkey]", "Test1[badtestkey]", "Test1[badtestkey]", "eq", "eq")
	AssertDeepError(t, errs, "Test.Test1[badtestkey]", "Test.Test1[badtestkey]", "Test1[badtestkey]", "Test1[badtestkey]", "eq", "eq")
	AssertDeepError(t, errs, "Test.Test2[10]", "Test.Test2[10]", "Test2[10]", "Test2[10]", "eq", "eq")
	AssertDeepError(t, errs, "Test.Test2[10]", "Test.Test2[10]", "Test2[10]", "Test2[10]", "eq", "eq")

	type Test2 struct {
		NestedKeys map[[1]string]string `validate:"gt=0,dive,keys,dive,eq=innertestkey,endkeys,eq=outertestval"`
	}

	var tst2 Test2

	err = validate.Struct(tst2)
	NotEqual(t, err, nil)
	Equal(t, len(err.(ValidationErrors)), 1)
	AssertError(t, err.(ValidationErrors), "Test2.NestedKeys", "Test2.NestedKeys", "NestedKeys", "NestedKeys", "gt")

	tst2.NestedKeys = map[[1]string]string{
		{"innertestkey"}: "outertestval",
	}

	err = validate.Struct(tst2)
	Equal(t, err, nil)

	tst2.NestedKeys[[1]string{"badtestkey"}] = "badtestvalue"

	err = validate.Struct(tst2)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)

	Equal(t, len(errs), 2)
	AssertDeepError(t, errs, "Test2.NestedKeys[[badtestkey]][0]", "Test2.NestedKeys[[badtestkey]][0]", "NestedKeys[[badtestkey]][0]", "NestedKeys[[badtestkey]][0]", "eq", "eq")
	AssertDeepError(t, errs, "Test2.NestedKeys[[badtestkey]]", "Test2.NestedKeys[[badtestkey]]", "NestedKeys[[badtestkey]]", "NestedKeys[[badtestkey]]", "eq", "eq")

	// test bad tag definitions

	PanicMatches(t, func() { _ = validate.Var(map[string]string{"key": "val"}, "endkeys,dive,eq=val") }, "'endkeys' tag encountered without a corresponding 'keys' tag")
	PanicMatches(t, func() { _ = validate.Var(1, "keys,eq=1,endkeys") }, "'keys' tag must be immediately preceded by the 'dive' tag")

	// test custom tag name
	validate = New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})

	err = validate.Struct(tst)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)

	Equal(t, len(errs), 4)

	AssertDeepError(t, errs, "Test.test1[badtestkey]", "Test.Test1[badtestkey]", "test1[badtestkey]", "Test1[badtestkey]", "eq", "eq")
	AssertDeepError(t, errs, "Test.test1[badtestkey]", "Test.Test1[badtestkey]", "test1[badtestkey]", "Test1[badtestkey]", "eq", "eq")
	AssertDeepError(t, errs, "Test.test2[10]", "Test.Test2[10]", "test2[10]", "Test2[10]", "eq", "eq")
	AssertDeepError(t, errs, "Test.test2[10]", "Test.Test2[10]", "test2[10]", "Test2[10]", "eq", "eq")
}

// Thanks @adrian-sgn specific test for your specific scenario
func TestKeysCustomValidation(t *testing.T) {
	type LangCode string
	type Label map[LangCode]string

	type TestMapStructPtr struct {
		Label Label `validate:"dive,keys,lang_code,endkeys,required"`
	}

	validate := New()
	err := validate.RegisterValidation("lang_code", func(fl FieldLevel) bool {
		validLangCodes := map[LangCode]struct{}{
			"en": {},
			"es": {},
			"pt": {},
		}

		_, ok := validLangCodes[fl.Field().Interface().(LangCode)]
		return ok
	})
	Equal(t, err, nil)

	label := Label{
		"en":  "Good morning!",
		"pt":  "",
		"es":  "¡Buenos días!",
		"xx":  "Bad key",
		"xxx": "",
	}

	err = validate.Struct(TestMapStructPtr{label})
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)
	Equal(t, len(errs), 4)

	AssertDeepError(t, errs, "TestMapStructPtr.Label[xx]", "TestMapStructPtr.Label[xx]", "Label[xx]", "Label[xx]", "lang_code", "lang_code")
	AssertDeepError(t, errs, "TestMapStructPtr.Label[pt]", "TestMapStructPtr.Label[pt]", "Label[pt]", "Label[pt]", "required", "required")
	AssertDeepError(t, errs, "TestMapStructPtr.Label[xxx]", "TestMapStructPtr.Label[xxx]", "Label[xxx]", "Label[xxx]", "lang_code", "lang_code")
	AssertDeepError(t, errs, "TestMapStructPtr.Label[xxx]", "TestMapStructPtr.Label[xxx]", "Label[xxx]", "Label[xxx]", "required", "required")

	// find specific error

	var e FieldError
	for _, e = range errs {
		if e.Namespace() == "TestMapStructPtr.Label[xxx]" {
			break
		}
	}

	Equal(t, e.Param(), "")
	Equal(t, e.Value().(LangCode), LangCode("xxx"))

	for _, e = range errs {
		if e.Namespace() == "TestMapStructPtr.Label[xxx]" && e.Tag() == "required" {
			break
		}
	}

	Equal(t, e.Param(), "")
	Equal(t, e.Value().(string), "")
}

func TestKeyOrs(t *testing.T) {
	type Test struct {
		Test1 map[string]string `validate:"gt=0,dive,keys,eq=testkey|eq=testkeyok,endkeys,eq=testval" json:"test1"`
	}

	var tst Test

	validate := New()
	err := validate.Struct(tst)
	NotEqual(t, err, nil)
	Equal(t, len(err.(ValidationErrors)), 1)
	AssertError(t, err.(ValidationErrors), "Test.Test1", "Test.Test1", "Test1", "Test1", "gt")

	tst.Test1 = map[string]string{
		"testkey": "testval",
	}

	err = validate.Struct(tst)
	Equal(t, err, nil)

	tst.Test1["badtestkey"] = "badtestval"

	err = validate.Struct(tst)
	NotEqual(t, err, nil)

	errs := err.(ValidationErrors)

	Equal(t, len(errs), 2)

	AssertDeepError(t, errs, "Test.Test1[badtestkey]", "Test.Test1[badtestkey]", "Test1[badtestkey]", "Test1[badtestkey]", "eq=testkey|eq=testkeyok", "eq=testkey|eq=testkeyok")
	AssertDeepError(t, errs, "Test.Test1[badtestkey]", "Test.Test1[badtestkey]", "Test1[badtestkey]", "Test1[badtestkey]", "eq", "eq")

	validate.RegisterAlias("okkey", "eq=testkey|eq=testkeyok")

	type Test2 struct {
		Test1 map[string]string `validate:"gt=0,dive,keys,okkey,endkeys,eq=testval" json:"test1"`
	}

	var tst2 Test2

	err = validate.Struct(tst2)
	NotEqual(t, err, nil)
	Equal(t, len(err.(ValidationErrors)), 1)
	AssertError(t, err.(ValidationErrors), "Test2.Test1", "Test2.Test1", "Test1", "Test1", "gt")

	tst2.Test1 = map[string]string{
		"testkey": "testval",
	}

	err = validate.Struct(tst2)
	Equal(t, err, nil)

	tst2.Test1["badtestkey"] = "badtestval"

	err = validate.Struct(tst2)
	NotEqual(t, err, nil)

	errs = err.(ValidationErrors)

	Equal(t, len(errs), 2)

	AssertDeepError(t, errs, "Test2.Test1[badtestkey]", "Test2.Test1[badtestkey]", "Test1[badtestkey]", "Test1[badtestkey]", "okkey", "eq=testkey|eq=testkeyok")
	AssertDeepError(t, errs, "Test2.Test1[badtestkey]", "Test2.Test1[badtestkey]", "Test1[badtestkey]", "Test1[badtestkey]", "eq", "eq")
}

func TestStructLevelValidationsPointerPassing(t *testing.T) {
	v1 := New()
	v1.RegisterStructValidation(StructValidationTestStruct, &TestStruct{})

	tst := &TestStruct{
		String: "good value",
	}

	errs := v1.Struct(tst)
	NotEqual(t, errs, nil)
	AssertError(t, errs, "TestStruct.StringVal", "TestStruct.String", "StringVal", "String", "badvalueteststruct")
}

func TestDirValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		title    string
		param    string
		expected bool
	}{
		{"existing dir", "testdata", true},
		{"existing self dir", ".", true},
		{"existing parent dir", "..", true},
		{"empty dir", "", false},
		{"missing dir", "non_existing_testdata", false},
		{"a file not a directory", filepath.Join("testdata", "a.go"), false},
	}

	for _, test := range tests {
		errs := validate.Var(test.param, "dir")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "dir")
	}, "Bad field type int")
}

func TestDirPathValidation(t *testing.T) {
	validate := New()

	tests := []struct {
		title    string
		param    string
		expected bool
	}{
		{"empty dirpath", "", false},
		{"valid dirpath - exists", "testdata", true},
		{"valid dirpath - explicit", "testdatanoexist" + string(os.PathSeparator), true},
		{"invalid dirpath", "testdata\000" + string(os.PathSeparator), false},
		{"file, not a dirpath", filepath.Join("testdata", "a.go"), false},
	}

	for _, test := range tests {
		errs := validate.Var(test.param, "dirpath")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Test: '%s' failed Error: %s", test.title, errs)
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(6, "filepath")
	}, "Bad field type int")
}

func TestStartsWithValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"startswith=(/^ヮ^)/*:・ﾟ✧"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "(/^ヮ^)/*:・ﾟ✧ glitter", Tag: "startswith=(/^ヮ^)/*:・ﾟ✧", ExpectedNil: true},
		{Value: "abcd", Tag: "startswith=(/^ヮ^)/*:・ﾟ✧", ExpectedNil: false},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestEndsWithValidation(t *testing.T) {
	tests := []struct {
		Value       string `validate:"endswith=(/^ヮ^)/*:・ﾟ✧"`
		Tag         string
		ExpectedNil bool
	}{
		{Value: "glitter (/^ヮ^)/*:・ﾟ✧", Tag: "endswith=(/^ヮ^)/*:・ﾟ✧", ExpectedNil: true},
		{Value: "(/^ヮ^)/*:・ﾟ✧ glitter", Tag: "endswith=(/^ヮ^)/*:・ﾟ✧", ExpectedNil: false},
	}

	validate := New()

	for i, s := range tests {
		errs := validate.Var(s.Value, s.Tag)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}

		errs = validate.Struct(s)

		if (s.ExpectedNil && errs != nil) || (!s.ExpectedNil && errs == nil) {
			t.Fatalf("Index: %d failed Error: %s", i, errs)
		}
	}
}

func TestRequiredIf(t *testing.T) {
	type Inner struct {
		Field *string
	}

	fieldVal := "test"
	test := struct {
		Inner   *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_if=FieldE test" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_if=Field1 test" json:"field_2"`
		Field3  map[string]string `validate:"required_if=Field2 test" json:"field_3"`
		Field4  interface{}       `validate:"required_if=Field3 1" json:"field_4"`
		Field5  int               `validate:"required_if=Inner.Field test" json:"field_5"`
		Field6  uint              `validate:"required_if=Field5 1" json:"field_6"`
		Field7  float32           `validate:"required_if=Field6 1" json:"field_7"`
		Field8  float64           `validate:"required_if=Field7 1.0" json:"field_8"`
		Field9  Inner             `validate:"required_if=Field1 test" json:"field_9"`
		Field10 *Inner            `validate:"required_if=Field1 test" json:"field_10"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: 2,
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner   *Inner
		Inner2  *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_if=FieldE test" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_if=Field1 test" json:"field_2"`
		Field3  map[string]string `validate:"required_if=Field2 test" json:"field_3"`
		Field4  interface{}       `validate:"required_if=Field2 test" json:"field_4"`
		Field5  string            `validate:"required_if=Field3 1" json:"field_5"`
		Field6  string            `validate:"required_if=Inner.Field test" json:"field_6"`
		Field7  string            `validate:"required_if=Inner2.Field test" json:"field_7"`
		Field8  Inner             `validate:"required_if=Field2 test" json:"field_8"`
		Field9  *Inner            `validate:"required_if=Field2 test" json:"field_9"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field2: &fieldVal,
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 5)
	AssertError(t, errs, "Field3", "Field3", "Field3", "Field3", "required_if")
	AssertError(t, errs, "Field4", "Field4", "Field4", "Field4", "required_if")
	AssertError(t, errs, "Field6", "Field6", "Field6", "Field6", "required_if")
	AssertError(t, errs, "Field8", "Field8", "Field8", "Field8", "required_if")
	AssertError(t, errs, "Field9", "Field9", "Field9", "Field9", "required_if")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test3 should have panicked!")
		}
	}()

	test3 := struct {
		Inner  *Inner
		Field1 string `validate:"required_if=Inner.Field" json:"field_1"`
	}{
		Inner: &Inner{Field: &fieldVal},
	}
	_ = validate.Struct(test3)
}

func TestRequiredUnless(t *testing.T) {
	type Inner struct {
		Field *string
	}

	fieldVal := "test"
	test := struct {
		Inner   *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_unless=FieldE test" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_unless=Field1 test" json:"field_2"`
		Field3  map[string]string `validate:"required_unless=Field2 test" json:"field_3"`
		Field4  interface{}       `validate:"required_unless=Field3 1" json:"field_4"`
		Field5  int               `validate:"required_unless=Inner.Field test" json:"field_5"`
		Field6  uint              `validate:"required_unless=Field5 2" json:"field_6"`
		Field7  float32           `validate:"required_unless=Field6 0" json:"field_7"`
		Field8  float64           `validate:"required_unless=Field7 0.0" json:"field_8"`
		Field9  bool              `validate:"omitempty" json:"field_9"`
		Field10 string            `validate:"required_unless=Field9 true" json:"field_10"`
		Field11 Inner             `validate:"required_unless=Field9 true" json:"field_11"`
		Field12 *Inner            `validate:"required_unless=Field9 true" json:"field_12"`
	}{
		FieldE: "test",
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: 2,
		Field9: true,
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner   *Inner
		Inner2  *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_unless=FieldE test" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_unless=Field1 test" json:"field_2"`
		Field3  map[string]string `validate:"required_unless=Field2 test" json:"field_3"`
		Field4  interface{}       `validate:"required_unless=Field2 test" json:"field_4"`
		Field5  string            `validate:"required_unless=Field3 0" json:"field_5"`
		Field6  string            `validate:"required_unless=Inner.Field test" json:"field_6"`
		Field7  string            `validate:"required_unless=Inner2.Field test" json:"field_7"`
		Field8  bool              `validate:"omitempty" json:"field_8"`
		Field9  string            `validate:"required_unless=Field8 true" json:"field_9"`
		Field10 Inner             `validate:"required_unless=Field9 true" json:"field_10"`
		Field11 *Inner            `validate:"required_unless=Field9 true" json:"field_11"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		FieldE: "test",
		Field1: "test",
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 6)
	AssertError(t, errs, "Field3", "Field3", "Field3", "Field3", "required_unless")
	AssertError(t, errs, "Field4", "Field4", "Field4", "Field4", "required_unless")
	AssertError(t, errs, "Field7", "Field7", "Field7", "Field7", "required_unless")
	AssertError(t, errs, "Field9", "Field9", "Field9", "Field9", "required_unless")
	AssertError(t, errs, "Field10", "Field10", "Field10", "Field10", "required_unless")
	AssertError(t, errs, "Field11", "Field11", "Field11", "Field11", "required_unless")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("test3 should have panicked!")
		}
	}()

	test3 := struct {
		Inner  *Inner
		Field1 string `validate:"required_unless=Inner.Field" json:"field_1"`
	}{
		Inner: &Inner{Field: &fieldVal},
	}
	_ = validate.Struct(test3)
}

func TestSkipUnless(t *testing.T) {
	type Inner struct {
		Field *string
	}

	fieldVal := "test1"
	test := struct {
		Inner   *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"skip_unless=FieldE test" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"skip_unless=Field1 test" json:"field_2"`
		Field3  map[string]string `validate:"skip_unless=Field2 test" json:"field_3"`
		Field4  interface{}       `validate:"skip_unless=Field3 1" json:"field_4"`
		Field5  int               `validate:"skip_unless=Inner.Field test" json:"field_5"`
		Field6  uint              `validate:"skip_unless=Field5 2" json:"field_6"`
		Field7  float32           `validate:"skip_unless=Field6 1" json:"field_7"`
		Field8  float64           `validate:"skip_unless=Field7 1.0" json:"field_8"`
		Field9  bool              `validate:"omitempty" json:"field_9"`
		Field10 string            `validate:"skip_unless=Field9 false" json:"field_10"`
		Field11 Inner             `validate:"skip_unless=Field9 false" json:"field_11"`
		Field12 *Inner            `validate:"skip_unless=Field9 false" json:"field_12"`
	}{
		FieldE: "test1",
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: 3,
		Field9: true,
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner   *Inner
		Inner2  *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"skip_unless=FieldE test" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"skip_unless=Field1 test" json:"field_2"`
		Field3  map[string]string `validate:"skip_unless=Field2 test" json:"field_3"`
		Field4  interface{}       `validate:"skip_unless=Field2 test" json:"field_4"`
		Field5  string            `validate:"skip_unless=Field3 0" json:"field_5"`
		Field6  string            `validate:"skip_unless=Inner.Field test" json:"field_6"`
		Field7  string            `validate:"skip_unless=Inner2.Field test" json:"field_7"`
		Field8  bool              `validate:"omitempty" json:"field_8"`
		Field9  string            `validate:"skip_unless=Field8 true" json:"field_9"`
		Field10 Inner             `validate:"skip_unless=Field8 false" json:"field_10"`
		Field11 *Inner            `validate:"skip_unless=Field8 false" json:"field_11"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		FieldE: "test1",
		Field1: "test1",
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 3)
	AssertError(t, errs, "Field5", "Field5", "Field5", "Field5", "skip_unless")
	AssertError(t, errs, "Field10", "Field10", "Field10", "Field10", "skip_unless")
	AssertError(t, errs, "Field11", "Field11", "Field11", "Field11", "skip_unless")

	test3 := struct {
		Inner  *Inner
		Field1 string `validate:"skip_unless=Inner.Field" json:"field_1"`
	}{
		Inner: &Inner{Field: &fieldVal},
	}
	PanicMatches(t, func() {
		_ = validate.Struct(test3)
	}, "Bad param number for skip_unless Field1")

	test4 := struct {
		Inner  *Inner
		Field1 string `validate:"skip_unless=Inner.Field test1" json:"field_1"`
	}{
		Inner: &Inner{Field: &fieldVal},
	}
	errs = validate.Struct(test4)
	NotEqual(t, errs, nil)

	ve = errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "Field1", "Field1", "Field1", "Field1", "skip_unless")
}

func TestRequiredWith(t *testing.T) {
	type Inner struct {
		Field *string
	}

	fieldVal := "test"
	test := struct {
		Inner   *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_with=FieldE" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_with=Field1" json:"field_2"`
		Field3  map[string]string `validate:"required_with=Field2" json:"field_3"`
		Field4  interface{}       `validate:"required_with=Field3" json:"field_4"`
		Field5  string            `validate:"required_with=Field" json:"field_5"`
		Field6  Inner             `validate:"required_with=Field2" json:"field_6"`
		Field7  *Inner            `validate:"required_with=Field2" json:"field_7"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: Inner{Field: &fieldVal},
		Field7: &Inner{Field: &fieldVal},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner   *Inner
		Inner2  *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_with=FieldE" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_with=Field1" json:"field_2"`
		Field3  map[string]string `validate:"required_with=Field2" json:"field_3"`
		Field4  interface{}       `validate:"required_with=Field2" json:"field_4"`
		Field5  string            `validate:"required_with=Field3" json:"field_5"`
		Field6  string            `validate:"required_with=Inner.Field" json:"field_6"`
		Field7  string            `validate:"required_with=Inner2.Field" json:"field_7"`
		Field8  Inner             `validate:"required_with=Field2" json:"field_8"`
		Field9  *Inner            `validate:"required_with=Field2" json:"field_9"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field2: &fieldVal,
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 5)
	AssertError(t, errs, "Field3", "Field3", "Field3", "Field3", "required_with")
	AssertError(t, errs, "Field4", "Field4", "Field4", "Field4", "required_with")
	AssertError(t, errs, "Field6", "Field6", "Field6", "Field6", "required_with")
	AssertError(t, errs, "Field8", "Field8", "Field8", "Field8", "required_with")
	AssertError(t, errs, "Field9", "Field9", "Field9", "Field9", "required_with")
}

func TestExcludedWith(t *testing.T) {
	type Inner struct {
		FieldE string
		Field  *string
	}

	fieldVal := "test"
	test := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_with=FieldE" json:"field_1"`
		Field2 *string           `validate:"excluded_with=FieldE" json:"field_2"`
		Field3 map[string]string `validate:"excluded_with=FieldE" json:"field_3"`
		Field4 interface{}       `validate:"excluded_with=FieldE" json:"field_4"`
		Field5 string            `validate:"excluded_with=Inner.FieldE" json:"field_5"`
		Field6 string            `validate:"excluded_with=Inner2.FieldE" json:"field_6"`
		Field7 Inner             `validate:"excluded_with=FieldE" json:"field_7"`
		Field8 *Inner            `validate:"excluded_with=FieldE" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: "test",
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_with=Field" json:"field_1"`
		Field2 *string           `validate:"excluded_with=Field" json:"field_2"`
		Field3 map[string]string `validate:"excluded_with=Field" json:"field_3"`
		Field4 interface{}       `validate:"excluded_with=Field" json:"field_4"`
		Field5 string            `validate:"excluded_with=Inner.Field" json:"field_5"`
		Field6 string            `validate:"excluded_with=Inner2.Field" json:"field_6"`
		Field7 Inner             `validate:"excluded_with=Field" json:"field_7"`
		Field8 *Inner            `validate:"excluded_with=Field" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field:  "populated",
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: "test",
		Field7: Inner{FieldE: "potato"},
		Field8: &Inner{FieldE: "potato"},
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 7)
	for i := 1; i <= 7; i++ {
		// accounting for field 7 & 8 failures, 6 skipped because no failure
		if i > 5 {
			i++
		}
		name := fmt.Sprintf("Field%d", i)
		AssertError(t, errs, name, name, name, name, "excluded_with")
	}

	test3 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_with=FieldE" json:"field_1"`
		Field2 *string           `validate:"excluded_with=FieldE" json:"field_2"`
		Field3 map[string]string `validate:"excluded_with=FieldE" json:"field_3"`
		Field4 interface{}       `validate:"excluded_with=FieldE" json:"field_4"`
		Field5 string            `validate:"excluded_with=Inner.FieldE" json:"field_5"`
		Field6 string            `validate:"excluded_with=Inner2.FieldE" json:"field_6"`
		Field7 Inner             `validate:"excluded_with=FieldE" json:"field_7"`
		Field8 *Inner            `validate:"excluded_with=FieldE" json:"field_8"`
	}{
		Inner:  &Inner{FieldE: "populated"},
		Inner2: &Inner{FieldE: "populated"},
		FieldE: "populated",
	}

	validate = New()

	errs = validate.Struct(test3)
	Equal(t, errs, nil)
}

func TestExcludedWithout(t *testing.T) {
	type Inner struct {
		FieldE string
		Field  *string
	}

	fieldVal := "test"
	test := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_without=Field" json:"field_1"`
		Field2 *string           `validate:"excluded_without=Field" json:"field_2"`
		Field3 map[string]string `validate:"excluded_without=Field" json:"field_3"`
		Field4 interface{}       `validate:"excluded_without=Field" json:"field_4"`
		Field5 string            `validate:"excluded_without=Inner.Field" json:"field_5"`
		Field6 Inner             `validate:"excluded_without=Field" json:"field_6"`
		Field7 *Inner            `validate:"excluded_without=Field" json:"field_7"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field:  "populated",
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: Inner{FieldE: "potato"},
		Field7: &Inner{FieldE: "potato"},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_without=FieldE" json:"field_1"`
		Field2 *string           `validate:"excluded_without=FieldE" json:"field_2"`
		Field3 map[string]string `validate:"excluded_without=FieldE" json:"field_3"`
		Field4 interface{}       `validate:"excluded_without=FieldE" json:"field_4"`
		Field5 string            `validate:"excluded_without=Inner.FieldE" json:"field_5"`
		Field6 string            `validate:"excluded_without=Inner2.FieldE" json:"field_6"`
		Field7 Inner             `validate:"excluded_without=FieldE" json:"field_7"`
		Field8 *Inner            `validate:"excluded_without=FieldE" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: "test",
		Field7: Inner{FieldE: "potato"},
		Field8: &Inner{FieldE: "potato"},
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 8)
	for i := 1; i <= 8; i++ {
		name := fmt.Sprintf("Field%d", i)
		AssertError(t, errs, name, name, name, name, "excluded_without")
	}

	test3 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_without=Field" json:"field_1"`
		Field2 *string           `validate:"excluded_without=Field" json:"field_2"`
		Field3 map[string]string `validate:"excluded_without=Field" json:"field_3"`
		Field4 interface{}       `validate:"excluded_without=Field" json:"field_4"`
		Field5 string            `validate:"excluded_without=Inner.Field" json:"field_5"`
		Field6 Inner             `validate:"excluded_without=Field" json:"field_6"`
		Field7 *Inner            `validate:"excluded_without=Field" json:"field_7"`
	}{
		Inner: &Inner{Field: &fieldVal},
		Field: "populated",
	}

	validate = New()

	errs = validate.Struct(test3)
	Equal(t, errs, nil)
}

func TestExcludedWithAll(t *testing.T) {
	type Inner struct {
		FieldE string
		Field  *string
	}

	fieldVal := "test"
	test := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_with_all=FieldE Field" json:"field_1"`
		Field2 *string           `validate:"excluded_with_all=FieldE Field" json:"field_2"`
		Field3 map[string]string `validate:"excluded_with_all=FieldE Field" json:"field_3"`
		Field4 interface{}       `validate:"excluded_with_all=FieldE Field" json:"field_4"`
		Field5 string            `validate:"excluded_with_all=Inner.FieldE" json:"field_5"`
		Field6 string            `validate:"excluded_with_all=Inner2.FieldE" json:"field_6"`
		Field7 Inner             `validate:"excluded_with_all=FieldE Field" json:"field_7"`
		Field8 *Inner            `validate:"excluded_with_all=FieldE Field" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field:  fieldVal,
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: "test",
		Field7: Inner{FieldE: "potato"},
		Field8: &Inner{FieldE: "potato"},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_with_all=Field FieldE" json:"field_1"`
		Field2 *string           `validate:"excluded_with_all=Field FieldE" json:"field_2"`
		Field3 map[string]string `validate:"excluded_with_all=Field FieldE" json:"field_3"`
		Field4 interface{}       `validate:"excluded_with_all=Field FieldE" json:"field_4"`
		Field5 string            `validate:"excluded_with_all=Inner.Field" json:"field_5"`
		Field6 string            `validate:"excluded_with_all=Inner2.Field" json:"field_6"`
		Field7 Inner             `validate:"excluded_with_all=Field FieldE" json:"field_7"`
		Field8 *Inner            `validate:"excluded_with_all=Field FieldE" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field:  "populated",
		FieldE: "populated",
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: "test",
		Field7: Inner{FieldE: "potato"},
		Field8: &Inner{FieldE: "potato"},
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 7)
	for i := 1; i <= 7; i++ {
		// accounting for no err for field 6
		if i > 5 {
			i++
		}
		name := fmt.Sprintf("Field%d", i)
		AssertError(t, errs, name, name, name, name, "excluded_with_all")
	}

	test3 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_with_all=FieldE Field" json:"field_1"`
		Field2 *string           `validate:"excluded_with_all=FieldE Field" json:"field_2"`
		Field3 map[string]string `validate:"excluded_with_all=FieldE Field" json:"field_3"`
		Field4 interface{}       `validate:"excluded_with_all=FieldE Field" json:"field_4"`
		Field5 string            `validate:"excluded_with_all=Inner.FieldE" json:"field_5"`
		Field6 string            `validate:"excluded_with_all=Inner2.FieldE" json:"field_6"`
		Field7 Inner             `validate:"excluded_with_all=Field FieldE" json:"field_7"`
		Field8 *Inner            `validate:"excluded_with_all=Field FieldE" json:"field_8"`
	}{
		Inner:  &Inner{FieldE: "populated"},
		Inner2: &Inner{FieldE: "populated"},
		Field:  "populated",
		FieldE: "populated",
	}

	validate = New()

	errs = validate.Struct(test3)
	Equal(t, errs, nil)
}

func TestExcludedWithoutAll(t *testing.T) {
	type Inner struct {
		FieldE string
		Field  *string
	}

	fieldVal := "test"
	test := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_without_all=Field FieldE" json:"field_1"`
		Field2 *string           `validate:"excluded_without_all=Field FieldE" json:"field_2"`
		Field3 map[string]string `validate:"excluded_without_all=Field FieldE" json:"field_3"`
		Field4 interface{}       `validate:"excluded_without_all=Field FieldE" json:"field_4"`
		Field5 string            `validate:"excluded_without_all=Inner.Field Inner2.Field" json:"field_5"`
		Field6 Inner             `validate:"excluded_without_all=Field FieldE" json:"field_6"`
		Field7 *Inner            `validate:"excluded_without_all=Field FieldE" json:"field_7"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Inner2: &Inner{Field: &fieldVal},
		Field:  "populated",
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: Inner{FieldE: "potato"},
		Field7: &Inner{FieldE: "potato"},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_without_all=FieldE Field" json:"field_1"`
		Field2 *string           `validate:"excluded_without_all=FieldE Field" json:"field_2"`
		Field3 map[string]string `validate:"excluded_without_all=FieldE Field" json:"field_3"`
		Field4 interface{}       `validate:"excluded_without_all=FieldE Field" json:"field_4"`
		Field5 string            `validate:"excluded_without_all=Inner.FieldE" json:"field_5"`
		Field6 string            `validate:"excluded_without_all=Inner2.FieldE" json:"field_6"`
		Field7 Inner             `validate:"excluded_without_all=Field FieldE" json:"field_7"`
		Field8 *Inner            `validate:"excluded_without_all=Field FieldE" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field1: fieldVal,
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: "test",
		Field7: Inner{FieldE: "potato"},
		Field8: &Inner{FieldE: "potato"},
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 8)
	for i := 1; i <= 8; i++ {
		name := fmt.Sprintf("Field%d", i)
		AssertError(t, errs, name, name, name, name, "excluded_without_all")
	}

	test3 := struct {
		Inner  *Inner
		Inner2 *Inner
		Field  string            `validate:"omitempty" json:"field"`
		FieldE string            `validate:"omitempty" json:"field_e"`
		Field1 string            `validate:"excluded_without_all=Field FieldE" json:"field_1"`
		Field2 *string           `validate:"excluded_without_all=Field FieldE" json:"field_2"`
		Field3 map[string]string `validate:"excluded_without_all=Field FieldE" json:"field_3"`
		Field4 interface{}       `validate:"excluded_without_all=Field FieldE" json:"field_4"`
		Field5 string            `validate:"excluded_without_all=Inner.Field Inner2.Field" json:"field_5"`
		Field6 Inner             `validate:"excluded_without_all=Field FieldE" json:"field_6"`
		Field7 *Inner            `validate:"excluded_without_all=Field FieldE" json:"field_7"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Inner2: &Inner{Field: &fieldVal},
		Field:  "populated",
		FieldE: "populated",
	}

	validate = New()

	errs = validate.Struct(test3)
	Equal(t, errs, nil)
}

func TestRequiredWithAll(t *testing.T) {
	type Inner struct {
		Field *string
	}

	fieldVal := "test"
	test := struct {
		Inner   *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_with_all=FieldE" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_with_all=Field1" json:"field_2"`
		Field3  map[string]string `validate:"required_with_all=Field2" json:"field_3"`
		Field4  interface{}       `validate:"required_with_all=Field3" json:"field_4"`
		Field5  string            `validate:"required_with_all=Inner.Field" json:"field_5"`
		Field6  Inner             `validate:"required_with_all=Field1 Field2" json:"field_6"`
		Field7  *Inner            `validate:"required_with_all=Field1 Field2" json:"field_7"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field1: "test_field1",
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: Inner{Field: &fieldVal},
		Field7: &Inner{Field: &fieldVal},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner   *Inner
		Inner2  *Inner
		FieldE  string            `validate:"omitempty" json:"field_e"`
		FieldER string            `validate:"required_with_all=FieldE" json:"field_er"`
		Field1  string            `validate:"omitempty" json:"field_1"`
		Field2  *string           `validate:"required_with_all=Field1" json:"field_2"`
		Field3  map[string]string `validate:"required_with_all=Field2" json:"field_3"`
		Field4  interface{}       `validate:"required_with_all=Field1 FieldE" json:"field_4"`
		Field5  string            `validate:"required_with_all=Inner.Field Field2" json:"field_5"`
		Field6  string            `validate:"required_with_all=Inner2.Field Field2" json:"field_6"`
		Field7  Inner             `validate:"required_with_all=Inner.Field Field2" json:"field_7"`
		Field8  *Inner            `validate:"required_with_all=Inner.Field Field2" json:"field_8"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field2: &fieldVal,
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 4)
	AssertError(t, errs, "Field3", "Field3", "Field3", "Field3", "required_with_all")
	AssertError(t, errs, "Field5", "Field5", "Field5", "Field5", "required_with_all")
	AssertError(t, errs, "Field7", "Field7", "Field7", "Field7", "required_with_all")
	AssertError(t, errs, "Field8", "Field8", "Field8", "Field8", "required_with_all")
}

func TestRequiredWithout(t *testing.T) {
	type Inner struct {
		Field *string
	}

	fieldVal := "test"
	test := struct {
		Inner  *Inner
		Field1 string            `validate:"omitempty" json:"field_1"`
		Field2 *string           `validate:"required_without=Field1" json:"field_2"`
		Field3 map[string]string `validate:"required_without=Field2" json:"field_3"`
		Field4 interface{}       `validate:"required_without=Field3" json:"field_4"`
		Field5 string            `validate:"required_without=Field3" json:"field_5"`
		Field6 Inner             `validate:"required_without=Field1" json:"field_6"`
		Field7 *Inner            `validate:"required_without=Field1" json:"field_7"`
	}{
		Inner:  &Inner{Field: &fieldVal},
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: Inner{Field: &fieldVal},
		Field7: &Inner{Field: &fieldVal},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Inner   *Inner
		Inner2  *Inner
		Field1  string            `json:"field_1"`
		Field2  *string           `validate:"required_without=Field1" json:"field_2"`
		Field3  map[string]string `validate:"required_without=Field2" json:"field_3"`
		Field4  interface{}       `validate:"required_without=Field3" json:"field_4"`
		Field5  string            `validate:"required_without=Field3" json:"field_5"`
		Field6  string            `validate:"required_without=Field1" json:"field_6"`
		Field7  string            `validate:"required_without=Inner.Field" json:"field_7"`
		Field8  string            `validate:"required_without=Inner.Field" json:"field_8"`
		Field9  Inner             `validate:"required_without=Field1" json:"field_9"`
		Field10 *Inner            `validate:"required_without=Field1" json:"field_10"`
	}{
		Inner:  &Inner{},
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
	}

	errs = validate.Struct(&test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 6)
	AssertError(t, errs, "Field2", "Field2", "Field2", "Field2", "required_without")
	AssertError(t, errs, "Field6", "Field6", "Field6", "Field6", "required_without")
	AssertError(t, errs, "Field7", "Field7", "Field7", "Field7", "required_without")
	AssertError(t, errs, "Field8", "Field8", "Field8", "Field8", "required_without")
	AssertError(t, errs, "Field9", "Field9", "Field9", "Field9", "required_without")
	AssertError(t, errs, "Field10", "Field10", "Field10", "Field10", "required_without")

	test3 := struct {
		Field1 *string `validate:"required_without=Field2,omitempty,min=1" json:"field_1"`
		Field2 *string `validate:"required_without=Field1,omitempty,min=1" json:"field_2"`
	}{
		Field1: &fieldVal,
	}

	errs = validate.Struct(&test3)
	Equal(t, errs, nil)
}

func TestRequiredWithoutAll(t *testing.T) {
	type nested struct {
		value string
	}
	fieldVal := "test"
	test := struct {
		Field1 string            `validate:"omitempty" json:"field_1"`
		Field2 *string           `validate:"required_without_all=Field1" json:"field_2"`
		Field3 map[string]string `validate:"required_without_all=Field2" json:"field_3"`
		Field4 interface{}       `validate:"required_without_all=Field3" json:"field_4"`
		Field5 string            `validate:"required_without_all=Field3" json:"field_5"`
		Field6 nested            `validate:"required_without_all=Field1" json:"field_6"`
		Field7 *nested           `validate:"required_without_all=Field1" json:"field_7"`
	}{
		Field1: "",
		Field2: &fieldVal,
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
		Field6: nested{"potato"},
		Field7: &nested{"potato"},
	}

	validate := New()

	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		Field1 string            `validate:"omitempty" json:"field_1"`
		Field2 *string           `validate:"required_without_all=Field1" json:"field_2"`
		Field3 map[string]string `validate:"required_without_all=Field2" json:"field_3"`
		Field4 interface{}       `validate:"required_without_all=Field3" json:"field_4"`
		Field5 string            `validate:"required_without_all=Field3" json:"field_5"`
		Field6 string            `validate:"required_without_all=Field1 Field3" json:"field_6"`
		Field7 nested            `validate:"required_without_all=Field1" json:"field_7"`
		Field8 *nested           `validate:"required_without_all=Field1" json:"field_8"`
	}{
		Field3: map[string]string{"key": "val"},
		Field4: "test",
		Field5: "test",
	}

	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)

	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 3)
	AssertError(t, errs, "Field2", "Field2", "Field2", "Field2", "required_without_all")
	AssertError(t, errs, "Field7", "Field7", "Field7", "Field7", "required_without_all")
	AssertError(t, errs, "Field8", "Field8", "Field8", "Field8", "required_without_all")
}

func TestExcludedIf(t *testing.T) {
	validate := New()
	type Inner struct {
		Field *string
	}

	shouldExclude := "exclude"
	shouldNotExclude := "dontExclude"

	test1 := struct {
		FieldE  string  `validate:"omitempty" json:"field_e"`
		FieldER *string `validate:"excluded_if=FieldE exclude" json:"field_er"`
	}{
		FieldE: shouldExclude,
	}
	errs := validate.Struct(test1)
	Equal(t, errs, nil)

	test2 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_if=FieldE exclude" json:"field_er"`
	}{
		FieldE:  shouldExclude,
		FieldER: "set",
	}
	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)
	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "FieldER", "FieldER", "FieldER", "FieldER", "excluded_if")

	test3 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldF  string `validate:"omitempty" json:"field_f"`
		FieldER string `validate:"excluded_if=FieldE exclude FieldF exclude" json:"field_er"`
	}{
		FieldE:  shouldExclude,
		FieldF:  shouldExclude,
		FieldER: "set",
	}
	errs = validate.Struct(test3)
	NotEqual(t, errs, nil)
	ve = errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "FieldER", "FieldER", "FieldER", "FieldER", "excluded_if")

	test4 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldF  string `validate:"omitempty" json:"field_f"`
		FieldER string `validate:"excluded_if=FieldE exclude FieldF exclude" json:"field_er"`
	}{
		FieldE:  shouldExclude,
		FieldF:  shouldNotExclude,
		FieldER: "set",
	}
	errs = validate.Struct(test4)
	Equal(t, errs, nil)

	test5 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_if=FieldE exclude" json:"field_er"`
	}{
		FieldE: shouldNotExclude,
	}
	errs = validate.Struct(test5)
	Equal(t, errs, nil)

	test6 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_if=FieldE exclude" json:"field_er"`
	}{
		FieldE:  shouldNotExclude,
		FieldER: "set",
	}
	errs = validate.Struct(test6)
	Equal(t, errs, nil)

	test7 := struct {
		Inner  *Inner
		FieldE string `validate:"omitempty" json:"field_e"`
		Field1 int    `validate:"excluded_if=Inner.Field exclude" json:"field_1"`
	}{
		Inner: &Inner{Field: &shouldExclude},
	}
	errs = validate.Struct(test7)
	Equal(t, errs, nil)

	test8 := struct {
		Inner  *Inner
		FieldE string `validate:"omitempty" json:"field_e"`
		Field1 int    `validate:"excluded_if=Inner.Field exclude" json:"field_1"`
	}{
		Inner:  &Inner{Field: &shouldExclude},
		Field1: 1,
	}
	errs = validate.Struct(test8)
	NotEqual(t, errs, nil)
	ve = errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "Field1", "Field1", "Field1", "Field1", "excluded_if")

	test9 := struct {
		Inner  *Inner
		FieldE string `validate:"omitempty" json:"field_e"`
		Field1 int    `validate:"excluded_if=Inner.Field exclude" json:"field_1"`
	}{
		Inner: &Inner{Field: &shouldNotExclude},
	}
	errs = validate.Struct(test9)
	Equal(t, errs, nil)

	test10 := struct {
		Inner  *Inner
		FieldE string `validate:"omitempty" json:"field_e"`
		Field1 int    `validate:"excluded_if=Inner.Field exclude" json:"field_1"`
	}{
		Inner:  &Inner{Field: &shouldNotExclude},
		Field1: 1,
	}
	errs = validate.Struct(test10)
	Equal(t, errs, nil)

	test11 := struct {
		Field1 bool
  		Field2 *string `validate:"excluded_if=Field1 false"`
	}{
		Field1: false,
		Field2: nil,
	}
	errs = validate.Struct(test11)
	Equal(t, errs, nil)

	test12 := struct {
		Field1 bool
		Field2 *string `validate:"excluded_if=Field1 !Field1"`
	}{
		Field1: true,
		Field2: nil,
	}
	errs = validate.Struct(test12)
	Equal(t, errs, nil)
	// Checks number of params in struct tag is correct
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panicTest should have panicked!")
		}
	}()
	fieldVal := "panicTest"
	panicTest := struct {
		Inner  *Inner
		Field1 string `validate:"excluded_if=Inner.Field" json:"field_1"`
	}{
		Inner: &Inner{Field: &fieldVal},
	}
	_ = validate.Struct(panicTest)
}

func TestExcludedUnless(t *testing.T) {
	validate := New()
	type Inner struct {
		Field *string
	}

	fieldVal := "test"
	test := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_unless=FieldE test" json:"field_er"`
	}{
		FieldE:  "test",
		FieldER: "filled",
	}
	errs := validate.Struct(test)
	Equal(t, errs, nil)

	test2 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_unless=FieldE test" json:"field_er"`
	}{
		FieldE:  "notest",
		FieldER: "filled",
	}
	errs = validate.Struct(test2)
	NotEqual(t, errs, nil)
	ve := errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "FieldER", "FieldER", "FieldER", "FieldER", "excluded_unless")

	// test5 and test6: excluded_unless has no effect if FieldER is left blank
	test5 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_unless=FieldE test" json:"field_er"`
	}{
		FieldE: "test",
	}
	errs = validate.Struct(test5)
	Equal(t, errs, nil)

	test6 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_unless=FieldE test" json:"field_er"`
	}{
		FieldE: "notest",
	}
	errs = validate.Struct(test6)
	Equal(t, errs, nil)

	shouldError := "notest"
	test3 := struct {
		Inner  *Inner
		Field1 string `validate:"excluded_unless=Inner.Field test" json:"field_1"`
	}{
		Inner:  &Inner{Field: &shouldError},
		Field1: "filled",
	}
	errs = validate.Struct(test3)
	NotEqual(t, errs, nil)
	ve = errs.(ValidationErrors)
	Equal(t, len(ve), 1)
	AssertError(t, errs, "Field1", "Field1", "Field1", "Field1", "excluded_unless")

	shouldPass := "test"
	test4 := struct {
		Inner  *Inner
		FieldE string `validate:"omitempty" json:"field_e"`
		Field1 string `validate:"excluded_unless=Inner.Field test" json:"field_1"`
	}{
		Inner:  &Inner{Field: &shouldPass},
		Field1: "filled",
	}
	errs = validate.Struct(test4)
	Equal(t, errs, nil)

	// test7 and test8: excluded_unless has no effect if FieldER is left blank
	test7 := struct {
		Inner   *Inner
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_unless=Inner.Field test" json:"field_er"`
	}{
		FieldE: "test",
	}
	errs = validate.Struct(test7)
	Equal(t, errs, nil)

	test8 := struct {
		FieldE  string `validate:"omitempty" json:"field_e"`
		FieldER string `validate:"excluded_unless=Inner.Field test" json:"field_er"`
	}{
		FieldE: "test",
	}
	errs = validate.Struct(test8)
	Equal(t, errs, nil)

	// Checks number of params in struct tag is correct
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("panicTest should have panicked!")
		}
	}()
	panicTest := struct {
		Inner  *Inner
		Field1 string `validate:"excluded_unless=Inner.Field" json:"field_1"`
	}{
		Inner: &Inner{Field: &fieldVal},
	}
	_ = validate.Struct(panicTest)
}

func TestLookup(t *testing.T) {
	type Lookup struct {
		FieldA *string `json:"fieldA,omitempty" validate:"required_without=FieldB"`
		FieldB *string `json:"fieldB,omitempty" validate:"required_without=FieldA"`
	}

	fieldAValue := "1232"
	lookup := Lookup{
		FieldA: &fieldAValue,
		FieldB: nil,
	}
	Equal(t, New().Struct(lookup), nil)
}

func TestAbilityToValidateNils(t *testing.T) {
	type TestStruct struct {
		Test *string `validate:"nil"`
	}

	ts := TestStruct{}
	val := New()
	fn := func(fl FieldLevel) bool {
		return fl.Field().Kind() == reflect.Ptr && fl.Field().IsNil()
	}

	err := val.RegisterValidation("nil", fn, true)
	Equal(t, err, nil)

	errs := val.Struct(ts)
	Equal(t, errs, nil)

	str := "string"
	ts.Test = &str

	errs = val.Struct(ts)
	NotEqual(t, errs, nil)
}

func TestRequiredWithoutPointers(t *testing.T) {
	type Lookup struct {
		FieldA *bool `json:"fieldA,omitempty" validate:"required_without=FieldB"`
		FieldB *bool `json:"fieldB,omitempty" validate:"required_without=FieldA"`
	}

	b := true
	lookup := Lookup{
		FieldA: &b,
		FieldB: nil,
	}

	val := New()
	errs := val.Struct(lookup)
	Equal(t, errs, nil)

	b = false
	lookup = Lookup{
		FieldA: &b,
		FieldB: nil,
	}
	errs = val.Struct(lookup)
	Equal(t, errs, nil)
}

func TestRequiredWithoutAllPointers(t *testing.T) {
	type Lookup struct {
		FieldA *bool `json:"fieldA,omitempty" validate:"required_without_all=FieldB"`
		FieldB *bool `json:"fieldB,omitempty" validate:"required_without_all=FieldA"`
	}

	b := true
	lookup := Lookup{
		FieldA: &b,
		FieldB: nil,
	}

	val := New()
	errs := val.Struct(lookup)
	Equal(t, errs, nil)

	b = false
	lookup = Lookup{
		FieldA: &b,
		FieldB: nil,
	}
	errs = val.Struct(lookup)
	Equal(t, errs, nil)
}

func TestGetTag(t *testing.T) {
	var tag string

	type Test struct {
		String string `validate:"mytag"`
	}

	val := New()
	_ = val.RegisterValidation("mytag", func(fl FieldLevel) bool {
		tag = fl.GetTag()
		return true
	})

	var test Test
	errs := val.Struct(test)
	Equal(t, errs, nil)
	Equal(t, tag, "mytag")
}

func TestJSONValidation(t *testing.T) {
	tests := []struct {
		param    interface{}
		expected bool
	}{
		{`foo`, false},
		{`}{`, false},
		{`{]`, false},
		{`{}`, true},
		{`{"foo":"bar"}`, true},
		{`{"foo":"bar","bar":{"baz":["qux"]}}`, true},
		{`{"foo": 3 "bar": 4}`, false},
		{`{"foo": 3 ,"bar": 4`, false},
		{`{foo": 3, "bar": 4}`, false},
		{`foo`, false},
		{`1`, true},
		{`true`, true},
		{`null`, true},
		{`"null"`, true},
		{json.RawMessage(`foo`), false},
		{json.RawMessage(`}{`), false},
		{json.RawMessage(`{]`), false},
		{json.RawMessage(`{}`), true},
		{json.RawMessage(`{"foo":"bar"}`), true},
		{json.RawMessage(`{"foo":"bar","bar":{"baz":["qux"]}}`), true},
		{json.RawMessage(`{"foo": 3 "bar": 4}`), false},
		{json.RawMessage(`{"foo": 3 ,"bar": 4`), false},
		{json.RawMessage(`{foo": 3, "bar": 4}`), false},
		{json.RawMessage(`foo`), false},
		{json.RawMessage(`1`), true},
		{json.RawMessage(`true`), true},
		{json.RawMessage(`null`), true},
		{json.RawMessage(`"null"`), true},
		{[]byte(`foo`), false},
		{[]byte(`}{`), false},
		{[]byte(`{]`), false},
		{[]byte(`{}`), true},
		{[]byte(`{"foo":"bar"}`), true},
		{[]byte(`{"foo":"bar","bar":{"baz":["qux"]}}`), true},
		{[]byte(`{"foo": 3 "bar": 4}`), false},
		{[]byte(`{"foo": 3 ,"bar": 4`), false},
		{[]byte(`{foo": 3, "bar": 4}`), false},
		{[]byte(`foo`), false},
		{[]byte(`1`), true},
		{[]byte(`true`), true},
		{[]byte(`null`), true},
		{[]byte(`"null"`), true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "json")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d json failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d json failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "json" {
					t.Fatalf("Index: %d json failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "json")
	}, "Bad field type int")
}

func TestJWTValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{"eyJhbGciOiJIUzI1NiJ9.eyJuYW1lIjoiZ29waGVyIn0.O_bROM_szPq9qBql-XDHMranHwP48ODdoLICWzqBr_U", true},
		{"acb123-_.def456-_.ghi789-_", true},
		{"eyJhbGciOiJOT05FIn0.e30.", true},
		{"eyJhbGciOiJOT05FIn0.e30.\n", false},
		{"\x00.\x00.\x00", false},
		{"", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "jwt")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d jwt failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d jwt failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "jwt" {
					t.Fatalf("Index: %d jwt failed Error: %s", i, errs)
				}
			}
		}
	}
}

func Test_hostnameport_validator(t *testing.T) {
	type Host struct {
		Addr string `validate:"hostname_port"`
	}

	type testInput struct {
		data     string
		expected bool
	}
	testData := []testInput{
		{"bad..domain.name:234", false},
		{"extra.dot.com.", false},
		{"localhost:1234", true},
		{"192.168.1.1:1234", true},
		{":1234", true},
		{"domain.com:1334", true},
		{"this.domain.com:234", true},
		{"domain:75000", false},
		{"missing.port", false},
	}
	for _, td := range testData {
		h := Host{Addr: td.data}
		v := New()
		err := v.Struct(h)
		if td.expected != (err == nil) {
			t.Fatalf("Test failed for data: %v Error: %v", td.data, err)
		}
	}
}

func Test_port_validator(t *testing.T) {
	type Host struct {
		Port uint32 `validate:"port"`
	}

	type testInput struct {
		data     uint32
		expected bool
	}
	testData := []testInput{
		{0, false},
		{1, true},
		{65535, true},
		{65536, false},
		{65538, false},
	}
	for _, td := range testData {
		h := Host{Port: td.data}
		v := New()
		err := v.Struct(h)
		if td.expected != (err == nil) {
			t.Fatalf("Test failed for data: %v Error: %v", td.data, err)
		}
	}
}

func TestLowercaseValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{`abcdefg`, true},
		{`Abcdefg`, false},
		{"", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "lowercase")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d lowercase failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d lowercase failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "lowercase" {
					t.Fatalf("Index: %d lowercase failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "lowercase")
	}, "Bad field type int")
}

func TestUppercaseValidation(t *testing.T) {
	tests := []struct {
		param    string
		expected bool
	}{
		{`ABCDEFG`, true},
		{`aBCDEFG`, false},
		{"", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.param, "uppercase")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d uppercase failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d uppercase failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "uppercase" {
					t.Fatalf("Index: %d uppercase failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "uppercase")
	}, "Bad field type int")
}

func TestDatetimeValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"datetime=2006-01-02"`
		tag      string
		expected bool
	}{
		{"2008-02-01", `datetime=2006-01-02`, true},
		{"2008-Feb-01", `datetime=2006-01-02`, false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d datetime failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d datetime failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "datetime" {
					t.Fatalf("Index: %d datetime failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "datetime")
	}, "Bad field type int")
}

func TestIsIso3166Alpha2Validation(t *testing.T) {
	tests := []struct {
		value    string `validate:"iso3166_1_alpha2"`
		expected bool
	}{
		{"PL", true},
		{"POL", false},
		{"AA", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_1_alpha2")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha2 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha2 failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso3166Alpha2EUValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"iso3166_1_alpha2_eu"`
		expected bool
	}{
		{"SE", true},
		{"UK", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_1_alpha2_eu")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha2_eu failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha2_eu failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso31662Validation(t *testing.T) {
	tests := []struct {
		value    string `validate:"iso3166_2"`
		expected bool
	}{
		{"US-FL", true},
		{"US-F", false},
		{"US", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_2")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_2 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_2 failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso3166Alpha3Validation(t *testing.T) {
	tests := []struct {
		value    string `validate:"iso3166_1_alpha3"`
		expected bool
	}{
		{"POL", true},
		{"PL", false},
		{"AAA", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_1_alpha3")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha3 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha3 failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso3166Alpha3EUValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"iso3166_1_alpha3_eu"`
		expected bool
	}{
		{"POL", true},
		{"SWE", true},
		{"UNK", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_1_alpha3_eu")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha3_eu failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha3_eu failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso3166AlphaNumericValidation(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{248, true},
		{"248", true},
		{0, false},
		{1, false},
		{"1", false},
		{"invalid_int", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_1_alpha_numeric")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha_numeric failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha_numeric failed Error: %s", i, errs)
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var([]string{"1"}, "iso3166_1_alpha_numeric")
	}, "Bad field type []string")
}

func TestIsIso3166AlphaNumericEUValidation(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{752, true}, //Sweden
		{"752", true},
		{826, false}, // UK
		{"826", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso3166_1_alpha_numeric_eu")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha_numeric_eu failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso3166_1_alpha_numeric_eu failed Error: %s", i, errs)
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var([]string{"1"}, "iso3166_1_alpha_numeric_eu")
	}, "Bad field type []string")
}

func TestCountryCodeValidation(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{248, true},
		{0, false},
		{1, false},
		{"POL", true},
		{"NO", true},
		{"248", true},
		{"1", false},
		{"0", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "country_code")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d country_code failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d country_code failed Error: %s", i, errs)
			}
		}
	}
}

func TestEUCountryCodeValidation(t *testing.T) {
	tests := []struct {
		value    interface{}
		expected bool
	}{
		{724, true},
		{0, false},
		{1, false},
		{"POL", true},
		{"NO", false},
		{"724", true},
		{"1", false},
		{"0", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "eu_country_code")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d eu_country_code failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d eu_country_code failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso4217Validation(t *testing.T) {
	tests := []struct {
		value    string `validate:"iso4217"`
		expected bool
	}{
		{"TRY", true},
		{"EUR", true},
		{"USA", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso4217")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso4217 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso4217 failed Error: %s", i, errs)
			}
		}
	}
}

func TestIsIso4217NumericValidation(t *testing.T) {
	tests := []struct {
		value    int `validate:"iso4217_numeric"`
		expected bool
	}{
		{8, true},
		{12, true},
		{13, false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, "iso4217_numeric")

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso4217 failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d iso4217 failed Error: %s", i, errs)
			}
		}
	}
}

func TestTimeZoneValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"timezone"`
		tag      string
		expected bool
	}{
		// systems may have different time zone database, some systems time zone are case insensitive
		{"America/New_York", `timezone`, true},
		{"UTC", `timezone`, true},
		{"", `timezone`, false},
		{"Local", `timezone`, false},
		{"Unknown", `timezone`, false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d time zone failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d time zone failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "timezone" {
					t.Fatalf("Index: %d time zone failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "timezone")
	}, "Bad field type int")
}

func TestDurationType(t *testing.T) {
	tests := []struct {
		name    string
		s       interface{} // struct
		success bool
	}{
		{
			name: "valid duration string pass",
			s: struct {
				Value time.Duration `validate:"gte=500ns"`
			}{
				Value: time.Second,
			},
			success: true,
		},
		{
			name: "valid duration int pass",
			s: struct {
				Value time.Duration `validate:"gte=500"`
			}{
				Value: time.Second,
			},
			success: true,
		},
	}

	validate := New()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			errs := validate.Struct(tc.s)
			if tc.success {
				Equal(t, errs, nil)
				return
			}
			NotEqual(t, errs, nil)
		})
	}
}

func TestBCP47LanguageTagValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"bcp47_language_tag"`
		tag      string
		expected bool
	}{
		{"en-US", "bcp47_language_tag", true},
		{"en_GB", "bcp47_language_tag", true},
		{"es", "bcp47_language_tag", true},
		{"English", "bcp47_language_tag", false},
		{"ESES", "bcp47_language_tag", false},
		{"az-Cyrl-AZ", "bcp47_language_tag", true},
		{"en-029", "bcp47_language_tag", true},
		{"xog", "bcp47_language_tag", true},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d locale failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d locale failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "bcp47_language_tag" {
					t.Fatalf("Index: %d locale failed Error: %s", i, errs)
				}
			}
		}
	}

	PanicMatches(t, func() {
		_ = validate.Var(2, "bcp47_language_tag")
	}, "Bad field type int")
}

func TestBicIsoFormatValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"bic"`
		tag      string
		expected bool
	}{
		{"SBICKEN1345", "bic", true},
		{"SBICKEN1", "bic", true},
		{"SBICKENY", "bic", true},
		{"SBICKEN1YYP", "bic", true},
		{"SBIC23NXXX", "bic", false},
		{"S23CKENXXXX", "bic", false},
		{"SBICKENXX", "bic", false},
		{"SBICKENXX9", "bic", false},
		{"SBICKEN13458", "bic", false},
		{"SBICKEN", "bic", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d bic failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d bic failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "bic" {
					t.Fatalf("Index: %d bic failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestSemverFormatValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"semver"`
		tag      string
		expected bool
	}{
		{"1.2.3", "semver", true},
		{"10.20.30", "semver", true},
		{"1.1.2-prerelease+meta", "semver", true},
		{"1.1.2+meta", "semver", true},
		{"1.1.2+meta-valid", "semver", true},
		{"1.0.0-alpha", "semver", true},
		{"1.0.0-alpha.1", "semver", true},
		{"1.0.0-alpha.beta", "semver", true},
		{"1.0.0-alpha.beta.1", "semver", true},
		{"1.0.0-alpha0.valid", "semver", true},
		{"1.0.0-alpha.0valid", "semver", true},
		{"1.0.0-alpha-a.b-c-somethinglong+build.1-aef.1-its-okay", "semver", true},
		{"1.0.0-rc.1+build.1", "semver", true},
		{"1.0.0-rc.1+build.123", "semver", true},
		{"1.2.3-beta", "semver", true},
		{"1.2.3-DEV-SNAPSHOT", "semver", true},
		{"1.2.3-SNAPSHOT-123", "semver", true},
		{"2.0.0+build.1848", "semver", true},
		{"2.0.1-alpha.1227", "semver", true},
		{"1.0.0-alpha+beta", "semver", true},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12+788", "semver", true},
		{"1.2.3----R-S.12.9.1--.12+meta", "semver", true},
		{"1.2.3----RC-SNAPSHOT.12.9.1--.12", "semver", true},
		{"1.0.0+0.build.1-rc.10000aaa-kk-0.1", "semver", true},
		{"99999999999999999999999.999999999999999999.99999999999999999", "semver", true},
		{"1.0.0-0A.is.legal", "semver", true},
		{"1", "semver", false},
		{"1.2", "semver", false},
		{"1.2.3-0123", "semver", false},
		{"1.2.3-0123.0123", "semver", false},
		{"1.1.2+.123", "semver", false},
		{"+invalid", "semver", false},
		{"-invalid", "semver", false},
		{"-invalid+invalid", "semver", false},
		{"alpha", "semver", false},
		{"alpha.beta.1", "semver", false},
		{"alpha.1", "semver", false},
		{"1.0.0-alpha_beta", "semver", false},
		{"1.0.0-alpha_beta", "semver", false},
		{"1.0.0-alpha...1", "semver", false},
		{"01.1.1", "semver", false},
		{"1.01.1", "semver", false},
		{"1.1.01", "semver", false},
		{"1.2", "semver", false},
		{"1.2.Dev", "semver", false},
		{"1.2.3.Dev", "semver", false},
		{"1.2-SNAPSHOT", "semver", false},
	}

	validate := New()

	for i, test := range tests {

		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d semver failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d semver failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "semver" {
					t.Fatalf("Index: %d semver failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestCveFormatValidation(t *testing.T) {

	tests := []struct {
		value    string `validate:"cve"`
		tag      string
		expected bool
	}{
		{"CVE-1999-0001", "cve", true},
		{"CVE-1998-0001", "cve", false},
		{"CVE-2000-0001", "cve", true},
		{"CVE-2222-0001", "cve", true},
		{"2222-0001", "cve", false},
		{"-2222-0001", "cve", false},
		{"CVE22220001", "cve", false},
		{"CVE-2222-000001", "cve", false},
		{"CVE-2222-100001", "cve", true},
		{"CVE-2222-99999999999", "cve", true},
		{"CVE-3000-0001", "cve", false},
		{"CVE-1999-0000", "cve", false},
		{"CVE-2099-0000", "cve", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d cve failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d cve failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "cve" {
					t.Fatalf("Index: %d cve failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestRFC1035LabelFormatValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"dns_rfc1035_label"`
		tag      string
		expected bool
	}{
		{"abc", "dns_rfc1035_label", true},
		{"abc-", "dns_rfc1035_label", false},
		{"abc-123", "dns_rfc1035_label", true},
		{"ABC", "dns_rfc1035_label", false},
		{"ABC-123", "dns_rfc1035_label", false},
		{"abc-abc", "dns_rfc1035_label", true},
		{"ABC-ABC", "dns_rfc1035_label", false},
		{"123-abc", "dns_rfc1035_label", false},
		{"", "dns_rfc1035_label", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d dns_rfc1035_label failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d dns_rfc1035_label failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "dns_rfc1035_label" {
					t.Fatalf("Index: %d dns_rfc1035_label failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestPostCodeByIso3166Alpha2(t *testing.T) {
	tests := map[string][]struct {
		value    string
		expected bool
	}{
		"VN": {
			{"ABC", false},
			{"700000", true},
			{"A1", false},
		},
		"GB": {
			{"EC1A 1BB", true},
			{"CF10 1B1H", false},
		},
		"VI": {
			{"00803", true},
			{"1234567", false},
		},
		"LC": {
			// not support regexp for post code
			{"123456", false},
		},
		"XX": {
			// not support country
			{"123456", false},
		},
	}

	validate := New()

	for cc, ccTests := range tests {
		for i, test := range ccTests {
			errs := validate.Var(test.value, fmt.Sprintf("postcode_iso3166_alpha2=%s", cc))

			if test.expected {
				if !IsEqual(errs, nil) {
					t.Fatalf("Index: %d postcode_iso3166_alpha2=%s failed Error: %s", i, cc, errs)
				}
			} else {
				if IsEqual(errs, nil) {
					t.Fatalf("Index: %d postcode_iso3166_alpha2=%s failed Error: %s", i, cc, errs)
				}
			}
		}
	}
}

func TestPostCodeByIso3166Alpha2Field(t *testing.T) {
	tests := []struct {
		Value       string `validate:"postcode_iso3166_alpha2_field=CountryCode"`
		CountryCode interface{}
		expected    bool
	}{
		{"ABC", "VN", false},
		{"700000", "VN", true},
		{"A1", "VN", false},
		{"EC1A 1BB", "GB", true},
		{"CF10 1B1H", "GB", false},
		{"00803", "VI", true},
		{"1234567", "VI", false},
		{"123456", "LC", false}, // not support regexp for post code
		{"123456", "XX", false}, // not support country
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Struct(test)
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d postcode_iso3166_alpha2_field=CountryCode failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d postcode_iso3166_alpha2_field=CountryCode failed Error: %s", i, errs)
			}
		}
	}
}

func TestPostCodeByIso3166Alpha2Field_WrongField(t *testing.T) {
	type test struct {
		Value        string `validate:"postcode_iso3166_alpha2_field=CountryCode"`
		CountryCode1 interface{}
		expected     bool
	}

	errs := New().Struct(test{"ABC", "VN", false})
	NotEqual(t, nil, errs)
}

func TestPostCodeByIso3166Alpha2Field_MissingParam(t *testing.T) {
	type test struct {
		Value        string `validate:"postcode_iso3166_alpha2_field="`
		CountryCode1 interface{}
		expected     bool
	}

	errs := New().Struct(test{"ABC", "VN", false})
	NotEqual(t, nil, errs)
}

func TestPostCodeByIso3166Alpha2Field_InvalidKind(t *testing.T) {
	type test struct {
		Value       string `validate:"postcode_iso3166_alpha2_field=CountryCode"`
		CountryCode interface{}
		expected    bool
	}
	defer func() { _ = recover() }()

	_ = New().Struct(test{"ABC", 123, false})
	t.Errorf("Didn't panic as expected")
}

func TestValidate_ValidateMapCtx(t *testing.T) {

	type args struct {
		data  map[string]interface{}
		rules map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test nested map in slice",
			args: args{
				data: map[string]interface{}{
					"Test_A": map[string]interface{}{
						"Test_B": "Test_B",
						"Test_C": []map[string]interface{}{
							{
								"Test_D": "Test_D",
							},
						},
						"Test_E": map[string]interface{}{
							"Test_F": "Test_F",
						},
					},
				},
				rules: map[string]interface{}{
					"Test_A": map[string]interface{}{
						"Test_B": "min=2",
						"Test_C": map[string]interface{}{
							"Test_D": "min=2",
						},
						"Test_E": map[string]interface{}{
							"Test_F": "min=2",
						},
					},
				},
			},
			want: 0,
		},

		{
			name: "test nested map error",
			args: args{
				data: map[string]interface{}{
					"Test_A": map[string]interface{}{
						"Test_B": "Test_B",
						"Test_C": []interface{}{"Test_D"},
						"Test_E": map[string]interface{}{
							"Test_F": "Test_F",
						},
						"Test_G": "Test_G",
						"Test_I": []map[string]interface{}{
							{
								"Test_J": "Test_J",
							},
						},
					},
				},
				rules: map[string]interface{}{
					"Test_A": map[string]interface{}{
						"Test_B": "min=2",
						"Test_C": map[string]interface{}{
							"Test_D": "min=2",
						},
						"Test_E": map[string]interface{}{
							"Test_F": "min=100",
						},
						"Test_G": map[string]interface{}{
							"Test_H": "min=2",
						},
						"Test_I": map[string]interface{}{
							"Test_J": "min=100",
						},
					},
				},
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validate := New()
			if got := validate.ValidateMapCtx(context.Background(), tt.args.data, tt.args.rules); len(got) != tt.want {
				t.Errorf("ValidateMapCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoDBObjectIDFormatValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"mongodb"`
		tag      string
		expected bool
	}{
		{"507f191e810c19729de860ea", "mongodb", true},
		{"507f191e810c19729de860eG", "mongodb", false},
		{"M07f191e810c19729de860eG", "mongodb", false},
		{"07f191e810c19729de860ea", "mongodb", false},
		{"507f191e810c19729de860e", "mongodb", false},
		{"507f191e810c19729de860ea4", "mongodb", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d mongodb failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d mongodb failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "mongodb" {
					t.Fatalf("Index: %d mongodb failed Error: %s", i, errs)
				}
			}
		}
	}
}
func TestMongoDBConnectionStringFormatValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"mongodb_connection_string"`
		tag      string
		expected bool
	}{
		{"mongodb://username:password@server.example.com:20017/database?replicaSet=test&connectTimeoutMS=300000&ssl=true", "mongodb_connection_string", true},
		{"mongodb+srv://username:password@server.example.com:20017/database?replicaSet=test&connectTimeoutMS=300000&ssl=true", "mongodb_connection_string", true},
		{"mongodb+srv://username:password@server.example.com:20017,server.example.com,server.example.com:20017/database?replicaSet=test&connectTimeoutMS=300000&ssl=true", "mongodb_connection_string", true},
		{"mongodb+srv://username:password@server.example.com:20017,server.example.com,server.example.com:20017?replicaSet=test&connectTimeoutMS=300000&ssl=true", "mongodb_connection_string", true},
		{"mongodb://username:password@server.example.com:20017,server.example.com,server.example.com:20017/database?replicaSet=test&connectTimeoutMS=300000&ssl=true", "mongodb_connection_string", true},
		{"mongodb://localhost", "mongodb_connection_string", true},
		{"mongodb://localhost:27017", "mongodb_connection_string", true},
		{"localhost", "mongodb_connection_string", false},
		{"mongodb://", "mongodb_connection_string", false},
		{"mongodb+srv://", "mongodb_connection_string", false},
		{"mongodbsrv://localhost", "mongodb_connection_string", false},
		{"mongodb+srv://localhost", "mongodb_connection_string", true},
		{"mongodb+srv://localhost:27017", "mongodb_connection_string", true},
		{"mongodb+srv://localhost?replicaSet=test", "mongodb_connection_string", true},
		{"mongodb+srv://localhost:27017?replicaSet=test", "mongodb_connection_string", true},
		{"mongodb+srv://localhost:27017?", "mongodb_connection_string", false},
		{"mongodb+srv://localhost:27017?replicaSet", "mongodb_connection_string", false},
		{"mongodb+srv://localhost/database", "mongodb_connection_string", true},
		{"mongodb+srv://localhost:27017/database", "mongodb_connection_string", true},
		{"mongodb+srv://username@localhost", "mongodb_connection_string", false},
		{"mongodb+srv://username:password@localhost", "mongodb_connection_string", true},
		{"mongodb+srv://username:password@localhost:27017", "mongodb_connection_string", true},
		{"mongodb+srv://username:password@localhost:27017,192.0.0.7,192.0.0.9:27018,server.example.com", "mongodb_connection_string", true},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d mongodb_connection_string failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d mongodb_connection_string failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "mongodb_connection_string" {
					t.Fatalf("Index: %d mongodb_connection_string failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestSpiceDBValueFormatValidation(t *testing.T) {
	tests := []struct {
		value    string
		tag      string
		expected bool
	}{
		// Must be an asterisk OR a string containing alphanumeric characters and a restricted set a special symbols: _ | / - = +
		{"*", "spicedb=id", true},
		{`azAZ09_|/-=+`, "spicedb=id", true},
		{`a*`, "spicedb=id", false},
		{`/`, "spicedb=id", true},
		{"*", "spicedb", true},

		// Must begin and end with a lowercase letter, may also contain numbers and underscores between, min length 3, max length 64
		{"a", "spicedb=permission", false},
		{"1", "spicedb=permission", false},
		{"a1", "spicedb=permission", false},
		{"a_b", "spicedb=permission", true},
		{"A_b", "spicedb=permission", false},
		{"a_B", "spicedb=permission", false},
		{"abcdefghijklmnopqrstuvwxyz_0123456789_abcdefghijklmnopqrstuvwxyz", "spicedb=permission", true},
		{"abcdefghijklmnopqrstuvwxyz_01234_56789_abcdefghijklmnopqrstuvwxyz", "spicedb=permission", false},

		// Object types follow the same rules as permissions for the type name plus an optional prefix up to 63 characters with a /
		{"a", "spicedb=type", false},
		{"1", "spicedb=type", false},
		{"a1", "spicedb=type", false},
		{"a_b", "spicedb=type", true},
		{"A_b", "spicedb=type", false},
		{"a_B", "spicedb=type", false},
		{"abcdefghijklmnopqrstuvwxyz_0123456789_abcdefghijklmnopqrstuvwxyz", "spicedb=type", true},
		{"abcdefghijklmnopqrstuvwxyz_01234_56789_abcdefghijklmnopqrstuvwxyz", "spicedb=type", false},

		{`a_b/a`, "spicedb=type", false},
		{`a_b/1`, "spicedb=type", false},
		{`a_b/a1`, "spicedb=type", false},
		{`a_b/a_b`, "spicedb=type", true},
		{`a_b/A_b`, "spicedb=type", false},
		{`a_b/a_B`, "spicedb=type", false},
		{`a_b/abcdefghijklmnopqrstuvwxyz_0123456789_abcdefghijklmnopqrstuvwxyz`, "spicedb=type", true},
		{`a_b/abcdefghijklmnopqrstuvwxyz_01234_56789_abcdefghijklmnopqrstuvwxyz`, "spicedb=type", false},

		{`a/a_b`, "spicedb=type", false},
		{`1/a_b`, "spicedb=type", false},
		{`a1/a_b`, "spicedb=type", false},
		{`a_b/a_b`, "spicedb=type", true},
		{`A_b/a_b`, "spicedb=type", false},
		{`a_B/a_b`, "spicedb=type", false},
		{`abcdefghijklmnopqrstuvwxyz_0123456789_abcdefghijklmnopqrstuvwxy/a_b`, "spicedb=type", true},
		{`abcdefghijklmnopqrstuvwxyz_0123456789_abcdefghijklmnopqrstuvwxyz/a_b`, "spicedb=type", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d spicedb failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d spicedb - expected error but there was none.", i)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "spicedb" {
					t.Fatalf("Index: %d spicedb failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestCreditCardFormatValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"credit_card"`
		tag      string
		expected bool
	}{
		{"586824160825533338", "credit_card", true},
		{"586824160825533328", "credit_card", false},
		{"4624748233249780", "credit_card", true},
		{"4624748233349780", "credit_card", false},
		{"378282246310005", "credit_card", true},
		{"378282146310005", "credit_card", false},
		{"4624 7482 3324 9780", "credit_card", true},
		{"4624 7482 3324  9780", "credit_card", false},
		{"4624 7482 3324 978A", "credit_card", false},
		{"4624 7482 332", "credit_card", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)

		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d credit_card failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d credit_card failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "credit_card" {
					t.Fatalf("Index: %d credit_card failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestLuhnChecksumValidation(t *testing.T) {
	testsUint := []struct {
		value    interface{} `validate:"luhn_checksum"` // the type is interface{} because the luhn_checksum works on both strings and numbers
		tag      string
		expected bool
	}{
		{uint64(586824160825533338), "luhn_checksum", true}, // credit card numbers are just special cases of numbers with luhn checksum
		{586824160825533338, "luhn_checksum", true},
		{"586824160825533338", "luhn_checksum", true},
		{uint64(586824160825533328), "luhn_checksum", false},
		{586824160825533328, "luhn_checksum", false},
		{"586824160825533328", "luhn_checksum", false},
		{10000000116, "luhn_checksum", true}, // but there may be shorter numbers (11 digits)
		{"10000000116", "luhn_checksum", true},
		{10000000117, "luhn_checksum", false},
		{"10000000117", "luhn_checksum", false},
		{uint64(12345678123456789011), "luhn_checksum", true}, // or longer numbers (19 digits)
		{"12345678123456789011", "luhn_checksum", true},
		{1, "luhn_checksum", false}, // single digits (checksum only) are not allowed
		{"1", "luhn_checksum", false},
		{-10, "luhn_checksum", false}, // negative ints are not allowed
		{"abcdefghijklmnop", "luhn_checksum", false},
	}

	validate := New()

	for i, test := range testsUint {
		errs := validate.Var(test.value, test.tag)
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d luhn_checksum failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d luhn_checksum failed Error: %s", i, errs)
			} else {
				val := getError(errs, "", "")
				if val.Tag() != "luhn_checksum" {
					t.Fatalf("Index: %d luhn_checksum failed Error: %s", i, errs)
				}
			}
		}
	}
}

func TestMultiOrOperatorGroup(t *testing.T) {
	tests := []struct {
		Value    int `validate:"eq=1|gte=5,eq=1|lt=7"`
		expected bool
	}{
		{1, true}, {2, false}, {5, true}, {6, true}, {8, false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Struct(test)
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf("Index: %d multi_group_of_OR_operators failed Error: %s", i, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf("Index: %d multi_group_of_OR_operators should have errs", i)
			}
		}
	}
}

func TestCronExpressionValidation(t *testing.T) {
	tests := []struct {
		value    string `validate:"cron"`
		tag      string
		expected bool
	}{
		{"0 0 12 * * ?", "cron", true},
		{"0 15 10 ? * *", "cron", true},
		{"0 15 10 * * ?", "cron", true},
		{"0 15 10 * * ? 2005", "cron", true},
		{"0 15 10 ? * 6L", "cron", true},
		{"0 15 10 ? * 6L 2002-2005", "cron", true},
		{"*/20 * * * *", "cron", true},
		{"0 15 10 ? * MON-FRI", "cron", true},
		{"0 15 10 ? * 6#3", "cron", true},
		{"0 */15 * * *", "cron", true},
		{"wrong", "cron", false},
	}

	validate := New()

	for i, test := range tests {
		errs := validate.Var(test.value, test.tag)
		if test.expected {
			if !IsEqual(errs, nil) {
				t.Fatalf(`Index: %d cron "%s" failed Error: %s`, i, test.value, errs)
			}
		} else {
			if IsEqual(errs, nil) {
				t.Fatalf(`Index: %d cron "%s" should have errs`, i, test.value)
			}
		}
	}
}

func TestNestedStructValidation(t *testing.T) {
	validator := New(WithRequiredStructEnabled())

	t.Run("required", func(t *testing.T) {
		type (
			value struct {
				Field string
			}
			topLevel struct {
				Nested value `validate:"required"`
			}
		)

		var validationErrs ValidationErrors
		if errs := validator.Struct(topLevel{}); errs != nil {
			validationErrs = errs.(ValidationErrors)
		}

		Equal(t, 1, len(validationErrs))
		AssertError(t, validationErrs, "topLevel.Nested", "topLevel.Nested", "Nested", "Nested", "required")

		Equal(t, validator.Struct(topLevel{value{"potato"}}), nil)
	})

	t.Run("omitempty", func(t *testing.T) {
		type (
			value struct {
				Field string
			}
			topLevel struct {
				Nested value `validate:"omitempty,required"`
			}
		)

		errs := validator.Struct(topLevel{})
		Equal(t, errs, nil)
	})

	t.Run("excluded_if", func(t *testing.T) {
		type (
			value struct {
				Field string
			}
			topLevel struct {
				Field  string
				Nested value `validate:"excluded_if=Field potato"`
			}
		)

		errs := validator.Struct(topLevel{Field: "test", Nested: value{"potato"}})
		Equal(t, errs, nil)

		errs = validator.Struct(topLevel{Field: "potato"})
		Equal(t, errs, nil)

		errs = validator.Struct(topLevel{Field: "potato", Nested: value{"potato"}})
		AssertError(t, errs, "topLevel.Nested", "topLevel.Nested", "Nested", "Nested", "excluded_if")
	})

	t.Run("excluded_unless", func(t *testing.T) {
		type (
			value struct {
				Field string
			}
			topLevel struct {
				Field  string
				Nested value `validate:"excluded_unless=Field potato"`
			}
		)

		errs := validator.Struct(topLevel{Field: "test"})
		Equal(t, errs, nil)

		errs = validator.Struct(topLevel{Field: "potato", Nested: value{"potato"}})
		Equal(t, errs, nil)

		errs = validator.Struct(topLevel{Field: "test", Nested: value{"potato"}})
		AssertError(t, errs, "topLevel.Nested", "topLevel.Nested", "Nested", "Nested", "excluded_unless")
	})

	t.Run("nonComparableField", func(t *testing.T) {
		type (
			value struct {
				Field []string
			}
			topLevel struct {
				Nested value `validate:"required"`
			}
		)

		errs := validator.Struct(topLevel{value{[]string{}}})
		Equal(t, errs, nil)
	})

	type (
		veggyBasket struct {
			Root   string
			Squash string `validate:"required"`
		}
		testErr struct {
			path string
			tag  string
		}
		test struct {
			name  string
			err   testErr
			value veggyBasket
		}
	)

	if err := validator.RegisterValidation("veggy", func(f FieldLevel) bool {
		v, ok := f.Field().Interface().(veggyBasket)
		if !ok || v.Root != "potato" {
			return false
		}
		return true
	}); err != nil {
		t.Fatal(fmt.Errorf("failed to register potato tag: %w", err))
	}

	tests := []test{
		{
			name:  "valid",
			value: veggyBasket{"potato", "zucchini"},
		}, {
			name:  "failedCustomTag",
			value: veggyBasket{"zucchini", "potato"},
			err:   testErr{"topLevel.VeggyBasket", "veggy"},
		}, {
			name:  "failedInnerField",
			value: veggyBasket{"potato", ""},
			err:   testErr{"topLevel.VeggyBasket.Squash", "required"},
		}, {
			name:  "customTagFailurePriorityCheck",
			value: veggyBasket{"zucchini", ""},
			err:   testErr{"topLevel.VeggyBasket", "veggy"},
		},
	}

	var evaluateTest = func(tt test, errs error) {
		if tt.err != (testErr{}) && errs != nil {
			Equal(t, len(errs.(ValidationErrors)), 1)

			segments := strings.Split(tt.err.path, ".")
			fieldName := segments[len(segments)-1]
			AssertError(t, errs, tt.err.path, tt.err.path, fieldName, fieldName, tt.err.tag)
		}

		shouldFail := tt.err != (testErr{})
		hasFailed := errs != nil
		if shouldFail != hasFailed {
			t.Fatalf("expected failure %v, got: %v with errs: %v", shouldFail, hasFailed, errs)
		}
	}

	for _, tt := range tests {
		type topLevel struct {
			VeggyBasket veggyBasket `validate:"veggy"`
		}

		t.Run(tt.name, func(t *testing.T) {
			evaluateTest(tt, validator.Struct(topLevel{tt.value}))
		})
	}

	// Also test on struct pointers
	for _, tt := range tests {
		type topLevel struct {
			VeggyBasket *veggyBasket `validate:"veggy"`
		}

		t.Run(tt.name+"Ptr", func(t *testing.T) {
			evaluateTest(tt, validator.Struct(topLevel{&tt.value}))
		})
	}
}

func TestTimeRequired(t *testing.T) {
	validate := New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	type TestTime struct {
		Time time.Time `validate:"required"`
	}

	var testTime TestTime

	err := validate.Struct(&testTime)
	NotEqual(t, err, nil)
	AssertError(t, err.(ValidationErrors), "TestTime.Time", "TestTime.Time", "Time", "Time", "required")
}

func TestOmitNilAndRequired(t *testing.T) {
	type (
		OmitEmpty struct {
			Str    string  `validate:"omitempty,required,min=10"`
			StrPtr *string `validate:"omitempty,required,min=10"`
			Inner  *OmitEmpty
		}
		OmitNil struct {
			Str    string  `validate:"omitnil,required,min=10"`
			StrPtr *string `validate:"omitnil,required,min=10"`
			Inner  *OmitNil
		}
	)

	var (
		validate = New(WithRequiredStructEnabled())
		valid    = "this is the long string to pass the validation rule"
	)

	t.Run("compare using valid data", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{Str: valid, StrPtr: &valid, Inner: &OmitEmpty{Str: valid, StrPtr: &valid}})
		err2 := validate.Struct(OmitNil{Str: valid, StrPtr: &valid, Inner: &OmitNil{Str: valid, StrPtr: &valid}})

		Equal(t, err1, nil)
		Equal(t, err2, nil)
	})

	t.Run("compare fully empty omitempty and omitnil", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{})
		err2 := validate.Struct(OmitNil{})

		Equal(t, err1, nil)
		AssertError(t, err2, "OmitNil.Str", "OmitNil.Str", "Str", "Str", "required")
	})

	t.Run("validate in deep", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{Str: valid, Inner: &OmitEmpty{}})
		err2 := validate.Struct(OmitNil{Str: valid, Inner: &OmitNil{}})

		Equal(t, err1, nil)
		AssertError(t, err2, "OmitNil.Inner.Str", "OmitNil.Inner.Str", "Str", "Str", "required")
	})
}

func TestOmitZero(t *testing.T) {
	type (
		OmitEmpty struct {
			Str    string  `validate:"omitempty,min=10"`
			StrPtr *string `validate:"omitempty,min=10"`
		}
		OmitZero struct {
			Str    string  `validate:"omitzero,min=10"`
			StrPtr *string `validate:"omitzero,min=10"`
		}
	)

	var (
		validate = New()
		valid    = "this is the long string to pass the validation rule"
		empty    = ""
	)

	t.Run("compare using valid data", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{Str: valid, StrPtr: &valid})
		err2 := validate.Struct(OmitZero{Str: valid, StrPtr: &valid})

		Equal(t, err1, nil)
		Equal(t, err2, nil)
	})

	t.Run("compare fully empty omitempty and omitzero", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{})
		err2 := validate.Struct(OmitZero{})

		Equal(t, err1, nil)
		Equal(t, err2, nil)
	})

	t.Run("compare with zero value", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{Str: "", StrPtr: nil})
		err2 := validate.Struct(OmitZero{Str: "", StrPtr: nil})

		Equal(t, err1, nil)
		Equal(t, err2, nil)
	})

	t.Run("compare with empty value", func(t *testing.T) {
		err1 := validate.Struct(OmitEmpty{Str: empty, StrPtr: &empty})
		err2 := validate.Struct(OmitZero{Str: empty, StrPtr: &empty})

		AssertError(t, err1, "OmitEmpty.StrPtr", "OmitEmpty.StrPtr", "StrPtr", "StrPtr", "min")
		Equal(t, err2, nil)
	})
}

func TestPrivateFieldsStruct(t *testing.T) {
	type tc struct {
		stct     interface{}
		errorNum int
	}

	tcs := []tc{
		{
			stct: &struct {
				f1 int8  `validate:"required"`
				f2 int16 `validate:"required"`
				f3 int32 `validate:"required"`
				f4 int64 `validate:"required"`
			}{},
			errorNum: 4,
		},
		{
			stct: &struct {
				f1 uint8  `validate:"required"`
				f2 uint16 `validate:"required"`
				f3 uint32 `validate:"required"`
				f4 uint64 `validate:"required"`
			}{},
			errorNum: 4,
		},
		{
			stct: &struct {
				f1 complex64  `validate:"required"`
				f2 complex128 `validate:"required"`
			}{},
			errorNum: 2,
		},
		{
			stct: &struct {
				f1 float32 `validate:"required"`
				f2 float64 `validate:"required"`
			}{},
			errorNum: 2,
		},
		{
			stct: struct {
				f1 int8  `validate:"required"`
				f2 int16 `validate:"required"`
				f3 int32 `validate:"required"`
				f4 int64 `validate:"required"`
			}{},
			errorNum: 4,
		},
		{
			stct: struct {
				f1 uint8  `validate:"required"`
				f2 uint16 `validate:"required"`
				f3 uint32 `validate:"required"`
				f4 uint64 `validate:"required"`
			}{},
			errorNum: 4,
		},
		{
			stct: struct {
				f1 complex64  `validate:"required"`
				f2 complex128 `validate:"required"`
			}{},
			errorNum: 2,
		},
		{
			stct: struct {
				f1 float32 `validate:"required"`
				f2 float64 `validate:"required"`
			}{},
			errorNum: 2,
		},
		{
			stct: struct {
				f1 *int `validate:"required"`
				f2 struct {
					f3 int `validate:"required"`
				}
			}{},
			errorNum: 2,
		},
		{
			stct: &struct {
				f1 *int `validate:"required"`
				f2 struct {
					f3 int `validate:"required"`
				}
			}{},
			errorNum: 2,
		},
	}

	validate := New(WithPrivateFieldValidation())

	for _, tc := range tcs {
		err := validate.Struct(tc.stct)
		NotEqual(t, err, nil)

		errs := err.(ValidationErrors)
		Equal(t, len(errs), tc.errorNum)
	}
}
