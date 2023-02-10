package internal

import (
	"applicationDesignTest/pkg/rooms"
	"context"
)

type RoomsUsecase struct{}

func NewRoomsUsecase() *RoomsUsecase {
	return &RoomsUsecase{}
}

func (s *RoomsUsecase) FindRoom(ctx context.Context, in rooms.FindRoomInDTO) (out rooms.FindRoomOutDTO, err error) {
	if in.FindParams.RoomID == 5 {
		out.Room.ID = 5
		out.Room.Name = "qwe"
		return
	}
	err = rooms.ErrRoomNotExists
	return
}
