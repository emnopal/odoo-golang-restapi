package routes

import (
	indexController "github.com/emnopal/odoo-golang-restapi/controllers"
	noRouteAndMethodController "github.com/emnopal/odoo-golang-restapi/controllers/handlerNoRouteAndMethod"
	resPartnerController "github.com/emnopal/odoo-golang-restapi/controllers/resPartner"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) *gin.Engine {

	public := r.Group("api/v1")
	// protected := r.Group("api/v1")
	private := r.Group("api/v1")

	resPartner := &resPartnerController.ResPartnerController{}
	private.GET("/", resPartner.GetResPartner)
	private.POST("/", resPartner.CreateResPartner)
	private.GET("/:id", resPartner.GetResPartnerById)
	private.PATCH("/:id", resPartner.UpdateResPartner)
	private.DELETE("/:id", resPartner.DeleteResPartner)

	index := &indexController.IndexController{}
	public.GET("/contoh", index.Contoh)
	public.POST("/contoh", index.Contoh)
	public.GET("/test", index.Index)
	public.POST("/test", index.Index)
	public.PUT("/test", index.Index)
	public.PATCH("/test", index.Index)
	public.DELETE("/test", index.Index)

	handlerNoRoute := &noRouteAndMethodController.NoRouteController{}
	r.NoRoute(handlerNoRoute.NoRouteHandler)

	handlerNoMethod := &noRouteAndMethodController.NoMethodController{}
	r.NoMethod(handlerNoMethod.NoMethodHandler)

	return r
}
