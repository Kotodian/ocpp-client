package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatSN(t *testing.T) {
	sn := "T1641735"
	assert.Equal(t, "T1641735000", formatSN(sn))
	sn = "T164173500000000000"
	assert.Equal(t, "T1641735000", formatSN(sn))
}

func TestCreateMultipleChargeStation(t *testing.T) {
	for i := 0; i < 10; i++ {
		sn := "T16417352" + fmt.Sprint(i)
		fmt.Println(formatSN(sn))
	}
}
