package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

var serverPort string
var validLevels = map[string]bool{
	"trace": true,
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start a FastHTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		handler := loggingRequest(func(ctx *fasthttp.RequestCtx) {
			fmt.Fprintf(ctx, "Hello from FastHTTP!")
		})
		addr := fmt.Sprintf(":%s", serverPort)
		if !validLevels[logLevel] {
			log.Warn().Msgf("Incorrect --log-level=%s specified. The default 'info' level will be used.", logLevel)
		}
		log.Info().Msgf("Starting FastHTTP server on port %s with '%s' logging level", addr, logLevel)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	viper.BindEnv("port")
	viper.BindEnv("default_log_level")

	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&serverPort, "port", viper.GetString("port"), "Port to run the server on")
	serverCmd.Flags().StringVar(&logLevel, "log-level", viper.GetString("default_log_level"), "Set log level for FastHTTP server: trace, debug, info, warn, error")
}

func loggingRequest(request fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		request(ctx)

		duration := time.Since(start)

		level := parseLogLevel(logLevel)
		configureLogger(level)

		log.Info().
			Str("method", string(ctx.Method())).
			Str("uri", ctx.URI().String()).
			Int("status", ctx.Response.StatusCode()).
			Str("remote_ip", ctx.RemoteAddr().String()).
			Dur("duration", duration).
			Msg("HTTP request")
	}
}
