package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/config"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/conn"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	systemCtr "github.com/monstar-lab-bd/golang-starter-rest-api/system_check/controller"
	systemRepo "github.com/monstar-lab-bd/golang-starter-rest-api/system_check/repository"
	systemUseCase "github.com/monstar-lab-bd/golang-starter-rest-api/system_check/service"

	userCtr "github.com/monstar-lab-bd/golang-starter-rest-api/user/controller"
	userRepo "github.com/monstar-lab-bd/golang-starter-rest-api/user/repository"
	userService "github.com/monstar-lab-bd/golang-starter-rest-api/user/service"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
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
	useRepo := userRepo.NewUserRepository(conn.Db())
	userUC := userService.NewUserService(useRepo)
	userCtr.NewUserController(grp, userUC)
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
