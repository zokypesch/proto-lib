package core

import (
	"math/rand"
	"strings"
	"time"
)

const letterBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// GeneratorCode struct info for generator
type GeneratorCode struct{}

// GeneratorCodeService for services
type GeneratorCodeService interface {
	Do(n int) string
}

var gen GeneratorCodeService

// NewServiceGenerator for new service generator
func NewServiceGenerator() GeneratorCodeService {
	if gen == nil {
		gen = &GeneratorCode{}
	}
	return gen
}

var randInt = map[int]string{
	0: "XBT-",
	1: "KLB-",
	2: "BCX-",
	3: "KTT-",
	4: "NCD-",
	5: "BMA-",
	6: "FWD-",
	7: "PSD-",
	8: "PDS-",
	9: "BBR-",
}

var src = rand.NewSource(time.Now().UnixNano())

// Do service for do action
func (g *GeneratorCode) Do(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return randInt[rand.Intn(9)] + sb.String()

}
