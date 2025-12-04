package main

import (
	"path/filepath"

	"backend/handler"
	"backend/repository"
	"backend/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Repository層の初期化
	repo, err := repository.NewFileTodoRepository(filepath.Join("tmp", "todos.json"))
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Service層の初期化
	svc := service.NewTodoService(repo)

	// Handler層の初期化とルーティング登録
	todoHandler := handler.NewTodoHandler(svc)
	todoHandler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
