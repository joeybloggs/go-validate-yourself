package es

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	spanish "github.com/go-playground/locales/es"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestTranslations(t *testing.T) {

	spa := spanish.New()
	uni := ut.New(spa, spa)
	trans, _ := uni.GetTranslator("es")

	validate := validator.New()

	err := RegisterDefaultTranslations(validate, trans)
	Equal(t, err, nil)

	type Inner struct {
		EqCSFieldString      string
		NeCSFieldString      string
		GtCSFieldString      string
		GteCSFieldString     string
		LtCSFieldString      string
		LteCSFieldString     string
		RequiredWithoutField string
	}

	type Test struct {
		Inner             Inner
		RequiredString    string            `validate:"required"`
		RequiredNumber    int               `validate:"required"`
		RequiredMultiple  []string          `validate:"required"`
		RequiredIf        string            `validate:"required_if=Inner.RequiredIf abcd"`
		RequiredWithout   string            `validate:"required_without=Inner.RequiredWithoutField"`
		LenString         string            `validate:"len=1"`
		LenNumber         float64           `validate:"len=1113.00"`
		LenMultiple       []string          `validate:"len=7"`
		MinString         string            `validate:"min=1"`
		MinNumber         float64           `validate:"min=1113.00"`
		MinMultiple       []string          `validate:"min=7"`
		MaxString         string            `validate:"max=3"`
		MaxNumber         float64           `validate:"max=1113.00"`
		MaxMultiple       []string          `validate:"max=7"`
		EqString          string            `validate:"eq=3"`
		EqNumber          float64           `validate:"eq=2.33"`
		EqMultiple        []string          `validate:"eq=7"`
		NeString          string            `validate:"ne="`
		NeNumber          float64           `validate:"ne=0.00"`
		NeMultiple        []string          `validate:"ne=0"`
		LtString          string            `validate:"lt=3"`
		LtNumber          float64           `validate:"lt=5.56"`
		LtMultiple        []string          `validate:"lt=2"`
		LtTime            time.Time         `validate:"lt"`
		LteString         string            `validate:"lte=3"`
		LteNumber         float64           `validate:"lte=5.56"`
		LteMultiple       []string          `validate:"lte=2"`
		LteTime           time.Time         `validate:"lte"`
		GtString          string            `validate:"gt=3"`
		GtNumber          float64           `validate:"gt=5.56"`
		GtMultiple        []string          `validate:"gt=2"`
		GtTime            time.Time         `validate:"gt"`
		GteString         string            `validate:"gte=3"`
		GteNumber         float64           `validate:"gte=5.56"`
		GteMultiple       []string          `validate:"gte=2"`
		GteTime           time.Time         `validate:"gte"`
		EqFieldString     string            `validate:"eqfield=MaxString"`
		EqCSFieldString   string            `validate:"eqcsfield=Inner.EqCSFieldString"`
		NeCSFieldString   string            `validate:"necsfield=Inner.NeCSFieldString"`
		GtCSFieldString   string            `validate:"gtcsfield=Inner.GtCSFieldString"`
		GteCSFieldString  string            `validate:"gtecsfield=Inner.GteCSFieldString"`
		LtCSFieldString   string            `validate:"ltcsfield=Inner.LtCSFieldString"`
		LteCSFieldString  string            `validate:"ltecsfield=Inner.LteCSFieldString"`
		NeFieldString     string            `validate:"nefield=EqFieldString"`
		GtFieldString     string            `validate:"gtfield=MaxString"`
		GteFieldString    string            `validate:"gtefield=MaxString"`
		LtFieldString     string            `validate:"ltfield=MaxString"`
		LteFieldString    string            `validate:"ltefield=MaxString"`
		AlphaString       string            `validate:"alpha"`
		AlphanumString    string            `validate:"alphanum"`
		NumericString     string            `validate:"numeric"`
		NumberString      string            `validate:"number"`
		HexadecimalString string            `validate:"hexadecimal"`
		HexColorString    string            `validate:"hexcolor"`
		RGBColorString    string            `validate:"rgb"`
		RGBAColorString   string            `validate:"rgba"`
		HSLColorString    string            `validate:"hsl"`
		HSLAColorString   string            `validate:"hsla"`
		Email             string            `validate:"email"`
		URL               string            `validate:"url"`
		URI               string            `validate:"uri"`
		Base64            string            `validate:"base64"`
		Contains          string            `validate:"contains=purpose"`
		ContainsAny       string            `validate:"containsany=!@#$"`
		Excludes          string            `validate:"excludes=text"`
		ExcludesAll       string            `validate:"excludesall=!@#$"`
		ExcludesRune      string            `validate:"excludesrune=☻"`
		ISBN              string            `validate:"isbn"`
		ISBN10            string            `validate:"isbn10"`
		ISBN13            string            `validate:"isbn13"`
		ISSN              string            `validate:"issn"`
		UUID              string            `validate:"uuid"`
		UUID3             string            `validate:"uuid3"`
		UUID4             string            `validate:"uuid4"`
		UUID5             string            `validate:"uuid5"`
		ULID              string            `validate:"ulid"`
		ASCII             string            `validate:"ascii"`
		PrintableASCII    string            `validate:"printascii"`
		MultiByte         string            `validate:"multibyte"`
		DataURI           string            `validate:"datauri"`
		Latitude          string            `validate:"latitude"`
		Longitude         string            `validate:"longitude"`
		SSN               string            `validate:"ssn"`
		IP                string            `validate:"ip"`
		IPv4              string            `validate:"ipv4"`
		IPv6              string            `validate:"ipv6"`
		CIDR              string            `validate:"cidr"`
		CIDRv4            string            `validate:"cidrv4"`
		CIDRv6            string            `validate:"cidrv6"`
		TCPAddr           string            `validate:"tcp_addr"`
		TCPAddrv4         string            `validate:"tcp4_addr"`
		TCPAddrv6         string            `validate:"tcp6_addr"`
		UDPAddr           string            `validate:"udp_addr"`
		UDPAddrv4         string            `validate:"udp4_addr"`
		UDPAddrv6         string            `validate:"udp6_addr"`
		IPAddr            string            `validate:"ip_addr"`
		IPAddrv4          string            `validate:"ip4_addr"`
		IPAddrv6          string            `validate:"ip6_addr"`
		UinxAddr          string            `validate:"unix_addr"` // can't fail from within Go's net package currently, but maybe in the future
		MAC               string            `validate:"mac"`
		IsColor           string            `validate:"iscolor"`
		StrPtrMinLen      *string           `validate:"min=10"`
		StrPtrMaxLen      *string           `validate:"max=1"`
		StrPtrLen         *string           `validate:"len=2"`
		StrPtrLt          *string           `validate:"lt=1"`
		StrPtrLte         *string           `validate:"lte=1"`
		StrPtrGt          *string           `validate:"gt=10"`
		StrPtrGte         *string           `validate:"gte=10"`
		OneOfString       string            `validate:"oneof=red green"`
		OneOfInt          int               `validate:"oneof=5 63"`
		UniqueSlice       []string          `validate:"unique"`
		UniqueArray       [3]string         `validate:"unique"`
		UniqueMap         map[string]string `validate:"unique"`
		Image			  string			`validate:"image"`
	}

	var test Test

	test.Inner.EqCSFieldString = "1234"
	test.Inner.GtCSFieldString = "1234"
	test.Inner.GteCSFieldString = "1234"

	test.MaxString = "1234"
	test.MaxNumber = 2000
	test.MaxMultiple = make([]string, 9)

	test.LtString = "1234"
	test.LtNumber = 6
	test.LtMultiple = make([]string, 3)
	test.LtTime = time.Now().Add(time.Hour * 24)

	test.LteString = "1234"
	test.LteNumber = 6
	test.LteMultiple = make([]string, 3)
	test.LteTime = time.Now().Add(time.Hour * 24)

	test.LtFieldString = "12345"
	test.LteFieldString = "12345"

	test.LtCSFieldString = "1234"
	test.LteCSFieldString = "1234"

	test.AlphaString = "abc3"
	test.AlphanumString = "abc3!"
	test.NumericString = "12E.00"
	test.NumberString = "12E"

	test.Excludes = "this is some test text"
	test.ExcludesAll = "This is Great!"
	test.ExcludesRune = "Love it ☻"

	test.ASCII = "ｶﾀｶﾅ"
	test.PrintableASCII = "ｶﾀｶﾅ"

	test.MultiByte = "1234feerf"

	s := "toolong"
	test.StrPtrMaxLen = &s
	test.StrPtrLen = &s

	test.UniqueSlice = []string{"1234", "1234"}
	test.UniqueMap = map[string]string{"key1": "1234", "key2": "1234"}

	err = validate.Struct(test)
	NotEqual(t, err, nil)

	errs, ok := err.(validator.ValidationErrors)
	Equal(t, ok, true)

	tests := []struct {
		ns       string
		expected string
	}{
		{
			ns:       "Test.IsColor",
			expected: "IsColor debe ser un color válido",
		},
		{
			ns:       "Test.MAC",
			expected: "MAC debe contener una dirección MAC válida",
		},
		{
			ns:       "Test.IPAddr",
			expected: "IPAddr debe ser una dirección IP resoluble",
		},
		{
			ns:       "Test.IPAddrv4",
			expected: "IPAddrv4 debe ser una dirección IPv4 resoluble",
		},
		{
			ns:       "Test.IPAddrv6",
			expected: "IPAddrv6 debe ser una dirección IPv6 resoluble",
		},
		{
			ns:       "Test.UDPAddr",
			expected: "UDPAddr debe ser una dirección UDP válida",
		},
		{
			ns:       "Test.UDPAddrv4",
			expected: "UDPAddrv4 debe ser una dirección IPv4 UDP válida",
		},
		{
			ns:       "Test.UDPAddrv6",
			expected: "UDPAddrv6 debe ser una dirección IPv6 UDP válida",
		},
		{
			ns:       "Test.TCPAddr",
			expected: "TCPAddr debe ser una dirección TCP válida",
		},
		{
			ns:       "Test.TCPAddrv4",
			expected: "TCPAddrv4 debe ser una dirección IPv4 TCP válida",
		},
		{
			ns:       "Test.TCPAddrv6",
			expected: "TCPAddrv6 debe ser una dirección IPv6 TCP válida",
		},
		{
			ns:       "Test.CIDR",
			expected: "CIDR debe contener una anotación válida del CIDR",
		},
		{
			ns:       "Test.CIDRv4",
			expected: "CIDRv4 debe contener una notación CIDR válida para una dirección IPv4",
		},
		{
			ns:       "Test.CIDRv6",
			expected: "CIDRv6 debe contener una notación CIDR válida para una dirección IPv6",
		},
		{
			ns:       "Test.SSN",
			expected: "SSN debe ser un número válido de SSN",
		},
		{
			ns:       "Test.IP",
			expected: "IP debe ser una dirección IP válida",
		},
		{
			ns:       "Test.IPv4",
			expected: "IPv4 debe ser una dirección IPv4 válida",
		},
		{
			ns:       "Test.IPv6",
			expected: "IPv6 debe ser una dirección IPv6 válida",
		},
		{
			ns:       "Test.DataURI",
			expected: "DataURI debe contener un URI de datos válido",
		},
		{
			ns:       "Test.Latitude",
			expected: "Latitude debe contener coordenadas de latitud válidas",
		},
		{
			ns:       "Test.Longitude",
			expected: "Longitude debe contener unas coordenadas de longitud válidas",
		},
		{
			ns:       "Test.MultiByte",
			expected: "MultiByte debe contener caracteres multibyte",
		},
		{
			ns:       "Test.ASCII",
			expected: "ASCII debe contener sólo caracteres ascii",
		},
		{
			ns:       "Test.PrintableASCII",
			expected: "PrintableASCII debe contener sólo caracteres ASCII imprimibles",
		},
		{
			ns:       "Test.UUID",
			expected: "UUID debe ser un UUID válido",
		},
		{
			ns:       "Test.UUID3",
			expected: "UUID3 debe ser una versión válida 3 UUID",
		},
		{
			ns:       "Test.UUID4",
			expected: "UUID4 debe ser una versión válida 4 UUID",
		},
		{
			ns:       "Test.UUID5",
			expected: "UUID5 debe ser una versión válida 5 UUID",
		},
		{
			ns:       "Test.ULID",
			expected: "ULID debe ser un ULID válido",
		},
		{
			ns:       "Test.ISBN",
			expected: "ISBN debe ser un número ISBN válido",
		},
		{
			ns:       "Test.ISBN10",
			expected: "ISBN10 debe ser un número ISBN-10 válido",
		},
		{
			ns:       "Test.ISBN13",
			expected: "ISBN13 debe ser un número ISBN-13 válido",
		},
		{
			ns:       "Test.ISSN",
			expected: "ISSN debe ser un número ISSN válido",
		},
		{
			ns:       "Test.Excludes",
			expected: "Excludes no puede contener el texto 'text'",
		},
		{
			ns:       "Test.ExcludesAll",
			expected: "ExcludesAll no puede contener ninguno de los siguientes caracteres '!@#$'",
		},
		{
			ns:       "Test.ExcludesRune",
			expected: "ExcludesRune no puede contener lo siguiente '☻'",
		},
		{
			ns:       "Test.ContainsAny",
			expected: "ContainsAny debe contener al menos uno de los siguientes caracteres '!@#$'",
		},
		{
			ns:       "Test.Contains",
			expected: "Contains debe contener el texto 'purpose'",
		},
		{
			ns:       "Test.Base64",
			expected: "Base64 debe ser una cadena de Base64 válida",
		},
		{
			ns:       "Test.Email",
			expected: "Email debe ser una dirección de correo electrónico válida",
		},
		{
			ns:       "Test.URL",
			expected: "URL debe ser un URL válido",
		},
		{
			ns:       "Test.URI",
			expected: "URI debe ser una URI válida",
		},
		{
			ns:       "Test.RGBColorString",
			expected: "RGBColorString debe ser un color RGB válido",
		},
		{
			ns:       "Test.RGBAColorString",
			expected: "RGBAColorString debe ser un color RGBA válido",
		},
		{
			ns:       "Test.HSLColorString",
			expected: "HSLColorString debe ser un color HSL válido",
		},
		{
			ns:       "Test.HSLAColorString",
			expected: "HSLAColorString debe ser un color HSL válido",
		},
		{
			ns:       "Test.HexadecimalString",
			expected: "HexadecimalString debe ser un hexadecimal válido",
		},
		{
			ns:       "Test.HexColorString",
			expected: "HexColorString debe ser un color HEX válido",
		},
		{
			ns:       "Test.NumberString",
			expected: "NumberString debe ser un número válido",
		},
		{
			ns:       "Test.NumericString",
			expected: "NumericString debe ser un valor numérico válido",
		},
		{
			ns:       "Test.AlphanumString",
			expected: "AlphanumString sólo puede contener caracteres alfanuméricos",
		},
		{
			ns:       "Test.AlphaString",
			expected: "AlphaString sólo puede contener caracteres alfabéticos",
		},
		{
			ns:       "Test.LtFieldString",
			expected: "LtFieldString debe ser menor que MaxString",
		},
		{
			ns:       "Test.LteFieldString",
			expected: "LteFieldString debe ser menor o igual a MaxString",
		},
		{
			ns:       "Test.GtFieldString",
			expected: "GtFieldString debe ser mayor que MaxString",
		},
		{
			ns:       "Test.GteFieldString",
			expected: "GteFieldString debe ser mayor o igual a MaxString",
		},
		{
			ns:       "Test.NeFieldString",
			expected: "NeFieldString no puede ser igual a EqFieldString",
		},
		{
			ns:       "Test.LtCSFieldString",
			expected: "LtCSFieldString debe ser menor que Inner.LtCSFieldString",
		},
		{
			ns:       "Test.LteCSFieldString",
			expected: "LteCSFieldString debe ser menor o igual a Inner.LteCSFieldString",
		},
		{
			ns:       "Test.GtCSFieldString",
			expected: "GtCSFieldString debe ser mayor que Inner.GtCSFieldString",
		},
		{
			ns:       "Test.GteCSFieldString",
			expected: "GteCSFieldString debe ser mayor o igual a Inner.GteCSFieldString",
		},
		{
			ns:       "Test.NeCSFieldString",
			expected: "NeCSFieldString no puede ser igual a Inner.NeCSFieldString",
		},
		{
			ns:       "Test.EqCSFieldString",
			expected: "EqCSFieldString debe ser igual a Inner.EqCSFieldString",
		},
		{
			ns:       "Test.EqFieldString",
			expected: "EqFieldString debe ser igual a MaxString",
		},
		{
			ns:       "Test.GteString",
			expected: "GteString debe tener al menos 3 caracteres de longitud",
		},
		{
			ns:       "Test.GteNumber",
			expected: "GteNumber debe ser 5,56 o mayor",
		},
		{
			ns:       "Test.GteMultiple",
			expected: "GteMultiple debe contener al menos 2 elementos",
		},
		{
			ns:       "Test.GteTime",
			expected: "GteTime debe ser después o durante la fecha y hora actuales",
		},
		{
			ns:       "Test.GtString",
			expected: "GtString debe ser mayor que 3 caracteres en longitud",
		},
		{
			ns:       "Test.GtNumber",
			expected: "GtNumber debe ser mayor que 5,56",
		},
		{
			ns:       "Test.GtMultiple",
			expected: "GtMultiple debe contener más de 2 elementos",
		},
		{
			ns:       "Test.GtTime",
			expected: "GtTime debe ser después de la fecha y hora actual",
		},
		{
			ns:       "Test.LteString",
			expected: "LteString debe tener un máximo de 3 caracteres de longitud",
		},
		{
			ns:       "Test.LteNumber",
			expected: "LteNumber debe ser 5,56 o menos",
		},
		{
			ns:       "Test.LteMultiple",
			expected: "LteMultiple debe contener como máximo 2 elementos",
		},
		{
			ns:       "Test.LteTime",
			expected: "LteTime debe ser antes o durante la fecha y hora actual",
		},
		{
			ns:       "Test.LtString",
			expected: "LtString debe tener menos de 3 caracteres de longitud",
		},
		{
			ns:       "Test.LtNumber",
			expected: "LtNumber debe ser menos de 5,56",
		},
		{
			ns:       "Test.LtMultiple",
			expected: "LtMultiple debe contener menos de 2 elementos",
		},
		{
			ns:       "Test.LtTime",
			expected: "LtTime debe ser antes de la fecha y hora actual",
		},
		{
			ns:       "Test.NeString",
			expected: "NeString no debería ser igual a ",
		},
		{
			ns:       "Test.NeNumber",
			expected: "NeNumber no debería ser igual a 0.00",
		},
		{
			ns:       "Test.NeMultiple",
			expected: "NeMultiple no debería ser igual a 0",
		},
		{
			ns:       "Test.EqString",
			expected: "EqString no es igual a 3",
		},
		{
			ns:       "Test.EqNumber",
			expected: "EqNumber no es igual a 2.33",
		},
		{
			ns:       "Test.EqMultiple",
			expected: "EqMultiple no es igual a 7",
		},
		{
			ns:       "Test.MaxString",
			expected: "MaxString debe tener un máximo de 3 caracteres de longitud",
		},
		{
			ns:       "Test.MaxNumber",
			expected: "MaxNumber debe ser 1.113,00 o menos",
		},
		{
			ns:       "Test.MaxMultiple",
			expected: "MaxMultiple debe contener como máximo 7 elementos",
		},
		{
			ns:       "Test.MinString",
			expected: "MinString debe tener al menos 1 carácter de longitud",
		},
		{
			ns:       "Test.MinNumber",
			expected: "MinNumber debe ser 1.113,00 o más",
		},
		{
			ns:       "Test.MinMultiple",
			expected: "MinMultiple debe contener al menos 7 elementos",
		},
		{
			ns:       "Test.LenString",
			expected: "LenString debe tener 1 carácter de longitud",
		},
		{
			ns:       "Test.LenNumber",
			expected: "LenNumber debe ser igual a 1.113,00",
		},
		{
			ns:       "Test.LenMultiple",
			expected: "LenMultiple debe contener 7 elementos",
		},
		{
			ns:       "Test.RequiredString",
			expected: "RequiredString es un campo requerido",
		},
		{
			ns:       "Test.RequiredNumber",
			expected: "RequiredNumber es un campo requerido",
		},
		{
			ns:       "Test.RequiredMultiple",
			expected: "RequiredMultiple es un campo requerido",
		},
		{
			ns:       "Test.StrPtrMinLen",
			expected: "StrPtrMinLen debe tener al menos 10 caracteres de longitud",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen debe tener un máximo de 1 carácter de longitud",
		},
		{
			ns:       "Test.StrPtrLen",
			expected: "StrPtrLen debe tener 2 caracteres de longitud",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt debe tener menos de 1 carácter de longitud",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte debe tener un máximo de 1 carácter de longitud",
		},
		{
			ns:       "Test.StrPtrGt",
			expected: "StrPtrGt debe ser mayor que 10 caracteres en longitud",
		},
		{
			ns:       "Test.StrPtrGte",
			expected: "StrPtrGte debe tener al menos 10 caracteres de longitud",
		},
		{
			ns:       "Test.OneOfString",
			expected: "OneOfString debe ser uno de [red green]",
		},
		{
			ns:       "Test.OneOfInt",
			expected: "OneOfInt debe ser uno de [5 63]",
		},
		{
			ns:       "Test.UniqueSlice",
			expected: "UniqueSlice debe contener valores únicos",
		},
		{
			ns:       "Test.UniqueArray",
			expected: "UniqueArray debe contener valores únicos",
		},
		{
			ns:       "Test.UniqueMap",
			expected: "UniqueMap debe contener valores únicos",
		},
		{
			ns:       "Test.RequiredWithout",
			expected: "RequiredWithout es un campo requerido cuando los campos [Inner.RequiredWithoutField] no estan presentes",
    },
    {
			ns: "Test.Image",
			expected: "Image debe ser una imagen válida",
		},
	}

	for _, tt := range tests {

		var fe validator.FieldError

		for _, e := range errs {
			if tt.ns == e.Namespace() {
				fe = e
				break
			}
		}

		NotEqual(t, fe, nil)
		Equal(t, tt.expected, fe.Translate(trans))
	}

}
