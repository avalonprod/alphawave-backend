package model

import "time"

type Tokens struct {
	AccessToken  string `json:"accessToken" bson:"accessToken"`
	RefreshToken string `json:"refreshToken" bson:"refreshToken"`
}

type Session struct {
	RefreshToken string    `json:"refreshToken" bson:"refreshToken"`
	ExpiresTime  time.Time `json:"expiresTime" bson:"expiresTime"`
}
