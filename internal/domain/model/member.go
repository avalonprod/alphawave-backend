package model

type Member struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	Email       string   `json:"email" bson:"email"`
	TeamID      string   `json:"teamID" bson:"teamID"`
	UserID      string   `json:"userID" bson:"userID"`
	VerifyToken string   `json:"verifyToken" bson:"verifyToken"`
	Status      string   `json:"status" bson:"status"`
	Roles       []string `json:"roles" bson:"roles"`
}
