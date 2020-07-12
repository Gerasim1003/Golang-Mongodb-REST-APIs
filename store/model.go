package store

//Hero ..
type Hero struct {
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Signed bool   `json:"signed"`
}

//User ...
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//JwtToken ...
type JwtToken struct {
	Token string `json:"token"`
}

//Exception ,...
type Exception struct {
	Message string `json:"message"`
}
