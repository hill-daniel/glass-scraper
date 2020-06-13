package colly

import (
	"strconv"
	"strings"
)

// Purge removes leading and trailing white space.
func Purge(s string) string {
	return strings.TrimSpace(s)
}

// ParseFloat parses glassdoor specific number format. E.g. Numbers > 999 are abreaviated with number/1000 + k; e.g 7500 -> 7.5k.
func ParseFloat(s string) float32 {
	multiplier := 1.0
	f := strings.Replace(s, ",", ".", -1)
	if strings.HasSuffix(f, "Tsd") || strings.HasSuffix(f, "k") {
		f = strings.TrimSuffix(f, "Tsd")
		multiplier = 1000.0
	}
	f = Purge(f)
	v, err := strconv.ParseFloat(f, 2)
	if err != nil {
		return 0.0
	}
	return float32(v * multiplier)
}
