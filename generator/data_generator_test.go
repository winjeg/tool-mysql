package generator

import (
	"testing"
)

func TestIntGenerator(t *testing.T) {
	IntRander.Get(0)
	IntRander.Get(TypeUInt)
	IntRander.Get(TypeUBigInt)
	IntRander.Get(TypeUTinyInt)
	IntRander.Get(TypeUSmallInt)
	IntRander.Get(TypeUMediumInt)
	IntRander.Get(TypeInt)
	IntRander.Get(TypeTinyInt)
	IntRander.Get(TypeSmallInt)
	IntRander.Get(TypeMediumInt)
	IntRander.Get(TypeBigInt)
}

func TestBitGenerator(t *testing.T) {
	BitRander.Get(10)
	BitRander.Get(100)
}

func TestFloatGenerator(t *testing.T) {
	FloatRander.Get(3, 5)
}

func TestTimeGenerator(t *testing.T) {
	TimeRander.Get()
}
func TestStringGenerator(t *testing.T) {
	StringRander.Get(100)
	WrapString(nil)
	s := `"a"`
	WrapString(&s)
}
