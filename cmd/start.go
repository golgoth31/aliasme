package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"time"

	"github.com/golgoth31/aliasme/internal/database"
	"github.com/golgoth31/aliasme/internal/email"
	"github.com/golgoth31/aliasme/internal/logger"
	"github.com/golgoth31/aliasme/internal/ovh"
	"github.com/golgoth31/aliasme/internal/user"
	aliasme "github.com/golgoth31/aliasme/pkg/proto"
	"github.com/golgoth31/aliasme/pkg/static"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the AliasMe server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func startServer() error {
	// Initialize database
	db, err := database.New(&database.Config{
		Path: viper.GetString("database.path"),
	})
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize OVH client
	ovhClient, err := ovh.New(&ovh.Config{
		Endpoint:          viper.GetString("ovh.endpoint"),
		ApplicationKey:    viper.GetString("ovh.application_key"),
		ApplicationSecret: viper.GetString("ovh.application_secret"),
		ConsumerKey:       viper.GetString("ovh.consumer_key"),
	})
	if err != nil {
		// return fmt.Errorf("failed to initialize OVH client: %w", err)
		log.Error().Err(err).Msgf("failed to initialize OVH client: %w", err)
	}

	// Initialize email service
	emailService := email.New(db, email.Config{
		SMTPHost:     viper.GetString("smtp.host"),
		SMTPPort:     viper.GetString("smtp.port"),
		SMTPUsername: viper.GetString("smtp.username"),
		SMTPPassword: viper.GetString("smtp.password"),
		FromEmail:    viper.GetString("smtp.from_email"),
	})

	// Initialize user service
	userService := user.New(db)

	// Initialize email service implementation
	emailServiceImpl := email.NewEmailService(db, ovhClient, emailService)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	aliasme.RegisterUserServiceServer(grpcServer, userService)
	aliasme.RegisterEmailServiceServer(grpcServer, emailServiceImpl)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("grpc.port")))
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to listen for gRPC")
		}
		log.Info().Int("port", viper.GetInt("grpc.port")).Msg("Starting gRPC server")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC")
		}
	}()

	// Create HTTP server with gRPC-Gateway
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endpoint := fmt.Sprintf("localhost:%d", viper.GetInt("grpc.port"))

	if err := aliasme.RegisterUserServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("failed to register user service handler: %w", err)
	}
	if err := aliasme.RegisterEmailServiceHandlerFromEndpoint(ctx, mux, endpoint, opts); err != nil {
		return fmt.Errorf("failed to register email service handler: %w", err)
	}

	// Create Echo instance
	e := echo.New()
	s := &http2.Server{
		MaxConcurrentStreams: 250,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          10 * time.Second,
	}

	// Configure Echo
	e.HideBanner = true
	e.HidePort = true
	e.Use(logger.EchoLogger())
	e.Use(echoprometheus.NewMiddleware("aliasme")) // adds middleware to gather metrics
	e.GET("/metrics", echoprometheus.NewHandler()) // adds route to serve gathered metrics

	// Add gRPC-Gateway handler
	e.Any("/api/*", func(c echo.Context) error {
		mux.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	// Serve Swagger UI from embedded files
	swaggerFS, err := fs.Sub(static.Swagger, "swagger")
	if err != nil {
		return fmt.Errorf("failed to get swagger subdirectory: %w", err)
	}

	// Serve Swagger UI static files
	e.GET("/swagger/*", echo.WrapHandler(http.StripPrefix("/swagger", http.FileServer(http.FS(swaggerFS)))))

	// Start HTTP server
	log.Info().Int("port", viper.GetInt("http.port")).Msg("Starting HTTP server")
	return e.StartH2CServer(fmt.Sprintf(":%d", viper.GetInt("http.port")), s)
}
