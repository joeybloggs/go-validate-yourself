package tr

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
			translation: "{0} zorunlu bir alandır",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {

				if err = ut.Add("len-string", "{0} uzunluğu {1} olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0}, {1} değerine eşit olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0}, {1} içermelidir", false); err != nil {
					return
				}
				if err = ut.AddCardinal("len-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
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

				if err = ut.Add("min-string", "{0} en az {1} uzunluğunda olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0}, {1} veya daha büyük olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} en az {1} içermelidir", false); err != nil {
					return
				}
				if err = ut.AddCardinal("min-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
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

				if err = ut.Add("max-string", "{0} uzunluğu en fazla {1} olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0}, {1} veya daha az olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} maksimum {1} içermelidir", false); err != nil {
					return
				}
				if err = ut.AddCardinal("max-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
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
			translation: "{0}, {1} değerine eşit değil",
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
			translation: "{0}, {1} değerine eşit olmamalıdır",
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

				if err = ut.Add("lt-string", "{0}, {1} uzunluğundan daha az olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0}, {1} değerinden küçük olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0}, {1}den daha az içermelidir", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} geçerli Tarih ve Saatten daha az olmalıdır", false); err != nil {
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

				if err = ut.Add("lte-string", "{0} en fazla {1} uzunluğunda olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0}, {1} veya daha az olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0}, maksimum {1} içermelidir", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} geçerli Tarih ve Saate eşit veya daha küçük olmalıdır", false); err != nil {
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

				if err = ut.Add("gt-string", "{0}, {1} uzunluğundan fazla olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0}, {1} değerinden büyük olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0}, {1}den daha fazla içermelidir", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} geçerli Tarih ve Saatten büyük olmalıdır", false); err != nil {
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

				if err = ut.Add("gte-string", "{0} en az {1} uzunluğunda olmalıdır", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} karakter", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} karakter", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0}, {1} veya daha büyük olmalıdır", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} en az {1} içermelidir", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} öğe", locales.PluralRuleOne, false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} öğe", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} geçerli Tarih ve Saatten büyük veya ona eşit olmalıdır", false); err != nil {
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
			translation: "{0}, {1} değerine eşit olmalıdır",
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
			translation: "{0}, {1} değerine eşit olmalıdır",
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
			translation: "{0}, {1} değerine eşit olmamalıdır",
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
			translation: "{0}, {1} değerinden büyük olmalıdır",
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
			translation: "{0}, {1} değerinden küçük veya ona eşit olmalıdır",
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
			translation: "{0}, {1} değerinden küçük olmalıdır",
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
			translation: "{0}, {1} değerinden küçük veya ona eşit olmalıdır",
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
			translation: "{0}, {1} değerine eşit olmamalıdır",
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
			translation: "{0}, {1} değerinden büyük olmalıdır",
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
			translation: "{0}, {1} değerinden büyük veya ona eşit olmalıdır",
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
			translation: "{0}, {1} değerinden küçük olmalıdır",
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
			translation: "{0}, {1} değerinden küçük veya ona eşit olmalıdır",
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
			translation: "{0} yalnızca alfabetik karakterler içerebilir",
			override:    false,
		},
		{
			tag:         "alphaunicode",
			translation: "{0} yalnızca alfabetik unicode karakterler içerebilir",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0} yalnızca alfanümerik karakterler içerebilir",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0} geçerli bir sayısal değer olmalıdır",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0} geçerli bir sayı olmalıdır",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0} geçerli bir onaltılık olmalıdır",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0} geçerli bir HEX rengi olmalıdır",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0} geçerli bir RGB rengi olmalıdır",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0} geçerli bir RGBA rengi olmalıdır",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0} geçerli bir HSL rengi olmalıdır",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0} geçerli bir HSLA rengi olmalıdır",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0} geçerli bir e-posta adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0} geçerli bir URL olmalıdır",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0} geçerli bir URI olmalıdır",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0} geçerli bir Base64 karakter dizesi olmalıdır",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0}, '{1}' metnini içermelidir",
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
			translation: "{0}, '{1}' karakterlerinden en az birini içermelidir",
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
			translation: "{0}, '{1}' metnini içeremez",
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
			translation: "{0}, '{1}' karakterlerinden hiçbirini içeremez",
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
			translation: "{0}, '{1}' ifadesini içeremez",
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
			translation: "{0} geçerli bir ISBN numarası olmalıdır",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0} geçerli bir ISBN-10 numarası olmalıdır",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0} geçerli bir ISBN-13 numarası olmalıdır",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "{0} geçerli bir ISSN numarası olmalıdır",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0} geçerli bir UUID olmalıdır",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0} geçerli bir sürüm 3 UUID olmalıdır",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0} geçerli bir sürüm 4 UUID olmalıdır",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0} geçerli bir sürüm 5 UUID olmalıdır",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "{0} geçerli bir ULID olmalıdır",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0} yalnızca ascii karakterler içermelidir",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0} yalnızca yazdırılabilir ascii karakterleri içermelidir",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0} çok baytlı karakterler içermelidir",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0} geçerli bir Veri URI içermelidir",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0} geçerli bir enlem koordinatı içermelidir",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0} geçerli bir boylam koordinatı içermelidir",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0} geçerli bir SSN numarası olmalıdır",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0} geçerli bir IPv4 adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0} geçerli bir IPv6 adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0} geçerli bir IP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0} geçerli bir CIDR gösterimi içermelidir",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0} bir IPv4 adresi için geçerli bir CIDR gösterimi içermelidir",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0} bir IPv6 adresi için geçerli bir CIDR gösterimi içermelidir",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0} geçerli bir TCP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} geçerli bir IPv4 TCP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} geçerli bir IPv6 TCP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0} geçerli bir UDP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0} geçerli bir IPv4 UDP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0} geçerli bir IPv6 UDP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0} çözülebilir bir IP adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0} çözülebilir bir IPv4 adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0} çözülebilir bir IPv6 adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0} çözülebilir bir UNIX adresi olmalıdır",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0} geçerli bir MAC adresi içermelidir",
			override:    false,
		},
		{
			tag:         "unique",
			translation: "{0} benzersiz değerler içermelidir",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0} geçerli bir renk olmalıdır",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0}, [{1}]'dan biri olmalıdır",
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
			translation: "{0} geçerli bir resim olmalıdır",
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
