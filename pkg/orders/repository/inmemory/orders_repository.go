package inmemory

import (
	"applicationDesignTest/pkg/orders/repository"
	"context"
	"sync"
)

var _ repository.OrdersRepository = (*OrdersRepository)(nil)

type OrdersRepository struct {
	orders  map[uint64]repository.OrderDTO
	rwMutex *sync.RWMutex
	lastID  uint64
}

func NewOrdersRepository() *OrdersRepository {
	return &OrdersRepository{
		orders:  make(map[uint64]repository.OrderDTO),
		rwMutex: &sync.RWMutex{},
	}
}

func (s *OrdersRepository) SaveOrder(ctx context.Context, in repository.SaveOrderInDTO) error {
	s.rwMutex.Lock()
	defer s.rwMutex.Unlock()

	for _, order := range s.orders {
		if in.Order.EqualForSave(order) {
			return repository.ErrOrderAlreadyExists
		}
	}

	s.lastID++
	in.Order.ID = s.lastID
	s.orders[in.Order.ID] = in.Order.Copy()

	return nil
}

func (s *OrdersRepository) FindOrder(ctx context.Context, in repository.FindOrderInDTO) (out repository.FindOrderOutDTO, err error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	filter := ordersFiltersFromFindOrdersParams(in.FindParams)

	for _, order := range s.orders {
		if filter(order) {
			out.Order = order.Copy()
			return
		}
	}

	err = repository.ErrOrderNotExists
	return
}

func (s *OrdersRepository) FindOrders(ctx context.Context, in repository.FindOrdersInDTO) (out repository.FindOrdersOutDTO, err error) {
	s.rwMutex.RLock()
	defer s.rwMutex.RUnlock()

	filter := ordersFiltersFromFindOrdersParams(in.FindParams)
	limit := in.PaginateParams.Limit
	offet := in.PaginateParams.Offset
	total := uint64(0)
	orders := repository.ListOrdersDTO{}

	for _, order := range s.orders {
		if !filter(order) {
			continue
		}
		total++
		if offet > 0 {
			offet--
			continue
		}
		if limit > 0 {
			limit--
			orders = append(orders, order.Copy())
		}
	}

	out.Orders = orders
	out.Total = total
	return
}

type OrderFilters []OrderFilter

func (fs OrderFilters) Check(order repository.OrderDTO) bool {
	for _, filter := range fs {
		if !filter(order) {
			return false
		}
	}
	return true
}

type OrderFilter func(order repository.OrderDTO) bool

func ordersFiltersFromFindOrdersParams(params repository.FindOrdersParamsInDTO) OrderFilter {
	var filters OrderFilters
	if params.Email != "" {
		filters = append(filters, func(order repository.OrderDTO) bool {
			return order.Email == params.Email
		})
	}
	if len(params.OneOfStatuses) > 0 {
		filters = append(filters, func(order repository.OrderDTO) bool {
			for _, status := range params.OneOfStatuses {
				if order.Status == status {
					return true
				}
			}
			return false
		})
	}
	if !params.From.IsZero() {
		filters = append(filters, func(order repository.OrderDTO) bool {
			return order.From.Equal(params.From)
		})
	}
	if !params.To.IsZero() {
		filters = append(filters, func(order repository.OrderDTO) bool {
			return order.To.Equal(params.To)
		})
	}
	return filters.Check
}
