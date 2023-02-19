package model

import "time"

type ParticipantGraph struct {
	UserId    string    `json:"userId"`
	SessionId string    `json:"sessionId"`
	Role      string    `json:"role"`
	Video     string    `json:"video"`
	Audio     string    `json:"audio"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Status    string    `json:"status"`
}
