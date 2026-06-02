package part

type PartService struct {
	repository PartRepository
}

func NewPartService(repository PartRepository) *PartService {
	return &PartService{
		repository: repository,
	}
}
