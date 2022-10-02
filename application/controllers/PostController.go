package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sanzharanarbay/golang-elastic-search/application/models"
	"github.com/sanzharanarbay/golang-elastic-search/application/requests"
	"github.com/sanzharanarbay/golang-elastic-search/application/services"
	"net/http"
	"strconv"
)

type PostController struct {
	postService *services.PostService
}

func NewPostController(postService *services.PostService) *PostController {
	return &PostController{
		postService: postService,
	}
}

func (h *PostController) GetPost(c *gin.Context) {
	var post *models.Post
	var err error
	param := c.Param("id")
	idInt, err := strconv.Atoi(param)
	post, err = h.postService.GetSinglePost(idInt)
	if err != nil {
		fmt.Printf("ERROR - %s", err)
	}
	if post != nil {
		c.JSON(http.StatusOK, post)
	} else {
		c.JSON(http.StatusNotFound, post)
	}

	return
}

func (h *PostController) GetPostList(c *gin.Context) {
	var postList *[]models.Post
	var err error
	postList, err = h.postService.GetAllPosts()
	if err != nil {
		fmt.Printf("ERROR - %s", err)
	}
	c.JSON(http.StatusOK, postList)
}

func (h *PostController) CreatePost(c *gin.Context) {
	var postToCreate models.Post

	if err := c.ShouldBindJSON(&postToCreate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	validErrs := requests.NewPostRequest(&postToCreate).Validate()
	 if len(validErrs) > 0{
		 errors := map[string]interface{}{"errors": validErrs}
		 c.JSON(http.StatusUnprocessableEntity, errors)
		 return
	 }
	created, err := h.postService.InsertPost(&postToCreate)
	if err != nil {
		fmt.Printf("ERROR - %s", err)
	}
	if created == true {
		fmt.Println("Saved Post Successfully")
	}
	c.JSON(http.StatusCreated, created)
}

func (h *PostController) DeletePost(c *gin.Context) {
	param := c.Param("id")
	idInt, err := strconv.Atoi(param)
	deleted, err := h.postService.DeletePost(idInt)
	if err != nil {
		fmt.Printf("ERROR - %s", err)
	}
	if deleted == true {
		fmt.Println("Deleted Post Successfully")
	}

	c.JSON(http.StatusOK, deleted)
}

func (h *PostController) UpdatePost(c *gin.Context) {
	var postToUpdate models.Post
	param := c.Param("id")
	idInt, err := strconv.Atoi(param)

	if err := c.ShouldBindJSON(&postToUpdate); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}

	validErrs := requests.NewPostRequest(&postToUpdate).Validate()
	if len(validErrs) > 0{
		errors := map[string]interface{}{"errors": validErrs}
		c.JSON(http.StatusUnprocessableEntity, errors)
		return
	}

	updated, err := h.postService.UpdatePost(&postToUpdate, idInt)
	if err != nil {
		fmt.Printf("ERROR - %s", err)
	}
	if updated == true {
		fmt.Println("Updated Post Successfully")
	}

	c.JSON(http.StatusCreated, updated)
}
