package converter

import (
	"strconv"
	"strings"
)

// ToIDRNumber converts number to indonesian format X.XXX
func ToIDRNumber(x interface{}) string {
	value := ToInt64(x)

	sign := ""
	if value < 0 {
		sign = "-"
		value = 0 - value
	}

	parts := []string{"", "", "", "", "", "", ""}
	j := len(parts) - 1

	for value > 999 {
		parts[j] = strconv.FormatInt(value%1000, 10)
		switch len(parts[j]) {
		case 2:
			parts[j] = "0" + parts[j]
		case 1:
			parts[j] = "00" + parts[j]
		}
		value = value / 1000
		j--
	}
	parts[j] = strconv.Itoa(int(value))
	return sign + strings.Join(parts[j:], ".")
}

// ToIDR converts number to Rp X.XXX style
func ToIDR(x interface{}) string {
	switch v := x.(type) {
	case int, int64, float64:
		return "Rp " + ToIDRNumber(v)
	default:
		return ""
	}
}
