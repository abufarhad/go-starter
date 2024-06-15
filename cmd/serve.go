package cmd

import (
	"context"
	"fmt"

	"github.com/abufarhad/golang-starter-rest-api/auth/middlewares"
	"github.com/abufarhad/golang-starter-rest-api/internal/config"
	"github.com/abufarhad/golang-starter-rest-api/internal/conn"
	"github.com/abufarhad/golang-starter-rest-api/internal/logger"
	systemCtr "github.com/abufarhad/golang-starter-rest-api/system_check/controller"
	systemRepo "github.com/abufarhad/golang-starter-rest-api/system_check/repository"
	systemUseCase "github.com/abufarhad/golang-starter-rest-api/system_check/service"
	"github.com/gin-gonic/gin"

	userCtr "github.com/abufarhad/golang-starter-rest-api/user/controller"
	userRepo "github.com/abufarhad/golang-starter-rest-api/user/repository"
	userService "github.com/abufarhad/golang-starter-rest-api/user/service"

	clubsCtr "github.com/abufarhad/golang-starter-rest-api/clubs/controller"
	clubsRepo "github.com/abufarhad/golang-starter-rest-api/clubs/repository"
	clubsService "github.com/abufarhad/golang-starter-rest-api/clubs/service"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	authCtr "github.com/abufarhad/golang-starter-rest-api/auth/controller"
	authService "github.com/abufarhad/golang-starter-rest-api/auth/service"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve will serve the system_check apis",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := conn.ConnectDb(config.Db()); err != nil {
			log.Println(err)
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		g := gin.Default()
		middlewares.Attach(g)

		appCfg := config.App()
		ApisToServe(g)

		server := &http.Server{
			Addr:    ":" + appCfg.Port,
			Handler: g,
		}
		/// start http server
		go func() {
			// service connections
			if err := server.ListenAndServe(); err != nil {
				log.Printf("listen: %s\n", err)
			}
		}()
		fmt.Println("server listening on port : ", appCfg.Port)

		// graceful shutdown
		GracefulShutdown(server)
		return nil
	},
}

func ApisToServe(g *gin.Engine) {
	grp := g.Group("api")

	//system_check pkg
	sysRepo := systemRepo.NewSystemRepository(conn.Db())
	sysUC := systemUseCase.NewSystemService(sysRepo)
	systemCtr.NewSystemController(grp, sysUC)

	//user pkg
	usrRepo := userRepo.NewUserRepository(conn.Db())
	userSvc := userService.NewUserService(usrRepo)
	userCtr.NewUserController(grp, userSvc)

	//auth pkg
	authSvc := authService.NewAuthService(usrRepo)
	authCtr.NewAuthController(grp, authSvc)

	//club pkg
	clubRepo := clubsRepo.NewClubRepository(conn.Db())
	clubsSvc := clubsService.NewClubService(clubRepo)
	clubsCtr.NewClubController(grp, clubsSvc)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// server will gracefully shutdown within 5 sec
func GracefulShutdown(srv *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	logger.Info("server shutdowns gracefully")
}
