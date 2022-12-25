package lis

type Session struct {
	GroupId   uint64  `json:"group_id"`
	Id        string  `json:"id"`
	LastLogin string  `json:"last_login"`
	UserId    uint64  `json:"user_id"`
	Password  *string `json:password,omitempty`
}

type SessionRequest struct {
	Groupname string `json:"groupname"`
	Password  string `json:"password"`
	Username  string `json:"username"`
}
