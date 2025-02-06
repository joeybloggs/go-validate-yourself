package th

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
	thai "github.com/go-playground/locales/th"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func TestTranslations(t *testing.T) {
	th := thai.New()
	uni := ut.New(th, th)
	trans, _ := uni.GetTranslator("th")

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
		PostCodeByField   string `validate:"postcode_iso3166_alpha2_field=PostCodeCountry"`
		BooleanString     string `validate:"boolean"`
		Image             string `validate:"image"`
		CveString         string `validate:"cve"`
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
			expected: "IsColor ต้องเป็นเลขสี",
		},
		{
			ns:       "Test.MAC",
			expected: "MAC ต้องเป็น MAC address",
		},
		{
			ns:       "Test.FQDN",
			expected: "FQDN ต้องเป็น FQDN",
		},
		{
			ns:       "Test.IPAddr",
			expected: "IPAddr ต้องเป็น IP address ที่เข้าถึงได้",
		},
		{
			ns:       "Test.IPAddrv4",
			expected: "IPAddrv4 ต้องเป็น IPv4 address ที่เข้าถึงได้",
		},
		{
			ns:       "Test.IPAddrv6",
			expected: "IPAddrv6 ต้องเป็น IPv6 address ที่เข้าถึงได้",
		},
		{
			ns:       "Test.UDPAddr",
			expected: "UDPAddr ต้องเป็น UDP address",
		},
		{
			ns:       "Test.UDPAddrv4",
			expected: "UDPAddrv4 ต้องเป็น IPv4 UDP address",
		},
		{
			ns:       "Test.UDPAddrv6",
			expected: "UDPAddrv6 ต้องเป็น IPv6 UDP address",
		},
		{
			ns:       "Test.TCPAddr",
			expected: "TCPAddr ต้องเป็น TCP address",
		},
		{
			ns:       "Test.TCPAddrv4",
			expected: "TCPAddrv4 ต้องเป็น IPv4 TCP address",
		},
		{
			ns:       "Test.TCPAddrv6",
			expected: "TCPAddrv6 ต้องเป็น IPv6 TCP address",
		},
		{
			ns:       "Test.CIDR",
			expected: "CIDR ต้องเป็น CIDR notation",
		},
		{
			ns:       "Test.CIDRv4",
			expected: "CIDRv4 ต้องเป็น CIDR notation สำหรับ an IPv4 address",
		},
		{
			ns:       "Test.CIDRv6",
			expected: "CIDRv6 ต้องเป็น CIDR notation สำหรับ an IPv6 address",
		},
		{
			ns:       "Test.SSN",
			expected: "SSN ต้องเป็นตัวเลข SSN",
		},
		{
			ns:       "Test.IP",
			expected: "IP ต้องเป็น IP address",
		},
		{
			ns:       "Test.IPv4",
			expected: "IPv4 ต้องเป็น IPv4 address",
		},
		{
			ns:       "Test.IPv6",
			expected: "IPv6 ต้องเป็น IPv6 address",
		},
		{
			ns:       "Test.DataURI",
			expected: "DataURI must contain a valid Data URI",
		},
		{
			ns:       "Test.Latitude",
			expected: "Latitude ต้องเป็นละติจูด",
		},
		{
			ns:       "Test.Longitude",
			expected: "Longitude ต้องเป็นลองจิจูด",
		},
		{
			ns:       "Test.UUID",
			expected: "UUID ต้องเป็น UUID",
		},
		{
			ns:       "Test.UUID3",
			expected: "UUID3 ต้องเป็น version 3 UUID",
		},
		{
			ns:       "Test.UUID4",
			expected: "UUID4 ต้องเป็น version 4 UUID",
		},
		{
			ns:       "Test.UUID5",
			expected: "UUID5 ต้องเป็น version 5 UUID",
		},
		{
			ns:       "Test.ULID",
			expected: "ULID must be a valid ULID",
		},
		{
			ns:       "Test.ISBN",
			expected: "ISBN ต้องเป็นตัวเลข ISBN",
		},
		{
			ns:       "Test.ISBN10",
			expected: "ISBN10 ต้องเป็นตัวเลข ISBN-10",
		},
		{
			ns:       "Test.ISBN13",
			expected: "ISBN13 ต้องเป็นตัวเลข ISBN-13",
		},
		{
			ns:       "Test.ISSN",
			expected: "ISSN ต้องเป็นตัวเลข ISSN",
		},
		{
			ns:       "Test.Excludes",
			expected: "Excludes ต้องไม่มี 'text'",
		},
		{
			ns:       "Test.ExcludesAll",
			expected: "ExcludesAll ต้องไม่มีอักขระ '!@#$' ทั้งหมด",
		},
		{
			ns:       "Test.ExcludesRune",
			expected: "ExcludesRune ต้องไม่มี '☻'",
		},
		{
			ns:       "Test.Contains",
			expected: "Contains ต้องมี 'purpose'",
		},
		{
			ns:       "Test.Base64",
			expected: "Base64 ต้องเป็น Base64 เท่านั้น",
		},
		{
			ns:       "Test.Email",
			expected: "Email ต้องเป็นอีเมลเท่านั้น",
		},
		{
			ns:       "Test.URL",
			expected: "URL ต้องเป็น URL เท่านั้น",
		},
		{
			ns:       "Test.URI",
			expected: "URI ต้องเป็น URI เท่านั้น",
		},
		{
			ns:       "Test.RGBColorString",
			expected: "RGBColorString ต้องเป็นเลขสี RGB เท่านั้น",
		},
		{
			ns:       "Test.RGBAColorString",
			expected: "RGBAColorString ต้องเป็นเลขสี RGBA เท่านั้น",
		},
		{
			ns:       "Test.HSLColorString",
			expected: "HSLColorString ต้องเป็นเลขสี HSL เท่านั้น",
		},
		{
			ns:       "Test.HSLAColorString",
			expected: "HSLAColorString ต้องเป็นเลขสี HSLA เท่านั้น",
		},
		{
			ns:       "Test.HexadecimalString",
			expected: "HexadecimalString ต้องเป็นค่าตัวเลขฐาน 16 เท่านั้น",
		},
		{
			ns:       "Test.HexColorString",
			expected: "HexColorString ต้องเป็นเลขสีฐาน 16 เท่านั้น",
		},
		{
			ns:       "Test.NumberString",
			expected: "NumberString ต้องเป็นตัวเลขเท่านั้น",
		},
		{
			ns:       "Test.NumericString",
			expected: "NumericString ต้องเป็นค่าตัวเลขเท่านั้น",
		},
		{
			ns:       "Test.LtFieldString",
			expected: "LtFieldString ต้องมีค่าน้อยกว่า MaxString",
		},
		{
			ns:       "Test.LteFieldString",
			expected: "LteFieldString ต้องมีค่าน้อยกว่าหรือเท่ากับ MaxString",
		},
		{
			ns:       "Test.GtFieldString",
			expected: "GtFieldString ต้องมีค่ามากกว่า MaxString",
		},
		{
			ns:       "Test.GteFieldString",
			expected: "GteFieldString ต้องมีค่ามากกว่าหรือเท่ากับ MaxString",
		},
		{
			ns:       "Test.NeFieldString",
			expected: "NeFieldString ต้องไม่เท่ากับ EqFieldString",
		},
		{
			ns:       "Test.LtCSFieldString",
			expected: "LtCSFieldString ต้องมีค่าน้อยกว่า Inner.LtCSFieldString",
		},
		{
			ns:       "Test.LteCSFieldString",
			expected: "LteCSFieldString ต้องมีค่าน้อยกว่าหรือเท่ากับ Inner.LteCSFieldString",
		},
		{
			ns:       "Test.GtCSFieldString",
			expected: "GtCSFieldString ต้องมีค่ามากกว่า Inner.GtCSFieldString",
		},
		{
			ns:       "Test.GteCSFieldString",
			expected: "GteCSFieldString ต้องมีค่ามากกว่าหรือเท่ากับ Inner.GteCSFieldString",
		},
		{
			ns:       "Test.NeCSFieldString",
			expected: "NeCSFieldString ต้องไม่เท่ากับ Inner.NeCSFieldString",
		},
		{
			ns:       "Test.EqCSFieldString",
			expected: "EqCSFieldString ต้องเท่ากับ Inner.EqCSFieldString",
		},
		{
			ns:       "Test.EqFieldString",
			expected: "EqFieldString ต้องเท่ากับ MaxString",
		},
		{
			ns:       "Test.GteNumber",
			expected: "GteNumber ต้องมีค่ามากกว่า 5.56",
		},
		{
			ns:       "Test.GteMultiple",
			expected: "GteMultiple ต้องมีอย่างน้อย 2 ไอเทม",
		},
		{
			ns:       "Test.GteTime",
			expected: "GteTime ต้องมีค่ามากกว่าหรือเท่ากับเวลาปัจจุบัน",
		},
		{
			ns:       "Test.GtNumber",
			expected: "GtNumber ต้องมีค่ามากกว่า 5.56",
		},
		{
			ns:       "Test.GtMultiple",
			expected: "GtMultiple ต้องมีมากกว่า 2 ไอเทม",
		},
		{
			ns:       "Test.GtTime",
			expected: "GtTime ต้องมีค่ามากกว่าเวลาปัจจุบัน",
		},
		{
			ns:       "Test.LteNumber",
			expected: "LteNumber ต้องมีค่าน้อยกว่าหรือเท่ากับ 5.56",
		},
		{
			ns:       "Test.LteMultiple",
			expected: "LteMultiple ต้องมีไม่เกิน 2 ไอเทม",
		},
		{
			ns:       "Test.LteTime",
			expected: "LteTime ต้องมีค่าน้อยกว่าหรือเท่ากับเวลาปัจจุบัน",
		},
		{
			ns:       "Test.LtNumber",
			expected: "LtNumber ต้องมีค่าน้อยกว่า 5.56",
		},
		{
			ns:       "Test.LtMultiple",
			expected: "LtMultiple ต้องมีน้อยกว่า 2 ไอเทม",
		},
		{
			ns:       "Test.LtTime",
			expected: "LtTime ต้องมีค่าน้อยกว่าเวลาปัจจุบัน",
		},
		{
			ns:       "Test.NeString",
			expected: "NeString ต้องไม่เท่ากับ ",
		},
		{
			ns:       "Test.NeNumber",
			expected: "NeNumber ต้องไม่เท่ากับ 0.00",
		},
		{
			ns:       "Test.NeMultiple",
			expected: "NeMultiple ต้องไม่เท่ากับ 0",
		},
		{
			ns:       "Test.EqString",
			expected: "EqString ไม่เท่ากับ 3",
		},
		{
			ns:       "Test.EqNumber",
			expected: "EqNumber ไม่เท่ากับ 2.33",
		},
		{
			ns:       "Test.EqMultiple",
			expected: "EqMultiple ไม่เท่ากับ 7",
		},
		{
			ns:       "Test.MaxNumber",
			expected: "MaxNumber ต้องมีค่าน้อยกว่าหรือเท่ากับ 1,113.00",
		},
		{
			ns:       "Test.MaxMultiple",
			expected: "MaxMultiple ต้องมีไม่เกิน 7 ไอเทม",
		},
		{
			ns:       "Test.MinString",
			expected: "MinString ต้องมีความยาวอย่างน้อย 1 ตัวอักษร",
		},
		{
			ns:       "Test.MinNumber",
			expected: "MinNumber ต้องมีค่ามากกว่า 1,113.00",
		},
		{
			ns:       "Test.MinMultiple",
			expected: "MinMultiple ต้องมีอย่างน้อย 7 ไอเทม",
		},
		{
			ns:       "Test.LenString",
			expected: "LenString ต้องมีความยาว 1 ตัวอักษร",
		},
		{
			ns:       "Test.LenNumber",
			expected: "LenNumber ต้องเท่ากับ 1,113.00",
		},
		{
			ns:       "Test.LenMultiple",
			expected: "LenMultiple must contain 7 ไอเทม",
		},
		{
			ns:       "Test.RequiredString",
			expected: "โปรดระบุ RequiredString",
		},
		{
			ns:       "Test.RequiredIf",
			expected: "โปรดระบุ RequiredIf",
		},
		{
			ns:       "Test.RequiredNumber",
			expected: "โปรดระบุ RequiredNumber",
		},
		{
			ns:       "Test.RequiredMultiple",
			expected: "โปรดระบุ RequiredMultiple",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen ต้องมีความยาวไม่เกิน 1 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt ต้องมีความยาวน้อยกว่า 1 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte ต้องมีความยาวไม่เกิน 1 ตัวอักษร",
		},
		{
			ns:       "Test.OneOfString",
			expected: "OneOfString ต้องอยู่ใน [red green]",
		},
		{
			ns:       "Test.OneOfInt",
			expected: "OneOfInt ต้องอยู่ใน [5 63]",
		},
		{
			ns:       "Test.UniqueSlice",
			expected: "UniqueSlice ต้องมีข้อมูลไม่ซ้ำ",
		},
		{
			ns:       "Test.UniqueArray",
			expected: "UniqueArray ต้องมีข้อมูลไม่ซ้ำ",
		},
		{
			ns:       "Test.UniqueMap",
			expected: "UniqueMap ต้องมีข้อมูลไม่ซ้ำ",
		},
		{
			ns:       "Test.JSONString",
			expected: "JSONString ต้องเป็น json string",
		},
		{
			ns:       "Test.JWTString",
			expected: "JWTString ต้องเป็น jwt string",
		},
		{
			ns:       "Test.LowercaseString",
			expected: "LowercaseString ต้องเป็นตัวพิมพ์เล็ก",
		},
		{
			ns:       "Test.UppercaseString",
			expected: "UppercaseString ต้องเป็นตัวพิมพ์ใหญ่",
		},
		{
			ns:       "Test.Datetime",
			expected: "Datetime ไม่ตรงกับรูปแบบ 2006-01-02",
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
			expected: "BooleanString ต้องเป็น boolean",
		},
		{
			ns:       "Test.Image",
			expected: "Image ต้องเป็นรูปภาพ",
		},
		{
			ns:       "Test.CveString",
			expected: "CveString ต้องเป็นรูปแบบ cve",
		},
		{
			ns:       "Test.StrPtrMinLen",
			expected: "StrPtrMinLen ต้องมีความยาวอย่างน้อย 10 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrMaxLen",
			expected: "StrPtrMaxLen ต้องมีความยาวไม่เกิน 1 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrLen",
			expected: "StrPtrLen ต้องมีความยาว 2 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrLt",
			expected: "StrPtrLt ต้องมีความยาวน้อยกว่า 1 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrLte",
			expected: "StrPtrLte ต้องมีความยาวไม่เกิน 1 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrGt",
			expected: "StrPtrGt ต้องมีความยาวมากกว่า 10 ตัวอักษร",
		},
		{
			ns:       "Test.StrPtrGte",
			expected: "StrPtrGte ต้องมีความยาวอย่างน้อย 10 ตัวอักษร",
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
