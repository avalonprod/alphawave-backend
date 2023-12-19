package model

type Feature struct {
	Name string `json:"name" bson:"name"`
	// Limit int    `json:"limit" bson:"limit"`
}

type Package struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	StripePriceId string    `json:"stripePriceId" bson:"stripePriceId"`
	Name          string    `json:"name" bson:"name"`
	Description   string    `json:"description" bson:"description"`
	Features      []Feature `json:"features" bson:"features"`
	Price         uint      `json:"price" bson:"price"`
	Currency      string    `json:"currency" bson:"currency"`
}
