package apperror

import (
	"encoding/json"
)

var (
	ErrBadID     = NewAppError(nil, "bad user id", "")
	ErrNotFound  = NewAppError(nil, "not found", "")
	ErrCreate    = NewAppError(nil, "failed to create user", "")
	ErrUpdate    = NewAppError(nil, "failed to update user", "")
	ErrDelete    = NewAppError(nil, "failed to delete user", "")
)

type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message, developerMessage string) *AppError {
	return &AppError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func systemError(err error) *AppError {
	return NewAppError(err, "internal system error", err.Error())
}
