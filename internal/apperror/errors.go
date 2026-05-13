package apperror

import (
	"errors"
)

var (
	ErrNotFound             = errors.New("entity not found")
	ErrDB                   = errors.New("database error")
	ErrUserNameExist        = errors.New("email already exist")
	ErrNotFoundVerify       = errors.New("verification not found")
	ErrAlreadyUsedVerify    = errors.New("verification already used")
	ErrUserVerify           = errors.New("user email is already verified")
	ErrExpiredVerify        = errors.New("verification link has expired")
	ErrUserEmailNotVerified = errors.New("user exists but email is not verified")
	ErrEmailNotVerified     = errors.New("email is not verified, please check your inbox and verify your email")
	ErrNotFoundResetPass    = errors.New("reset pass not found")
	ErrExpiredResetPass     = errors.New("reset password link has expired")
	ErrAlreadyUsedResetPass = errors.New("reset password already used")
	ErrTooManyRequests      = errors.New("please wait before requesting again")
)
