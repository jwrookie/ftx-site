package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/foxdex/ftx-site/pkg/db"
	"github.com/hashicorp/go-multierror"
	"go.uber.org/zap"

	"github.com/foxdex/ftx-site/config"
	"github.com/foxdex/ftx-site/middleware"
	"github.com/foxdex/ftx-site/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var ApiCmd = &cobra.Command{
	Use:  "api",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		setup()
	},
}

// WaitForSignal creates a channel and wait until desired signal
// arriving. It's a block call so be sure you're using it correctly.
func WaitForSignal(callback func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	sig := <-sigCh
	log.Log.Info("signal arrived",
		zap.String("signal", sig.String()),
	)

	callback()
}

func setup() {
	initApp()

	app := config.GetConfig().App
	if app.Model == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	middleware.NewRoute(r)

	httpSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.Port),
		Handler: r,
	}
	go func() {
		WaitForSignal(func() {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			if err := httpSrv.Shutdown(ctx); err != nil {
				log.Log.Error("failed to shutdown http server",
					zap.Error(err),
				)
			}
			if err := resourceRelease(); err != nil {
				log.Log.Error("failed to close context used resource connections",
					zap.Error(err),
				)
			}

		})
	}()

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Log.Fatal("failed to run status server",
					zap.Error(err),
				)
			}
		}
	}()
}

func initApp() {
	log.InitLog()
	db.NewMysql()
}

func resourceRelease() error {
	var err error

	if err1 := db.DisconnectMysql(); err1 != nil {
		err = multierror.Append(err, err1)
	}

	return err
}
