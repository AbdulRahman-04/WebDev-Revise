package routes

import (
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/ai"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/controllers/private"
	"github.com/AbdulRahman-04/GoProjects/EventManagement/server/middleware"
	"github.com/gin-gonic/gin"
)

func PrivateRoutes(r *gin.Engine) {
	privateGroup := r.Group("/api/private")
	privateGroup.Use(middleware.AuthMiddleware())
	privateGroup.Use(middleware.RateLimitMiddleware(10))

	{
		// ==========================
		// 🔹 EVENT ROUTES (USER)
		// ==========================
		privateGroup.POST("/events/create", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.CreateEvent)
		privateGroup.GET("/events/all", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10), private.GetAllEvents)
		privateGroup.GET("/events/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10), private.GetOneEvent)
		privateGroup.PUT("/events/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.EditEventApi)
		privateGroup.DELETE("/events/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.DeleteOneEvent)
		privateGroup.DELETE("/events", middleware.OnlyUsers(), middleware.RateLimitMiddleware(2), private.DeleteAllEvents)

		// ==========================
		// 🔹 FUNCTION ROUTES (USER)
		// ==========================
		privateGroup.POST("/functions/create", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.CreateFunction)
		privateGroup.GET("/functions/all", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10), private.GetAllFunctions)
		privateGroup.GET("/functions/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10), private.GetOneFunction)
		privateGroup.PUT("/functions/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.EditFunction)
		privateGroup.DELETE("/functions/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.DeleteOneFunction)
		privateGroup.DELETE("/functions", middleware.OnlyUsers(), middleware.RateLimitMiddleware(2), private.DeleteAllFunctions)

		// ==========================
		// 🔹 USER ROUTES
		// ==========================
		privateGroup.GET("/users/all", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10), private.GetAllUsers)
		privateGroup.GET("/users/me", middleware.OnlyUsers(), middleware.RateLimitMiddleware(10), private.GetMyProfile)
		privateGroup.PUT("/users/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.EditUser)
		privateGroup.DELETE("/users/:id", middleware.OnlyUsers(), middleware.RateLimitMiddleware(5), private.DeleteOneUser)
		privateGroup.POST("/users/logout", middleware.OnlyUsers(), private.UserLogout)

		// ==========================
		// 🔹 ADMIN ROUTES
		// ==========================
		privateGroup.GET("/admin/events", middleware.OnlyAdmins(), private.AdminGetAllEvents)
		privateGroup.GET("/admin/events/:id", middleware.OnlyAdmins(), private.GetOneEventAdmin)
		privateGroup.GET("/admin/users", middleware.OnlyAdmins(), private.GetAllUsersAdmin)
		privateGroup.GET("/admin/users/:id", middleware.OnlyAdmins(), private.GetOneUserAdmin)
		privateGroup.GET("/admin/functions", middleware.OnlyAdmins(), private.GetAllFunctionsAdmin)
		privateGroup.GET("/admin/functions/:id", middleware.OnlyAdmins(), private.GetOneFunctionAdmin)
		privateGroup.POST("/admin/logout", middleware.OnlyAdmins(), private.AdminLogout)

		// ==========================
		// 🤖 AI ROUTES
		// ==========================
		privateGroup.POST("/ai/event-desc", middleware.OnlyUsers(), middleware.RateLimitMiddleware(3), ai.GenerateDescription)
		privateGroup.POST("/ai/function-desc", middleware.OnlyUsers(), middleware.RateLimitMiddleware(3), ai.GenerateFunctionDesc)
		privateGroup.POST("/ai/assistant", middleware.OnlyUsers(), middleware.RateLimitMiddleware(3), ai.RecommendAI)

		// ==========================
		// 🧩 EVENT JOIN APIs
		// ==========================
		privateGroup.POST("/events/:id/join", middleware.OnlyUsers(), private.JoinEvent)
		privateGroup.PUT("/events/requests/:id/approve", middleware.OnlyUsers(), private.ApproveJoinRequest)
		privateGroup.PUT("/events/requests/:id/reject", middleware.OnlyUsers(), private.RejectJoinRequest)
		privateGroup.GET("/events/requests/pending", middleware.OnlyUsers(), private.ViewPendingRequests)

		// ==========================
		// 🎉 FUNCTION JOIN APIs
		// ==========================
		privateGroup.POST("/functions/:id/join", middleware.OnlyUsers(), private.JoinFunction)
		privateGroup.PUT("/functions/requests/:id/approve", middleware.OnlyUsers(), private.ApproveFunctionJoinRequest)
		privateGroup.PUT("/functions/requests/:id/reject", middleware.OnlyUsers(), private.RejectFunctionJoinRequest)
		privateGroup.GET("/functions/requests/pending", middleware.OnlyUsers(), private.ViewPendingFunctionRequests)
	}
}
