package def

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
