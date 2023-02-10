package emails

type StatusDTO uint8

const (
	InProcessConfirmation StatusDTO = iota
	ConfirmationIsExpired
	ConfirmationIsFailed
	Confirmed
)

type EmailDTO struct {
	Email  string
	Status StatusDTO
}
