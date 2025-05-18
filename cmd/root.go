package cmd

import (
	"fmt"
	"os"

	"github.com/golgoth31/aliasme/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "aliasme",
		Short: "AliasMe - Email alias management service",
		Long:  `A gRPC service for managing email aliases using OVH provider.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().Int("grpc-port", 9090, "The gRPC server port")
	rootCmd.PersistentFlags().Int("http-port", 8080, "The HTTP server port")

	if err := viper.BindPFlag("grpc.port", rootCmd.PersistentFlags().Lookup("grpc-port")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding grpc port flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("http.port", rootCmd.PersistentFlags().Lookup("http-port")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding http port flag: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("ALIASME")

	// Logging configuration
	viper.SetDefault("logging.format", "console")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.time_format", "2006-01-02 15:04:05")

	// Database configuration
	viper.SetDefault("database.path", "aliasme.db")

	// OVH configuration
	viper.SetDefault("ovh.endpoint", "")
	viper.SetDefault("ovh.application_key", "")
	viper.SetDefault("ovh.application_secret", "")
	viper.SetDefault("ovh.consumer_key", "")

	// SMTP configuration
	viper.SetDefault("smtp.host", "")
	viper.SetDefault("smtp.port", "")
	viper.SetDefault("smtp.username", "")
	viper.SetDefault("smtp.password", "")
	viper.SetDefault("smtp.from_email", "")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Initialize logger
	logger.Configure()
}
