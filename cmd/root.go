package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nlowe/euchre/hosting"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "euchred",
		Short: "Euchre Game Server",
		Long:  "A hosting hosting tables of Euchre played in your browser",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt)

			s, err := hosting.NewServer()
			if err != nil {
				return err
			}

			go func() {
				logrus.Info("Starting Up")
				if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.WithError(err).Fatal("Error serving requests")
				}
			}()

			<-c
			logrus.Info("Shutting Down")

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			return s.Shutdown(ctx)
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringP("verbosity", "v", "info", "The log level to use [panic, fatal, error, warning, info, debug, trace]")

	_ = viper.BindPFlags(flags)

	cobra.OnInitialize(func() {
		level, err := logrus.ParseLevel(viper.GetString("verbosity"))
		if err != nil {
			logrus.WithError(err).Fatal("Failed to parse verbosity")
		}

		logrus.SetLevel(level)
	})

	return cmd
}
