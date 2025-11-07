package routes

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/public"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/gin-gonic/gin"
)

func PublicRoutes(r *gin.Engine) {
	publicGroup := r.Group("/api/public")
	publicGroup.Use(middleware.RateLimitMiddleware(5)) // global rate limit on public routes

	{
		// ============================
		// 🔹 USER AUTH ROUTES
		// ============================
		publicGroup.POST("/users/signup", public.UserSignUp)
		publicGroup.POST("/users/signin", public.UserSignIn)
		publicGroup.GET("/users/emailverify/:token", public.EmailVerifyUser)
		publicGroup.POST("/users/forgot-password", public.ForgotPass)
		publicGroup.POST("/users/refreshtoken", public.RefreshToken)
		// optional future endpoint (for later): change password
		// publicGroup.POST("/users/change-password", public.UserChangePass)

		// ============================
		// 🔹 ADMIN AUTH ROUTES
		// ============================
		publicGroup.POST("/admins/signup", public.AdminSignUp)
		publicGroup.POST("/admins/signin", public.AdminSignIn)
		publicGroup.GET("/admins/emailverify/:token", public.EmailVerifyAdmin)
		publicGroup.POST("/admins/forgot-password", public.AdminForgotPass)
		// optional future endpoint (for later): change password
		// publicGroup.POST("/admins/change-password", public.AdminChangePass)
	}
}
