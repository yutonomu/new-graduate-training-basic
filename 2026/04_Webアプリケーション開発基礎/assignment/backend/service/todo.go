package service

import (
	"errors"

	"backend/model"
	"backend/repository"
)

var ErrNotFound = errors.New("todo not found")

// TodoService はTODOに関するビジネスロジックを定義するインターフェース
type TodoService interface {
	List() []model.Todo
	Create(title string) (model.Todo, error)
	Update(id int, title *string, completed *bool) (model.Todo, error)
	Delete(id int) error
}

// todoService はTodoServiceの実装
type todoService struct {
	repo repository.TodoRepository
}

// NewTodoService は新しいTodoServiceを作成する
func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

// List は全てのTODOを取得する
func (s *todoService) List() []model.Todo {
	return s.repo.List()
}

// Create は新しいTODOを作成する
func (s *todoService) Create(title string) (model.Todo, error) {
	todo, err := model.NewTodo(s.repo.NextID(), title, false)
	if err != nil {
		return model.Todo{}, err
	}

	if err := s.repo.Create(todo); err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

// Update は既存のTODOを更新する
func (s *todoService) Update(id int, title *string, completed *bool) (model.Todo, error) {
	todo, err := s.repo.FindByID(id)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Todo{}, ErrNotFound
	}
	if err != nil {
		return model.Todo{}, err
	}

	if title != nil {
		todo, err = todo.WithTitle(*title)
		if err != nil {
			return model.Todo{}, err
		}
	}
	if completed != nil {
		todo, err = todo.WithCompleted(*completed)
		if err != nil {
			return model.Todo{}, err
		}
	}

	if err := s.repo.Update(todo); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Todo{}, ErrNotFound
		}
		return model.Todo{}, err
	}

	return todo, nil
}

// Delete は指定IDのTODOを削除する
func (s *todoService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
