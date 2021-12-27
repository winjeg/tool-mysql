package generator

import (
	"fmt"
	"math"
)

// generate most frequently used data types

const (
	maxBitNum = 64

	// types
	_       = iota
	TypeInt = iota
	TypeTinyInt
	TypeSmallInt
	TypeMediumInt
	TypeBigInt
	TypeUInt
	TypeUTinyInt
	TypeUSmallInt
	TypeUMediumInt
	TypeUBigInt
)

var (
	IntRander = IntGenerator{
		RandomInt:  RandomInt,
		RandomUInt: RandomUInt,
	}
	StringRander = StringGenerator{
		GenerateString: RandomReadable,
	}

	BitRander = BitGenerator{}

	TimeRander = TimeGenerator{}

	FloatRander = FloatGenerator{
		RandomFloat: RandomFloat,
	}
)

// IntGenerator generate all kinds of int types
type IntGenerator struct {
	RandomInt  func(int64) int64
	RandomUInt func(uint64) uint64
}

func (g *IntGenerator) TinyInt() int {
	return int(g.RandomInt(math.MaxInt8))
}

func (g *IntGenerator) SmallInt() int {
	return int(g.RandomInt(math.MaxInt16))
}

func (g *IntGenerator) MediumInt() int {
	return int(g.RandomInt(1<<23 - 1))
}

func (g *IntGenerator) Int() int {
	return int(g.RandomInt(math.MaxInt32))
}

func (g *IntGenerator) BigInt() int64 {
	return int64(g.RandomInt(math.MaxInt64))
}

func (g *IntGenerator) UTinyInt() int {
	return int(g.RandomUInt(math.MaxInt8))
}

func (g *IntGenerator) USmallInt() uint {
	return uint(g.RandomUInt(math.MaxInt16))
}

func (g *IntGenerator) UMediumInt() uint {
	return uint(g.RandomUInt(1<<23 - 1))
}

func (g *IntGenerator) UInt() uint {
	return uint(g.RandomUInt(math.MaxInt32))
}

func (g *IntGenerator) UBigInt() uint64 {
	return uint64(g.RandomUInt(math.MaxInt64))
}

// Get get random int values from the generator
func (g *IntGenerator) Get(t int) string {
	switch t {
	case TypeTinyInt:
		return formatInt(g.TinyInt())
	case TypeSmallInt:
		return formatInt(g.SmallInt())
	case TypeMediumInt:
		return formatInt(g.MediumInt())
	case TypeInt:
		return formatInt(g.Int())
	case TypeBigInt:
		return formatInt(g.BigInt())
	case TypeUTinyInt:
		return formatInt(g.UTinyInt())
	case TypeUSmallInt:
		return formatInt(g.USmallInt())
	case TypeUMediumInt:
		return formatInt(g.UMediumInt())
	case TypeUInt:
		return formatInt(g.UInt())
	case TypeUBigInt:
		return formatInt(g.UBigInt())
	default:
		return formatInt(g.Int())
	}
}

// format int to mysql required format
func formatInt(i interface{}) string {
	return fmt.Sprintf("'%d'", i)
}

// FloatGenerator generate float numbers
type FloatGenerator struct {
	RandomFloat func(int, int) string
}

// GetFloat get float with n digits of int and f for digits of fraction
func (fg *FloatGenerator) GetFloat(n, f int) string {
	return fg.RandomFloat(n, f)
}

// Get get float with n digits of int and f for digits of fraction
func (fg *FloatGenerator) Get(n, f int) string {
	return formatFloat(fg.GetFloat(n, f))
}

// format float to mysql required format
func formatFloat(i interface{}) string {
	return fmt.Sprintf("'%s'", i)
}

// StringGenerator string generator
type StringGenerator struct {
	GenerateString func(int) string
}

// Get get string wrapped by quote
func (sg *StringGenerator) Get(n int) string {
	rawStr := sg.GenerateString(n)
	return fmt.Sprintf("'%s'", WrapString(&rawStr))
}

// BitGenerator generate bits, 0 and 1
type BitGenerator struct {
}

// Get get n bits of 0 or 1
func (bg *BitGenerator) Get(n int) string {
	if n > maxBitNum {
		n = maxBitNum
	}
	return fmt.Sprintf("b'%s'", RandomBits(n))
}

// TimeGenerator generator for time
type TimeGenerator struct {
}

// Get get legal time from 1970 to 2038
func (tg *TimeGenerator) Get() string {
	return fmt.Sprintf("'%s'", RandomTime().Format(TimeFormat))
}

const (
	// special chars
	backslash   = 92
	quote       = 34
	singleQuote = 39
)

// WrapString add backslash to special chars like ", ', \
func WrapString(raw *string) string {
	if raw == nil {
		return ""
	}
	utf8Chars := []rune(*raw)
	result := make([]rune, 0, len(*raw))
	for i := range utf8Chars {
		if utf8Chars[i] == quote || utf8Chars[i] == singleQuote || utf8Chars[i] == backslash {
			result = append(result, backslash, utf8Chars[i])
		} else {
			result = append(result, utf8Chars[i])
		}
	}
	return string(result)
}
