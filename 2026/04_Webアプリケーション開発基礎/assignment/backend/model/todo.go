package model

import "errors"

// バリデーションエラー
var (
	ErrInvalidID  = errors.New("id must be positive")
	ErrEmptyTitle = errors.New("title is required")
)

// Todo はTODO項目を表す値オブジェクト
// フィールドは非公開にし、コンストラクタでのみ生成可能
type Todo struct {
	id        int
	title     string
	completed bool
}

// NewTodo は完全コンストラクタ - 必要な値を全て受け取りバリデーション後にTodoを生成
func NewTodo(id int, title string, completed bool) (Todo, error) {
	if id <= 0 {
		return Todo{}, ErrInvalidID
	}
	if title == "" {
		return Todo{}, ErrEmptyTitle
	}
	return Todo{
		id:        id,
		title:     title,
		completed: completed,
	}, nil
}

// ゲッターメソッド（セッターは用意しない = 不変性）

func (t Todo) ID() int {
	return t.id
}

func (t Todo) Title() string {
	return t.title
}

func (t Todo) Completed() bool {
	return t.completed
}

// 値を変更した新しいTodoを返すメソッド（不変性を保つ）

func (t Todo) WithTitle(title string) (Todo, error) {
	return NewTodo(t.id, title, t.completed)
}

func (t Todo) WithCompleted(completed bool) (Todo, error) {
	return NewTodo(t.id, t.title, completed)
}
