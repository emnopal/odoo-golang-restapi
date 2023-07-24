package routes

import (
	indexController "github.com/emnopal/odoo-golang-restapi/controllers"
	authController "github.com/emnopal/odoo-golang-restapi/controllers/authHandler"
	noRouteAndMethodController "github.com/emnopal/odoo-golang-restapi/controllers/handlerNoRouteAndMethod"
	resPartnerController "github.com/emnopal/odoo-golang-restapi/controllers/resPartner"
	"github.com/emnopal/odoo-golang-restapi/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) *gin.Engine {

	public := r.Group("api/v1")
	// protected := r.Group("api/v1")
	private := r.Group("api/v1")
	private.Use(middleware.JwtAuthMiddleware())

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

	auth := &authController.AuthController{}
	public.POST("/login", auth.Login)

	handlerNoRoute := &noRouteAndMethodController.NoRouteController{}
	r.NoRoute(handlerNoRoute.NoRouteHandler)

	handlerNoMethod := &noRouteAndMethodController.NoMethodController{}
	r.NoMethod(handlerNoMethod.NoMethodHandler)

	return r
}
