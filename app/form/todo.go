package form

type ToDoForm struct {
	Name   string `json:"name" binding:"required,min=5"`
	UserId string `json:"-"`
}
