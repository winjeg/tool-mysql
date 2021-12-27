package generator

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// using a rand.Seed(n) function slows down all random functions
// so in every function,  we don't use a seed

const (
	TimeFormat    = "2006-01-02 15:04:05"
	timeStartDate = "1970-01-01 00:00:01"
	timeEndDate   = "2038-01-19 03:14:07"

	// cjk section of unicode charset
	cjkStart = 0x4E00
	cjkStop  = 0x9FBF
)

var (
	// time related vars, max, min time and duration
	timeStart, _ = time.Parse(TimeFormat, timeStartDate)
	timeEnd, _   = time.Parse(TimeFormat, timeEndDate)
	maxDuration  = uint64(timeEnd.Sub(timeStart).Nanoseconds())
)

// RandomInt random int, returns possibly negative value or positive value
// including the limit, and negative values, the arg limit should be positive
func RandomInt(limit int64) int64 {
	result := rd.Int63n(limit) - 2
	if RandomBool() {
		return -result
	}
	return result
}

// RandomUInt random unsigned int, positive value only
func RandomUInt(limit uint64) uint64 {
	r := rd.Uint64()
	if r > limit {
		return r % limit
	}
	return r
}

// RandomBool random bool
func RandomBool() bool {
	result := rd.Intn(2) == 1
	return result
}

// RandomFloat return a random double, limit is the max value
// part return 0 for illegal values
func RandomFloat(n, f int) string {
	if n <= 0 {
		n = 0
	}
	if f <= 0 {
		f = 0
	}
	partN := make([]byte, n)
	partF := make([]byte, f)
	for i := 0; i < n; i++ {
		partN[i] = byte(rd.Intn(10) + 48)
	}
	for i := 0; i < f; i++ {
		partF[i] = byte(rd.Intn(10) + 48)
	}
	if n == 0 && f == 0 {
		return "0.0"
	}
	if n <= 0 {
		return "0." + string(partF)
	}
	if f <= 0 {
		return string(partN) + ".0"
	}
	return string(partN) + "." + string(partF)
}

// RandomTime return a random legal time both for mysql type datetime and timestamp
func RandomTime() time.Time {
	return timeStart.Add(time.Duration(RandomUInt(maxDuration)))
}

// RandomBits random bits,  n should at most be 64
func RandomBits(n int) string {
	num := RandomUInt(math.MaxUint64) + math.MaxUint64>>2
	bits := fmt.Sprintf("%b", num)
	if len(bits) > n {
		return bits[:n]
	}
	return bits[:rd.Intn(len(bits))]
}

// RandomCJK random CJK in UTF8
// for utf-8 is a mutable length charset
// so the result length may not actually equals to the size specified
func RandomCJK(size int) string {
	result := make([]rune, size)
	for i := range result {
		result[i] = rune(RandIntSection(cjkStart, cjkStop))
	}
	return string(result)
}

// RandomASCII random ascii
func RandomASCII(size int) string {
	result := make([]byte, size)
	for i := range result {
		result[i] = byte(RandIntSection(0, 128))
	}
	return string(result)
}

// RandomReadable random readable string
func RandomReadable(size int) string {
	return string(RandStr(size, KindAllWithSpecial))
}

// RandIntSection random int section
func RandIntSection(min, max int64) int64 {
	return min + rd.Int63n(max-min)
}

func RandStr(size int, kind int) []byte {
	kinds := [][]int{{10, 48}, {26, 97}, {26, 65}}
	specialChars := []byte{95, 45, 46, 35, 36, 37, 38}
	specialCharLen := len(specialChars)
	iKind, result := kind, make([]byte, size)
	isAll := kind == 3
	for i := 0; i < size; i++ {

		// random iKind
		if isAll {
			iKind = rand.Intn(3)
		}
		if kind == KindAllWithSpecial {
			iKind = rand.Intn(4)
		}
		if iKind == 3 {
			result[i] = specialChars[rand.Intn(specialCharLen)]
		} else {
			scope, base := kinds[iKind][0], kinds[iKind][1]
			result[i] = uint8(base + rand.Intn(scope))
		}
	}
	return result
}

const (
	KindAllWithSpecial = 4
)

// the rand instance, to improve efficiency
var rd = NewRand()

func NewRand() *rand.Rand {
	return rand.New(&LockedSource{src: rand.NewSource(time.Now().UnixNano()).(Source64)})
}
