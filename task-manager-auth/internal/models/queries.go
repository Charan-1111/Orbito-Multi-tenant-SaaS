package models

type Queries struct {
	Create map[string]string `json:"create"`
	Insert Insert            `json:"insert"`
	Fetch  Fetch             `json:"fetch"`
}

type Create struct {
	TaskManagerAuth string `json:"taskManagerAuth"`
}

type Insert struct {
	RegisterUser string `json:"registerUser"`
}

type Fetch struct {
	CheckUserExistance string `json:"checkUserExistance"`
}
