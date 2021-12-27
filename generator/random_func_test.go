package generator

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"unicode/utf8"
)

const (
	charNum = 255
)

func BenchmarkRandomBool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomBool()
	}
	b.ReportAllocs()
}

func BenchmarkRandomInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomInt(math.MaxInt32)
	}
	b.ReportAllocs()
}

func BenchmarkRandomUInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomUInt(math.MaxInt32)
	}
	b.ReportAllocs()
}

func BenchmarkRandomTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomTime()
	}
	b.ReportAllocs()
}

func BenchmarkRandomDouble(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomFloat(10, 10)
	}
	b.ReportAllocs()
}

func BenchmarkRandomASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomASCII(charNum)
	}
	b.ReportAllocs()
}

func BenchmarkRandomReadable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomReadable(charNum)
	}
	b.ReportAllocs()
}
func BenchmarkRandomCJK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomCJK(charNum)
	}
	b.ReportAllocs()
}

func BenchmarkRandomBits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomBits(10)
	}
	b.ReportAllocs()
}

func TestRandom(t *testing.T) {
	RandomBool()
	assert.True(t, RandomInt(math.MaxInt64) < math.MaxInt64)
	assert.True(t, RandomUInt(math.MaxInt64) < math.MaxInt64)
	assert.True(t, timeEnd.Sub(RandomTime()).Nanoseconds() > 0)
	assert.True(t, RandomTime().Sub(timeStart).Nanoseconds() > 0)
	assert.True(t, len(RandomFloat(1, 1)) == 3)
	assert.True(t, len(RandomFloat(0, 1)) == 3)
	assert.True(t, len(RandomFloat(1, -1)) == 3)
	assert.True(t, len(RandomFloat(-1, -1)) == 3)
	assert.True(t, len(RandomReadable(charNum)) == charNum)
	assert.True(t, len(RandomASCII(charNum)) == charNum)
	assert.True(t, utf8.RuneCountInString(RandomCJK(charNum)) == charNum)
	assert.True(t, len(RandomBits(10)) == 10)
}
