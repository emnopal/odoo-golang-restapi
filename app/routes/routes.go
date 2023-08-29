package routes

import (
	indexController "github.com/emnopal/odoo-golang-restapi/app/controllers"
	authController "github.com/emnopal/odoo-golang-restapi/app/controllers/authHandler"
	noRouteAndMethodController "github.com/emnopal/odoo-golang-restapi/app/controllers/handlerNoRouteAndMethod"
	helpdeskTicketController "github.com/emnopal/odoo-golang-restapi/app/controllers/helpdesk"
	resPartnerController "github.com/emnopal/odoo-golang-restapi/app/controllers/resPartner"
	"github.com/emnopal/odoo-golang-restapi/app/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) *gin.Engine {

	public := r.Group("api/v1")
	// protected := r.Group("api/v1")
	private := r.Group("api/v1")
	private.Use(middleware.JwtAuthMiddleware())

	resPartner := &resPartnerController.ResPartnerController{}
	private.GET("/contact", resPartner.GetResPartner)
	private.POST("/contact", resPartner.CreateResPartner)
	private.GET("/contact/:id", resPartner.GetResPartnerById)
	private.PATCH("/contact/:id", resPartner.UpdateResPartner)
	private.DELETE("/contact/:id", resPartner.DeleteResPartner)

	helpdeskTicket := &helpdeskTicketController.HelpdeskController{}
	private.GET("/helpdesk", helpdeskTicket.GetHelpdeskTicket)
	private.GET("/helpdesk/:id", helpdeskTicket.GetHelpdeskTicketById)
	private.GET("/helpdesk/stage", helpdeskTicket.GetHelpdeskTicketStage)
	private.GET("/helpdesk/stage/:id", helpdeskTicket.GetHelpdeskTicketStageById)

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
	private.GET("/user", auth.Profile)
	private.GET("/user/:param", auth.ProfileBy)
	private.GET("/me", auth.Profile)
	private.GET("/profile", auth.Profile)
	private.GET("/logout", auth.Logout)

	handlerNoRoute := &noRouteAndMethodController.NoRouteController{}
	r.NoRoute(handlerNoRoute.NoRouteHandler)

	handlerNoMethod := &noRouteAndMethodController.NoMethodController{}
	r.NoMethod(handlerNoMethod.NoMethodHandler)

	return r
}
