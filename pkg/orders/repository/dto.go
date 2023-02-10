package repository

import "time"

type ListOrdersDTO []OrderDTO

type OrderStatusDTO uint8

const (
	OrderInProcessConfirmationByOwner OrderStatusDTO = iota
	OrderConfirmedByOwner
	OrderWaitingPayment
	OrderInProcessPayment
	OrderCanceledByOwnerBeforePayment
	OrderCanceledByCustomerBeforePayment
	OrderPiad
	OrderRefund
)

type OrderDTO struct {
	ID        uint64
	RoomID    uint64
	Status    OrderStatusDTO
	Email     string
	From      time.Time
	To        time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (o1 OrderDTO) EqualForSave(o2 OrderDTO) bool {
	return o1.Email == o2.Email && o1.From.Equal(o2.From) && o1.To.Equal(o2.To)
}

func (o OrderDTO) Copy() OrderDTO {
	var deletedAt *time.Time
	if o.DeletedAt != nil {
		_deletedAt := *o.DeletedAt
		deletedAt = &_deletedAt
	}
	return OrderDTO{
		ID:        o.ID,
		RoomID:    o.RoomID,
		Status:    o.Status,
		Email:     o.Email,
		From:      o.From,
		To:        o.To,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		DeletedAt: deletedAt,
	}
}
