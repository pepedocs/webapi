package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type swatchTimeTest struct {
	switzerlandTime time.Time
	swatchTime      string
}

func TestGetInternetTime(t *testing.T) {
	// To test Swatch.GetInternetTime we fix the biel time otherwise we
	// would get different results when we call getBMTTime. We then test
	// the return values to pre-calculated internet time.

	loc, _ := time.LoadLocation("Africa/Casablanca")
	tests := []swatchTimeTest{
		{
			switzerlandTime: time.Date(2000, time.January, 1, 8, 0, 0, 0, loc),
			swatchTime:      "2000-01-01@333.33",
		},
		{
			switzerlandTime: time.Date(2024, time.June, 5, 17, 1, 1, 0, loc),
			swatchTime:      "2024-06-05@709.04",
		},
	}
	for _, test := range tests {
		sw := Swatch{&test.switzerlandTime}
		swatchTime, err := sw.GetInternetTime()
		assert.NoError(t, err)
		assert.Equal(t, test.swatchTime, swatchTime)
	}
}
