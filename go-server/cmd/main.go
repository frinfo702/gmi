package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/frinfo702/fixer/delivery"
	"github.com/frinfo702/fixer/infrastructure"
	"github.com/frinfo702/fixer/usecase"
	"github.com/labstack/echo/v4"
)

func main() {
	binaryPath := os.Getenv("SEARCH_ROOT")
	if binaryPath == "" {
		log.Println("SEARCH_ROOT env variable is empty, so default value was set")
		exePath, _ := os.Executable()
		exeDir := filepath.Dir(exePath)
		binaryPath = filepath.Join(exeDir, "..", "rust-search", "target", "release", "rust_search")
	}
	seacher := infrastructure.NewRustSearchAdapter(binaryPath)
	searchInteractor := usecase.NewSearchInteractor(seacher)
	handler := delivery.NewHTTPHandler(searchInteractor)

	e := echo.New()
	handler.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
