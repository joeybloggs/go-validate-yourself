package ko

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
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "required_if",
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "required_unless",
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "required_with",
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "required_with_all",
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "required_without",
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "required_without_all",
			translation: "{0}은(는) 필수 필드입니다.",
			override:    false,
		},
		{
			tag:         "excluded_if",
			translation: "{0}은(는) 제외된 필드입니다.",
			override:    false,
		},
		{
			tag:         "excluded_unless",
			translation: "{0}은(는) 제외된 필드입니다.",
			override:    false,
		},
		{
			tag:         "excluded_with",
			translation: "{0}은(는) 제외된 필드입니다.",
			override:    false,
		},
		{
			tag:         "excluded_with_all",
			translation: "{0}은(는) 제외된 필드입니다.",
			override:    false,
		},
		{
			tag:         "excluded_without",
			translation: "{0}은(는) 제외된 필드입니다.",
			override:    false,
		},
		{
			tag:         "excluded_without_all",
			translation: "{0}은(는) 제외된 필드입니다.",
			override:    false,
		},
		{
			tag:         "isdefault",
			translation: "{0}은(는) 기본값이어야 합니다.",
			override:    false,
		},
		{
			tag: "len",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("len-string", "{0}의 길이는 {1}여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("len-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("len-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("len-number", "{0}은(는) {1}와(과) 같아야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("len-items", "{0}은(는) {1}을 포함해야 합니다.", false); err != nil {
					return
				}
				// if err = ut.AddCardinal("len-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("len-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
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
				if err = ut.Add("min-string", "{0}의 길이는 최소 {1}여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("min-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("min-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("min-number", "{0}은(는) {1} 이상여야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("min-items", "{0}은(는) 최소 {1}을 포함해야 합니다.", false); err != nil {
					return
				}
				// if err = ut.AddCardinal("min-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("min-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
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
				if err = ut.Add("max-string", "{0}의 길이는 최대 {1}여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("max-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("max-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("max-number", "{0}은(는) {1} 이하여야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("max-items", "{0}은(는) 최대 {1}여야 합니다.", false); err != nil {
					return
				}
				// if err = ut.AddCardinal("max-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("max-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
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
			translation: "{0}은(는) {1}와(과) 같아야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					fmt.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ne",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("ne-items", "{0}의 항목 수는 {1}와(과) 달라야 합니다.", false); err != nil {
					fmt.Printf("ne customRegisFunc #1 error because of %v\n", err)
					return
				}
				// if err = ut.AddCardinal("ne-items-item", "{0}개", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("ne-items-item", "{0}개", locales.PluralRuleOther, false); err != nil {
					return
				}
				if err = ut.Add("ne", "{0}은(는) {1}와(과) 달라야 합니다.", false); err != nil {
					fmt.Printf("ne customRegisFunc #2 error because of %v\n", err)
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
				case reflect.Slice:
					var c string
					err = fn()
					if err != nil {
						goto END
					}

					c, err = ut.C("ne-items-item", f64, digits, ut.FmtNumber(f64, digits))
					if err != nil {
						goto END
					}
					t, err = ut.T("ne-items", fe.Field(), c)
				default:
					t, err = ut.T("ne", fe.Field(), fe.Param())
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
			tag: "lt",
			customRegisFunc: func(ut ut.Translator) (err error) {
				if err = ut.Add("lt-string", "{0}의 길이는 {1}보다 작아야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("lt-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("lt-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-number", "{0}은(는) {1}보다 작아야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("lt-items", "{0}은(는) {1}보다 적은 항목여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("lt-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("lt-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lt-datetime", "{0}은(는) 현재 시간보다 이전이어야 합니다.", false); err != nil {
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
				if err = ut.Add("lte-string", "{0}의 길이는 최대 {1}여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("lte-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("lte-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-number", "{0}은(는) {1} 이하여야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("lte-items", "{0}은(는) 최대 {1}여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("lte-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("lte-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("lte-datetime", "{0}은(는) 현재 시간보다 이전이어야 합니다.", false); err != nil {
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
				if err = ut.Add("gt-string", "{0}의 길이는 {1}보다 길어야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("gt-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("gt-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-number", "{0}은(는) {1}보다 커야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("gt-items", "{0}은(는) {1}보다 많은 항목을 포함해야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("gt-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("gt-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gt-datetime", "{0}은(는) 현재 시간 이후이어야 합니다.", false); err != nil {
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
				if err = ut.Add("gte-string", "{0}의 길이는 최소 {1} 이상여야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("gte-string-character", "{0}자", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("gte-string-character", "{0}자", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-number", "{0}은(는) {1} 이상여야 합니다.", false); err != nil {
					return
				}

				if err = ut.Add("gte-items", "{0}은(는) 최소 {1}을 포함해야 합니다.", false); err != nil {
					return
				}

				// if err = ut.AddCardinal("gte-items-item", "{0}개의 항목", locales.PluralRuleOne, false); err != nil {
				// 	return
				// }

				if err = ut.AddCardinal("gte-items-item", "{0}개의 항목", locales.PluralRuleOther, false); err != nil {
					return
				}

				if err = ut.Add("gte-datetime", "{0}은(는) 현재 시간 이후이어야 합니다.", false); err != nil {
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
			translation: "{0}은(는) {1}와(과) 같아야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "eqcsfield",
			translation: "{0}은(는) {1}와(과) 같아야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "necsfield",
			translation: "{0}은(는) {1}와(과) 달라야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtcsfield",
			translation: "{0}은(는) {1}보다 커야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtecsfield",
			translation: "{0}은(는) {1} 이상여야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltcsfield",
			translation: "{0}은(는) {1}보다 작아야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltecsfield",
			translation: "{0}은(는) {1} 이하여야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "nefield",
			translation: "{0}은(는) {1}와(과) 달라야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtfield",
			translation: "{0}은(는) {1}보다 커야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "gtefield",
			translation: "{0}은(는) {1} 이상여야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltfield",
			translation: "{0}은(는) {1}보다 작아야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "ltefield",
			translation: "{0}은(는) {1} 이하여야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "alpha",
			translation: "{0}은(는) 알파벳만 포함할 수 있습니다.",
			override:    false,
		},
		{
			tag:         "alphanum",
			translation: "{0}은(는) 알파벳과 숫자만 포함할 수 있습니다.",
			override:    false,
		},
		{
			tag:         "numeric",
			translation: "{0}은(는) 올바른 숫자여야 합니다.",
			override:    false,
		},
		{
			tag:         "number",
			translation: "{0}은(는) 올바른 수여야 합니다.",
			override:    false,
		},
		{
			tag:         "hexadecimal",
			translation: "{0}은(는) 올바른 16진수 표기여야 합니다.",
			override:    false,
		},
		{
			tag:         "hexcolor",
			translation: "{0}은(는) 올바른 HEX 색상 코드여야 합니다.",
			override:    false,
		},
		{
			tag:         "rgb",
			translation: "{0}은(는) 올바른 RGB 색상 코드여야 합니다.",
			override:    false,
		},
		{
			tag:         "rgba",
			translation: "{0}은(는) 올바른 RGBA 색상 코드여야 합니다.",
			override:    false,
		},
		{
			tag:         "hsl",
			translation: "{0}은(는) 올바른 HSL 색상 코드여야 합니다.",
			override:    false,
		},
		{
			tag:         "hsla",
			translation: "{0}은(는) 올바른 HSLA 색상 코드여야 합니다.",
			override:    false,
		},
		{
			tag:         "e164",
			translation: "{0}은(는) 유효한 E.164 형식의 전화번호여야 합니다.",
			override:    false,
		},
		{
			tag:         "email",
			translation: "{0}은(는) 올바른 이메일 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "url",
			translation: "{0}은(는) 올바른 URL여야 합니다.",
			override:    false,
		},
		{
			tag:         "uri",
			translation: "{0}은(는) 올바른 URI여야 합니다.",
			override:    false,
		},
		{
			tag:         "base64",
			translation: "{0}은(는) 올바른 Base64 문자열여야 합니다.",
			override:    false,
		},
		{
			tag:         "contains",
			translation: "{0}은(는) '{1}'을(를) 포함해야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "containsany",
			translation: "{0}은(는) '{1}' 중 최소 하나를 포함해야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludes",
			translation: "{0}에는 '{1}'라는 텍스트를 포함할 수 없습니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesall",
			translation: "{0}에는 '{1}' 중 어느 것도 포함할 수 없습니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "excludesrune",
			translation: "{0}에는 '{1}'을(를) 포함할 수 없습니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "isbn",
			translation: "{0}은(는) 올바른 ISBN 번호여야 합니다.",
			override:    false,
		},
		{
			tag:         "isbn10",
			translation: "{0}은(는) 올바른 ISBN-10 번호여야 합니다.",
			override:    false,
		},
		{
			tag:         "isbn13",
			translation: "{0}은(는) 올바른 ISBN-13 번호여야 합니다.",
			override:    false,
		},
		{
			tag:         "issn",
			translation: "{0}은(는) 올바른 ISSN 번호여야 합니다.",
			override:    false,
		},
		{
			tag:         "uuid",
			translation: "{0}은(는) 올바른 UUID여야 합니다.",
			override:    false,
		},
		{
			tag:         "uuid3",
			translation: "{0}은(는) 버전 3의 올바른 UUID여야 합니다.",
			override:    false,
		},
		{
			tag:         "uuid4",
			translation: "{0}은(는) 버전 4의 올바른 UUID여야 합니다.",
			override:    false,
		},
		{
			tag:         "uuid5",
			translation: "{0}은(는) 버전 5의 올바른 UUID여야 합니다.",
			override:    false,
		},
		{
			tag:         "ulid",
			translation: "{0}은(는) 올바른 ULID여야 합니다.",
			override:    false,
		},
		{
			tag:         "ascii",
			translation: "{0}은(는) ASCII 문자만 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "printascii",
			translation: "{0}은(는) 인쇄 가능한 ASCII 문자만 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "multibyte",
			translation: "{0}은(는) 멀티바이트 문자를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "datauri",
			translation: "{0}은(는) 올바른 데이터 URI를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "latitude",
			translation: "{0}은(는) 올바른 위도 좌표를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "longitude",
			translation: "{0}은(는) 올바른 경도 좌표를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "ssn",
			translation: "{0}은(는) 올바른 사회 보장 번호여야 합니다.",
			override:    false,
		},
		{
			tag:         "ipv4",
			translation: "{0}은(는) 올바른 IPv4 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "ipv6",
			translation: "{0}은(는) 올바른 IPv6 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "ip",
			translation: "{0}은(는) 올바른 IP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "cidr",
			translation: "{0}은(는) 올바른 CIDR 표기를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "cidrv4",
			translation: "{0}은(는) IPv4 주소의 올바른 CIDR 표기를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "cidrv6",
			translation: "{0}은(는) IPv6 주소의 올바른 CIDR 표기를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "tcp_addr",
			translation: "{0}은(는) 올바른 TCP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "tcp4_addr",
			translation: "{0}은(는) 올바른 IPv4의 TCP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "tcp6_addr",
			translation: "{0}은(는) 올바른 IPv6의 TCP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "udp_addr",
			translation: "{0}은(는) 올바른 UDP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "udp4_addr",
			translation: "{0}은(는) 올바른 IPv4의 UDP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "udp6_addr",
			translation: "{0}은(는) 올바른 IPv6의 UDP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "ip_addr",
			translation: "{0}은(는) 해석 가능한 IP 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "ip4_addr",
			translation: "{0}은(는) 해석 가능한 IPv4 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "ip6_addr",
			translation: "{0}은(는) 해석 가능한 IPv6 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "unix_addr",
			translation: "{0}은(는) 해석 가능한 UNIX 주소여야 합니다.",
			override:    false,
		},
		{
			tag:         "mac",
			translation: "{0}은(는) 올바른 MAC 주소를 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "fqdn",
			translation: "{0}은(는) 유효한 FQDN이어야 합니다.",
			override:    false,
		},
		{
			tag:         "unique",
			translation: "{0}은(는) 고유한 값만 포함해야 합니다.",
			override:    false,
		},
		{
			tag:         "iscolor",
			translation: "{0}은(는) 올바른 색이여야 합니다.",
			override:    false,
		},
		{
			tag:         "cron",
			translation: "{0}은(는) 유효한 cron 표현식이어야 합니다.",
			override:    false,
		},
		{
			tag:         "oneof",
			translation: "{0}은(는) [{1}] 중 하나여야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				s, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}
				return s
			},
		},
		{
			tag:         "json",
			translation: "{0}은(는) 올바른 JSON 문자열여야 합니다.",
			override:    false,
		},
		{
			tag:         "jwt",
			translation: "{0}은(는) 올바른 JWT 문자열여야 합니다.",
			override:    false,
		},
		{
			tag:         "lowercase",
			translation: "{0}은(는) 소문자여야 합니다.",
			override:    false,
		},
		{
			tag:         "uppercase",
			translation: "{0}은(는) 대문자여야 합니다.",
			override:    false,
		},
		{
			tag:         "datetime",
			translation: "{0}은(는) {1} 형식과 일치해야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "postcode_iso3166_alpha2",
			translation: "{0}은(는) 국가 코드 {1}의 우편번호 형식과 일치해야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "postcode_iso3166_alpha2_field",
			translation: "{0}은(는) {1} 필드에 지정된 국가 코드의 우편번호 형식과 일치해야 합니다.",
			override:    false,
			customTransFunc: func(ut ut.Translator, fe validator.FieldError) string {
				t, err := ut.T(fe.Tag(), fe.Field(), fe.Param())
				if err != nil {
					log.Printf("경고: FieldError 번역 중 오류 발생: %#v", fe)
					return fe.(error).Error()
				}

				return t
			},
		},
		{
			tag:         "boolean",
			translation: "{0}은(는) 올바른 부울 값여야 합니다.",
			override:    false,
		},
		{
			tag:         "image",
			translation: "{0}은(는) 유효한 이미지여야 합니다.",
			override:    false,
		},
		{
			tag:         "cve",
			translation: "{0}은(는) 유효한 CVE 식별자여야 합니다.",
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
