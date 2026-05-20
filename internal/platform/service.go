package platform

type Service interface {
	Start() error
	Stop() error
	Restart() error
	Status() error
}

func NewService() (Service, error) {
	return newPlatformService()
}
