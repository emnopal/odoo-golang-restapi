package schemas

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AfterLogin struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
