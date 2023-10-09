package models

import "time"

type Config struct {
	ServerSetting struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server_setting"`
	DbSetting struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		Database string `json:"database"`
	} `json:"db_setting"`
}
type User struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	Password     string `json:"password"`
	PersonalData struct {
		FirstName  string `json:"first_name"`
		SecondName string `json:"second_name"`
		Patronymic string `json:"patronymic"`
		Phone      string `json:"phone"`
		Address    string `json:"address"`
		GenderId   int    `json:"gender_id"`
		Company    string `json:"company"`
	} `json:"personal_data"`
}
type Token struct {
	ID             int       `json:"id"`
	Token          string    `json:"token"`
	ExpirationTime time.Time `json:"expiration_time"`
	UserID         int       `json:"user_id"`
	CreatedAt      time.Time `json:"created_at"`
	Active         bool      `json:"active"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
}
type Check struct {
	Id         int       `json:"id"`
	Login      string    `json:"login"`
	Password   string    `json:"password"`
	PersonalId int       `json:"personal_id"`
	CreatedAt  time.Time `json:"created_at"`
	Active     bool      `json:"active"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
type Vacancies struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Terms          string    `json:"terms"`
	Duration       int       `json:"duration"`
	Fee            float64   `json:"fee"`
	UserID         int       `json:"user_id"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreatedAt      time.Time `json:"created_at"`
	Active         bool      `json:"active"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      time.Time `json:"deleted_at"`
	CategoryID     int       `json:"category"`
}
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
type Notification struct {
	ID        int    `json:"id"`
	Comment   string `json:"comment"`
	Status    bool   `json:"status"`
	OwnerID   int    `json:"owner_id"`
	SenderID  int    `json:"sender_id"`
	VacancyID int    `json:"vacancy_id"`
}
