package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/arbha1erao/cohereDB/pb/db_manager"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var (
		managerAddr = flag.String("addr", "127.0.0.1:9090", "DB Manager address")
		operation   = flag.String("op", "", "Operation: get, set, delete")
		key         = flag.String("key", "", "Key")
		value       = flag.String("value", "", "Value (for set operation)")
	)
	flag.Parse()

	if *operation == "" || *key == "" {
		fmt.Println("Usage:")
		fmt.Println("  Set: ./client -op=set -key=mykey -value=myvalue")
		fmt.Println("  Get: ./client -op=get -key=mykey")
		fmt.Println("  Delete: ./client -op=delete -key=mykey")
		os.Exit(1)
	}

	// Connect to DB Manager
	conn, err := grpc.NewClient(*managerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect to DB Manager: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := db_manager.NewDBManagerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch *operation {
	case "set":
		if *value == "" {
			fmt.Println("Value is required for set operation")
			os.Exit(1)
		}
		
		resp, err := client.Set(ctx, &db_manager.SetRequest{
			Key:   *key,
			Value: *value,
		})
		if err != nil {
			fmt.Printf("Set operation failed: %v\n", err)
			os.Exit(1)
		}
		
		if resp.Success {
			fmt.Printf("Successfully set key '%s' = '%s'\n", *key, *value)
		} else {
			fmt.Printf("Failed to set key '%s'\n", *key)
		}

	case "get":
		resp, err := client.Get(ctx, &db_manager.GetRequest{
			Key: *key,
		})
		if err != nil {
			fmt.Printf("Get operation failed: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Key '%s' = '%s'\n", *key, resp.Value)

	case "delete":
		resp, err := client.Delete(ctx, &db_manager.DeleteRequest{
			Key: *key,
		})
		if err != nil {
			fmt.Printf("Delete operation failed: %v\n", err)
			os.Exit(1)
		}
		
		if resp.Success {
			fmt.Printf("Successfully deleted key '%s'\n", *key)
		} else {
			fmt.Printf("Failed to delete key '%s'\n", *key)
		}

	default:
		fmt.Printf("Unknown operation: %s\n", *operation)
		fmt.Println("Supported operations: get, set, delete")
		os.Exit(1)
	}
}
