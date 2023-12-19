package types

import "time"

type TeamsDTO struct {
	ID         string `json:"id"`
	TeamName   string `json:"teamName"`
	JobTitle   string `json:"jobTitle"`
	OwnerID    string `json:"ownerID"`
	CustomerId string `json:"customerId"`
}

type CreateTeamsDTO struct {
	TeamName string `json:"teamName"`
	JobTitle string `json:"jobTitle"`
	// OwnerID  string `json:"ownerID"`
}

type MemberDTO struct {
	MemberID      string    `json:"memberID"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	LastVisitTime time.Time `json:"lastVisitTime"`
	Roles         []string  `json:"roles"`
	Status        string    `json:"status"`
}

type UpdateTeamSettingsDTO struct {
	LogoURL               *string
	UserActivityIndicator *bool
	DisplayLinkPreview    *bool
	DisplayFilePreview    *bool
	EnableGifs            *bool
	ShowWeekends          *bool
	FirstDayOfWeek        *string
}
