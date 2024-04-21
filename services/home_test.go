package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomeService(t *testing.T) {
	svc := HomeService{}
	assert.Equal(t, "Hello, world!", svc.Greet())
}
