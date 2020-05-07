package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)


var db *gorm.DB

func init()  {
	//open db connection
	var err error
	db,err = gorm.Open("mysql", "root:root@/go_api?charset=utf8&parseTime=True&loc=Local")
	if err!=nil{
		panic("failed to connect database")
	}

	//migrate the schema
	db.AutoMigrate(&todoModel{})
}
func main()  {
	//define router
	router := gin.Default()

	path:= router.Group("api/todo")
	{
		path.POST("/",createTodo)
		path.GET("/",fetchAllTodo)
		path.GET("/:id",fetchsingleTodo)
		path.PUT("/:id",updateTodo)
		path.DELETE("/:id",deleteTodo)

	}

	router.Run()
	
}

type (
	//describe a todo model type
	todoModel struct {
		gorm.Model
		Title string `json:"title"`
		Description string `json:"description"`
	}

	//represented data format for todo
	tarnsformedTodo struct {
		ID uint `json:"id"`
		Title string `json:"title"`
		Description string `json:"description"`
	}
)

//add todo information
func createTodo(c *gin.Context)  {
	todo := todoModel{
		Title: c.PostForm("title"),
		Description: c.PostForm("description")}
	//save to database
	db.Save(&todo)
	c.JSON(http.StatusCreated,gin.H{
		"status":http.StatusCreated,
		"message": "todo item created successfully"})
}

//fetch all information
func fetchAllTodo(c *gin.Context)  {
	var todos []todoModel
	var _todos []tarnsformedTodo

	//find data
	db.Find(&todos)

	if len(todos) <=0 {
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"message":"No todo found!"})
		return
	}

	for _, item:= range todos{
		_todos= append(_todos,tarnsformedTodo{
			ID: item.ID,
			Title:item.Title,
			Description:item.Description})
	}

	c.JSON(http.StatusOK , gin.H{
		"status":http.StatusOK,
		"data":_todos})
}

//fetch single information
func fetchsingleTodo(c *gin.Context)  {
	var todo todoModel
	todoID := c.Param("id")

	//find data
	db.First(&todo,todoID)

	//check data
	if todo.ID == 0{
		c.JSON(http.StatusNotFound, gin.H{
			"status":http.StatusNotFound,
			"message":"No data found!"})
		return
	}

	_todo := tarnsformedTodo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}

	c.JSON(http.StatusOK,gin.H{
		"status":http.StatusOK,
		"data": _todo,
	})
}

//update data

func updateTodo(c *gin.Context)  {
	var todo todoModel
	todoID := c.Param("id")

	//find data
	db.First(&todo, todoID)

	//check data
	if todo.ID == 0{
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"message":"No data found",
		})
		return
	}

	db.Model(&todo).Update("title",c.PostForm("title"))
	db.Model(&todo).Update("description", c.PostForm("description"))

	c.JSON(http.StatusOK,gin.H{
		"status":http.StatusOK,
		"message":"Todo updated successfully",
	})
}

//delete information
func deleteTodo(c *gin.Context)  {

	var todo todoModel
	todoID:= c.Param("id")

	//find data
	db.First(&todo,todoID)

	//check data
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound,gin.H{
			"status":http.StatusNotFound,
			"message":"No data found",
		})
		return
	}

	//delete data
	db.Delete(&todo)
	c.JSON(http.StatusOK,gin.H{
		"status":http.StatusOK,
		"message":"todo deleted successfully",
	})

}