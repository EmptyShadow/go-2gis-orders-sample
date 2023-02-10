package http

import (
	"applicationDesignTest/pkg/orders"
	"applicationDesignTest/utils/date"
	"applicationDesignTest/utils/httputils"
	"encoding/json"
	"io"
	"net/url"
	"strconv"
)

const (
	RoomIDField = "room_id"
	EmailField  = "email"
	FromField   = "from"
	ToField     = "to"
	LimitField  = "limit"
	OffsetField = "offset"
)

type CreateOrderRequest struct {
	RoomID uint64 `json:"room_id"`
	Email  string `json:"email"`
	From   string `json:"from"`
	To     string `json:"to"`
}

func CreateOrderRequestFromReader(r io.Reader) (reqBody CreateOrderRequest, err error) {
	if err = json.NewDecoder(r).Decode(&reqBody); err != nil {
		err = httputils.NewInvalidBodyFormatAPIError()
	}
	return
}

func (r CreateOrderRequest) CreateOrderInDTO() (_ orders.CreateOrderInDTO, err error) {
	var badArgs []httputils.BadArgument

	from, err := date.FromString(r.From)
	if err != nil {
		badArgs = append(badArgs, httputils.BadArgument{
			Field:  FromField,
			Reason: err.Error(),
		})
	}

	to, err := date.FromString(r.To)
	if err != nil {
		badArgs = append(badArgs, httputils.BadArgument{
			Field:  ToField,
			Reason: err.Error(),
		})
	}

	if len(badArgs) > 0 {
		err = httputils.NewBadArgumentsAPIError(badArgs...)
		return
	}

	inDTO := orders.CreateOrderInDTO{
		RoomID: r.RoomID,
		Email:  r.Email,
		From:   from,
		To:     to,
	}

	return inDTO, nil
}

type CreateOrderResponse struct {
	Order Order `json:"order"`
}

func CreateOrderResponseFromOutDTO(outDTO orders.CreateOrderOutDTO) CreateOrderResponse {
	return CreateOrderResponse{
		Order: OrderFromDTO(outDTO.Order),
	}
}

type FindOrdersRequest struct {
	Email  string `json:"-"`
	Limit  uint64 `json:"-"`
	Offset uint64 `json:"-"`
}

func FindOrdersRequestFromQuery(query url.Values) (_ FindOrdersRequest, err error) {
	var badArgs []httputils.BadArgument

	email := query.Get(EmailField)
	if email == "" {
		badArgs = append(badArgs, httputils.BadArgument{
			Field:  EmailField,
			Reason: "requared",
		})
	}

	var limit uint64
	limitStr := query.Get(LimitField)
	if limitStr != "" {
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			badArgs = append(badArgs, httputils.BadArgument{
				Field:  LimitField,
				Reason: err.Error(),
			})
		}
	}

	var offset uint64
	offsetStr := query.Get(OffsetField)
	if offsetStr != "" {
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			badArgs = append(badArgs, httputils.BadArgument{
				Field:  OffsetField,
				Reason: err.Error(),
			})
		}
	}

	if len(badArgs) > 0 {
		err = httputils.NewBadArgumentsAPIError(badArgs...)
		return
	}

	reqBody := FindOrdersRequest{
		Email:  email,
		Limit:  limit,
		Offset: offset,
	}

	return reqBody, nil
}

func (r FindOrdersRequest) FindOrdersInDTO() orders.FindOrdersInDTO {
	return orders.FindOrdersInDTO{
		FindParams: orders.FindParamsInDTO{
			Email: r.Email,
		},
		PaginateParams: orders.PaginateParamsInDTO{
			Limit:  r.Limit,
			Offset: r.Offset,
		},
	}
}

type FindOrdersResponse struct {
	Orders ListOrders `json:"orders"`
}

func FindOrdersResponseFromDTO(dto orders.FindOrdersOutDTO) FindOrdersResponse {
	return FindOrdersResponse{
		Orders: ListOrderFromDTO(dto.Orders),
	}
}
