package handler

import "backend/model"

// TodoResponse はAPIレスポンス用のDTO
type TodoResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// NewTodoResponse はmodel.TodoからTodoResponseを生成する
func NewTodoResponse(t model.Todo) TodoResponse {
	return TodoResponse{
		ID:        t.ID(),
		Title:     t.Title(),
		Completed: t.Completed(),
	}
}

// NewTodoResponses はmodel.Todoのスライスから[]TodoResponseを生成する
func NewTodoResponses(todos []model.Todo) []TodoResponse {
	result := make([]TodoResponse, len(todos))
	for i, t := range todos {
		result[i] = NewTodoResponse(t)
	}
	return result
}

// CreateTodoRequest は新規TODO作成リクエスト
type CreateTodoRequest struct {
	Title string `json:"title"`
}

// UpdateTodoRequest はTODO更新リクエスト
type UpdateTodoRequest struct {
	Title     *string `json:"title"`
	Completed *bool   `json:"completed"`
}
