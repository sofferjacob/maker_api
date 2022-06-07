package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sofferjacob/maker_api/conf"
	"github.com/sofferjacob/maker_api/db"
	"github.com/sofferjacob/maker_api/middleware"
	"github.com/sofferjacob/maker_api/routes"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	conf.Load()
	db.Client.Connect()
	defer db.Client.Close()
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Type", "Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	id := r.Group("/id")
	{
		id.POST("/register", routes.Register)
		id.POST("/login", routes.Login)
		id.GET("/profile", middleware.RequireAuth(), routes.Profile)
		id.PUT("/profile", middleware.RequireAuth(), routes.UpdateProfile)
	}

	collections := r.Group("/collections", middleware.RequireAuth())
	{
		collections.POST("/", routes.CreateCollection)
		collections.PUT("/", routes.UpdateCollection)
		collections.GET("/u/:uid", routes.GetUserCollections)
		collections.GET("/:id", routes.GetCollection)
		collections.GET("/query", routes.QueryCollections)
		collections.DELETE("/:id", routes.DeleteCollection)
		collections.POST("/level", routes.LinkLevel)
		collections.DELETE("/level", routes.UnlinkLevel)
		collections.GET("/levels/:id", routes.GetCollectionLevels)
		collections.GET("/trending", routes.TrendingCollections)
	}

	drafts := r.Group("/drafts", middleware.RequireAuth())
	{
		drafts.POST("/", routes.CreateDraft)
		drafts.PUT("/", routes.UpdateDraft)
		drafts.GET("/:id", routes.GetDraft)
		drafts.GET("/level/:id", routes.GetLevelDraft)
		drafts.GET("/u", routes.GetUserDrafts)
		drafts.DELETE("/:id", routes.DeleteDraft)
	}

	levels := r.Group("/levels", middleware.RequireAuth())
	{
		levels.POST("/fromDraft", routes.CreateLevelFromDraft)
		levels.POST("/", routes.CreateLevel)
		levels.GET("/info/:id", routes.GetLevelInfo)
		levels.GET("/:id", routes.GetLevel)
		levels.PUT("/fromDraft", routes.UpdateLevelFromDraft)
		levels.PUT("/", routes.UpdateLevel)
		levels.DELETE("/:id", routes.DeleteLevel)
		levels.GET("/query", routes.QueryLevels)
		levels.GET("/trending", routes.TrendingLevels)
		levels.GET("/leaderboard/:id", routes.Leaderboard)
	}

	users := r.Group("/u", middleware.RequireAuth())
	{
		users.GET("/:id", routes.GetUser)
		users.GET("/query", routes.QueryUsers)
	}

	transport := r.Group("/t", middleware.RequireAuth())
	{
		transport.POST("/", routes.PostEvent)
	}

	r.Run(fmt.Sprintf(":%v", os.Getenv("PORT")))
}
