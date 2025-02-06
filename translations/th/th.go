package th

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
			translation: "โปรดระบุ {0}",
			override:    false,
		},
		{
			tag:         "required_if",
			translation: "โปรดระบุ {0}",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("len-string", "{0} ต้องมีความยาว {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0} ต้องเท่ากับ {1}", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0} must contain {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("len-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
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
				if err = ut.Add("min-string", "{0} ต้องมีความยาวอย่างน้อย {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0} ต้องมีค่ามากกว่า {1}", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0} ต้องมีอย่างน้อย {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("min-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
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
				if err = ut.Add("max-string", "{0} ต้องมีความยาวไม่เกิน {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0} ต้องมีค่าน้อยกว่าหรือเท่ากับ {1}", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0} ต้องมีไม่เกิน {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("max-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
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
			translation: "{0} ไม่เท่ากับ {1}",
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
			translation: "{0} ต้องไม่เท่ากับ {1}",
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
				if err = ut.Add("lt-string", "{0} ต้องมีความยาวน้อยกว่า {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0} ต้องมีค่าน้อยกว่า {1}", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0} ต้องมีน้อยกว่า {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lt-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0} ต้องมีค่าน้อยกว่าเวลาปัจจุบัน", false); err != nil {
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
				if err = ut.Add("lte-string", "{0} ต้องมีความยาวไม่เกิน {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0} ต้องมีค่าน้อยกว่าหรือเท่ากับ {1}", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0} ต้องมีไม่เกิน {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("lte-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0} ต้องมีค่าน้อยกว่าหรือเท่ากับเวลาปัจจุบัน", false); err != nil {
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
				if err = ut.Add("gt-string", "{0} ต้องมีความยาวมากกว่า {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0} ต้องมีค่ามากกว่า {1}", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0} ต้องมีมากกว่า {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gt-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0} ต้องมีค่ามากกว่าเวลาปัจจุบัน", false); err != nil {
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
				if err = ut.Add("gte-string", "{0} ต้องมีความยาวอย่างน้อย {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-string-character", "{0} ตัวอักษร", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0} ต้องมีค่ามากกว่า {1}", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0} ต้องมีอย่างน้อย {1}", false); err != nil {
					return
				}

				if err = ut.AddCardinal("gte-items-item", "{0} ไอเทม", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0} ต้องมีค่ามากกว่าหรือเท่ากับเวลาปัจจุบัน", false); err != nil {
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
			translation: "{0} ต้องเท่ากับ {1}",
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
			translation: "{0} ต้องเท่ากับ {1}",
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
			translation: "{0} ต้องไม่เท่ากับ {1}",
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
			translation: "{0} ต้องมีค่ามากกว่า {1}",
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
			translation: "{0} ต้องมีค่ามากกว่าหรือเท่ากับ {1}",
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
			translation: "{0} ต้องมีค่าน้อยกว่า {1}",
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
			translation: "{0} ต้องมีค่าน้อยกว่าหรือเท่ากับ {1}",
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
			translation: "{0} ต้องไม่เท่ากับ {1}",
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
			translation: "{0} ต้องมีค่ามากกว่า {1}",
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
			translation: "{0} ต้องมีค่ามากกว่าหรือเท่ากับ {1}",
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
			translation: "{0} ต้องมีค่าน้อยกว่า {1}",
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
			translation: "{0} ต้องมีค่าน้อยกว่าหรือเท่ากับ {1}",
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
			translation: "{0} ต้องเป็นตัวอักษรเท่านั้น",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0} ต้องเป็นตัวอักษรหรือตัวเลขเท่านั้น",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0} ต้องเป็นค่าตัวเลขเท่านั้น",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0} ต้องเป็นตัวเลขเท่านั้น",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0} ต้องเป็นค่าตัวเลขฐาน 16 เท่านั้น",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0} ต้องเป็นเลขสีฐาน 16 เท่านั้น",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0} ต้องเป็นเลขสี RGB เท่านั้น",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0} ต้องเป็นเลขสี RGBA เท่านั้น",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0} ต้องเป็นเลขสี HSL เท่านั้น",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0} ต้องเป็นเลขสี HSLA เท่านั้น",
			override:    false,
		},
		{
			tag:         "e164",
			translation: "{0} ต้องเป็นเบอร์โทรศัพท์รูปแบบ E.164 เท่านั้น",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0} ต้องเป็นอีเมลเท่านั้น",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0} ต้องเป็น URL เท่านั้น",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0} ต้องเป็น URI เท่านั้น",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0} ต้องเป็น Base64 เท่านั้น",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0} ต้องมี '{1}'",
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
			translation: "{0} ต้องมีอย่างน้อยอักขระใน '{1}' อย่างน้อย 1 ตัว",
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
			translation: "{0} ต้องไม่มี '{1}'",
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
			translation: "{0} ต้องไม่มีอักขระ '{1}' ทั้งหมด",
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
			translation: "{0} ต้องไม่มี '{1}'",
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
			translation: "{0} ต้องเป็นตัวเลข ISBN",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0} ต้องเป็นตัวเลข ISBN-10",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0} ต้องเป็นตัวเลข ISBN-13",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "{0} ต้องเป็นตัวเลข ISSN",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0} ต้องเป็น UUID",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0} ต้องเป็น version 3 UUID",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0} ต้องเป็น version 4 UUID",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0} ต้องเป็น version 5 UUID",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "{0} must be a valid ULID",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0} ต้องมีแค่อักขระ ascii",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0} ต้องมีแค่อักขระ ascii ที่แสดงผลได้",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0} ต้องมี multibyte character",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0} must contain a valid Data URI",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0} ต้องเป็นละติจูด",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0} ต้องเป็นลองจิจูด",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0} ต้องเป็นตัวเลข SSN",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0} ต้องเป็น IPv4 address",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0} ต้องเป็น IPv6 address",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0} ต้องเป็น IP address",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0} ต้องเป็น CIDR notation",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0} ต้องเป็น CIDR notation สำหรับ an IPv4 address",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0} ต้องเป็น CIDR notation สำหรับ an IPv6 address",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0} ต้องเป็น TCP address",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0} ต้องเป็น IPv4 TCP address",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0} ต้องเป็น IPv6 TCP address",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0} ต้องเป็น UDP address",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0} ต้องเป็น IPv4 UDP address",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0} ต้องเป็น IPv6 UDP address",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0} ต้องเป็น IP address ที่เข้าถึงได้",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0} ต้องเป็น IPv4 address ที่เข้าถึงได้",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0} ต้องเป็น IPv6 address ที่เข้าถึงได้",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0} ต้องเป็น UNIX address ที่เข้าถึงได้",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0} ต้องเป็น MAC address",
			override:    false,
		},
		{
			tag:         "fqdn",
			translation: "{0} ต้องเป็น FQDN",
			override:    false,
		},
		{
			tag:         "unique",
			translation: "{0} ต้องมีข้อมูลไม่ซ้ำ",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0} ต้องเป็นเลขสี",
			override:    false,
		},
		{
			tag:         "cron",
			translation: "{0} ต้องเป็น cron expression",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0} ต้องอยู่ใน [{1}]",
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
			translation: "{0} ต้องเป็น json string",
			override:    false,
		},
		{
			tag:         "jwt",
			translation: "{0} ต้องเป็น jwt string",
			override:    false,
		},
		{
			tag:         "lowercase",
			translation: "{0} ต้องเป็นตัวพิมพ์เล็ก",
			override:    false,
		},
		{
			tag:         "uppercase",
			translation: "{0} ต้องเป็นตัวพิมพ์ใหญ่",
			override:    false,
		},
		{
			tag:         "datetime",
			translation: "{0} ไม่ตรงกับรูปแบบ {1}",
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
			translation: "{0} does not match postcode format of {1} country",
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
			translation: "{0} does not match postcode format of country in {1} field",
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
			translation: "{0} ต้องเป็น boolean",
			override:    false,
		},
		{
			tag:         "image",
			translation: "{0} ต้องเป็นรูปภาพ",
			override:    false,
		},
		{
			tag:         "cve",
			translation: "{0} ต้องเป็นรูปแบบ cve",
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
