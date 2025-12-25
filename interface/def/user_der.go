package def

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

type LoginResp struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	Code     string `json:"code"` //邮箱验证码
}

type RegisterResp struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
