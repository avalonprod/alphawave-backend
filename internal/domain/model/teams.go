package model

type Team struct {
	ID         string       `json:"id" bson:"_id,omitempty"`
	TeamName   string       `json:"teamName" bson:"teamName"`
	JobTitle   string       `json:"jobTitle" bson:"jobTitle"`
	OwnerID    string       `json:"ownerID" bson:"ownerID"`
	CustomerId string       `json:"customerId" bson:"customerId"`
	Settings   TeamSettings `json:"settings" bson:"settings"`
}

type TeamSettings struct {
	LogoURL               string `json:"logoUrl" bson:"logoUrl"`
	UserActivityIndicator bool   `json:"userActivityIndicator" bson:"userActivityIndicator"`
	DisplayLinkPreview    bool   `json:"displayLinkPreview" bson:"displayLinkPreview"`
	DisplayFilePreview    bool   `json:"displayFilePreview" bson:"displayFilePreview"`
	EnableGifs            bool   `json:"enableGifs" bson:"enableGifs"`
	ShowWeekends          bool   `json:"showWeekends" bson:"showWeekends"`
	FirstDayOfWeek        string `json:"firstDayOfWeek" bson:"firstDayOfWeek"`
}

type UpdateTeamSettingsInput struct {
	LogoURL               *string
	UserActivityIndicator *bool
	DisplayLinkPreview    *bool
	DisplayFilePreview    *bool
	EnableGifs            *bool
	ShowWeekends          *bool
	FirstDayOfWeek        *string
}
