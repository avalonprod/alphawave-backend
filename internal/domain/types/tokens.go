package types

type Tokens struct {
	AccessToken     string
	RefreshToken    string
	MattermostToken string
}

type AuthPayload struct {
	UserId   string
	UserInfo UserDTO
	Tokens   Tokens
}
