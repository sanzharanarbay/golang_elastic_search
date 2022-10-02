package services

import (
	"github.com/sanzharanarbay/golang-elastic-search/application/models"
	"github.com/sanzharanarbay/golang-elastic-search/application/repositories"
)

type PostService struct {
	postRepository *repositories.PostRepository
}

func NewPostService(postRepository *repositories.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (s *PostService) IsPostAvailable(id int) (bool, error) {
	post, err := s.postRepository.GetPostById(id)
	if err == nil && post != nil {
		return true, nil
	}
	return false, err
}

func (s *PostService) GetSinglePost(ID int) (*models.Post, error) {
	post, err := s.postRepository.GetPostById(ID)
	return post, err
}

func (s *PostService) GetAllPosts() (*[]models.Post, error) {
	postList, err := s.postRepository.GetAllPosts()
	return postList, err
}

func (s *PostService) InsertPost(post *models.Post) (bool, error) {
	state, err := s.postRepository.SavePost(post)
	return state, err
}

func (s *PostService) DeletePost(id int) (bool, error) {
	var err error
	found, err := s.IsPostAvailable(id)
	if found == false {
		return false, err
	}
	state, err := s.postRepository.DeletePost(id)
	return state, err
}

func (s *PostService) UpdatePost(post *models.Post, ID int) (bool, error) {
	found, err := s.IsPostAvailable(ID)
	if found == false {
		return false, err
	}
	state, err := s.postRepository.UpdatePost(post, ID)
	return state, err
}

