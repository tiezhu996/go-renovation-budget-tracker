package main

import (
	"fmt"

	"renobudget/internal/config"
	"renobudget/internal/service"
	"renobudget/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st, cfg)
	_ = svc
	fmt.Println("renobudget ready")
}
