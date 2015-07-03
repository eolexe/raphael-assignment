package db

import (

	//Imports the mysql driver

	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func InitDB(config string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", config)

	if err != nil {
		return nil, err
	}

	// Get database connection handle
	db.DB()

	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)

	return &db, nil
}

//Task is the model.
type Task struct {
	Id          int64
	Title       string
	Description string
	Priority    int
	CreatedAt   int64
	UpdatedAt   int64
	CompletedAt int64
	IsDeleted   bool
	IsCompleted bool
}

type TaskManager struct {
	Db *gorm.DB
}

func NewTaskManager(db *gorm.DB) TaskManager {
	return TaskManager{Db: db}
}

func (m *TaskManager) Close() error {
	return m.Db.Close()
}

func (m *TaskManager) All() []Task {
	var tasks []Task
	if err := m.Db.Find(&tasks).Error; err != nil {
		log.Println(err)
	}

	return tasks
}

func (m *TaskManager) Get(id int) (task Task, err error) {
	err = m.Db.First(&task, id).Error

	return
}

func (m *TaskManager) Create(task *Task) (errs []error) {
	// errs = m.validate(task)
	// if len(errs) > 0 {
	// 	return
	// }

	now := time.Now().UnixNano()
	task.CreatedAt = now
	task.UpdatedAt = now

	m.Db.Create(task)
	return
}

func (m *TaskManager) Update(task *Task) (errs []error) {
	//  errs = m.validate(task)
	if task.Id <= 0 {
		errs = append(errs, errors.New("Invalid id"))
	}

	now := time.Now().UnixNano()
	task.UpdatedAt = now

	var taskModel Task
	errs = append(errs, m.Db.Model(&taskModel).Updates(task).Error)

	return
}

func (m *TaskManager) Delete(id int) (err error) {
	task, err := m.Get(id)

	if err != nil {
		return
	}

	task.IsDeleted = true
	m.Db.Save(&task)

	return
}

// Validator

// func (m *TaskManager) validate(task *Task) (errs []error) {
//
// }
