package ru

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/locales"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// RegisterDefaultTranslations registers a set of default translations
// for all built in tag's in validator; you may add your own as desired.
func RegisterDefaultTranslations(v *validator.Validate, trans ut.Translator) (err error) {

	translations := []struct {
		tag             string
		translation     string
		override        bool
		customRegisFunc validator.RegisterTranslationsFunc
		customTransFunc validator.TranslationFunc
	}{
		{
			tag:         "required",
			translation: "{0} обязательное поле",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("len-string", "{0} должен быть длиной в {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} символа", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0} должен быть равен {1}", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0} должен содержать {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} элемента", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("len-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("len-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("len-items", fe.Field(), c)

				default:
					t, err = ut.T("len-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "min",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("min-string", "{0} должен содержать минимум {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} символа", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0} должен быть больше или равно {1}", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} должен содержать минимум {1}", false); err != nil {
					return
				}
				if err = ut.AddCardinal("min-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} элемента", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("min-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("min-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("min-items", fe.Field(), c)

				default:
					t, err = ut.T("min-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "max",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("max-string", "{0} должен содержать максимум {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} символа", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0} должен быть меньше или равно {1}", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} должен содержать максимум {1}", false); err != nil {
					return
				}
				if err = ut.AddCardinal("max-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} элемента", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string

				var digits uint64
				var kind reflect.Kind

				if idx := strings.Index(fe.Param(), "."); idx != -1 {
					digits = uint64(len(fe.Param()[idx+1:]))
				}

				f64, err := strconv.ParseFloat(fe.Param(), 64)
				if err != nil {
					goto END
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					c, err = ut.C("max-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					c, err = ut.C("max-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("max-items", fe.Field(), c)

				default:
					t, err = ut.T("max-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eq",
			translation: "{0} не равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ne",
			translation: "{0} должен быть не равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lt-string", "{0} должен иметь менее {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} символов", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0} должен быть менее {1}", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0} должен содержать менее {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} элементов", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} must be less than the current Date & Time", false); err != nil {
					return
				}

				return

			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("lt-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "lte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("lte-string", "{0} должен содержать максимум {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} символа", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0} должен быть менее или равен {1}", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0} должен содержать максимум {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} элемента", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} must be less than or equal to the current Date & Time", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("lte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("lte-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("lte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gt",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gt-string", "{0} должен быть длиннее {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} символов", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0} должен быть больше {1}", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0} должен содержать более {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} элементов", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} должна быть позже текущего момента", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gt-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("gt-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gt-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag: "gte",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("gte-string", "{0} должен содержать минимум {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} символ", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} символа", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} символов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} символы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0} должен быть больше или равно {1}", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} должен содержать минимум {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} элемент", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} элемента", locales.PluralRuleFew, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} элементов", locales.PluralRuleMany, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} элементы", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} должна быть позже или равна текущему моменту", false); err != nil {
					return
				}

				return
			},
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				var err error
				var t string
				var f64 float64
				var digits uint64
				var kind reflect.Kind

				fn := func() (err error) {

					if idx := strings.Index(fe.Param(), "."); idx != -1 {
						digits = uint64(len(fe.Param()[idx+1:]))
					}

					f64, err = strconv.ParseFloat(fe.Param(), 64)

					return
				}

				kind = fe.Kind()
				if kind == reflect.Ptr {
					kind = fe.Type().Elem().Kind()
				}

				switch kind {
				case reflect.String:

					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-string-character", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-string", fe.Field(), c)

				case reflect.Slice, reflect.Map, reflect.Array:
					var c string

					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("gte-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-items", fe.Field(), c)

				case reflect.Struct:
					if fe.Type() != reflect.TypeOf(time.Time{}) {
						err = fmt.Errorf("tag '%s' cannot be used on a struct type", fe.Tag())
						goto END
					}

					t, err = ut.T("gte-datetime", fe.Field())

				default:
					err = fn()
					if err != nil {
						goto END
					}

					t, err = ut.T("gte-number", fe.Field(), ut.FmtNumber(f64, digits))
				}

			END:
				if err != nil {
					fmt.Printf("warning: error translating FieldError: %s", err)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqfield",
			translation: "{0} должен быть равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "{0} должен быть равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0} не должен быть равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "{0} должен быть больше {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "{0} должен быть больше или равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "{0} должен быть менее {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "{0} должен быть менее или равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0} не должен быть равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtfield",
			translation: "{0} должен быть больше {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "{0} должен быть больше или равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "{0} должен быть менее {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "{0} должен быть менее или равен {1}",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "alpha",
			translation: "{0} должен содержать только буквы",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0} должен содержать только буквы и цифры",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0} должен быть цифровым значением",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0} должен быть цифрой",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0} должен быть шестнадцатеричной строкой",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0} должен быть HEX цветом",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0} должен быть RGB цветом",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0} должен быть RGBA цветом",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0} должен быть HSL цветом",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0} должен быть HSLA цветом",
			override:    false,
		},
		{
			tag:         "e164",
			translation: "{0} должен быть действительным номером телефона в формате E.164",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0} должен быть email адресом",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0} должен быть URL",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0} должен быть URI",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0} должен быть Base64 строкой",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0} должен содержать текст '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "{0} должен содержать минимум один из символов '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "{0} не должен содержать текст '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "{0} не должен содержать символы '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "{0} не должен содержать '{1}'",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {

				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "{0} должен быть ISBN номером",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0} должен быть ISBN-10 номером",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0} должен быть ISBN-13 номером",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "{0} должен быть ISSN номером",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0} должен быть UUID",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0} должен быть UUID 3 версии",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0} должен быть UUID 4 версии",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0} должен быть UUID 5 версии",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "{0} должен быть ULID",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0} должен содержать только ascii символы",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0} должен содержать только доступные для печати ascii символы",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0} должен содержать мультибайтные символы",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0} должен содержать Data URI",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0} должен содержать координаты широты",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0} должен содержать координаты долготы",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0} должен быть SSN номером",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0} должен быть IPv4 адресом",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0} должен быть IPv6 адресом",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0} должен быть IP адресом",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0} должен содержать CIDR обозначения",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0} должен содержать CIDR обозначения для IPv4 адреса",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0} должен содержать CIDR обозначения для IPv6 адреса",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0} должен быть TCP адресом",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} должен быть IPv4 TCP адресом",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} должен быть IPv6 TCP адресом",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0} должен быть UDP адресом",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0} должен быть IPv4 UDP адресом",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0} должен быть IPv6 UDP адресом",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0} должен быть распознаваемым IP адресом",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0} должен быть распознаваемым IPv4 адресом",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0} должен быть распознаваемым IPv6 адресом",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0} должен быть распознаваемым UNIX адресом",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0} должен содержать MAC адрес",
			override:    false,
		},
		{
			tag:         "unique",
			translation: "{0} должен содержать уникальные значения",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0} должен быть цветом",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0} должен быть одним из [{1}]",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("warning: error translating FieldError: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "image",
			translation: "{0} должно быть допустимым изображением",
			override:    false,
		},
	}

	for _, t := range translations {

		if t.customTransFunc != nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, t.customTransFunc)

		} else if t.customTransFunc != nil && t.customRegisFunc == nil {

			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), t.customTransFunc)

		} else if t.customTransFunc == nil && t.customRegisFunc != nil {

			err = v.RegisterTranslation(t.tag, trans, t.customRegisFunc, translateFunc)

		} else {
			err = v.RegisterTranslation(t.tag, trans, registrationFunc(t.tag, t.translation, t.override), translateFunc)
		}

		if err != nil {
			return
		}
	}

	return
}

func registrationFunc(tag string, translation string, override bool) validator.RegisterTranslationsFunc {

	return func(ut ut.Translator) (err error) {

		if err = ut.Add(tag, translation, override); err != nil {
			return
		}

		return

	}

}

func translateFunc(ut ut.Translator, fe validator.FieldError) string {

	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		log.Printf("warning: error translating FieldError: %#v", fe)
		return fe.(error).Error()
	}

	return t
}
