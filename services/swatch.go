package services

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const secondsPerBeat = 86.4

type iSwatch interface {
	getInternetTime() (string, error)
}

type Swatch struct {
	switzerlandTime *time.Time
}

type SwatchService struct {
	Swatch iSwatch
}

func NewSwatchService() SwatchService {
	return SwatchService{
		Swatch: &Swatch{},
	}
}

func (s SwatchService) GetInternetTime() (string, error) {
	internetTime, err := s.Swatch.getInternetTime()
	return internetTime, err
}

func (s *Swatch) getInternetTime() (string, error) {

	bmtTime := s.getBMTTime()

	if s.switzerlandTime != nil {
		bmtTime = *s.switzerlandTime
	}

	bmtMidnightTime := s.getPreviousMidnight(bmtTime)

	// Seconds since BMT midngiht
	secondsSinceBMTMidnight := bmtTime.Sub(bmtMidnightTime).Seconds()

	// Convert to beats
	numBeats := s.secondsToBeats(secondsSinceBMTMidnight)

	// Todo: Is there such thing as date format string in go?
	tokens := strings.Split(bmtTime.Format(time.RFC3339), "T")
	date := tokens[0]

	return fmt.Sprintf("%s@%v", date, numBeats), nil
}

func (s *Swatch) getBMTTime() time.Time {
	loc, _ := time.LoadLocation("Africa/Casablanca")
	return time.Now().In(loc)
}

func (s *Swatch) getPreviousMidnight(sometime time.Time) time.Time {
	return time.Date(
		sometime.Year(),
		sometime.Month(),
		sometime.Day(),
		0, 0, 0, 0,
		sometime.Location(),
	)
}

func (s *Swatch) secondsToBeats(seconds float64) float64 {
	beats := seconds / secondsPerBeat
	return math.Round(beats*100) / 100
}
