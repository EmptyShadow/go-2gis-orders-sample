package http

import (
	"applicationDesignTest/pkg/orders"
	"applicationDesignTest/utils/httputils"
	"context"
	"errors"
	"net/http"
)

type OrdersEndpoint struct {
	ordersService orders.Usecase
}

func NewOrdersEndpoint(ordersService orders.Usecase) *OrdersEndpoint {
	return &OrdersEndpoint{
		ordersService: ordersService,
	}
}

const (
	CreateOrderPath = "/order"
	ListOrdersPath  = "/orders"
)

func (e *OrdersEndpoint) Registrate(mux *http.ServeMux) {
	mux.Handle(CreateOrderPath, httputils.CheckMethod(http.MethodPost, httputils.Handler(e.CreateOrder)))
	mux.Handle(ListOrdersPath, httputils.CheckMethod(http.MethodGet, httputils.Handler(e.FindOrders)))
}

var (
	ErrInvalidPeriod      = httputils.NewAPIError(http.StatusBadRequest, orders.ErrInvalidPeriod.Error())
	ErrEmailNotFound      = httputils.NewAPIError(http.StatusNotFound, orders.ErrEmailNotFound.Error())
	ErrEmailNotConfirmed  = httputils.NewAPIError(http.StatusBadRequest, orders.ErrEmailNotConfirmed.Error())
	ErrRoomNotFound       = httputils.NewAPIError(http.StatusNotFound, orders.ErrRoomNotFound.Error())
	ErrRoomAlreadyBooked  = httputils.NewAPIError(http.StatusBadRequest, orders.ErrRoomAlreadyBooked.Error())
	ErrOrderAlreadyExists = httputils.NewAPIError(http.StatusConflict, orders.ErrOrderAlreadyExists.Error())
)

func (e *OrdersEndpoint) CreateOrder(ctx context.Context, req httputils.Request) (_ httputils.Response, err error) {
	reqBody, err := CreateOrderRequestFromReader(req.Body)
	if err != nil {
		return
	}

	inDTO, err := reqBody.CreateOrderInDTO()
	if err != nil {
		return
	}

	outDTO, err := e.ordersService.CreateOrder(ctx, inDTO)
	if errors.Is(err, orders.ErrInvalidPeriod) {
		err = ErrInvalidPeriod
	}
	if errors.Is(err, orders.ErrEmailNotFound) {
		err = ErrEmailNotFound
	}
	if errors.Is(err, orders.ErrEmailNotConfirmed) {
		err = ErrEmailNotConfirmed
	}
	if errors.Is(err, orders.ErrRoomNotFound) {
		err = ErrRoomNotFound
	}
	if errors.Is(err, orders.ErrRoomAlreadyBooked) {
		err = ErrRoomAlreadyBooked
	}
	if errors.Is(err, orders.ErrOrderAlreadyExists) {
		err = ErrOrderAlreadyExists
	}
	if err != nil {
		return
	}

	resp := httputils.Response{
		Body: CreateOrderResponseFromOutDTO(outDTO),
	}
	return resp, nil
}

func (e *OrdersEndpoint) FindOrders(ctx context.Context, req httputils.Request) (_ httputils.Response, err error) {
	reqBody, err := FindOrdersRequestFromQuery(req.Query)
	if err != nil {
		return
	}

	inDTO := reqBody.FindOrdersInDTO()

	outDTO, err := e.ordersService.FindOrders(ctx, inDTO)
	if err != nil {
		return
	}

	resp := httputils.Response{
		Body: FindOrdersResponseFromDTO(outDTO),
	}
	return resp, nil
}
