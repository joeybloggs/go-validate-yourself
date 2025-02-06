package de

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	german "github.com/go-playground/locales/de"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestTranslations(t *testing.T) {
	ger := german.New()
	uni := ut.New(ger, ger)
	trans, _ := uni.GetTranslator("de")

	validate := validator.New()

	err := RegisterDefaultTranslations(validate, trans)
	Equal(t, err, nil)

	type Inner struct {
		EqCSFieldString    string
		NeCSFieldString    string
		GtCSFieldString    string
		GteCSFieldString   string
		LtCSFieldString    string
		LteCSFieldString   string
		RequiredIf         string
		RequiredUnless     string
		RequiredWith       string
		RequiredWithAll    string
		RequiredWithout    string
		RequiredWithoutAll string
		ExcludedIf         string
		ExcludedUnless     string
		ExcludedWith       string
		ExcludedWithAll    string
		ExcludedWithout    string
		ExcludedWithoutAll string
	}

	type Test struct {
		Inner              Inner
		RequiredString     string            `validate:"required"`
		RequiredNumber     int               `validate:"required"`
		RequiredMultiple   []string          `validate:"required"`
		RequiredIf         string            `validate:"required_if=Inner.RequiredIf abcd"`
		RequiredUnless     string            `validate:"required_unless=Inner.RequiredUnless abcd"`
		RequiredWith       string            `validate:"required_with=Inner.RequiredWith"`
		RequiredWithAll    string            `validate:"required_with_all=Inner.RequiredWith Inner.RequiredWithAll"`
		RequiredWithout    string            `validate:"required_without=Inner.RequiredWithout"`
		RequiredWithoutAll string            `validate:"required_without_all=Inner.RequiredWithout Inner.RequiredWithoutAll"`
		ExcludedIf         string            `validate:"excluded_if=Inner.ExcludedIf abcd"`
		ExcludedUnless     string            `validate:"excluded_unless=Inner.ExcludedUnless abcd"`
		ExcludedWith       string            `validate:"excluded_with=Inner.ExcludedWith"`
		ExcludedWithout    string            `validate:"excluded_with_all=Inner.ExcludedWithAll"`
		ExcludedWithAll    string            `validate:"excluded_without=Inner.ExcludedWithout"`
		ExcludedWithoutAll string            `validate:"excluded_without_all=Inner.ExcludedWithoutAll"`
		IsDefault          string            `validate:"isdefault"`
		LenString          string            `validate:"len=1"`
		LenNumber          float64           `validate:"len=1113.00"`
		LenMultiple        []string          `validate:"len=7"`
		MinString          string            `validate:"min=1"`
		MinNumber          float64           `validate:"min=1113.00"`
		MinMultiple        []string          `validate:"min=7"`
		MaxString          string            `validate:"max=3"`
		MaxNumber          float64           `validate:"max=1113.00"`
		MaxMultiple        []string          `validate:"max=7"`
		EqString           string            `validate:"eq=3"`
		EqNumber           float64           `validate:"eq=2.33"`
		EqMultiple         []string          `validate:"eq=7"`
		NeString           string            `validate:"ne="`
		NeNumber           float64           `validate:"ne=0.00"`
		NeMultiple         []string          `validate:"ne=0"`
		LtString           string            `validate:"lt=3"`
		LtNumber           float64           `validate:"lt=5.56"`
		LtMultiple         []string          `validate:"lt=2"`
		LtTime             time.Time         `validate:"lt"`
		LteString          string            `validate:"lte=3"`
		LteNumber          float64           `validate:"lte=5.56"`
		LteMultiple        []string          `validate:"lte=2"`
		LteTime            time.Time         `validate:"lte"`
		GtString           string            `validate:"gt=3"`
		GtNumber           float64           `validate:"gt=5.56"`
		GtMultiple         []string          `validate:"gt=2"`
		GtTime             time.Time         `validate:"gt"`
		GteString          string            `validate:"gte=3"`
		GteNumber          float64           `validate:"gte=5.56"`
		GteMultiple        []string          `validate:"gte=2"`
		GteTime            time.Time         `validate:"gte"`
		EqFieldString      string            `validate:"eqfield=MaxString"`
		EqCSFieldString    string            `validate:"eqcsfield=Inner.EqCSFieldString"`
		NeCSFieldString    string            `validate:"necsfield=Inner.NeCSFieldString"`
		GtCSFieldString    string            `validate:"gtcsfield=Inner.GtCSFieldString"`
		GteCSFieldString   string            `validate:"gtecsfield=Inner.GteCSFieldString"`
		LtCSFieldString    string            `validate:"ltcsfield=Inner.LtCSFieldString"`
		LteCSFieldString   string            `validate:"ltecsfield=Inner.LteCSFieldString"`
		NeFieldString      string            `validate:"nefield=EqFieldString"`
		GtFieldString      string            `validate:"gtfield=MaxString"`
		GteFieldString     string            `validate:"gtefield=MaxString"`
		LtFieldString      string            `validate:"ltfield=MaxString"`
		LteFieldString     string            `validate:"ltefield=MaxString"`
		AlphaString        string            `validate:"alpha"`
		AlphanumString     string            `validate:"alphanum"`
		NumericString      string            `validate:"numeric"`
		NumberString       string            `validate:"number"`
		HexadecimalString  string            `validate:"hexadecimal"`
		HexColorString     string            `validate:"hexcolor"`
		RGBColorString     string            `validate:"rgb"`
		RGBAColorString    string            `validate:"rgba"`
		HSLColorString     string            `validate:"hsl"`
		HSLAColorString    string            `validate:"hsla"`
		Email              string            `validate:"email"`
		URL                string            `validate:"url"`
		URI                string            `validate:"uri"`
		Base64             string            `validate:"base64"`
		Contains           string            `validate:"contains=purpose"`
		ContainsAny        string            `validate:"containsany=!@#$"`
		Excludes           string            `validate:"excludes=text"`
		ExcludesAll        string            `validate:"excludesall=!@#$"`
		ExcludesRune       string            `validate:"excludesrune=☻"`
		ISBN               string            `validate:"isbn"`
		ISBN10             string            `validate:"isbn10"`
		ISBN13             string            `validate:"isbn13"`
		ISSN               string            `validate:"issn"`
		UUID               string            `validate:"uuid"`
		UUID3              string            `validate:"uuid3"`
		UUID4              string            `validate:"uuid4"`
		UUID5              string            `validate:"uuid5"`
		ULID               string            `validate:"ulid"`
		ASCII              string            `validate:"ascii"`
		PrintableASCII     string            `validate:"printascii"`
		MultiByte          string            `validate:"multibyte"`
		DataURI            string            `validate:"datauri"`
		Latitude           string            `validate:"latitude"`
		Longitude          string            `validate:"longitude"`
		SSN                string            `validate:"ssn"`
		IP                 string            `validate:"ip"`
		IPv4               string            `validate:"ipv4"`
		IPv6               string            `validate:"ipv6"`
		CIDR               string            `validate:"cidr"`
		CIDRv4             string            `validate:"cidrv4"`
		CIDRv6             string            `validate:"cidrv6"`
		TCPAddr            string            `validate:"tcp_addr"`
		TCPAddrv4          string            `validate:"tcp4_addr"`
		TCPAddrv6          string            `validate:"tcp6_addr"`
		UDPAddr            string            `validate:"udp_addr"`
		UDPAddrv4          string            `validate:"udp4_addr"`
		UDPAddrv6          string            `validate:"udp6_addr"`
		IPAddr             string            `validate:"ip_addr"`
		IPAddrv4           string            `validate:"ip4_addr"`
		IPAddrv6           string            `validate:"ip6_addr"`
		UinxAddr           string            `validate:"unix_addr"` // can't fail from within Go's net package currently, but maybe in the future
		MAC                string            `validate:"mac"`
		FQDN               string            `validate:"fqdn"`
		IsColor            string            `validate:"iscolor"`
		StrPtrMinLen       *string           `validate:"min=10"`
		StrPtrMaxLen       *string           `validate:"max=1"`
		StrPtrLen          *string           `validate:"len=2"`
		StrPtrLt           *string           `validate:"lt=1"`
		StrPtrLte          *string           `validate:"lte=1"`
		StrPtrGt           *string           `validate:"gt=10"`
		StrPtrGte          *string           `validate:"gte=10"`
		OneOfString        string            `validate:"oneof=red green"`
		OneOfInt           int               `validate:"oneof=5 63"`
		UniqueSlice        []string          `validate:"unique"`
		UniqueArray        [3]string         `validate:"unique"`
		UniqueMap          map[string]string `validate:"unique"`
		JSONString         string            `validate:"json"`
		JWTString          string            `validate:"jwt"`
		LowercaseString    string            `validate:"lowercase"`
		UppercaseString    string            `validate:"uppercase"`
		Datetime           string            `validate:"datetime=2006-01-02"`
		PostCode           string            `validate:"postcode_iso3166_alpha2=SG"`
		PostCodeCountry    string
		PostCodeByField    string `validate:"postcode_iso3166_alpha2_field=PostCodeCountry"`
		BooleanString      string `validate:"boolean"`
		Image              string `validate:"image"`
		CveString          string `validate:"cve"`
	}

	var test Test

	test.Inner.EqCSFieldString = "1234"
	test.Inner.GtCSFieldString = "1234"
	test.Inner.GteCSFieldString = "1234"
	test.Inner.RequiredUnless = "1234"
	test.Inner.RequiredWith = "1234"
	test.Inner.RequiredWithAll = "1234"
	test.Inner.ExcludedIf = "abcd"
	test.Inner.ExcludedUnless = "1234"
	test.Inner.ExcludedWith = "1234"
	test.Inner.ExcludedWithAll = "1234"

	test.ExcludedIf = "1234"
	test.ExcludedUnless = "1234"
	test.ExcludedWith = "1234"
	test.ExcludedWithAll = "1234"
	test.ExcludedWithout = "1234"
	test.ExcludedWithoutAll = "1234"

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
			expected: "IsColor muss eine gültige Farbe sein",
		},
		{
			ns:       "Test.MAC",
			expected: "MAC muss eine gültige MAC-Adresse sein",
		},
		{
			ns:       "Test.FQDN",
			expected: "FQDN muss eine gültige FQDN sein",
		},
		{
			ns:       "Test.IPAddr",
			expected: "IPAddr muss eine auflösbare IP-Adresse sein",
		},
		{
			ns:       "Test.IPAddrv4",
			expected: "IPAddrv4 muss eine auflösbare IPv4-Adresse sein",
		},
		{
			ns:       "Test.IPAddrv6",
			expected: "IPAddrv6 muss eine auflösbare IPv6-Adresse sein",
		},
		{
			ns:       "Test.UDPAddr",
			expected: "UDPAddr muss eine gültige UDP-Adresse sein",
		},
		{
			ns:       "Test.UDPAddrv4",
			expected: "UDPAddrv4 muss eine gültige IPv4-UDP-Adresse sein",
		},
		{
			ns:       "Test.UDPAddrv6",
			expected: "UDPAddrv6 muss eine gültige IPv6-UDP-Adresse sein",
		},
		{
			ns:       "Test.TCPAddr",
			expected: "TCPAddr muss eine gültige TCP-Adresse sein",
		},
		{
			ns:       "Test.TCPAddrv4",
			expected: "TCPAddrv4 muss eine gültige IPv4-TCP-Adresse sein",
		},
		{
			ns:       "Test.TCPAddrv6",
			expected: "TCPAddrv6 muss eine gültige IPv6-TCP-Adresse sein",
		},
		{
			ns:       "Test.CIDR",
			expected: "CIDR muss eine gültige CIDR-Notation enthalten",
		},
		{
			ns:       "Test.CIDRv4",
			expected: "CIDRv4 muss eine gültige CIDR-Notation für eine IPv4-Adresse enthalten",
		},
		{
			ns:       "Test.CIDRv6",
			expected: "CIDRv6 muss eine gültige CIDR-Notation für eine IPv6-Adresse enthalten",
		},
		{
			ns:       "Test.SSN",
			expected: "SSN muss eine gültige SSN-Nummer sein",
		},
		{
			ns:       "Test.IP",
			expected: "IP muss eine gültige IP-Adresse sein",
		},
		{
			ns:       "Test.IPv4",
			expected: "IPv4 muss eine gültige IPv4-Adresse sein",
		},
		{
			ns:       "Test.IPv6",
			expected: "IPv6 muss eine gültige IPv6-Adresse sein",
		},
		{
			ns:       "Test.DataURI",
			expected: "DataURI muss eine gültige Data-URI sein",
		},
		{
			ns:       "Test.Latitude",
			expected: "Latitude muss gültige Breitengradkoordinaten enthalten",
		},
		{
			ns:       "Test.Longitude",
			expected: "Longitude muss gültige Längengradkoordinaten enthalten",
		},
		{
			ns:       "Test.MultiByte",
			expected: "MultiByte darf nur Mehrbyte-Zeichen enthalten",
		},
		{
			ns:       "Test.ASCII",
			expected: "ASCII darf nur ASCII-Zeichen enthalten",
		},
		{
			ns:       "Test.PrintableASCII",
			expected: "PrintableASCII darf nur druckbare ASCII-Zeichen enthalten",
		},
		{
			ns:       "Test.UUID",
			expected: "UUID muss eine gültige UUID sein",
		},
		{
			ns:       "Test.UUID3",
			expected: "UUID3 muss eine gültige Version 3 UUID sein",
		},
		{
			ns:       "Test.UUID4",
			expected: "UUID4 muss eine gültige Version 4 UUID sein",
		},
		{
			ns:       "Test.UUID5",
			expected: "UUID5 muss eine gültige Version 5 UUID sein",
		},
		{
			ns:       "Test.ULID",
			expected: "ULID muss eine gültige ULID sein",
		},
		{
			ns:       "Test.ISBN",
			expected: "ISBN muss eine gültige ISBN-Nummer sein",
		},
		{
			ns:       "Test.ISBN10",
			expected: "ISBN10 muss eine gültige ISBN-10-Nummer sein",
		},
		{
			ns:       "Test.ISBN13",
			expected: "ISBN13 muss eine gültige ISBN-13-Nummer sein",
		},
		{
			ns:       "Test.ISSN",
			expected: "ISSN muss eine gültige ISSN-Nummer sein",
		},
		{
			ns:       "Test.Excludes",
			expected: "Excludes darf den Text 'text' nicht enthalten",
		},
		{
			ns:       "Test.ExcludesAll",
			expected: "ExcludesAll darf keines der folgenden Zeichen enthalten: '!@#$'",
		},
		{
			ns:       "Test.ExcludesRune",
			expected: "ExcludesRune darf die folgenden Runen nicht enthalten: '☻'",
		},
		{
			ns:       "Test.ContainsAny",
			expected: "ContainsAny muss mindestens eines der folgenden Zeichen enthalten: '!@#$'",
		},
		{
			ns:       "Test.Contains",
			expected: "Contains muss den Text 'purpose' enthalten",
		},
		{
			ns:       "Test.Base64",
			expected: "Base64 muss eine gültige Base64-Zeichenkette sein",
		},
		{
			ns:       "Test.Email",
			expected: "Email muss eine gültige E-Mail-Adresse sein",
		},
		{
			ns:       "Test.URL",
			expected: "URL muss eine gültige URL sein",
		},
		{
			ns:       "Test.URI",
			expected: "URI muss eine gültige URI sein",
		},
		{
			ns:       "Test.RGBColorString",
			expected: "RGBColorString muss eine gültige RGB-Farbe sein",
		},
		{
			ns:       "Test.RGBAColorString",
			expected: "RGBAColorString muss eine gültige RGBA-Farbe sein",
		},
		{
			ns:       "Test.HSLColorString",
			expected: "HSLColorString muss eine gültige HSL-Farbe sein",
		},
		{
			ns:       "Test.HSLAColorString",
			expected: "HSLAColorString muss eine gültige HSLA-Farbe sein",
		},
		{
			ns:       "Test.HexadecimalString",
			expected: "HexadecimalString muss eine gültige hexadezimale Zahl sein",
		},
		{
			ns:       "Test.HexColorString",
			expected: "HexColorString muss eine gültige Hexadezimalfarbe sein",
		},
		{
			ns:       "Test.NumberString",
			expected: "NumberString muss eine gültige Zahl sein",
		},
		{
			ns:       "Test.NumericString",
			expected: "NumericString muss eine gültige Zahl sein",
		},
		{
			ns:       "Test.AlphanumString",
			expected: "AlphanumString darf nur alphanumerische Zeichen enthalten",
		},
		{
			ns:       "Test.AlphaString",
			expected: "AlphaString darf nur alphabetische Zeichen enthalten",
		},
		{
			ns:       "Test.LtFieldString",
			expected: "LtFieldString muss kleiner als MaxString sein",
		},
		{
			ns:       "Test.LteFieldString",
			expected: "LteFieldString muss kleiner als oder gleich MaxString sein",
		},
		{
			ns:       "Test.GtFieldString",
			expected: "GtFieldString muss größer als MaxString sein",
		},
		{
			ns:       "Test.GteFieldString",
			expected: "GteFieldString muss größer als oder gleich MaxString sein",
		},
		{
			ns:       "Test.NeFieldString",
			expected: "NeFieldString darf nicht gleich EqFieldString sein",
		},
		{
			ns:       "Test.LtCSFieldString",
			expected: "LtCSFieldString muss kleiner als Inner.LtCSFieldString sein",
		},
		{
			ns:       "Test.LteCSFieldString",
			expected: "LteCSFieldString muss kleiner als oder gleich Inner.LteCSFieldString sein",
		},
		{
			ns:       "Test.GtCSFieldString",
			expected: "GtCSFieldString muss größer als Inner.GtCSFieldString sein",
		},
		{
			ns:       "Test.GteCSFieldString",
			expected: "GteCSFieldString muss größer als oder gleich Inner.GteCSFieldString sein",
		},
		{
			ns:       "Test.NeCSFieldString",
			expected: "NeCSFieldString darf nicht gleich Inner.NeCSFieldString sein",
		},
		{
			ns:       "Test.EqCSFieldString",
			expected: "EqCSFieldString muss gleich Inner.EqCSFieldString sein",
		},
		{
			ns:       "Test.EqFieldString",
			expected: "EqFieldString muss gleich MaxString sein",
		},
		{
			ns:       "Test.GteString",
			expected: "GteString muss mindestens 3 Zeichen lang sein",
		},
		{
			ns:       "Test.GteNumber",
			expected: "GteNumber muss 5,56 oder größer sein",
		},
		{
			ns:       "Test.GteMultiple",
			expected: "GteMultiple muss mindestens 2 Elemente enthalten",
		},
		{
			ns:       "Test.GteTime",
			expected: "GteTime muss vor dem aktuellen Datum und Uhrzeit liegen oder gleich sein",
		},
		{
			ns:       "Test.GtString",
			expected: "GtString muss größer als 3 Zeichen sein",
		},
		{
			ns:       "Test.GtNumber",
			expected: "GtNumber muss größer als 5,56 sein",
		},
		{
			ns:       "Test.GtMultiple",
			expected: "GtMultiple muss 2 Elemente oder mehr enthalten",
		},
		{
			ns:       "Test.GtTime",
			expected: "GtTime muss vor dem aktuellen Datum und Uhrzeit liegen",
		},
		{
			ns:       "Test.LteString",
			expected: "LteString darf maximal 3 Zeichen lang sein",
		},
		{
			ns:       "Test.LteNumber",
			expected: "LteNumber darf 5,56 oder weniger sein",
		},
		{
			ns:       "Test.LteMultiple",
			expected: "LteMultiple darf maximal 2 Elemente enthalten",
		},
		{
			ns:       "Test.LteTime",
			expected: "LteTime muss vor dem aktuellen Datum und Uhrzeit liegen oder gleich sein",
		},
		{
			ns:       "Test.LtString",
			expected: "LtString muss kleiner als 3 Zeichen sein",
		},
		{
			ns:       "Test.LtNumber",
			expected: "LtNumber muss kleiner als 5,56 sein",
		},
		{
			ns:       "Test.LtMultiple",
			expected: "LtMultiple muss 2 Elemente oder weniger enthalten",
		},
		{
			ns:       "Test.LtTime",
			expected: "LtTime muss vor dem aktuellen Datum und Uhrzeit liegen",
		},
		{
			ns:       "Test.NeString",
			expected: "NeString darf nicht gleich  sein",
		},
		{
			ns:       "Test.NeNumber",
			expected: "NeNumber darf nicht gleich 0.00 sein",
		},
		{
			ns:       "Test.NeMultiple",
			expected: "NeMultiple darf nicht gleich 0 sein",
		},
		{
			ns:       "Test.EqString",
			expected: "EqString ist nicht gleich 3",
		},
		{
			ns:       "Test.EqNumber",
			expected: "EqNumber ist nicht gleich 2.33",
		},
		{
			ns:       "Test.EqMultiple",
			expected: "EqMultiple ist nicht gleich 7",
		},
		{
			ns:       "Test.MaxString",
			expected: "MaxString darf maximal 3 Zeichen lang sein",
		},
		{
			ns:       "Test.MaxNumber",
			expected: "MaxNumber darf 1.113,00 oder weniger sein",
		},
		{
			ns:       "Test.MaxMultiple",
			expected: "MaxMultiple darf maximal 7 Elemente enthalten",
		},
		{
			ns:       "Test.MinString",
			expected: "MinString muss mindestens 1 Zeichen lang sein",
		},
		{
			ns:       "Test.MinNumber",
			expected: "MinNumber muss 1.113,00 oder größer sein",
		},
		{
			ns:       "Test.MinMultiple",
			expected: "MinMultiple muss mindestens 7 Elemente enthalten",
		},
		{
			ns:       "Test.LenString",
			expected: "LenString darf nur 1 Zeichen lang sein",
		},
		{
			ns:       "Test.LenNumber",
			expected: "LenNumber muss gleich 1.113,00 sein",
		},
		{
			ns:       "Test.LenMultiple",
			expected: "LenMultiple muss 7 Elemente enthalten",
		},
		{
			ns:       "Test.RequiredString",
			expected: "RequiredString ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredIf",
			expected: "RequiredIf ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredNumber",
			expected: "RequiredNumber ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredMultiple",
			expected: "RequiredMultiple ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredUnless",
			expected: "RequiredUnless ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredWith",
			expected: "RequiredWith ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredWithAll",
			expected: "RequiredWithAll ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredWithout",
			expected: "RequiredWithout ist ein Pflichtfeld",
		},
		{
			ns:       "Test.RequiredWithoutAll",
			expected: "RequiredWithoutAll ist ein Pflichtfeld",
		},
		{
			ns:       "Test.ExcludedIf",
			expected: "ExcludedIf ist ein ausgeschlossenes Feld",
		},
		{
			ns:       "Test.ExcludedUnless",
			expected: "ExcludedUnless ist ein ausgeschlossenes Feld",
		},
		{
			ns:       "Test.ExcludedWith",
			expected: "ExcludedWith ist ein ausgeschlossenes Feld",
		},
		{
			ns:       "Test.ExcludedWithAll",
			expected: "ExcludedWithAll ist ein ausgeschlossenes Feld",
		},
		{
			ns:       "Test.ExcludedWithout",
			expected: "ExcludedWithout ist ein ausgeschlossenes Feld",
		},
		{
			ns:       "Test.ExcludedWithoutAll",
			expected: "ExcludedWithoutAll ist ein ausgeschlossenes Feld",
		},
		{
			ns:       "Test.StrPtrMinLen",
			expected: "StrPtrMinLen muss mindestens 10 Zeichen lang sein",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen darf maximal 1 Zeichen lang sein",
		},
		{
			ns:       "Test.StrPtrLen",
			expected: "StrPtrLen darf nur 2 Zeichen lang sein",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt muss kleiner als 1 Zeichen sein",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte darf maximal 1 Zeichen lang sein",
		},
		{
			ns:       "Test.StrPtrGt",
			expected: "StrPtrGt muss größer als 10 Zeichen sein",
		},
		{
			ns:       "Test.StrPtrGte",
			expected: "StrPtrGte muss mindestens 10 Zeichen lang sein",
		},
		{
			ns:       "Test.OneOfString",
			expected: "OneOfString muss einer der folgenden sein: [red green]",
		},
		{
			ns:       "Test.OneOfInt",
			expected: "OneOfInt muss einer der folgenden sein: [5 63]",
		},
		{
			ns:       "Test.UniqueSlice",
			expected: "UniqueSlice darf nur einmal vorkommen",
		},
		{
			ns:       "Test.UniqueArray",
			expected: "UniqueArray darf nur einmal vorkommen",
		},
		{
			ns:       "Test.UniqueMap",
			expected: "UniqueMap darf nur einmal vorkommen",
		},
		{
			ns:       "Test.JSONString",
			expected: "JSONString muss eine gültige JSON-Zeichenkette sein",
		},
		{
			ns:       "Test.JWTString",
			expected: "JWTString muss eine gültige JWT-Zeichenkette sein",
		},
		{
			ns:       "Test.LowercaseString",
			expected: "LowercaseString darf nur Kleinbuchstaben enthalten",
		},
		{
			ns:       "Test.UppercaseString",
			expected: "UppercaseString darf nur Großbuchstaben enthalten",
		},
		{
			ns:       "Test.Datetime",
			expected: "Datetime entspricht nicht dem 2006-01-02-Format",
		},
		{
			ns:       "Test.PostCode",
			expected: "PostCode entspricht nicht dem Postleitzahlformat von SG",
		},
		{
			ns:       "Test.PostCodeByField",
			expected: "PostCodeByField entspricht nicht dem Postleitzahlformat des Feldes PostCodeCountry",
		},
		{
			ns:       "Test.BooleanString",
			expected: "BooleanString muss eine gültige Booleanwert sein",
		},
		{
			ns:       "Test.Image",
			expected: "Image muss ein Bild sein",
		},
		{
			ns:       "Test.CveString",
			expected: "CveString muss eine gültige CVE-Kennung sein",
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

		/*if fe == nil {
			t.Logf("no error found for %s", tt.ns)
		}*/

		NotEqual(t, fe, nil)
		/*if strings.Contains(tt.expected, "or greater") {
			got := fe.Translate(trans)
			Equal(t, tt.expected, got)
		}*/
		Equal(t, tt.expected, fe.Translate(trans))
	}
}
