package ko

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	ko_locale "github.com/go-playground/locales/ko"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestTranslations(t *testing.T) {
	korean := ko_locale.New()
	uni := ut.New(korean, korean)
	trans, _ := uni.GetTranslator("ko")

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

	test.ASCII = "가나다라"
	test.PrintableASCII = "가나다라"

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
			expected: "IsColor은(는) 올바른 색이여야 합니다.",
		},
		{
			ns:       "Test.MAC",
			expected: "MAC은(는) 올바른 MAC 주소를 포함해야 합니다.",
		},
		{
			ns:       "Test.FQDN",
			expected: "FQDN은(는) 유효한 FQDN이어야 합니다.",
		},
		{
			ns:       "Test.IPAddr",
			expected: "IPAddr은(는) 해석 가능한 IP 주소여야 합니다.",
		},
		{
			ns:       "Test.IPAddrv4",
			expected: "IPAddrv4은(는) 해석 가능한 IPv4 주소여야 합니다.",
		},
		{
			ns:       "Test.IPAddrv6",
			expected: "IPAddrv6은(는) 해석 가능한 IPv6 주소여야 합니다.",
		},
		{
			ns:       "Test.UDPAddr",
			expected: "UDPAddr은(는) 올바른 UDP 주소여야 합니다.",
		},
		{
			ns:       "Test.UDPAddrv4",
			expected: "UDPAddrv4은(는) 올바른 IPv4의 UDP 주소여야 합니다.",
		},
		{
			ns:       "Test.UDPAddrv6",
			expected: "UDPAddrv6은(는) 올바른 IPv6의 UDP 주소여야 합니다.",
		},
		{
			ns:       "Test.TCPAddr",
			expected: "TCPAddr은(는) 올바른 TCP 주소여야 합니다.",
		},
		{
			ns:       "Test.TCPAddrv4",
			expected: "TCPAddrv4은(는) 올바른 IPv4의 TCP 주소여야 합니다.",
		},
		{
			ns:       "Test.TCPAddrv6",
			expected: "TCPAddrv6은(는) 올바른 IPv6의 TCP 주소여야 합니다.",
		},
		{
			ns:       "Test.CIDR",
			expected: "CIDR은(는) 올바른 CIDR 표기를 포함해야 합니다.",
		},
		{
			ns:       "Test.CIDRv4",
			expected: "CIDRv4은(는) IPv4 주소의 올바른 CIDR 표기를 포함해야 합니다.",
		},
		{
			ns:       "Test.CIDRv6",
			expected: "CIDRv6은(는) IPv6 주소의 올바른 CIDR 표기를 포함해야 합니다.",
		},
		{
			ns:       "Test.SSN",
			expected: "SSN은(는) 올바른 사회 보장 번호여야 합니다.",
		},
		{
			ns:       "Test.IP",
			expected: "IP은(는) 올바른 IP 주소여야 합니다.",
		},
		{
			ns:       "Test.IPv4",
			expected: "IPv4은(는) 올바른 IPv4 주소여야 합니다.",
		},
		{
			ns:       "Test.IPv6",
			expected: "IPv6은(는) 올바른 IPv6 주소여야 합니다.",
		},
		{
			ns:       "Test.DataURI",
			expected: "DataURI은(는) 올바른 데이터 URI를 포함해야 합니다.",
		},
		{
			ns:       "Test.Latitude",
			expected: "Latitude은(는) 올바른 위도 좌표를 포함해야 합니다.",
		},
		{
			ns:       "Test.Longitude",
			expected: "Longitude은(는) 올바른 경도 좌표를 포함해야 합니다.",
		},
		{
			ns:       "Test.MultiByte",
			expected: "MultiByte은(는) 멀티바이트 문자를 포함해야 합니다.",
		},
		{
			ns:       "Test.ASCII",
			expected: "ASCII은(는) ASCII 문자만 포함해야 합니다.",
		},
		{
			ns:       "Test.PrintableASCII",
			expected: "PrintableASCII은(는) 인쇄 가능한 ASCII 문자만 포함해야 합니다.",
		},
		{
			ns:       "Test.UUID",
			expected: "UUID은(는) 올바른 UUID여야 합니다.",
		},
		{
			ns:       "Test.UUID3",
			expected: "UUID3은(는) 버전 3의 올바른 UUID여야 합니다.",
		},
		{
			ns:       "Test.UUID4",
			expected: "UUID4은(는) 버전 4의 올바른 UUID여야 합니다.",
		},
		{
			ns:       "Test.UUID5",
			expected: "UUID5은(는) 버전 5의 올바른 UUID여야 합니다.",
		},
		{
			ns:       "Test.ULID",
			expected: "ULID은(는) 올바른 ULID여야 합니다.",
		},
		{
			ns:       "Test.ISBN",
			expected: "ISBN은(는) 올바른 ISBN 번호여야 합니다.",
		},
		{
			ns:       "Test.ISBN10",
			expected: "ISBN10은(는) 올바른 ISBN-10 번호여야 합니다.",
		},
		{
			ns:       "Test.ISBN13",
			expected: "ISBN13은(는) 올바른 ISBN-13 번호여야 합니다.",
		},
		{
			ns:       "Test.ISSN",
			expected: "ISSN은(는) 올바른 ISSN 번호여야 합니다.",
		},
		{
			ns:       "Test.Excludes",
			expected: "Excludes에는 'text'라는 텍스트를 포함할 수 없습니다.",
		},
		{
			ns:       "Test.ExcludesAll",
			expected: "ExcludesAll에는 '!@#$' 중 어느 것도 포함할 수 없습니다.",
		},
		{
			ns:       "Test.ExcludesRune",
			expected: "ExcludesRune에는 '☻'을(를) 포함할 수 없습니다.",
		},
		{
			ns:       "Test.ContainsAny",
			expected: "ContainsAny은(는) '!@#$' 중 최소 하나를 포함해야 합니다.",
		},
		{
			ns:       "Test.Contains",
			expected: "Contains은(는) 'purpose'을(를) 포함해야 합니다.",
		},
		{
			ns:       "Test.Base64",
			expected: "Base64은(는) 올바른 Base64 문자열여야 합니다.",
		},
		{
			ns:       "Test.Email",
			expected: "Email은(는) 올바른 이메일 주소여야 합니다.",
		},
		{
			ns:       "Test.URL",
			expected: "URL은(는) 올바른 URL여야 합니다.",
		},
		{
			ns:       "Test.URI",
			expected: "URI은(는) 올바른 URI여야 합니다.",
		},
		{
			ns:       "Test.RGBColorString",
			expected: "RGBColorString은(는) 올바른 RGB 색상 코드여야 합니다.",
		},
		{
			ns:       "Test.RGBAColorString",
			expected: "RGBAColorString은(는) 올바른 RGBA 색상 코드여야 합니다.",
		},
		{
			ns:       "Test.HSLColorString",
			expected: "HSLColorString은(는) 올바른 HSL 색상 코드여야 합니다.",
		},
		{
			ns:       "Test.HSLAColorString",
			expected: "HSLAColorString은(는) 올바른 HSLA 색상 코드여야 합니다.",
		},
		{
			ns:       "Test.HexadecimalString",
			expected: "HexadecimalString은(는) 올바른 16진수 표기여야 합니다.",
		},
		{
			ns:       "Test.HexColorString",
			expected: "HexColorString은(는) 올바른 HEX 색상 코드여야 합니다.",
		},
		{
			ns:       "Test.NumberString",
			expected: "NumberString은(는) 올바른 수여야 합니다.",
		},
		{
			ns:       "Test.NumericString",
			expected: "NumericString은(는) 올바른 숫자여야 합니다.",
		},
		{
			ns:       "Test.AlphanumString",
			expected: "AlphanumString은(는) 알파벳과 숫자만 포함할 수 있습니다.",
		},
		{
			ns:       "Test.AlphaString",
			expected: "AlphaString은(는) 알파벳만 포함할 수 있습니다.",
		},
		{
			ns:       "Test.LtFieldString",
			expected: "LtFieldString은(는) MaxString보다 작아야 합니다.",
		},
		{
			ns:       "Test.LteFieldString",
			expected: "LteFieldString은(는) MaxString 이하여야 합니다.",
		},
		{
			ns:       "Test.GtFieldString",
			expected: "GtFieldString은(는) MaxString보다 커야 합니다.",
		},
		{
			ns:       "Test.GteFieldString",
			expected: "GteFieldString은(는) MaxString 이상여야 합니다.",
		},
		{
			ns:       "Test.NeFieldString",
			expected: "NeFieldString은(는) EqFieldString와(과) 달라야 합니다.",
		},
		{
			ns:       "Test.LtCSFieldString",
			expected: "LtCSFieldString은(는) Inner.LtCSFieldString보다 작아야 합니다.",
		},
		{
			ns:       "Test.LteCSFieldString",
			expected: "LteCSFieldString은(는) Inner.LteCSFieldString 이하여야 합니다.",
		},
		{
			ns:       "Test.GtCSFieldString",
			expected: "GtCSFieldString은(는) Inner.GtCSFieldString보다 커야 합니다.",
		},
		{
			ns:       "Test.GteCSFieldString",
			expected: "GteCSFieldString은(는) Inner.GteCSFieldString 이상여야 합니다.",
		},
		{
			ns:       "Test.NeCSFieldString",
			expected: "NeCSFieldString은(는) Inner.NeCSFieldString와(과) 달라야 합니다.",
		},
		{
			ns:       "Test.EqCSFieldString",
			expected: "EqCSFieldString은(는) Inner.EqCSFieldString와(과) 같아야 합니다.",
		},
		{
			ns:       "Test.EqFieldString",
			expected: "EqFieldString은(는) MaxString와(과) 같아야 합니다.",
		},
		{
			ns:       "Test.GteString",
			expected: "GteString의 길이는 최소 3자 이상여야 합니다.",
		},
		{
			ns:       "Test.GteNumber",
			expected: "GteNumber은(는) 5.56 이상여야 합니다.",
		},
		{
			ns:       "Test.GteMultiple",
			expected: "GteMultiple은(는) 최소 2개의 항목을 포함해야 합니다.",
		},
		{
			ns:       "Test.GteTime",
			expected: "GteTime은(는) 현재 시간 이후이어야 합니다.",
		},
		{
			ns:       "Test.GtString",
			expected: "GtString의 길이는 3자보다 길어야 합니다.",
		},
		{
			ns:       "Test.GtNumber",
			expected: "GtNumber은(는) 5.56보다 커야 합니다.",
		},
		{
			ns:       "Test.GtMultiple",
			expected: "GtMultiple은(는) 2개의 항목보다 많은 항목을 포함해야 합니다.",
		},
		{
			ns:       "Test.GtTime",
			expected: "GtTime은(는) 현재 시간 이후이어야 합니다.",
		},
		{
			ns:       "Test.LteString",
			expected: "LteString의 길이는 최대 3자여야 합니다.",
		},
		{
			ns:       "Test.LteNumber",
			expected: "LteNumber은(는) 5.56 이하여야 합니다.",
		},
		{
			ns:       "Test.LteMultiple",
			expected: "LteMultiple은(는) 최대 2개의 항목여야 합니다.",
		},
		{
			ns:       "Test.LteTime",
			expected: "LteTime은(는) 현재 시간보다 이전이어야 합니다.",
		},
		{
			ns:       "Test.LtString",
			expected: "LtString의 길이는 3자보다 작아야 합니다.",
		},
		{
			ns:       "Test.LtNumber",
			expected: "LtNumber은(는) 5.56보다 작아야 합니다.",
		},
		{
			ns:       "Test.LtMultiple",
			expected: "LtMultiple은(는) 2개의 항목보다 적은 항목여야 합니다.",
		},
		{
			ns:       "Test.LtTime",
			expected: "LtTime은(는) 현재 시간보다 이전이어야 합니다.",
		},
		{
			ns:       "Test.NeString",
			expected: "NeString은(는) 와(과) 달라야 합니다.",
		},
		{
			ns:       "Test.NeNumber",
			expected: "NeNumber은(는) 0.00와(과) 달라야 합니다.",
		},
		{
			ns:       "Test.NeMultiple",
			expected: "NeMultiple의 항목 수는 0개와(과) 달라야 합니다.",
		},
		{
			ns:       "Test.EqString",
			expected: "EqString은(는) 3와(과) 같아야 합니다.",
		},
		{
			ns:       "Test.EqNumber",
			expected: "EqNumber은(는) 2.33와(과) 같아야 합니다.",
		},
		{
			ns:       "Test.EqMultiple",
			expected: "EqMultiple은(는) 7와(과) 같아야 합니다.",
		},
		{
			ns:       "Test.MaxString",
			expected: "MaxString의 길이는 최대 3자여야 합니다.",
		},
		{
			ns:       "Test.MaxNumber",
			expected: "MaxNumber은(는) 1,113.00 이하여야 합니다.",
		},
		{
			ns:       "Test.MaxMultiple",
			expected: "MaxMultiple은(는) 최대 7개의 항목여야 합니다.",
		},
		{
			ns:       "Test.MinString",
			expected: "MinString의 길이는 최소 1자여야 합니다.",
		},
		{
			ns:       "Test.MinNumber",
			expected: "MinNumber은(는) 1,113.00 이상여야 합니다.",
		},
		{
			ns:       "Test.MinMultiple",
			expected: "MinMultiple은(는) 최소 7개의 항목을 포함해야 합니다.",
		},
		{
			ns:       "Test.LenString",
			expected: "LenString의 길이는 1자여야 합니다.",
		},
		{
			ns:       "Test.LenNumber",
			expected: "LenNumber은(는) 1,113.00와(과) 같아야 합니다.",
		},
		{
			ns:       "Test.LenMultiple",
			expected: "LenMultiple은(는) 7개의 항목을 포함해야 합니다.",
		},
		{
			ns:       "Test.RequiredString",
			expected: "RequiredString은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredIf",
			expected: "RequiredIf은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredNumber",
			expected: "RequiredNumber은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredMultiple",
			expected: "RequiredMultiple은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredUnless",
			expected: "RequiredUnless은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredWith",
			expected: "RequiredWith은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredWithAll",
			expected: "RequiredWithAll은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredWithout",
			expected: "RequiredWithout은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.RequiredWithoutAll",
			expected: "RequiredWithoutAll은(는) 필수 필드입니다.",
		},
		{
			ns:       "Test.ExcludedIf",
			expected: "ExcludedIf은(는) 제외된 필드입니다.",
		},
		{
			ns:       "Test.ExcludedUnless",
			expected: "ExcludedUnless은(는) 제외된 필드입니다.",
		},
		{
			ns:       "Test.ExcludedWith",
			expected: "ExcludedWith은(는) 제외된 필드입니다.",
		},
		{
			ns:       "Test.ExcludedWithAll",
			expected: "ExcludedWithAll은(는) 제외된 필드입니다.",
		},
		{
			ns:       "Test.ExcludedWithout",
			expected: "ExcludedWithout은(는) 제외된 필드입니다.",
		},
		{
			ns:       "Test.ExcludedWithoutAll",
			expected: "ExcludedWithoutAll은(는) 제외된 필드입니다.",
		},
		{
			ns:       "Test.StrPtrMinLen",
			expected: "StrPtrMinLen의 길이는 최소 10자여야 합니다.",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen의 길이는 최대 1자여야 합니다.",
		},
		{
			ns:       "Test.StrPtrLen",
			expected: "StrPtrLen의 길이는 2자여야 합니다.",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt의 길이는 1자보다 작아야 합니다.",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte의 길이는 최대 1자여야 합니다.",
		},
		{
			ns:       "Test.StrPtrGt",
			expected: "StrPtrGt의 길이는 10자보다 길어야 합니다.",
		},
		{
			ns:       "Test.StrPtrGte",
			expected: "StrPtrGte의 길이는 최소 10자 이상여야 합니다.",
		},
		{
			ns:       "Test.OneOfString",
			expected: "OneOfString은(는) [red green] 중 하나여야 합니다.",
		},
		{
			ns:       "Test.OneOfInt",
			expected: "OneOfInt은(는) [5 63] 중 하나여야 합니다.",
		},
		{
			ns:       "Test.UniqueSlice",
			expected: "UniqueSlice은(는) 고유한 값만 포함해야 합니다.",
		},
		{
			ns:       "Test.UniqueArray",
			expected: "UniqueArray은(는) 고유한 값만 포함해야 합니다.",
		},
		{
			ns:       "Test.UniqueMap",
			expected: "UniqueMap은(는) 고유한 값만 포함해야 합니다.",
		},
		{
			ns:       "Test.JSONString",
			expected: "JSONString은(는) 올바른 JSON 문자열여야 합니다.",
		},
		{
			ns:       "Test.JWTString",
			expected: "JWTString은(는) 올바른 JWT 문자열여야 합니다.",
		},
		{
			ns:       "Test.LowercaseString",
			expected: "LowercaseString은(는) 소문자여야 합니다.",
		},
		{
			ns:       "Test.UppercaseString",
			expected: "UppercaseString은(는) 대문자여야 합니다.",
		},
		{
			ns:       "Test.Datetime",
			expected: "Datetime은(는) 2006-01-02 형식과 일치해야 합니다.",
		},
		{
			ns:       "Test.PostCode",
			expected: "PostCode은(는) 국가 코드 SG의 우편번호 형식과 일치해야 합니다.",
		},
		{
			ns:       "Test.PostCodeByField",
			expected: "PostCodeByField은(는) PostCodeCountry 필드에 지정된 국가 코드의 우편번호 형식과 일치해야 합니다.",
		},
		{
			ns:       "Test.BooleanString",
			expected: "BooleanString은(는) 올바른 부울 값여야 합니다.",
		},
		{
			ns:       "Test.Image",
			expected: "Image은(는) 유효한 이미지여야 합니다.",
		},
		{
			ns:       "Test.CveString",
			expected: "CveString은(는) 유효한 CVE 식별자여야 합니다.",
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
