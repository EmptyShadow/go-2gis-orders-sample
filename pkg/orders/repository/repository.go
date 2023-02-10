package repository

import (
	"context"
	"errors"
	"time"
)

var (
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderNotExists     = errors.New("order not exists")
)

// OrdersRepository хранилище Order.
type OrdersRepository interface {
	// SaveOrder метод должен сохранять order и модифицировать в нем ID.
	//
	// Если order уже существует нужно вернуть ErrOrderAlreadyExists.
	SaveOrder(ctx context.Context, in SaveOrderInDTO) error

	// FindOrder метод должен произвести выборку по FindOrdersParamsInDTO и вернуть OrderDTO.
	//
	// Если order  нет вернуть ErrOrderNotExists.
	FindOrder(ctx context.Context, in FindOrderInDTO) (out FindOrderOutDTO, err error)

	// FindOrders метод должен произвести выборку по FindOrdersParamsInDTO и вернуть ListOrdersDTO.
	FindOrders(ctx context.Context, in FindOrdersInDTO) (out FindOrdersOutDTO, err error)
}

type SaveOrderInDTO struct {
	Order *OrderDTO
}

type FindOrderInDTO struct {
	FindParams FindOrdersParamsInDTO
}

type FindOrderOutDTO struct {
	Order OrderDTO
}

type FindOrdersInDTO struct {
	FindParams     FindOrdersParamsInDTO
	PaginateParams PaginateOrdersParamsInDTO
}

type FindOrdersParamsInDTO struct {
	Email         string
	OneOfStatuses []OrderStatusDTO
	From          time.Time
	To            time.Time
}

type PaginateOrdersParamsInDTO struct {
	Limit  uint64
	Offset uint64
}

type FindOrdersOutDTO struct {
	Orders ListOrdersDTO
	Total  uint64
}
