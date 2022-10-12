#Golang elastic-search

#Run commands below:
- git clone {repository}
- docker-compose build --no-cache
- docker-compose up -d
- docker-compose ps
- docker-compose logs -f {serviceName}
- application run's http://localhost:8080
- elastic-search run's http://localhost:9200
- phppgadmin run's http://localhost:8087
- kibana run's http://localhost:5601



### API Documentation

1. POST `http://localhost:8080/api/v1/dashboard/posts/create`

```
{
    "title":"Test 29",
    "content":"Test 29",
    "category_id":29
}
```

2. PUT `http://localhost:8080/api/v1/dashboard/posts/update/27`

```
{
    "title":"Test 29",
    "content":"Test 29",
    "category_id":29
}
```

3. DELETE `http://localhost:8080/api/v1/dashboard/posts/delete/28`

4. GET `http://localhost:8080/api/v1/dashboard/posts/all`

5. GET `http://localhost:8080/api/v1/dashboard/posts/27` 

6. POST `http://localhost:8080/api/v1/dashboard/posts/search` 

```
{
    "title":"Test 29", 
    "content":"Test 29", // not required
    "category_id":29 // not required
}
```
