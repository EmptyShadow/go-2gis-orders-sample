package http

import (
	"applicationDesignTest/pkg/orders"
	"applicationDesignTest/utils/date"
	"time"
)

type OrderStatus string

const (
	InProcessConfirmationByOwner    OrderStatus = "in_process_confirmation_by_owner"
	ConfirmedByOwner                OrderStatus = "confirmed_by_owner"
	WaitingPayment                  OrderStatus = "waiting_payment"
	InProcessPayment                OrderStatus = "in_process_payment"
	CanceledByOwnerBeforePayment    OrderStatus = "canceled_by_owner_before_payment"
	CanceledByCustomerBeforePayment OrderStatus = "canceled_by_customer_before_payment"
	Piad                            OrderStatus = "piad"
	Refund                          OrderStatus = "refund"
)

func OrderStatusFromDTO(dto orders.StatusDTO) OrderStatus {
	switch dto {
	case orders.InProcessConfirmationByOwner:
		return InProcessConfirmationByOwner
	case orders.ConfirmedByOwner:
		return ConfirmedByOwner
	case orders.WaitingPayment:
		return WaitingPayment
	case orders.InProcessPayment:
		return InProcessPayment
	case orders.CanceledByOwnerBeforePayment:
		return CanceledByOwnerBeforePayment
	case orders.CanceledByCustomerBeforePayment:
		return CanceledByCustomerBeforePayment
	case orders.Piad:
		return Piad
	case orders.Refund:
		return Refund
	}
	return ""
}

type Order struct {
	ID        uint64
	RoomID    uint64
	Status    OrderStatus
	Email     string
	From      string
	To        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func OrderFromDTO(dto orders.OrderDTO) Order {
	return Order{
		ID:        dto.ID,
		RoomID:    dto.RoomID,
		Status:    OrderStatusFromDTO(dto.Status),
		Email:     dto.Email,
		From:      date.ToString(dto.From),
		To:        date.ToString(dto.To),
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		DeletedAt: dto.DeletedAt,
	}
}

type ListOrders []Order

func ListOrderFromDTO(dto orders.ListOrdersDTO) ListOrders {
	list := make(ListOrders, len(dto))
	for i := range dto {
		list[i] = OrderFromDTO(dto[i])
	}
	return list
}
