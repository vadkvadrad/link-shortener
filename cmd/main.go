package main

import (
	"adv-go/api/configs"
	"adv-go/api/internal/auth"
	"adv-go/api/internal/link"
	"adv-go/api/internal/stat"
	"adv-go/api/internal/user"
	"adv-go/api/pkg/db"
	"adv-go/api/pkg/event"
	"adv-go/api/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.Load()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repository
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Services 
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus: eventBus,
		StatRepository: statRepository,
	})

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		Config: conf,
		EventBus: eventBus,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config: conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	// creating server
	server := http.Server{
		Addr: ":8081",
		Handler: stack(router),
	}

	// listening for statistic
	go statService.AddClick()

	// service start
	fmt.Println("service started on port", server.Addr)
	server.ListenAndServe()
}