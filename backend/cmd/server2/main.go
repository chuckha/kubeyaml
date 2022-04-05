package main

import (
	"github.com/cristifalcas/kubeyaml/backend/internal/adapters/web"
	"github.com/cristifalcas/kubeyaml/backend/internal/service/validation"
)

func main() {
	svc := validation.NewService()
	svr := web.NewServer(svc, web.WithDevMode(true))
	svr.Run()
}
