package en

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	english "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestTranslations(t *testing.T) {
	eng := english.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()

	err := RegisterDefaultTranslations(validate, trans)
	Equal(t, err, nil)

	type Inner struct {
		EqCSFieldString  string
		NeCSFieldString  string
		GtCSFieldString  string
		GteCSFieldString string
		LtCSFieldString  string
		LteCSFieldString string
		RequiredIf       string
	}

	type Test struct {
		Inner             Inner
		RequiredString    string            `validate:"required"`
		RequiredNumber    int               `validate:"required"`
		RequiredMultiple  []string          `validate:"required"`
		RequiredIf        string            `validate:"required_if=Inner.RequiredIf abcd"`
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
		FQDN              string            `validate:"fqdn"`
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
		JSONString        string            `validate:"json"`
		JWTString         string            `validate:"jwt"`
		LowercaseString   string            `validate:"lowercase"`
		UppercaseString   string            `validate:"uppercase"`
		Datetime          string            `validate:"datetime=2006-01-02"`
		PostCode          string            `validate:"postcode_iso3166_alpha2=SG"`
		PostCodeCountry   string
		PostCodeByField   string        `validate:"postcode_iso3166_alpha2_field=PostCodeCountry"`
		BooleanString     string        `validate:"boolean"`
		Image             string        `validate:"image"`
		CveString         string        `validate:"cve"`
		MinDuration       time.Duration `validate:"min=1h30m,max=2h"`
		MaxDuration       time.Duration `validate:"min=1h30m,max=2h"`
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

	test.LowercaseString = "ABCDEFG"
	test.UppercaseString = "abcdefg"

	s := "toolong"
	test.StrPtrMaxLen = &s
	test.StrPtrLen = &s

	test.UniqueSlice = []string{"1234", "1234"}
	test.UniqueMap = map[string]string{"key1": "1234", "key2": "1234"}
	test.Datetime = "2008-Feb-01"
	test.BooleanString = "A"
	test.CveString = "A"

	test.MinDuration = 1 * time.Hour
	test.MaxDuration = 3 * time.Hour

	test.Inner.RequiredIf = "abcd"

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
			expected: "IsColor must be a valid color",
		},
		{
			ns:       "Test.MAC",
			expected: "MAC must contain a valid MAC address",
		},
		{
			ns:       "Test.FQDN",
			expected: "FQDN must be a valid FQDN",
		},
		{
			ns:       "Test.IPAddr",
			expected: "IPAddr must be a resolvable IP address",
		},
		{
			ns:       "Test.IPAddrv4",
			expected: "IPAddrv4 must be a resolvable IPv4 address",
		},
		{
			ns:       "Test.IPAddrv6",
			expected: "IPAddrv6 must be a resolvable IPv6 address",
		},
		{
			ns:       "Test.UDPAddr",
			expected: "UDPAddr must be a valid UDP address",
		},
		{
			ns:       "Test.UDPAddrv4",
			expected: "UDPAddrv4 must be a valid IPv4 UDP address",
		},
		{
			ns:       "Test.UDPAddrv6",
			expected: "UDPAddrv6 must be a valid IPv6 UDP address",
		},
		{
			ns:       "Test.TCPAddr",
			expected: "TCPAddr must be a valid TCP address",
		},
		{
			ns:       "Test.TCPAddrv4",
			expected: "TCPAddrv4 must be a valid IPv4 TCP address",
		},
		{
			ns:       "Test.TCPAddrv6",
			expected: "TCPAddrv6 must be a valid IPv6 TCP address",
		},
		{
			ns:       "Test.CIDR",
			expected: "CIDR must contain a valid CIDR notation",
		},
		{
			ns:       "Test.CIDRv4",
			expected: "CIDRv4 must contain a valid CIDR notation for an IPv4 address",
		},
		{
			ns:       "Test.CIDRv6",
			expected: "CIDRv6 must contain a valid CIDR notation for an IPv6 address",
		},
		{
			ns:       "Test.SSN",
			expected: "SSN must be a valid SSN number",
		},
		{
			ns:       "Test.IP",
			expected: "IP must be a valid IP address",
		},
		{
			ns:       "Test.IPv4",
			expected: "IPv4 must be a valid IPv4 address",
		},
		{
			ns:       "Test.IPv6",
			expected: "IPv6 must be a valid IPv6 address",
		},
		{
			ns:       "Test.DataURI",
			expected: "DataURI must contain a valid Data URI",
		},
		{
			ns:       "Test.Latitude",
			expected: "Latitude must contain valid latitude coordinates",
		},
		{
			ns:       "Test.Longitude",
			expected: "Longitude must contain a valid longitude coordinates",
		},
		{
			ns:       "Test.MultiByte",
			expected: "MultiByte must contain multibyte characters",
		},
		{
			ns:       "Test.ASCII",
			expected: "ASCII must contain only ascii characters",
		},
		{
			ns:       "Test.PrintableASCII",
			expected: "PrintableASCII must contain only printable ascii characters",
		},
		{
			ns:       "Test.UUID",
			expected: "UUID must be a valid UUID",
		},
		{
			ns:       "Test.UUID3",
			expected: "UUID3 must be a valid version 3 UUID",
		},
		{
			ns:       "Test.UUID4",
			expected: "UUID4 must be a valid version 4 UUID",
		},
		{
			ns:       "Test.UUID5",
			expected: "UUID5 must be a valid version 5 UUID",
		},
		{
			ns:       "Test.ULID",
			expected: "ULID must be a valid ULID",
		},
		{
			ns:       "Test.ISBN",
			expected: "ISBN must be a valid ISBN number",
		},
		{
			ns:       "Test.ISBN10",
			expected: "ISBN10 must be a valid ISBN-10 number",
		},
		{
			ns:       "Test.ISBN13",
			expected: "ISBN13 must be a valid ISBN-13 number",
		},
		{
			ns:       "Test.ISSN",
			expected: "ISSN must be a valid ISSN number",
		},
		{
			ns:       "Test.Excludes",
			expected: "Excludes cannot contain the text 'text'",
		},
		{
			ns:       "Test.ExcludesAll",
			expected: "ExcludesAll cannot contain any of the following characters '!@#$'",
		},
		{
			ns:       "Test.ExcludesRune",
			expected: "ExcludesRune cannot contain the following '☻'",
		},
		{
			ns:       "Test.ContainsAny",
			expected: "ContainsAny must contain at least one of the following characters '!@#$'",
		},
		{
			ns:       "Test.Contains",
			expected: "Contains must contain the text 'purpose'",
		},
		{
			ns:       "Test.Base64",
			expected: "Base64 must be a valid Base64 string",
		},
		{
			ns:       "Test.Email",
			expected: "Email must be a valid email address",
		},
		{
			ns:       "Test.URL",
			expected: "URL must be a valid URL",
		},
		{
			ns:       "Test.URI",
			expected: "URI must be a valid URI",
		},
		{
			ns:       "Test.RGBColorString",
			expected: "RGBColorString must be a valid RGB color",
		},
		{
			ns:       "Test.RGBAColorString",
			expected: "RGBAColorString must be a valid RGBA color",
		},
		{
			ns:       "Test.HSLColorString",
			expected: "HSLColorString must be a valid HSL color",
		},
		{
			ns:       "Test.HSLAColorString",
			expected: "HSLAColorString must be a valid HSLA color",
		},
		{
			ns:       "Test.HexadecimalString",
			expected: "HexadecimalString must be a valid hexadecimal",
		},
		{
			ns:       "Test.HexColorString",
			expected: "HexColorString must be a valid HEX color",
		},
		{
			ns:       "Test.NumberString",
			expected: "NumberString must be a valid number",
		},
		{
			ns:       "Test.NumericString",
			expected: "NumericString must be a valid numeric value",
		},
		{
			ns:       "Test.AlphanumString",
			expected: "AlphanumString can only contain alphanumeric characters",
		},
		{
			ns:       "Test.AlphaString",
			expected: "AlphaString can only contain alphabetic characters",
		},
		{
			ns:       "Test.LtFieldString",
			expected: "LtFieldString must be less than MaxString",
		},
		{
			ns:       "Test.LteFieldString",
			expected: "LteFieldString must be less than or equal to MaxString",
		},
		{
			ns:       "Test.GtFieldString",
			expected: "GtFieldString must be greater than MaxString",
		},
		{
			ns:       "Test.GteFieldString",
			expected: "GteFieldString must be greater than or equal to MaxString",
		},
		{
			ns:       "Test.NeFieldString",
			expected: "NeFieldString cannot be equal to EqFieldString",
		},
		{
			ns:       "Test.LtCSFieldString",
			expected: "LtCSFieldString must be less than Inner.LtCSFieldString",
		},
		{
			ns:       "Test.LteCSFieldString",
			expected: "LteCSFieldString must be less than or equal to Inner.LteCSFieldString",
		},
		{
			ns:       "Test.GtCSFieldString",
			expected: "GtCSFieldString must be greater than Inner.GtCSFieldString",
		},
		{
			ns:       "Test.GteCSFieldString",
			expected: "GteCSFieldString must be greater than or equal to Inner.GteCSFieldString",
		},
		{
			ns:       "Test.NeCSFieldString",
			expected: "NeCSFieldString cannot be equal to Inner.NeCSFieldString",
		},
		{
			ns:       "Test.EqCSFieldString",
			expected: "EqCSFieldString must be equal to Inner.EqCSFieldString",
		},
		{
			ns:       "Test.EqFieldString",
			expected: "EqFieldString must be equal to MaxString",
		},
		{
			ns:       "Test.GteString",
			expected: "GteString must be at least 3 characters in length",
		},
		{
			ns:       "Test.GteNumber",
			expected: "GteNumber must be 5.56 or greater",
		},
		{
			ns:       "Test.GteMultiple",
			expected: "GteMultiple must contain at least 2 items",
		},
		{
			ns:       "Test.GteTime",
			expected: "GteTime must be greater than or equal to the current Date & Time",
		},
		{
			ns:       "Test.GtString",
			expected: "GtString must be greater than 3 characters in length",
		},
		{
			ns:       "Test.GtNumber",
			expected: "GtNumber must be greater than 5.56",
		},
		{
			ns:       "Test.GtMultiple",
			expected: "GtMultiple must contain more than 2 items",
		},
		{
			ns:       "Test.GtTime",
			expected: "GtTime must be greater than the current Date & Time",
		},
		{
			ns:       "Test.LteString",
			expected: "LteString must be at maximum 3 characters in length",
		},
		{
			ns:       "Test.LteNumber",
			expected: "LteNumber must be 5.56 or less",
		},
		{
			ns:       "Test.LteMultiple",
			expected: "LteMultiple must contain at maximum 2 items",
		},
		{
			ns:       "Test.LteTime",
			expected: "LteTime must be less than or equal to the current Date & Time",
		},
		{
			ns:       "Test.LtString",
			expected: "LtString must be less than 3 characters in length",
		},
		{
			ns:       "Test.LtNumber",
			expected: "LtNumber must be less than 5.56",
		},
		{
			ns:       "Test.LtMultiple",
			expected: "LtMultiple must contain less than 2 items",
		},
		{
			ns:       "Test.LtTime",
			expected: "LtTime must be less than the current Date & Time",
		},
		{
			ns:       "Test.NeString",
			expected: "NeString should not be equal to ",
		},
		{
			ns:       "Test.NeNumber",
			expected: "NeNumber should not be equal to 0.00",
		},
		{
			ns:       "Test.NeMultiple",
			expected: "NeMultiple should not be equal to 0",
		},
		{
			ns:       "Test.EqString",
			expected: "EqString is not equal to 3",
		},
		{
			ns:       "Test.EqNumber",
			expected: "EqNumber is not equal to 2.33",
		},
		{
			ns:       "Test.EqMultiple",
			expected: "EqMultiple is not equal to 7",
		},
		{
			ns:       "Test.MaxString",
			expected: "MaxString must be a maximum of 3 characters in length",
		},
		{
			ns:       "Test.MaxNumber",
			expected: "MaxNumber must be 1,113.00 or less",
		},
		{
			ns:       "Test.MaxMultiple",
			expected: "MaxMultiple must contain at maximum 7 items",
		},
		{
			ns:       "Test.MinString",
			expected: "MinString must be at least 1 character in length",
		},
		{
			ns:       "Test.MinNumber",
			expected: "MinNumber must be 1,113.00 or greater",
		},
		{
			ns:       "Test.MinMultiple",
			expected: "MinMultiple must contain at least 7 items",
		},
		{
			ns:       "Test.LenString",
			expected: "LenString must be 1 character in length",
		},
		{
			ns:       "Test.LenNumber",
			expected: "LenNumber must be equal to 1,113.00",
		},
		{
			ns:       "Test.LenMultiple",
			expected: "LenMultiple must contain 7 items",
		},
		{
			ns:       "Test.RequiredString",
			expected: "RequiredString is a required field",
		},
		{
			ns:       "Test.RequiredIf",
			expected: "RequiredIf is a required field",
		},
		{
			ns:       "Test.RequiredNumber",
			expected: "RequiredNumber is a required field",
		},
		{
			ns:       "Test.RequiredMultiple",
			expected: "RequiredMultiple is a required field",
		},
		{
			ns:       "Test.StrPtrMinLen",
			expected: "StrPtrMinLen must be at least 10 characters in length",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen must be a maximum of 1 character in length",
		},
		{
			ns:       "Test.StrPtrLen",
			expected: "StrPtrLen must be 2 characters in length",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt must be less than 1 character in length",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte must be at maximum 1 character in length",
		},
		{
			ns:       "Test.StrPtrGt",
			expected: "StrPtrGt must be greater than 10 characters in length",
		},
		{
			ns:       "Test.StrPtrGte",
			expected: "StrPtrGte must be at least 10 characters in length",
		},
		{
			ns:       "Test.OneOfString",
			expected: "OneOfString must be one of [red green]",
		},
		{
			ns:       "Test.OneOfInt",
			expected: "OneOfInt must be one of [5 63]",
		},
		{
			ns:       "Test.UniqueSlice",
			expected: "UniqueSlice must contain unique values",
		},
		{
			ns:       "Test.UniqueArray",
			expected: "UniqueArray must contain unique values",
		},
		{
			ns:       "Test.UniqueMap",
			expected: "UniqueMap must contain unique values",
		},
		{
			ns:       "Test.JSONString",
			expected: "JSONString must be a valid json string",
		},
		{
			ns:       "Test.JWTString",
			expected: "JWTString must be a valid jwt string",
		},
		{
			ns:       "Test.LowercaseString",
			expected: "LowercaseString must be a lowercase string",
		},
		{
			ns:       "Test.UppercaseString",
			expected: "UppercaseString must be an uppercase string",
		},
		{
			ns:       "Test.Datetime",
			expected: "Datetime does not match the 2006-01-02 format",
		},
		{
			ns:       "Test.PostCode",
			expected: "PostCode does not match postcode format of SG country",
		},
		{
			ns:       "Test.PostCodeByField",
			expected: "PostCodeByField does not match postcode format of country in PostCodeCountry field",
		},
		{
			ns:       "Test.BooleanString",
			expected: "BooleanString must be a valid boolean value",
		},
		{
			ns:       "Test.Image",
			expected: "Image must be a valid image",
		},
		{
			ns:       "Test.CveString",
			expected: "CveString must be a valid cve identifier",
		},
		{
			ns:       "Test.MinDuration",
			expected: "MinDuration must be 1h30m or greater",
		},
		{
			ns:       "Test.MaxDuration",
			expected: "MaxDuration must be 2h or less",
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
