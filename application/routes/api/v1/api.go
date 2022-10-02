package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sanzharanarbay/golang-elastic-search/application/configs/database"
	elasticSearch "github.com/sanzharanarbay/golang-elastic-search/application/configs/elastic"
	"github.com/sanzharanarbay/golang-elastic-search/application/controllers"
	"github.com/sanzharanarbay/golang-elastic-search/application/repositories"
	"github.com/sanzharanarbay/golang-elastic-search/application/services"
	"log"
)

func ApiRoutes(prefix string, router *gin.Engine) {
	db := database.ConnectDB()
	elastic := elasticSearch.NewElasticSearch()
	if err := elastic.CreateIndex(); err != nil {
		log.Fatalln(err)
	}

	apiGroup := router.Group(prefix)
	{
		dashboard := apiGroup.Group("/dashboard/posts")
		{
			postRepo := repositories.NewPostRepository(db)
			postService := services.NewPostService(postRepo)
			postController := controllers.NewPostController(postService)

			dashboard.GET("/all", postController.GetPostList)
			dashboard.GET("/:id", postController.GetPost)
			dashboard.POST("/create", postController.CreatePost)
			dashboard.PUT("/update/:id", postController.UpdatePost)
			dashboard.DELETE("/delete/:id", postController.DeletePost)
		}
	}
}
