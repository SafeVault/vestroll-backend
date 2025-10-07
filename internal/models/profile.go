package models

import "time"

// Account types supported for the onboarding profile
const (
    AccountTypeFreelancer string = "freelancer"
    AccountTypeContractorProfile string = "contractor"
)

// AccountTypeRequest sets the user's selected account type
// Note: user_id is required since we are not using auth yet
// binding: oneof constraint ensures value is freelancer or contractor
// We keep a separate constant name AccountTypeContractorProfile to avoid
// conflicting with the BusinessProfile model's constant.
type AccountTypeRequest struct {
    UserID      string `json:"user_id" binding:"required"`
    AccountType string `json:"account_type" binding:"required,oneof=freelancer contractor"`
}

// PersonalDetails holds the basic personal info captured during onboarding
// DateOfBirth uses the canonical format YYYY-MM-DD
// Phone should be digits without formatting; dial_code includes the country code prefix like +234
// Gender is optional but if present must be one of male,female,other
// We keep validation in the handler/service as well for stricter rules.
type PersonalDetails struct {
    FirstName   string `json:"first_name" binding:"required"`
    LastName    string `json:"last_name" binding:"required"`
    Gender      string `json:"gender" binding:"omitempty,oneof=male female other"`
    DateOfBirth string `json:"date_of_birth" binding:"required"`
    DialCode    string `json:"dial_code" binding:"required"`
    Phone       string `json:"phone" binding:"required,min=7,max=20"`
}

type PersonalDetailsRequest struct {
    UserID string          `json:"user_id" binding:"required"`
    Data   PersonalDetails `json:"data" binding:"required"`
}

// Address details for the user's profile
// Country can be ISO alpha-2 code or country name; kept as string with non-empty validation
// Postal code kept optional but validated for length/format in service.
type Address struct {
    Country    string `json:"country" binding:"required"`
    Street     string `json:"street" binding:"required"`
    City       string `json:"city" binding:"required"`
    PostalCode string `json:"postal_code" binding:"omitempty,min=3,max=12"`
}

type AddressRequest struct {
    UserID string  `json:"user_id" binding:"required"`
    Data   Address `json:"data" binding:"required"`
}

// UserProfile is the aggregate profile assembled across the steps
// CompletionPercent is computed based on filled sections; Completed when 100%.
type UserProfile struct {
    UserID            string           `json:"user_id"`
    AccountType       string           `json:"account_type"`
    Personal          *PersonalDetails `json:"personal,omitempty"`
    Address           *Address         `json:"address,omitempty"`
    Completed         bool             `json:"completed"`
    CompletionPercent int              `json:"completion_percent"`
    UpdatedAt         time.Time        `json:"updated_at"`
}

type ProfileResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Profile UserProfile `json:"profile"`
}