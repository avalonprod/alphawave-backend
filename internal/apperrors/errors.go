package apperrors

import "errors"

var (
	ErrIncorrectVerificationCode = errors.New("incorrect verificaion code")
	ErrUserAlreadyExists         = errors.New("user with such email already exists")
	ErrIncorrectUserData         = errors.New("first name, last name, company name ust be more than 2 characters long")
	ErrUserAlreadyVerifyed       = errors.New("user with such email aready verifyed")
	ErrUserNotVerifyed           = errors.New("user with such email not verifyed")
	ErrUserNotFound              = errors.New("user doesn't exists")
	ErrMemberNotFound            = errors.New("member doesn't exists")
	ErrVerificationCodeExpired   = errors.New("the code has expired")
	ErrInternalServerError       = errors.New("thre was an internal server bug, please try again later")
	ErrIncorrectEmailFormat      = errors.New("incorrect email format")
	ErrIncorrectPasswordFormat   = errors.New("password must be at least 8 characters long and contain at least one uppercase letter and one digit")
	ErrUserBlocked               = errors.New("user is blocked")
	ErrDocumentNotFound          = errors.New("document doesn't exists")
	ErrFileNameIsEmpty           = errors.New("file name can't be empty")
	ErrTeamNotFound              = errors.New("team dosn't exists")
	ErrRoleIsNotAvailable        = errors.New("not available role")
	ErrInvalidFileType           = errors.New("invalid file type")
	ErrInvalidIdFormat           = errors.New("invalid id format")
)
