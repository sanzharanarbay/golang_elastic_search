package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	elasticObj "github.com/sanzharanarbay/golang-elastic-search/application/configs/elastic"
	"github.com/sanzharanarbay/golang-elastic-search/application/models"
	"log"
	"strconv"
	"time"
)

type PostRepository struct {
	dbClient *sql.DB
	elastic  *elasticObj.ElasticSearch
}

func NewPostRepository(dbClient *sql.DB, elastic  *elasticObj.ElasticSearch) *PostRepository {
	return &PostRepository{
		dbClient: dbClient,
		elastic: elastic,
	}
}

type PostRepositoryInterface interface {
	GetPostById(ID int) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	SavePost(*models.Post) (bool, error)
	DeletePost(ID int) (bool, error)
	UpdatePost(*models.Post) (bool, error)
}

func (p *PostRepository) GetPostById(ID int) (*models.Post, error) {
	var post models.Post
	err := p.dbClient.QueryRow(`SELECT * FROM posts WHERE id=$1`, ID).Scan(&post.ID, &post.Title, &post.Content, &post.CategoryId,
		&post.CreatedAt, &post.UpdatedAt)
	switch err {
	case sql.ErrNoRows:
		log.Printf("No rows were returned!")
		return nil, nil
	case nil:
		return &post, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return &post, nil
}

func (p *PostRepository) GetAllPosts() (*[]models.Post, error) {
	rows, err := p.dbClient.Query("SELECT * FROM posts")
	if err != nil {
		fmt.Printf("ERROR SELECT QUERY - %s", err)
		return nil, err
	}
	var postList []models.Post
	for rows.Next() {
		var post models.Post
		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.CategoryId,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			fmt.Printf("ERROR QUERY SCAN - %s", err)
			return nil, err
		}
		postList = append(postList, post)
	}
	return &postList, nil
}

func (p *PostRepository) SavePost(post *models.Post) (bool, error) {
	post.CreatedAt = time.Now().Local()
	post.UpdatedAt = time.Now().Local()
	sqlStatement := `INSERT into posts (title,content, category_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	res, err := p.dbClient.Exec(sqlStatement, post.Title, post.Content, post.CategoryId, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		log.Printf("ERROR EXEC INSERT QUERY - %s", err)
		return false, err
	}

	documentID,_ :=res.LastInsertId()
	body, _ := json.Marshal(post)

	put1, err := p.elastic.Client.Index().
		Index(p.elastic.Index).
		Type("post").
		Id(strconv.Itoa(int(documentID))).
		BodyJson(body).
		Do(p.elastic.Context)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return true, nil
}

func (p *PostRepository) DeletePost(ID int) (bool, error) {
	_, err := p.dbClient.Exec(`DELETE FROM posts WHERE id=$1`, ID)
	if err != nil {
		log.Printf("ERROR EXEC DELETE QUERY - %s", err)
		return false, err
	}

	del1, err := p.elastic.Client.Delete().
		Index(p.elastic.Index).
		Type("post").
		Id(strconv.Itoa(ID)).
		Do(p.elastic.Context)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Indexed tweet %s to index %s, type %s\n", del1.Id, del1.Index, del1.Type)

	return true, nil
}

func (p *PostRepository) UpdatePost(post *models.Post, PostID int) (bool, error) {
	post.UpdatedAt = time.Now().Local()
	sqlStatement := `UPDATE posts SET title=$1, content=$2, category_id=$3, updated_at=$4 WHERE id=$5`
	_, err := p.dbClient.Exec(sqlStatement, post.Title, post.Content, post.CategoryId, post.UpdatedAt, PostID)
	if err != nil {
		fmt.Printf("ERROR EXEC UPDATE QUERY - %s", err)
	}

	body, _ := json.Marshal(post)

	put1, err := p.elastic.Client.Index().
		Index(p.elastic.Index).
		Type("post").
		Id(strconv.Itoa(PostID)).
		BodyJson(body).
		Do(p.elastic.Context)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return true, nil
}

