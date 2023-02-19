package repository

import "hls.streaming.com/src/internal/model"

type RoomRepository struct {
}

func (r RoomRepository) GetRoomPublisherIds(roomId string) ([]string, error) {
	return nil, nil
}

func (r RoomRepository) GetRoomViewerIds(roomId string) ([]string, error) {
	return nil, nil
}

func (r RoomRepository) getConnectionsDetails(ids []string, roomId string) {

}

func (r RoomRepository) QueryParticipant(roomId, first, after string) ([]*model.ParticipantGraph, error) {
	_, err := r.GetRoomPublisherIds(roomId)
	if err != nil {
		return nil, err
	}

	_, err = r.GetRoomViewerIds(roomId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
