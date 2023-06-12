package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	systemCtr "github.com/monstar-lab-bd/golang-starter-rest-api/system_check/controller"
	systemRepo "github.com/monstar-lab-bd/golang-starter-rest-api/system_check/repository"
	systemUseCase "github.com/monstar-lab-bd/golang-starter-rest-api/system_check/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/conn"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve will serve the system_check apis",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := conn.ConnectDb(); err != nil {
			log.Println(err)
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		g := gin.Default()

		ApisToServe(g)

		server := &http.Server{
			Addr:    ":" + os.Getenv("APP_PORT"),
			Handler: g,
		}
		/// start http server
		go func() {
			// service connections
			if err := server.ListenAndServe(); err != nil {
				log.Printf("listen: %s\n", err)
			}
		}()
		fmt.Println("server listening on port : ", os.Getenv("APP_PORT"))

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
