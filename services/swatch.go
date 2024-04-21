package services

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const SECONDS_PER_BEAT = 86.4

type ISwatch interface {
	GetInternetTime() (string, error)
}

type Swatch struct {
	switzerlandTime *time.Time
}

type SwatchService struct {
	Swatch ISwatch
}

func NewSwatchService() SwatchService {
	return SwatchService{
		Swatch: &Swatch{},
	}
}

func (s SwatchService) GetInternetTime() (string, error) {
	internetTime, err := s.Swatch.GetInternetTime()
	internetTime = fmt.Sprintf("It is currently %s", internetTime)
	return internetTime, err
}

// Gets the internet/swatch time
func (s *Swatch) GetInternetTime() (string, error) {

	bmtTime := s.GetBMTTime()

	if s.switzerlandTime != nil {
		bmtTime = *s.switzerlandTime
	}

	bmtMidnightTime := s.GetPreviousMidnight(bmtTime)

	// Seconds since BMT midngiht
	secondsSinceBMTMidnight := bmtTime.Sub(bmtMidnightTime).Seconds()

	// Convert to beats
	numBeats := s.SecondsToBeats(secondsSinceBMTMidnight)

	// Todo: Is there such thing as date format string in go?
	tokens := strings.Split(bmtTime.Format(time.RFC3339), "T")
	date := tokens[0]

	return fmt.Sprintf("%s@%v", date, numBeats), nil
}

func (s *Swatch) GetBMTTime() time.Time {
	loc, _ := time.LoadLocation("Africa/Casablanca")
	return time.Now().In(loc)
}

func (s *Swatch) GetPreviousMidnight(sometime time.Time) time.Time {
	return time.Date(
		sometime.Year(),
		sometime.Month(),
		sometime.Day(),
		0, 0, 0, 0,
		sometime.Location(),
	)
}

func (s *Swatch) SecondsToBeats(seconds float64) float64 {
	beats := seconds / SECONDS_PER_BEAT
	return math.Round(beats*100) / 100
}
