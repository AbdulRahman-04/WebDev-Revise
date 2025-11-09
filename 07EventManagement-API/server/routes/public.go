package routes

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/auth"
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

		// ============================
		// 🔹 ADMIN AUTH ROUTES
		// ============================
		publicGroup.POST("/admins/signup", public.AdminSignUp)
		publicGroup.POST("/admins/signin", public.AdminSignIn)
		publicGroup.GET("/admins/emailverify/:token", public.EmailVerifyAdmin)
		publicGroup.POST("/admins/forgot-password", public.AdminForgotPass)

		// ============================
		// 🔹 GOOGLE OAUTH ROUTES
		// ============================

		// 👤 User Google OAuth
		publicGroup.GET("/auth/google/user", auth.GoogleLoginUser)
		publicGroup.GET("/auth/google/user/callback", auth.GoogleCallbackUser)

		// 🛡️ Admin Google OAuth
		publicGroup.GET("/auth/google/admin", auth.GoogleLoginAdmin)
		publicGroup.GET("/auth/google/admin/callback", auth.GoogleCallbackAdmin)
	}
}
