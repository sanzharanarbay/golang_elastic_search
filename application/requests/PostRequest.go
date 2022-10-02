package requests

import (
	"github.com/sanzharanarbay/golang-elastic-search/application/models"
	"net/url"
)

type PostRequest struct {
	model *models.Post
}

func NewPostRequest(model *models.Post) *PostRequest {
	return &PostRequest{
		model: model,
	}
}

func(r *PostRequest) Validate() url.Values{
	errs := url.Values{}

	if r.model.Title == "" {
		errs.Add("title", "The title field is required!")
	}

	// check the title field is between 3 to 120 chars
	if len(r.model.Title) < 3 || len(r.model.Title) > 255 {
		errs.Add("title", "The title field must be between 3-255 chars!")
	}

	if r.model.Content == "" {
		errs.Add("content", "The content field is required!")
	}

	// check the title field is between 3 to 120 chars
	if len(r.model.Content) < 3 || len(r.model.Content) > 500 {
		errs.Add("content", "The content field must be between 3-500 chars!")
	}

	if r.model.CategoryId == 0 {
		errs.Add("category_id", "The category_id field is required!")
	}
	return errs
}