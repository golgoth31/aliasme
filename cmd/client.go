package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	aliasme "github.com/golgoth31/aliasme/pkg/proto"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Interact with the AliasMe gRPC service",
	Long:  `Client command to interact with the AliasMe gRPC service endpoints.`,
}

var createAliasCmd = &cobra.Command{
	Use:   "create-alias",
	Short: "Create a new email alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := grpc.NewClient(
			fmt.Sprintf("localhost:%d", viper.GetInt("grpc.port")),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return fmt.Errorf("failed to connect to gRPC server: %w", err)
		}
		defer conn.Close()

		client := aliasme.NewEmailServiceClient(conn)
		resp, err := client.CreateAlias(ctx, &aliasme.CreateAliasRequest{
			UserId:      viper.GetString("alias.user_id"),
			EmailId:     viper.GetString("alias.email_id"),
			AliasPrefix: viper.GetString("alias.source"),
		})
		if err != nil {
			return fmt.Errorf("failed to create alias: %w", err)
		}

		fmt.Printf("Successfully created alias: %s\n", resp.AliasAddress)
		return nil
	},
}

var listAliasesCmd = &cobra.Command{
	Use:   "list-aliases",
	Short: "List all aliases for a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		conn, err := grpc.NewClient(
			fmt.Sprintf("localhost:%d", viper.GetInt("grpc.port")),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return fmt.Errorf("failed to connect to gRPC server: %w", err)
		}
		defer conn.Close()

		client := aliasme.NewEmailServiceClient(conn)
		resp, err := client.ListAliases(ctx, &aliasme.ListAliasesRequest{
			UserId: viper.GetString("alias.user_id"),
		})
		if err != nil {
			return fmt.Errorf("failed to list aliases: %w", err)
		}

		for _, alias := range resp.Aliases {
			fmt.Printf("Alias: %s\n", alias.AliasAddress)
		}
		return nil
	},
}

var listUsersCmd = &cobra.Command{
	Use:   "list-users",
	Short: "List all users",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		client, err := grpc.NewClient(
			fmt.Sprintf("localhost:%d", viper.GetInt("grpc.port")),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return fmt.Errorf("failed to connect to gRPC server: %w", err)
		}
		defer client.Close()

		userClient := aliasme.NewUserServiceClient(client)
		resp, err := userClient.ListUsers(ctx, &aliasme.ListUsersRequest{})
		if err != nil {
			return fmt.Errorf("failed to list users: %w", err)
		}

		if len(resp.Users) == 0 {
			fmt.Println("No users found")
			return nil
		}

		fmt.Println("Users:")
		for _, user := range resp.Users {
			fmt.Printf("- id: %s, name: %s (Email: %s)\n", user.GetId(), user.GetUsername(), user.GetEmail())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(createAliasCmd, listAliasesCmd, listUsersCmd)

	// Common flags for all client commands
	clientCmd.PersistentFlags().String("user-id", "", "User ID")
	// clientCmd.MarkPersistentFlagRequired("user-id")

	if err := viper.BindPFlag("alias.user_id", clientCmd.PersistentFlags().Lookup("user-id")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding user-id flag: %v\n", err)
		os.Exit(1)
	}

	// Flags specific to create-alias command
	createAliasCmd.Flags().String("email-id", "", "Email ID")
	createAliasCmd.Flags().String("source", "", "Source email address")

	// createAliasCmd.MarkFlagRequired("email-id")
	// createAliasCmd.MarkFlagRequired("source")

	if err := viper.BindPFlag("alias.email_id", createAliasCmd.Flags().Lookup("email-id")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding email-id flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("alias.source", createAliasCmd.Flags().Lookup("source")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding source flag: %v\n", err)
		os.Exit(1)
	}
}
