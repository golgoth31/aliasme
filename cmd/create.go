package cmd

import (
	"fmt"
	"os"

	"github.com/golgoth31/aliasme/internal/ovh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an email alias directly in OVH",
	Long:  `Create an email alias directly in OVH without going through the gRPC service.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ovh.NewClient(
			viper.GetString("ovh.endpoint"),
			viper.GetString("ovh.application_key"),
			viper.GetString("ovh.application_secret"),
			viper.GetString("ovh.consumer_key"),
		)
		if err != nil {
			return fmt.Errorf("failed to create OVH client: %w", err)
		}

		alias, err := client.CreateAlias(
			viper.GetString("alias.domain"),
			viper.GetString("alias.prefix"),
			viper.GetString("alias.destination"),
		)
		if err != nil {
			return fmt.Errorf("failed to create alias: %w", err)
		}

		fmt.Printf("Successfully created alias: %s -> %s\n", alias.Source, alias.Destination)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String("domain", "aliasme.ovh", "Domain for the alias")
	createCmd.Flags().String("prefix", "", "Prefix for alias email address")
	createCmd.Flags().String("destination", "david.sabatie@notrenet.com", "Destination email address")

	// createCmd.MarkFlagRequired("prefix")
	// createCmd.MarkFlagRequired("destination")

	if err := viper.BindPFlag("alias.domain", createCmd.Flags().Lookup("domain")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding domain flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("alias.prefix", createCmd.Flags().Lookup("prefix")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding prefix flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("alias.destination", createCmd.Flags().Lookup("destination")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding destination flag: %v\n", err)
		os.Exit(1)
	}
}
