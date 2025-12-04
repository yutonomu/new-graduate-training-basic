package repository

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"backend/model"
)

var ErrNotFound = errors.New("todo not found")

// TodoRepository はTODOデータへのアクセスを抽象化するインターフェース
type TodoRepository interface {
	List() []model.Todo
	FindByID(id int) (model.Todo, error)
	Create(todo model.Todo) error
	Update(todo model.Todo) error
	Delete(id int) error
	NextID() int
}

// FileTodoRepository はJSONファイルベースのTodoRepository実装
type FileTodoRepository struct {
	sync.Mutex
	todos    []model.Todo
	nextID   int
	filePath string
}

// NewFileTodoRepository は新しいFileTodoRepositoryを作成する
func NewFileTodoRepository(filePath string) (*FileTodoRepository, error) {
	r := &FileTodoRepository{filePath: filePath, nextID: 1}
	if err := r.load(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *FileTodoRepository) load() error {
	r.Lock()
	defer r.Unlock()

	if err := os.MkdirAll(filepath.Dir(r.filePath), 0o755); err != nil {
		return err
	}

	data, err := os.ReadFile(r.filePath)
	if errors.Is(err, os.ErrNotExist) || len(data) == 0 {
		r.todos = []model.Todo{}
		r.nextID = 1
		return nil
	}
	if err != nil {
		return err
	}

	// DTOを使ってJSONをデシリアライズ
	var dtos []TodoDTO
	if err := json.Unmarshal(data, &dtos); err != nil {
		return err
	}

	// DTOからmodel.Todoに変換
	todos, err := ToModels(dtos)
	if err != nil {
		return err
	}
	r.todos = todos

	maxID := 0
	for _, t := range r.todos {
		if t.ID() > maxID {
			maxID = t.ID()
		}
	}
	r.nextID = maxID + 1

	return nil
}

func (r *FileTodoRepository) save() error {
	// model.TodoをDTOに変換してシリアライズ
	dtos := ToDTOs(r.todos)
	payload, err := json.MarshalIndent(dtos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, payload, 0o644)
}

func (r *FileTodoRepository) findIndexByID(id int) int {
	return slices.IndexFunc(r.todos, func(t model.Todo) bool {
		return t.ID() == id
	})
}

// List は全てのTODOを取得する
func (r *FileTodoRepository) List() []model.Todo {
	r.Lock()
	defer r.Unlock()

	copied := make([]model.Todo, len(r.todos))
	copy(copied, r.todos)
	return copied
}

// FindByID は指定IDのTODOを取得する
func (r *FileTodoRepository) FindByID(id int) (model.Todo, error) {
	r.Lock()
	defer r.Unlock()

	idx := r.findIndexByID(id)
	if idx == -1 {
		return model.Todo{}, ErrNotFound
	}
	return r.todos[idx], nil
}

// Create は新しいTODOを保存する
func (r *FileTodoRepository) Create(todo model.Todo) error {
	r.Lock()
	defer r.Unlock()

	r.todos = append(r.todos, todo)
	if todo.ID() >= r.nextID {
		r.nextID = todo.ID() + 1
	}
	return r.save()
}

// Update は既存のTODOを更新する
func (r *FileTodoRepository) Update(todo model.Todo) error {
	r.Lock()
	defer r.Unlock()

	idx := r.findIndexByID(todo.ID())
	if idx == -1 {
		return ErrNotFound
	}

	r.todos[idx] = todo
	return r.save()
}

// Delete は指定IDのTODOを削除する
func (r *FileTodoRepository) Delete(id int) error {
	r.Lock()
	defer r.Unlock()

	idx := r.findIndexByID(id)
	if idx == -1 {
		return ErrNotFound
	}

	r.todos = append(r.todos[:idx], r.todos[idx+1:]...)
	return r.save()
}

// NextID は次に使用可能なIDを返す
func (r *FileTodoRepository) NextID() int {
	r.Lock()
	defer r.Unlock()
	return r.nextID
}
