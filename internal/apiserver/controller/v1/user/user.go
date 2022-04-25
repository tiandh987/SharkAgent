package user

import "github.com/tiandh987/SharkAgent/internal/apiserver/store"

type UserController struct {
	srv srvv1.Service
}
func NewUserController(store store.Factory) *UserController {
	return &UserController{
		srv: srvv1.NewService(store),
	}
}