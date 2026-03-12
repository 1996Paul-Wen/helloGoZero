package types

type CreateUserReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
