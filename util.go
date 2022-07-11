package api

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"log"
	"math/big"
	"net/http"
	"time"
)

const (
	UpperEnglishLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var headerIPAddress = []string{"X-Forwarded-For", "X-Real-Ip", "HTTP_X_FORWARDED_FOR", "fwd"}

// GetIPAddresses returns comma separated list of remote IP addresses
func GetIPAddresses(r *http.Request) string {
	var ips []string

	for _, h := range headerIPAddress {
		addresses := r.Header.Get(h)
		log.Printf(" %s :%v", h, addresses)

		if len(addresses) == 0 {
			continue
		}
		log.Printf(" %v visit site", addresses)

		ips = append(ips, addresses)
	}
	if len(ips) > 0 {
		return ips[0]
	}
	return ""
}

// RandomString generates a random string consisting of characters in the provided alphabet
func RandomString(n int) string {
	if n <= 0 {
		return ""
	}

	r := []rune(UpperEnglishLetters)
	k := byte(len(r))

	s := make([]rune, n)
	b := make([]byte, n)

	_, _ = rand.Read(b)
	for i := range b {
		s[i] = r[b[i]%k]
	}

	return string(s)
}

// RandomNumberFromRange generates a random number
func RandomNumberFromRange(min int, max int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		panic(err)
	}
	return min + int(nBig.Int64())
}

func GetCurrentTimeStamp() time.Time {
	return time.Now().UTC()
}

func GetCurrentISOTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func ParseToISOTime(dt time.Time) string {
	return dt.UTC().Format(time.RFC3339)
}

func HashString(s string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%v", h.Sum32())
}
