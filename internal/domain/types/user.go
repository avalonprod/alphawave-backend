package types

import "time"

type UserDTO struct {
	FirstName      string       `json:"firstName"`
	LastName       string       `json:"lastName"`
	JobTitle       string       `json:"jobTitle"`
	Email          string       `json:"email"`
	Verification   bool         `json:"verification"`
	RegisteredTime time.Time    `json:"registeredTime"`
	LastVisitTime  time.Time    `json:"lastVisitTime"`
	Blocked        bool         `json:"blocked"`
	Settings       UserSettings `json:"settings"`
}

type UserImageDTO struct {
	Url              string    `json:"url"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
}

type UserSettings struct {
	UserIcon    UserImageDTO `json:"userIcon"`
	BannerImage UserImageDTO `json:"bannerImage"`
	TimeZone    string       `json:"timeZone"`
	DateFormat  string       `json:"dateFormat"`
	TimeFormat  string       `json:"timeFormat"`
}

type UserSignUpDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	JobTitle  string `json:"jobTitle"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UserSignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerificationCodeDTO struct {
	Email                       string        `json:"email"`
	VerificationCodeExpiresTime time.Duration `json:"verificationCodeExpiresTime"`
}

type ForgotPasswordPayloadDTO struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	ResultToken string `json:"resultToken"`
}

type UpdateUserInfoDTO struct {
	FirstName *string
	LastName  *string
	JobTitle  *string
	Email     *string
}

type UpdateUserImageDTO struct {
	Url              string    `json:"url"`
	Path             string    `json:"path"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
}

type UpdateUserSettingsDTO struct {
	UserIcon    *UpdateUserImageDTO
	BannerImage *UpdateUserImageDTO
	TimeZone    *string
	DateFormat  *string
	TimeFormat  *string
}

type CreateUserMattermostPayloadDTO struct {
	Email     string
	Username  string
	FirstName string
	LastName  string
	Password  string
}
