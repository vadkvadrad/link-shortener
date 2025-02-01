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

func App() http.Handler {
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

	// listening for statistic
	go statService.AddClick()

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	return stack(router)
}

func main() {
	app := App()

	// creating server
	server := http.Server{
		Addr: ":8081",
		Handler: app,
	}

	// service start
	fmt.Println("service started on port", server.Addr)
	server.ListenAndServe()
}