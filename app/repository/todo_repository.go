package repository

import (
	"go-gin/app/form"
	"go-gin/app/model"
	"go-gin/db"
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ToDoEntity IToDo

type toDoEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type IToDo interface {
	GetAll(userId string) ([]model.ToDo, int, error)
	CreateOne(todo form.ToDoForm) (model.ToDo, int, error)
	GetOneByID(userId string, id string) (*model.ToDo, int, error) // need return pointer
	Update(userId string, id string, todo form.ToDoForm) (model.ToDo, int, error)
}

//func NewToDoEntity
func NewToDoEntity(resource *db.Resource) IToDo {
	toDoRepo := resource.DB.Collection("todo")
	ToDoEntity = &toDoEntity{resource: resource, repo: toDoRepo}
	return ToDoEntity
}

func (entity *toDoEntity) GetAll(userId string) ([]model.ToDo, int, error) {
	toDoList := []model.ToDo{}
	ctx, cancel := initContext()
	defer cancel()
	cursor, err := entity.repo.Find(ctx, bson.M{"userId": userId})

	if err != nil {
		return []model.ToDo{}, 400, err
	}

	for cursor.Next(ctx) {
		var todo model.ToDo
		err = cursor.Decode(&todo)
		if err != nil {
			logrus.Print(err)
		}
		toDoList = append(toDoList, todo)
	}
	return toDoList, http.StatusOK, nil
}

func (entity *toDoEntity) CreateOne(todoForm form.ToDoForm) (model.ToDo, int, error) {
	todo := model.ToDo{
		Id:     primitive.NewObjectID(),
		UserId: todoForm.UserId,
		Name:   todoForm.Name,
	}
	ctx, cancel := initContext()
	defer cancel()
	_, err := entity.repo.InsertOne(ctx, todo)

	if err != nil {
		return model.ToDo{}, 400, err
	}

	return todo, http.StatusOK, nil
}

func (entity *toDoEntity) GetOneByID(userId string, id string) (*model.ToDo, int, error) {
	var todo model.ToDo
	ctx, cancel := initContext()
	defer cancel()
	logrus.Print(id)
	objID, _ := primitive.ObjectIDFromHex(id)

	err := entity.repo.FindOne(ctx, bson.M{"_id": objID, "userId": userId}).Decode(&todo)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	return &todo, http.StatusOK, nil
}

func (entity *toDoEntity) Update(userId string, id string, todoForm form.ToDoForm) (model.ToDo, int, error) {
	var todo *model.ToDo
	ctx, cancel := initContext()

	defer cancel()
	objID, _ := primitive.ObjectIDFromHex(id)

	todo, _, err := entity.GetOneByID(userId, id)
	if err != nil {
		return model.ToDo{}, http.StatusNotFound, nil
	}

	todoForm.UserId = userId

	err = copier.Copy(todo, todoForm) // this is why we need return a pointer: to copy value
	if err != nil {
		logrus.Error(err)
		return model.ToDo{}, getHTTPCode(err), err
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}
	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"_id": objID, "userId": userId}, bson.M{"$set": todo}, opts).Decode(&todo)

	return *todo, http.StatusOK, nil
}
