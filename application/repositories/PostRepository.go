package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	elasticLib "github.com/olivere/elastic/v7"
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
	GetAllPosts() (* []models.Post, error)
	SavePost(*models.Post) (bool, error)
	DeletePost(ID int) (bool, error)
	UpdatePost(post *models.Post, PostID int) (bool, error)
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
	lastInsertId := 0
	sqlStatement := `INSERT into posts (title,content, category_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := p.dbClient.QueryRow(sqlStatement, post.Title, post.Content, post.CategoryId, post.CreatedAt, post.UpdatedAt).Scan(&lastInsertId)
	if err != nil {
		panic(err)
	}
	object := models.PostElastic{
		Title: post.Title,
		Content: post.Content,
		CategoryId: post.CategoryId,
		CreatedAt: post.CreatedAt.Format("2006-01-02"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02"),
	}

	// with map
	//object := map[string]interface{}{
	//	"title":        post.Title,
	//	"content":       post.Content,
	//	"category_id":post.CategoryId,
	//	"created_at": post.CreatedAt.Format("2006-01-02"),
	//	"updated_at": post.UpdatedAt.Format("2006-01-02"),
	//}

	put1, err := p.elastic.Client.Index().
		Index(p.elastic.Index).
		Id(strconv.Itoa(lastInsertId)).
		BodyJson(object).
		Do(p.elastic.Context)
	if err != nil {
		// Handle error
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Printf("Created post %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return true, nil
}

func (p *PostRepository) DeletePost(ID int) (bool, error) {
	_, err := p.dbClient.Exec(`DELETE FROM posts WHERE id=$1`, ID)
	if err != nil {
		panic(err)
	}

	del1, err := p.elastic.Client.Delete().
		Index(p.elastic.Index).
		Id(strconv.Itoa(ID)).
		Do(p.elastic.Context)
	if err != nil {
		// Handle error
		panic(err)
	}

	p.elastic.Client.Refresh(p.elastic.Index)

	fmt.Printf("Deleted post %s to index %s, type %s\n", del1.Id, del1.Index, del1.Type)

	return true, nil
}

func (p *PostRepository) UpdatePost(post *models.Post, PostID int) (bool, error) {
	post.UpdatedAt = time.Now().Local()
	sqlStatement := `UPDATE posts SET title=$1, content=$2, category_id=$3, updated_at=$4 WHERE id=$5`
	_, err := p.dbClient.Exec(sqlStatement, post.Title, post.Content, post.CategoryId, post.UpdatedAt, PostID)
	if err != nil {
		fmt.Printf("ERROR EXEC UPDATE QUERY - %s", err)
	}

	object := models.PostElastic{
		Title: post.Title,
		Content: post.Content,
		CategoryId: post.CategoryId,
		UpdatedAt: post.UpdatedAt.Format("2006-01-02"),
	}

	put1, err := p.elastic.Client.Update().
		Index(p.elastic.Index).
		Id(strconv.Itoa(PostID)).
		Doc(object).
		Do(p.elastic.Context)
	if err != nil {
		// Handle error
		panic(err)
	}

	p.elastic.Client.Refresh(p.elastic.Index)

	fmt.Printf("Updated post %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	return true, nil
}

func (p *PostRepository) SearchPost(body string) (*[]models.PostElastic, error) {
	data := models.PostElastic{}
	json.Unmarshal([]byte(body), &data)
	var postList []models.PostElastic
	matchQuery := elasticLib.NewMatchQuery("title", data.Title)
	boolQuery := elasticLib.NewBoolQuery().Must(matchQuery)

	if data.Content != "" {
		boolQuery.Must(elasticLib.NewMatchQuery("content", data.Content))
	}
	if data.CategoryId > 0 {
		boolQuery.Must(elasticLib.NewMatchQuery("category_id", data.CategoryId))
	}

	searchResult, err := p.elastic.Client.Search().
		Index(p.elastic.Index).
		Query(boolQuery).
		From(0). // Starting from this result
		Size(100).  // Limit of responds
		Do(p.elastic.Context)         // execute

	if err != nil {
		panic(err)
	}


	for _, hit := range searchResult.Hits.Hits {
		var post models.PostElastic
		err := json.Unmarshal(hit.Source, &post)
		if err != nil {
			panic(err)
		}
		 postList= append(postList, post)
	}
	return &postList, nil
}

