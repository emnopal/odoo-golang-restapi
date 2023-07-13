package routes

import (
	indexController "github.com/emnopal/go_postgres/controllers"
	noRouteAndMethodController "github.com/emnopal/go_postgres/controllers/handlerNoRouteAndMethod"
	resPartnerController "github.com/emnopal/go_postgres/controllers/resPartner"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) *gin.Engine {

	resPartner := &resPartnerController.ResPartnerController{}
	r.GET("/", resPartner.GetResPartner)
	r.POST("/", resPartner.CreateResPartner)
	r.GET("/:id", resPartner.GetResPartnerById)

	index := &indexController.IndexController{}
	r.GET("/contoh", index.Index)
	r.POST("/contoh", index.Index)

	handlerNoRoute := &noRouteAndMethodController.NoRouteController{}
	r.NoRoute(handlerNoRoute.NoRouteHandler)

	handlerNoMethod := &noRouteAndMethodController.NoMethodController{}
	r.NoMethod(handlerNoMethod.NoMethodHandler)

	return r
}
