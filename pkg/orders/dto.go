package orders

import "time"

type ListOrdersDTO []OrderDTO

type StatusDTO uint8

const (
	// В процессе подтверждения владельцем room.
	InProcessConfirmationByOwner StatusDTO = iota

	// Подтвержден владельцем.
	ConfirmedByOwner

	// В ожидании оплаты заказчиком.
	WaitingPayment

	// В процессе оплаты заказчиком.
	InProcessPayment

	// Отменена владельцем room до оплаты order.
	CanceledByOwnerBeforePayment

	// Отменена заказчиком до оплаты order.
	CanceledByCustomerBeforePayment

	// Заказ оплачен.
	Piad

	// Возврат оплаты.
	Refund
)

type OrderDTO struct {
	ID        uint64    // идентификатор бронирования.
	RoomID    uint64    // комната которую хотят забронировать.
	Status    StatusDTO // статус order.
	Email     string    // email заказчика.
	From      time.Time // начало периода бронирования.
	To        time.Time // конец периода бронирования.
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
