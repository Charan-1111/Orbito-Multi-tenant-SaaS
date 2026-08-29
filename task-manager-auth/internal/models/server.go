package models

type Server struct {
	Port string `json:"port"`
}

type Jwt struct {
	AccessExpiryInMinutes  int `json:"accessExpiryInMinutes"`
	RefreshExpiryInMinutes int `json:"refreshExpiryInMinutes"`
}
