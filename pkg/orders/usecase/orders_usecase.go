package usecase

import (
	"applicationDesignTest/pkg/emails"
	"applicationDesignTest/pkg/orders"
	"applicationDesignTest/pkg/orders/repository"
	"applicationDesignTest/pkg/rooms"
	"context"
	"errors"
	"fmt"
	"time"
)

var _ orders.Usecase = (*OrdersUsecase)(nil)

type OrdersUsecase struct {
	emailsUsecase    emails.Usecase
	roomsUsecase     rooms.Usecase
	ordersRepository repository.OrdersRepository
}

func NewOrdersUsecase(
	emailsUsecase emails.Usecase,
	roomsUsecase rooms.Usecase,
	ordersRepository repository.OrdersRepository,
) *OrdersUsecase {
	return &OrdersUsecase{
		emailsUsecase:    emailsUsecase,
		roomsUsecase:     roomsUsecase,
		ordersRepository: ordersRepository,
	}
}

func (s *OrdersUsecase) CreateOrder(ctx context.Context, in orders.CreateOrderInDTO) (_ orders.CreateOrderOutDTO, err error) {
	// валидация периода бронирования.
	now := time.Now().UTC()
	if now.After(in.From) || now.After(in.To) || in.From.After(in.To) {
		err = orders.ErrInvalidPeriod
		return
	}

	// проверка что email валиден и подтвержден.
	findEmailOutDTO, err := s.emailsUsecase.FindEmail(ctx, emails.FindEmailInDTO{
		FindParams: emails.FindParamsInDTO{
			Email: in.Email,
		},
	})
	if errors.Is(err, emails.ErrEmailNotExists) {
		err = orders.ErrEmailNotFound
		return
	}
	if findEmailOutDTO.Email.Status != emails.Confirmed {
		err = orders.ErrEmailNotConfirmed
		return
	}
	if err != nil {
		err = fmt.Errorf("find email in emails service: %w", err)
		return
	}

	// проверка что room существует.
	_, err = s.roomsUsecase.FindRoom(ctx, rooms.FindRoomInDTO{
		FindParams: rooms.FindParamsInDTO{
			RoomID: in.RoomID,
		},
	})
	if errors.Is(err, rooms.ErrRoomNotExists) {
		err = orders.ErrRoomNotFound
		return
	}
	if err != nil {
		err = fmt.Errorf("find room in rooms service: %w", err)
		return
	}

	// надо найти order, который является блокирующим для создания нашего order.
	// кто то нас опередил и раньше смог сделать заказ и перейти на статусы из выборки.
	_, err = s.ordersRepository.FindOrder(ctx, repository.FindOrderInDTO{
		FindParams: repository.FindOrdersParamsInDTO{
			OneOfStatuses: []repository.OrderStatusDTO{
				repository.OrderConfirmedByOwner,
				repository.OrderWaitingPayment,
				repository.OrderInProcessPayment,
				repository.OrderPiad,
			},
			From: in.From,
			To:   in.To,
		},
	})
	if err != nil && !errors.Is(err, repository.ErrOrderNotExists) {
		err = fmt.Errorf("find order records from repository: %w", err)
		return
	}
	if !errors.Is(err, repository.ErrOrderNotExists) {
		err = orders.ErrRoomAlreadyBooked
		return
	}

	// формируем запись в хранилище об order.
	order := repository.OrderDTO{
		RoomID:    in.RoomID,
		Status:    repository.OrderInProcessConfirmationByOwner,
		Email:     in.Email,
		From:      in.From,
		To:        in.To,
		CreatedAt: now,
		UpdatedAt: now,
	}
	err = s.ordersRepository.SaveOrder(ctx, repository.SaveOrderInDTO{
		Order: &order,
	})
	if errors.Is(err, repository.ErrOrderAlreadyExists) {
		err = orders.ErrOrderAlreadyExists
		return
	}
	if err != nil {
		err = fmt.Errorf("save order record to repository: %w", err)
		return
	}

	// маппим на выход
	outDTO := orders.CreateOrderOutDTO{
		Order: MakeOrderDTO(order),
	}
	return outDTO, nil
}

func (s *OrdersUsecase) FindOrders(ctx context.Context, in orders.FindOrdersInDTO) (_ orders.FindOrdersOutDTO, err error) {
	findOrdersOutDTO, err := s.ordersRepository.FindOrders(ctx, repository.FindOrdersInDTO{
		FindParams: repository.FindOrdersParamsInDTO{
			Email: in.FindParams.Email,
		},
		PaginateParams: repository.PaginateOrdersParamsInDTO{
			Limit:  in.PaginateParams.Limit,
			Offset: in.PaginateParams.Offset,
		},
	})
	if err != nil {
		err = fmt.Errorf("find order records in repository: %w", err)
		return
	}

	outDTO := orders.FindOrdersOutDTO{
		Orders: MakeListOrdersDTO(findOrdersOutDTO.Orders),
		Total:  findOrdersOutDTO.Total,
	}
	return outDTO, nil
}
