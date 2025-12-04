package repository

import "backend/model"

// TodoDTO はファイル永続化用のDTO
type TodoDTO struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ToDTO はmodel.Todoを永続化用DTOに変換
func ToDTO(t model.Todo) TodoDTO {
	return TodoDTO{
		ID:        t.ID(),
		Title:     t.Title(),
		Completed: t.Completed(),
	}
}

// ToModel はDTOからmodel.Todoに変換
func (d TodoDTO) ToModel() (model.Todo, error) {
	return model.NewTodo(d.ID, d.Title, d.Completed)
}

// ToDTOs はmodel.Todoのスライスを永続化用DTOに変換
func ToDTOs(todos []model.Todo) []TodoDTO {
	result := make([]TodoDTO, len(todos))
	for i, t := range todos {
		result[i] = ToDTO(t)
	}
	return result
}

// ToModels はDTOのスライスからmodel.Todoに変換
func ToModels(dtos []TodoDTO) ([]model.Todo, error) {
	result := make([]model.Todo, len(dtos))
	for i, d := range dtos {
		t, err := d.ToModel()
		if err != nil {
			return nil, err
		}
		result[i] = t
	}
	return result, nil
}
