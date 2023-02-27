package main

import (
	"time"
)

type Config struct {
	BaseUrl string
	ApiKey  string
	Query   string
}

type Project struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Client struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"client"`
	Team struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Authors   []struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	} `json:"authors"`
	Owners []struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	} `json:"owners"`
	CustomFields []struct{} `json:"custom_fields"`
}

type Issue struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	Fields struct {
		Title          string `json:"Title"`
		RiskRating     string `json:"RiskRating"`
		ImpactRating   string `json:"ImpactRating"`
		OverallRating  string `json:"OverallRating"`
		Type           string `json:"Type"`
		Description    string `json:"Description"`
		Recommendation string `json:"Recommendation"`
		References     string `json:"References"`
	} `json:"fields"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []struct {
		Color       string `json:"color"`
		DisplayName string `json:"display_name"`
	} `json:"tags"`
}
