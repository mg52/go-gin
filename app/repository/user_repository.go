package repository

import (
	"errors"
	"mg52-gin/app/form"
	"mg52-gin/app/model"
	"mg52-gin/db"
	"mg52-gin/utils/bcrypt"
	"mg52-gin/utils/constant"
	"net/http"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserEntity IUser

type userEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type IUser interface {
	GetAll() ([]model.User, int, error)
	GetOneByUsername(username string) (*model.User, int, error)
	CreateOne(userForm form.User) (*model.User, int, error)
}

//func NewToDoEntity
func NewUserEntity(resource *db.Resource) IUser {
	userRepo := resource.DB.Collection("user")
	UserEntity = &userEntity{resource: resource, repo: userRepo}
	return UserEntity
}

func (entity *userEntity) GetAll() ([]model.User, int, error) {
	usersList := []model.User{}
	ctx, cancel := initContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{})

	if err != nil {
		logrus.Print(err)
		return []model.User{}, 400, err
	}

	for cursor.Next(ctx) {
		var user model.User
		err = cursor.Decode(&user)
		if err != nil {
			logrus.Print(err)
		}
		usersList = append(usersList, user)
	}
	return usersList, http.StatusOK, nil
}

func (entity *userEntity) GetOneByUsername(username string) (*model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	var user model.User
	err :=
		entity.repo.FindOne(ctx, bson.M{"username": username}).Decode(&user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}

func (entity *userEntity) CreateOne(userForm form.User) (*model.User, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	user := model.User{
		Id:       primitive.NewObjectID(),
		Username: userForm.Username,
		Password: bcrypt.HashPassword(userForm.Password),
		Roles:    []string{constant.ADMIN, constant.USER},
	}
	found, _, _ := entity.GetOneByUsername(user.Username)
	if found != nil {
		return nil, http.StatusBadRequest, errors.New("Username is taken")
	}
	_, err := entity.repo.InsertOne(ctx, user)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &user, http.StatusOK, nil
}
