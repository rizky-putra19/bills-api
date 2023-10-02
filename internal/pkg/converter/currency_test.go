package converter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToIDRNumber(t *testing.T) {
	assert.Equal(t, "1.000", ToIDRNumber(1000))
	assert.Equal(t, "-10.000.000", ToIDRNumber(-10000000))
}

func TestToIDR(t *testing.T) {
	assert.Equal(t, "Rp 1.000", ToIDR(1000))
	assert.Equal(t, "Rp 10.010", ToIDR(10010))
	assert.Equal(t, "Rp 1.000", ToIDR(1000))
	assert.Equal(t, "Rp 1.000", ToIDR(1000.00))
	assert.Equal(t, "", ToIDR(""))
}
