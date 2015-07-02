package db

//Task is the model.
type Task struct {
	Id          int64
	Title       string
	Description string
	Priority    int
	CreatedAt   int64
	UpdatedAt   int64
	CompletedAt bool
	IsDeleted   bool
	IsCompleted bool
}
