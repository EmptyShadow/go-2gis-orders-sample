package usecase

import (
	"applicationDesignTest/pkg/orders"
	"applicationDesignTest/pkg/orders/repository"
)

func MakeListOrdersDTO(list repository.ListOrdersDTO) orders.ListOrdersDTO {
	out := make(orders.ListOrdersDTO, len(list))
	for i := range list {
		out[i] = MakeOrderDTO(list[i])
	}
	return out
}

func MakeOrderDTO(order repository.OrderDTO) orders.OrderDTO {
	return orders.OrderDTO{
		ID:        order.ID,
		RoomID:    order.RoomID,
		Status:    MakeOrderStatus(order.Status),
		Email:     order.Email,
		From:      order.From,
		To:        order.To,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
		DeletedAt: order.DeletedAt,
	}
}

func MakeOrderStatus(status repository.OrderStatusDTO) orders.StatusDTO {
	switch status {
	case repository.OrderInProcessConfirmationByOwner:
		return orders.InProcessConfirmationByOwner
	case repository.OrderConfirmedByOwner:
		return orders.ConfirmedByOwner
	case repository.OrderWaitingPayment:
		return orders.WaitingPayment
	case repository.OrderInProcessPayment:
		return orders.InProcessPayment
	case repository.OrderCanceledByOwnerBeforePayment:
		return orders.CanceledByOwnerBeforePayment
	case repository.OrderCanceledByCustomerBeforePayment:
		return orders.CanceledByCustomerBeforePayment
	case repository.OrderPiad:
		return orders.Piad
	case repository.OrderRefund:
		return orders.Refund
	}
	return 0
}
