package main

import (
	"github.com/chuckha/kubeyaml.com/backend/internal/adapters/web"
	"github.com/chuckha/kubeyaml.com/backend/internal/service/validation"
)

func main() {
	svc := validation.NewService()
	svr := web.NewServer(svc, web.WithDevMode(true))
	svr.Run()
}
