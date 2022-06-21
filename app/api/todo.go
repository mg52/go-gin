package api

import (
	"fmt"
	"go-gin/app/form"
	"go-gin/app/repository"
	"go-gin/db"
	"go-gin/middlewares"
	"go-gin/utils/constant"
	err2 "go-gin/utils/err"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ApplyToDoAPI(app *gin.RouterGroup, resource *db.Resource) {
	toDoEntity := repository.NewToDoEntity(resource)
	toDoRoute := app.Group("/todo")

	toDoRoute.Use(middlewares.RequireAuthenticated())               // when need authentication
	toDoRoute.Use(middlewares.RequireAuthorization(constant.ADMIN)) // when need authorization
	toDoRoute.GET("", getAllToDo(toDoEntity))
	toDoRoute.GET("/:id", getToDoById(toDoEntity))
	toDoRoute.POST("", createToDo(toDoEntity))
	toDoRoute.PUT("/:id", updateToDo(toDoEntity))
}

// GetAllToDo godoc
// @Tags ToDoController
// @Summary Get all ToDo TEST
// @Description Get all ToDo TEST2
// @Accept  json
// @Produce  json
// @Success 200 {array} model.ToDo
// @Router /todo [get]
// @Security ApiKeyAuth
func getAllToDo(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		jwtToken := strings.Split(token, "Bearer ")
		userId := middlewares.GetUserIdFromToken(jwtToken[1])
		fmt.Println(userId)

		list, code, err := toDoEntity.GetAll(userId)
		response := map[string]interface{}{
			"todo": list,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// Create ToDo godoc
// @Tags ToDoController
// @Summary Create ToDo
// @Description Create ToDo
// @Accept  json
// @Produce  json
// @Param mediaType query string false "mediaType" Enums(json, xml)
// @Param ToDo body form.ToDoForm true "New ToDo"
// @Success 200 {array} model.ToDo
// @Router /todo [post]
// @Security ApiKeyAuth
func createToDo(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		jwtToken := strings.Split(token, "Bearer ")
		userId := middlewares.GetUserIdFromToken(jwtToken[1])
		fmt.Println(userId)
		fmt.Println("before validation")
		todoReq := form.ToDoForm{}
		if err := ctx.Bind(&todoReq); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		fmt.Println("after validation")
		fmt.Println(todoReq)

		todoReq.UserId = userId

		fmt.Println(todoReq)

		todo, code, err := toDoEntity.CreateOne(todoReq)
		response := map[string]interface{}{
			"todo": todo,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetToDo godoc
// @Summary Get a ToDo
// @Description Get a ToDo by Id
// @Tags ToDoController
// @Accept json
// @Produce json
// @Param mediaType query string false "mediaType" Enums(json, xml)
// @Param id path string true "ToDo ID"
// @Success 200 {object} model.ToDo
// @Router /todo/{id} [get]
// @Security ApiKeyAuth
func getToDoById(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		jwtToken := strings.Split(token, "Bearer ")
		userId := middlewares.GetUserIdFromToken(jwtToken[1])
		fmt.Println(userId)

		id := ctx.Param("id")
		todo, code, err := toDoEntity.GetOneByID(userId, id)
		response := map[string]interface{}{
			"todo": todo,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// UpdateToDo godoc
// @Summary Update a ToDo
// @Description Update a ToDo by Id
// @Tags ToDoController
// @Accept json
// @Produce json
// @Param mediaType query string false "mediaType" Enums(json, xml)
// @Param id path string true "ToDo ID"
// @Param ToDo body form.ToDoForm true "Update ToDo"
// @Success 200 {object} model.ToDo
// @Router /todo/{id} [put]
// @Security ApiKeyAuth
func updateToDo(toDoEntity repository.IToDo) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		jwtToken := strings.Split(token, "Bearer ")
		userId := middlewares.GetUserIdFromToken(jwtToken[1])
		fmt.Println(userId)

		id := ctx.Param("id")
		todoReq := form.ToDoForm{}
		if err := ctx.Bind(&todoReq); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		todo, code, err := toDoEntity.Update(userId, id, todoReq)
		response := map[string]interface{}{
			"todo": todo,
			"err":  err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
