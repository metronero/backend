package apierror

import "gitlab.com/metronero/backend/pkg/models"

var (
	ErrInvalidSession = &models.ApiError{Code: 1, Msg: "Invalid session.", Status: 403}
	ErrUnauthorized   = &models.ApiError{Code: 2, Msg: "Unknown username or password.", Status: 401}
	ErrRequired       = &models.ApiError{Code: 3, Msg: "Required field(s) can't be empty.", Status: 400}
	ErrHash           = &models.ApiError{Code: 4, Msg: "Failed to hash password.", Status: 500}
	ErrTokenIssue     = &models.ApiError{Code: 5, Msg: "Failed to issue token.", Status: 500}
	ErrUserExists     = &models.ApiError{Code: 6, Msg: "User already exists.", Status: 400}
	ErrNoId           = &models.ApiError{Code: 7, Msg: "Unknown resource ID.", Status: 400}
	ErrBadRequest     = &models.ApiError{Code: 8, Msg: "Invalid request body.", Status: 400}
	ErrTemplateSave   = &models.ApiError{Code: 9, Msg: "Failed to save template.", Status: 500}
	ErrTemplateLoad   = &models.ApiError{Code: 10, Msg: "Failed to load template.", Status: 500}
	ErrTemplateDelete = &models.ApiError{Code: 11, Msg: "Failed to delete template.", Status: 500}
	ErrDatabase       = &models.ApiError{Code: 12, Msg: "Database errors.", Status: 500}
	ErrMoneropay      = &models.ApiError{Code: 13, Msg: "MoneroPay errors.", Status: 500}
	ErrWithdraw       = &models.ApiError{Code: 14, Msg: "Withdrawal errors.", Status: 500}
	ErrNoFunds        = &models.ApiError{Code: 15, Msg: "No funds to withdraw.", Status: 500}
	ErrPassGen        = &models.ApiError{Code: 16, Msg: "Failed to generate password.", Status: 500}
	ErrSession        = &models.ApiError{Code: 17, Msg: "Failed to access session", Status: 500}
)
