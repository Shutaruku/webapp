package models

import "time"

type User struct {
	ID        int
	FullName  string
	Email     string
	Password  string
	Role      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Bungalow struct {
	ID           int
	BungalowName string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID         int
	FullName   string
	Email      string
	Phone      string
	StartDate  time.Time
	EndDate    time.Time
	BungalowID int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Bungalow   Bungalow
	Status     int
}

type BungalowRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	BungalowID    int
	ReservationID int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Bungalow      Bungalow
	Reservation   Reservation
	Restriction   Restriction
}

type MailData struct {
	To      string
	From    string
	Subject string
	Content string
}
