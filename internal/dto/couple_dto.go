package dto

import "time"

type CreateInviteRequest struct {
	AnniversaryDate *string `json:"anniversary_date" binding:"omitempty,datetime=2006-01-02"`
}

type JoinCoupleRequest struct {
	InviteCode string `json:"invite_code" binding:"required"`
}

type CoupleResponse struct {
	CoupleID        string        `json:"couple_id"`
	InviteCode      string        `json:"invite_code"`
	AnniversaryDate *time.Time    `json:"anniversary_date"`
	Status          string        `json:"status"`
	Partner         *UserResponse `json:"partner,omitempty"`
}