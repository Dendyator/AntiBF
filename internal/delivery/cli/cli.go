package cli

import (
	"context"
	"fmt"
	"github.com/Dendyator/AntiBF/internal/delivery/grpc/proto/pb"
	"github.com/Dendyator/AntiBF/pkg/logger"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcServerAddress = "localhost:50051"

var appLogger *logger.Logger

func RunCLI(logger *logger.Logger) {
	appLogger = logger

	var rootCmd = &cobra.Command{
		Use:   "antibruteforce-cli",
		Short: "CLI tool for Anti-Bruteforce service",
	}

	rootCmd.AddCommand(resetBucketCmd)
	rootCmd.AddCommand(addToWhitelistCmd)
	rootCmd.AddCommand(removeFromWhitelistCmd)
	rootCmd.AddCommand(addToBlacklistCmd)
	rootCmd.AddCommand(removeFromBlacklistCmd)

	if err := rootCmd.Execute(); err != nil {
		appLogger.Fatalf("CLI execution failed: %v", err)
	}
}

var resetBucketCmd = &cobra.Command{
	Use:   "reset [login] [ip]",
	Short: "Reset bucket for given login and IP",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		login := args[0]
		ip := args[1]
		conn, client := createGRPCClient()
		defer func() {
			if err := conn.Close(); err != nil {
				appLogger.Errorf("Failed to close GRPC connection: %v", err)
			}
		}()

		req := &pb.ResetRequest{Login: login, Ip: ip}
		res, err := client.ResetBucket(context.Background(), req)
		if err != nil {
			appLogger.Fatalf("Could not reset bucket: %v", err)
		}

		fmt.Printf("Reset bucket succeeded: %v\n", res.Success)
		appLogger.Infof("Bucket reset for login: %s, IP: %s", login, ip)
	},
}

var addToWhitelistCmd = &cobra.Command{
	Use:   "whitelist-add [subnet]",
	Short: "Add subnet to whitelist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subnet := args[0]
		conn, client := createGRPCClient()
		defer func() {
			if err := conn.Close(); err != nil {
				appLogger.Errorf("Failed to close GRPC connection: %v", err)
			}
		}()

		req := &pb.ListRequest{Subnet: subnet}
		res, err := client.AddToWhitelist(context.Background(), req)
		if err != nil {
			appLogger.Fatalf("Could not add to whitelist: %v", err)
		}

		fmt.Printf("Add to whitelist succeeded: %v\n", res.Success)
		appLogger.Infof("Subnet added to whitelist: %s", subnet)
	},
}

var removeFromWhitelistCmd = &cobra.Command{
	Use:   "whitelist-remove [subnet]",
	Short: "Remove subnet from whitelist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subnet := args[0]
		conn, client := createGRPCClient()
		defer func() {
			if err := conn.Close(); err != nil {
				appLogger.Errorf("Failed to close GRPC connection: %v", err)
			}
		}()

		req := &pb.ListRequest{Subnet: subnet}
		res, err := client.RemoveFromWhitelist(context.Background(), req)
		if err != nil {
			appLogger.Fatalf("Could not remove from whitelist: %v", err)
		}

		fmt.Printf("Remove from whitelist succeeded: %v\n", res.Success)
		appLogger.Infof("Subnet removed from whitelist: %s", subnet)
	},
}

var addToBlacklistCmd = &cobra.Command{
	Use:   "blacklist-add [subnet]",
	Short: "Add subnet to blacklist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subnet := args[0]
		conn, client := createGRPCClient()
		defer func() {
			if err := conn.Close(); err != nil {
				appLogger.Errorf("Failed to close GRPC connection: %v", err)
			}
		}()

		req := &pb.ListRequest{Subnet: subnet}
		res, err := client.AddToBlacklist(context.Background(), req)
		if err != nil {
			appLogger.Fatalf("Could not add to blacklist: %v", err)
		}

		fmt.Printf("Add to blacklist succeeded: %v\n", res.Success)
		appLogger.Infof("Subnet added to blacklist: %s", subnet)
	},
}

var removeFromBlacklistCmd = &cobra.Command{
	Use:   "blacklist-remove [subnet]",
	Short: "Remove subnet from blacklist",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subnet := args[0]
		conn, client := createGRPCClient()
		defer func() {
			if err := conn.Close(); err != nil {
				appLogger.Errorf("Failed to close GRPC connection: %v", err)
			}
		}()

		req := &pb.ListRequest{Subnet: subnet}
		res, err := client.RemoveFromBlacklist(context.Background(), req)
		if err != nil {
			appLogger.Fatalf("Could not remove from blacklist: %v", err)
		}

		fmt.Printf("Remove from blacklist succeeded: %v\n", res.Success)
		appLogger.Infof("Subnet removed from blacklist: %s", subnet)
	},
}

func createGRPCClient() (*grpc.ClientConn, pb.AntiBruteForceClient) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(grpcServerAddress, opts...)
	if err != nil {
		appLogger.Fatalf("Failed to create gRPC client: %v", err)
	}

	client := pb.NewAntiBruteForceClient(conn)
	return conn, client
}
