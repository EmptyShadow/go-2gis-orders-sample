package orders

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidPeriod      = errors.New("invalid period")
	ErrEmailNotFound      = errors.New("email not found")
	ErrEmailNotConfirmed  = errors.New("email not confirmed")
	ErrRoomNotFound       = errors.New("room not found")
	ErrRoomAlreadyBooked  = errors.New("room already booked")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

type Usecase interface {
	// CreateOrder метод должен создавать order на бронирование room.
	//
	// Если email не найден вернуть ErrEmailNotFound.
	// Если email не подтвержден вернуть ErrEmailNotConfirmed.
	// Есть room не найден вернуть ErrRoomNotFound.
	// Если room уже забронирована вернуть ErrRoomAlreadyBooked.
	// Если room уже кем то бронируется вернуть ErrRoomLockedForBooking.
	// Если такой order уже есть вернуть ErrOrderAlreadyExists.
	CreateOrder(ctx context.Context, in CreateOrderInDTO) (CreateOrderOutDTO, error)

	// FindOrders метод должен искать orders по переданным FindParamsDTO и
	// ограничевать выборку по PaginateParams.
	FindOrders(ctx context.Context, in FindOrdersInDTO) (FindOrdersOutDTO, error)
}

type CreateOrderInDTO struct {
	RoomID uint64
	Email  string
	From   time.Time
	To     time.Time
}

type CreateOrderOutDTO struct {
	Order OrderDTO
}

type FindOrdersInDTO struct {
	FindParams     FindParamsInDTO
	PaginateParams PaginateParamsInDTO
}

type FindParamsInDTO struct {
	Email string
}

type PaginateParamsInDTO struct {
	Limit  uint64
	Offset uint64
}

type FindOrdersOutDTO struct {
	Orders ListOrdersDTO
	Total  uint64
}
