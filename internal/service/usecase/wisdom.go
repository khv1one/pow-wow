package usecase

type WisdomRepository interface {
	ExtractWisdom() (string, error)
}

type WisdomUsecase struct {
	wisdomRepository WisdomRepository
}

func NewWisdomUsecase(wisdomRepository WisdomRepository) *WisdomUsecase {
	return &WisdomUsecase{wisdomRepository: wisdomRepository}
}

func (wuc *WisdomUsecase) ExtractWisdom() string {
	res, _ := wuc.wisdomRepository.ExtractWisdom()

	return res
}
