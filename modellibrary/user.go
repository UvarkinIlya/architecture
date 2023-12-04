package modellibrary

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserWithActions struct {
	User
	Actions map[string]struct{} `json:"actions"`
}
