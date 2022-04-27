package apiserver

import (
	"github.com/gin-gonic/gin"
	"github.com/tiandh987/SharkAgent/internal/apiserver/controller/v1/user"
	"github.com/tiandh987/SharkAgent/internal/apiserver/store/mysql"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {

}

func installController(g *gin.Engine) (*gin.Engine) {

	storeIns, _ := mysql.GetMySQLFactoryOr(nil)
	v1 := g.Group("/v1")
	{
		userController := user.NewUserController(storeIns)
		userv1 := v1.Group("/users")
		{
			userv1.POST("", userController.Create)
		}
	}

	return g
}