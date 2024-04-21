package services

type HomeService struct{}

func (s HomeService) Greet() string {
	return "Hello, world!"
}
