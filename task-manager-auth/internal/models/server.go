package models

type Server struct {
	Port string `json:"port"`
}

type Jwt struct {
	Secret                 string `json:"secret"`
	AccessExpiryInMinutes  int    `json:"accessExpiryInMinutes"`
	RefreshExpiryInMinutes int    `json:"refreshExpiryInMinutes"`
}
