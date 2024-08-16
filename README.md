# USERS
### 1) userRoutes.GET("/", controllers.GetAllUsers)
```console
GET http://localhost:8080/users
```

### 2) userRoutes.POST("/", controllers.CreateUser)
```console
POST http://localhost:8080/users
{
  "name": "Neci",
  "age": 28,
  "posts": []
}
```

### 3) userRoutes.PUT("/:name", controllers.UpdateUser)
```console
PUT http://localhost:8080/users/Dek
{
  "name": "Dek",
  "age": 31,
  "posts": []
}
```

### 4) userRoutes.DELETE("/:name", controllers.DeleteUser)
```console
DELETE http://localhost:8080/users/Dek
```

# POSTS
### 1) postRoutes.POST("/", controllers.CreatePost)
```console
POST http://localhost:8080/posts
{
  "title": "My Second Post",
  "message": "This is the content of my lllll post.",
  "userId": "66be64afb75a9d0532b8beca" 
}
```

### 2) postRoutes.GET("/:id", controllers.GetPost) , postRoutes.GET("/", controllers.GetAllPosts)
```console
GET http://localhost:8080/posts/66be66f426254240eac80182
```

### 3) postRoutes.GET("/user/:userId", controllers.GetPostsByUserID)
```console
GET http://localhost:8080/posts/user/66be64afb75a9d0532b8beca
```

### 4) postRoutes.PUT("/:id", controllers.UpdatePost)
```console
PUT http://localhost:8080/posts/66be66f426254240eac80182
{
  "title": "Updated Post Title",
  "message": "Updated content of the post."
}
```

### 5) postRoutes.DELETE("/:id", controllers.DeletePost)
```console
DELETE http://localhost:8080/posts/66be66f426254240eac80182
```

