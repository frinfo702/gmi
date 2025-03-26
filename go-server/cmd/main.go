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
	log.SetOutput(os.Stdout)
	log.SetPrefix("[RustySearch] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Starting RustySearch server...")

	binaryPath := os.Getenv("SEARCH_ROOT")
	if binaryPath == "" {
		log.Println("SEARCH_ROOT env variable is empty, so default value was set")
		exePath, _ := os.Executable()
		exeDir := filepath.Dir(exePath)
		binaryPath = filepath.Join(exeDir, "..", "rust-search", "target", "release", "rust_search")
	}

	absPath, err := filepath.Abs(binaryPath)
	if err == nil {
		log.Printf("Absolute path to search binary: %s", absPath)
		_, err := os.Stat(absPath)
		if err != nil {
			log.Printf("Warning: Search binary not found at path: %s, Error: %v", absPath, err)
		} else {
			log.Printf("Search binary exists at path: %s", absPath)
			binaryPath = absPath
		}
	}

	seacher := infrastructure.NewRustSearchAdapter(binaryPath)
	searchInteractor := usecase.NewSearchInteractor(seacher)
	handler := delivery.NewHTTPHandler(searchInteractor)

	e := echo.New()
	handler.RegisterRoutes(e)

	log.Println("Server is ready at http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
