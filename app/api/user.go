package api

import (
	"go-gin/app/form"
	"go-gin/app/repository"
	"go-gin/db"
	"go-gin/middlewares"
	"go-gin/utils/bcrypt"
	"go-gin/utils/constant"
	err2 "go-gin/utils/err"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ApplyUserAPI(app *gin.RouterGroup, resource *db.Resource) {
	userEntity := repository.NewUserEntity(resource)
	authRoute := app.Group("")
	authRoute.POST("/login", login(userEntity))
	authRoute.POST("/sign-up", signUp(userEntity))

	userRoute := app.Group("/user")
	userRoute.GET("/get-all", getAllUSer(userEntity))
	userRoute.Use(middlewares.RequireAuthenticated())               // when need authentication
	userRoute.Use(middlewares.RequireAuthorization(constant.ADMIN)) // when need authorization
	userRoute.GET("", getAllUSer(userEntity))
}

// SignUp godoc
// @Tags UserController
// @Summary Login
// @Description Login
// @Accept  json
// @Produce  json
// @Param mediaType query string false "mediaType" Enums(json, xml)
// @Param user body form.User true "User"
// @Success 200 {array} model.User
// @Router /login [post]
func login(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {

		userRequest := form.User{}
		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}

		user, code, _ := userEntity.GetOneByUsername(userRequest.Username)

		if (user == nil) || bcrypt.ComparePasswordAndHashedPassword(userRequest.Password, user.Password) != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"err": "Wrong username or password"})
			return
		}
		token := middlewares.GenerateJWTToken(*user)
		response := map[string]interface{}{
			"token": token,
			"error": nil,
		}
		ctx.JSON(code, response)
	}
}

// SignUp godoc
// @Tags UserController
// @Summary Sign Up
// @Description Sign Up
// @Accept  json
// @Produce  json
// @Param mediaType query string false "mediaType" Enums(json, xml)
// @Param user body form.User true "New User"
// @Success 200 {array} model.User
// @Router /sign-up [post]
func signUp(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		userRequest := form.User{}
		if err := ctx.Bind(&userRequest); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
			return
		}
		user, code, err := userEntity.CreateOne(userRequest)
		response := map[string]interface{}{
			"user":  user,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}

// GetAllUser godoc
// @Tags UserController
// @Summary Get all user
// @Description Get all user
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.User
// @Router /user [get]
func getAllUSer(userEntity repository.IUser) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		list, code, err := userEntity.GetAll()
		response := map[string]interface{}{
			"users": list,
			"error": err2.GetErrorMessage(err),
		}
		ctx.JSON(code, response)
	}
}
