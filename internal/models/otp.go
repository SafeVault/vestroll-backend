package models

import (
	"time"
)

type OTPType string

const (
	OTPTypeSMS   OTPType = "sms"
	OTPTypeEmail OTPType = "email"
)

type OTPRequest struct {
	Identifier string  `json:"identifier" binding:"required"`
	Type       OTPType `json:"type" binding:"required,oneof=sms email"`
}

type OTPVerificationRequest struct {
	Identifier string  `json:"identifier" binding:"required"`
	Code       string  `json:"code" binding:"required,len=6"`
	Type       OTPType `json:"type" binding:"required,oneof=sms email"`
}

type OTPData struct {
	Code      string    `json:"code"`
	Type      OTPType   `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Attempts  int       `json:"attempts"`
}

type OTPResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type OTPVerificationResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
