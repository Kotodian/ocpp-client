package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatSN(t *testing.T) {
	sn := "T1641735"
	assert.Equal(t, "T1641735000", formatSN(sn))
	sn = "T164173500000000000"
	assert.Equal(t, "T1641735000", formatSN(sn))
}
