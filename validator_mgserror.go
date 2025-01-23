package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Help valid struct with tag custom
const tagCustom = "errormgs"

func errorTagFunc[T interface{}](obj interface{}, snp string, fieldname, actualTag string) error {
	o := obj.(T)

	if !strings.Contains(snp, fieldname) {
		return nil
	}

	fieldArr := strings.Split(snp, ".")
	rsf := reflect.TypeOf(&o).Elem()

	for i := 1; i < len(fieldArr); i++ {
		field, found := rsf.FieldByName(fieldArr[i])
		if found {
			if fieldArr[i] == fieldname {
				customMessage := field.Tag.Get(tagCustom)
				if customMessage != "" {
					return fmt.Errorf("%s: %s (%s)", fieldname, customMessage, actualTag)
				}
				return nil
			} else {
				nestedFieldType := field.Type
				rsf = nestedFieldType
			}
		}
	}
	return nil
}

// ValidateFunc validates (error message with tag) the given object obj against the struct tags defined using
// the "validator v10" package. The type of the object is defined by the generic type T.
// It returns an error indicating the validation failure(s), if any.
// If a panic occurs during the validation process, the function recovers and returns
// an error indicating the panic. The function also uses the errorTagFunc function to
// generate error messages with custom messages defined in the struct tags.
// If no validation errors or panics occur, the function returns nil.
// https://dev.to/thanhphuchuynh/customizing-error-messages-in-struct-validation-using-tags-in-go-4k0j
func ValidateFunc[T interface{}](obj interface{}, validate *Validate) (errs error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in Validate:", r)
			errs = fmt.Errorf("can't validate %+v", r)
		}
	}()

	o := obj.(T)

	if err := validate.Struct(o); err != nil {
		errorValid := err.(ValidationErrors)
		for _, e := range errorValid {
			// snp  X.Y.Z
			snp := e.StructNamespace()
			errmgs := errorTagFunc[T](obj, snp, e.Field(), e.ActualTag())
			if errmgs != nil {
				errs = errors.Join(errs, fmt.Errorf("%w", errmgs))
			} else {
				errs = errors.Join(errs, fmt.Errorf("%w", e))
			}
		}
	}

	if errs != nil {
		return errs
	}

	return nil
}
