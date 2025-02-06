package de

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
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "required_if",
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "required_unless",
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "required_with",
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "required_with_all",
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "required_without",
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "required_without_all",
			translation: "{0} ist ein Pflichtfeld",
			override:    false,
		},
		{
			tag:         "excluded_if",
			translation: "{0} ist ein ausgeschlossenes Feld",
			override:    false,
		},
		{
			tag:         "excluded_unless",
			translation: "{0} ist ein ausgeschlossenes Feld",
			override:    false,
		},
		{
			tag:         "excluded_with",
			translation: "{0} ist ein ausgeschlossenes Feld",
			override:    false,
		},
		{
			tag:         "excluded_with_all",
			translation: "{0} ist ein ausgeschlossenes Feld",
			override:    false,
		},
		{
			tag:         "excluded_without",
			translation: "{0} ist ein ausgeschlossenes Feld",
			override:    false,
		},
		{
			tag:         "excluded_without_all",
			translation: "{0} ist ein ausgeschlossenes Feld",
			override:    false,
		},
		{
			tag:         "isdefault",
			translation: "{0} muss der Standardwert sein",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("len-string", "{0} darf nur {1} lang sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0} muss gleich {1} sein", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0} muss {1} enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
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
				if err = ut.Add("min-string", "{0} muss mindestens {1} lang sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0} muss {1} oder größer sein", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} muss mindestens {1} enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
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
				if err = ut.Add("max-string", "{0} darf maximal {1} lang sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0} darf {1} oder weniger sein", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} darf maximal {1} enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
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
			translation: "{0} ist nicht gleich {1}",
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
			translation: "{0} darf nicht gleich {1} sein",
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
				if err = ut.Add("lt-string", "{0} muss kleiner als {1} sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0} muss kleiner als {1} sein", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0} muss {1} oder weniger enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} muss vor dem aktuellen Datum und Uhrzeit liegen", false); err != nil {
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
				if err = ut.Add("lte-string", "{0} darf maximal {1} lang sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0} darf {1} oder weniger sein", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0} darf maximal {1} enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} muss vor dem aktuellen Datum und Uhrzeit liegen oder gleich sein", false); err != nil {
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
				if err = ut.Add("gt-string", "{0} muss größer als {1} sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0} muss größer als {1} sein", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0} muss {1} oder mehr enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} muss vor dem aktuellen Datum und Uhrzeit liegen", false); err != nil {
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
				if err = ut.Add("gte-string", "{0} muss mindestens {1} lang sein", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} Zeichen", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} Zeichen", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0} muss {1} oder größer sein", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} muss mindestens {1} enthalten", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} Element", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} Elemente", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} muss vor dem aktuellen Datum und Uhrzeit liegen oder gleich sein", false); err != nil {
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
			translation: "{0} muss gleich {1} sein",
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
			translation: "{0} muss gleich {1} sein",
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
			translation: "{0} darf nicht gleich {1} sein",
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
			translation: "{0} muss größer als {1} sein",
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
			translation: "{0} muss größer als oder gleich {1} sein",
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
			translation: "{0} muss kleiner als {1} sein",
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
			translation: "{0} muss kleiner als oder gleich {1} sein",
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
			translation: "{0} darf nicht gleich {1} sein",
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
			translation: "{0} muss größer als {1} sein",
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
			translation: "{0} muss größer als oder gleich {1} sein",
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
			translation: "{0} muss kleiner als {1} sein",
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
			translation: "{0} muss kleiner als oder gleich {1} sein",
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
			translation: "{0} darf nur alphabetische Zeichen enthalten",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0} darf nur alphanumerische Zeichen enthalten",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0} muss eine gültige Zahl sein",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0} muss eine gültige Zahl sein",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0} muss eine gültige hexadezimale Zahl sein",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0} muss eine gültige Hexadezimalfarbe sein",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0} muss eine gültige RGB-Farbe sein",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0} muss eine gültige RGBA-Farbe sein",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0} muss eine gültige HSL-Farbe sein",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0} muss eine gültige HSLA-Farbe sein",
			override:    false,
		},
		{
			tag:         "e164",
			translation: "{0} muss eine gültige E.164-Telefonnummer sein",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0} muss eine gültige E-Mail-Adresse sein",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0} muss eine gültige URL sein",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0} muss eine gültige URI sein",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0} muss eine gültige Base64-Zeichenkette sein",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0} muss den Text '{1}' enthalten",
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
			translation: "{0} muss mindestens eines der folgenden Zeichen enthalten: '{1}'",
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
			translation: "{0} darf den Text '{1}' nicht enthalten",
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
			translation: "{0} darf keines der folgenden Zeichen enthalten: '{1}'",
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
			translation: "{0} darf die folgenden Runen nicht enthalten: '{1}'",
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
			translation: "{0} muss eine gültige ISBN-Nummer sein",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0} muss eine gültige ISBN-10-Nummer sein",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0} muss eine gültige ISBN-13-Nummer sein",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "{0} muss eine gültige ISSN-Nummer sein",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0} muss eine gültige UUID sein",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0} muss eine gültige Version 3 UUID sein",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0} muss eine gültige Version 4 UUID sein",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0} muss eine gültige Version 5 UUID sein",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "{0} muss eine gültige ULID sein",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0} darf nur ASCII-Zeichen enthalten",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0} darf nur druckbare ASCII-Zeichen enthalten",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0} darf nur Mehrbyte-Zeichen enthalten",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0} muss eine gültige Data-URI sein",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0} muss gültige Breitengradkoordinaten enthalten",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0} muss gültige Längengradkoordinaten enthalten",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0} muss eine gültige SSN-Nummer sein",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0} muss eine gültige IPv4-Adresse sein",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0} muss eine gültige IPv6-Adresse sein",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0} muss eine gültige IP-Adresse sein",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0} muss eine gültige CIDR-Notation enthalten",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0} muss eine gültige CIDR-Notation für eine IPv4-Adresse enthalten",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0} muss eine gültige CIDR-Notation für eine IPv6-Adresse enthalten",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0} muss eine gültige TCP-Adresse sein",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} muss eine gültige IPv4-TCP-Adresse sein",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} muss eine gültige IPv6-TCP-Adresse sein",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0} muss eine gültige UDP-Adresse sein",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0} muss eine gültige IPv4-UDP-Adresse sein",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0} muss eine gültige IPv6-UDP-Adresse sein",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0} muss eine auflösbare IP-Adresse sein",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0} muss eine auflösbare IPv4-Adresse sein",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0} muss eine auflösbare IPv6-Adresse sein",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0} muss eine auflösbare UNIX-Adresse sein",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0} muss eine gültige MAC-Adresse sein",
			override:    false,
		},
		{
			tag:         "fqdn",
			translation: "{0} muss eine gültige FQDN sein",
			override:    false,
		},
		{
			tag:         "unique",
			translation: "{0} darf nur einmal vorkommen",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0} muss eine gültige Farbe sein",
			override:    false,
		},
		{
			tag:         "cron",
			translation: "{0} muss eine gültige Cron-Ausdruck sein",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0} muss einer der folgenden sein: [{1}]",
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
			tag:         "json",
			translation: "{0} muss eine gültige JSON-Zeichenkette sein",
			override:    false,
		},
		{
			tag:         "jwt",
			translation: "{0} muss eine gültige JWT-Zeichenkette sein",
			override:    false,
		},
		{
			tag:         "lowercase",
			translation: "{0} darf nur Kleinbuchstaben enthalten",
			override:    false,
		},
		{
			tag:         "uppercase",
			translation: "{0} darf nur Großbuchstaben enthalten",
			override:    false,
		},
		{
			tag:         "datetime",
			translation: "{0} entspricht nicht dem {1}-Format",
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
			tag:         "postcode_iso3166_alpha2",
			translation: "{0} entspricht nicht dem Postleitzahlformat von {1}",
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
			tag:         "postcode_iso3166_alpha2_field",
			translation: "{0} entspricht nicht dem Postleitzahlformat des Feldes {1}",
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
			tag:         "boolean",
			translation: "{0} muss eine gültige Booleanwert sein",
			override:    false,
		},
		{
			tag:         "image",
			translation: "{0} muss ein Bild sein",
			override:    false,
		},
		{
			tag:         "cve",
			translation: "{0} muss eine gültige CVE-Kennung sein",
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
