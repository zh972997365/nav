package routes

import (
	"nav/backend/controllers"
	"nav/backend/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	corsConfig.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	router.Use(cors.New(corsConfig))

	public := router.Group("/api")
	{
		public.POST("/login", controllers.Login)
		public.GET("/applications", controllers.GetApplications)
		public.GET("/menus", controllers.GetMenus)
		public.GET("/settings", controllers.GetSiteSettings)
		public.GET("/hitokoto", controllers.HitokotoHandler)
		public.GET("/hotrank", controllers.HotRankHandler)
		public.GET("/scrape-website", controllers.ScrapeWebsiteHandler)
		public.POST("/sync-icon", controllers.SyncIconHandler)
		public.POST("/save-token", controllers.SaveImageHostTokenHandler)
	}

	private := router.Group("/api")
	private.Use(middlewares.AuthMiddleware())
	{
		private.POST("/applications", controllers.CreateApplication)
		private.PUT("/applications/:id", controllers.UpdateApplication)
		private.DELETE("/applications/:id", controllers.DeleteApplication)
		private.POST("/applications/batch_delete", controllers.BatchDeleteApplications)
		private.GET("/applications/total", controllers.GetTotalApplications)
		private.GET("/applications/recent", controllers.GetRecentApplications)

		private.POST("/menus", controllers.CreateMenu)
		private.PUT("/menus/:id", controllers.UpdateMenu)
		private.DELETE("/menus/:id", controllers.DeleteMenu)
		private.POST("/menus/batch_delete", controllers.BatchDeleteMenus)
		private.GET("/menus/total", controllers.GetTotalMenus)

		private.GET("/users", controllers.GetUsers)
		private.GET("/user", controllers.GetCurrentUser)
		private.POST("/users", controllers.CreateUser)
		private.PUT("/users", controllers.UpdateUser)
		private.PUT("/users/password", controllers.UpdateUserPassword)
		private.PUT("/users/status", controllers.UpdateUserStatus)
		private.PUT("/users/admin", controllers.UpdateUserAdminStatus)
		private.DELETE("/users/:id", controllers.DeleteUser)
		private.GET("/users/total", controllers.GetTotalUsers)

		private.PUT("/settings", controllers.UpdateSiteSettings)
		private.POST("/logout", controllers.Logout)
	}

	return router
}
