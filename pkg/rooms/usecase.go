package rooms

import (
	"context"
	"errors"
)

var ErrRoomNotExists = errors.New("room not exists")

type Usecase interface {
	FindRoom(ctx context.Context, in FindRoomInDTO) (out FindRoomOutDTO, err error)
}

type FindRoomInDTO struct {
	FindParams FindParamsInDTO
}

type FindParamsInDTO struct {
	RoomID uint64
}

type FindRoomOutDTO struct {
	Room RoomDTO
}
