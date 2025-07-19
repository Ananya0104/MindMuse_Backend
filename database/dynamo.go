package database

import (
	"context"
	"sync"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	// client holds the singleton instance of the DynamoDB client (aws-sdk-go-v2).
	client *dynamodb.Client
	// initErr stores any error that occurred during the client's initialization.
	initErr error
	// once ensures that the initialization logic runs exactly once.
	once sync.Once
)

func init() {
	// Call GetClient internally to leverage sync.Once for thread-safe, single initialization.
	// We ignore the returned client here as it's assigned to the package-level 'client'.
	// The error is stored in initErr, and if present, causes a panic.
	_, err := GetClient()
	if err != nil {
		// Panicking is the common way to signal unrecoverable initialization failures
		// in an init function, as the application cannot proceed without a database.
		panic(fmt.Errorf("database package initialization failed: %w", err))
	}
}

// GetClient returns a singleton DynamoDB client (*dynamodb.Client) instance.
// It performs lazy initialization the first time it's called, using sync.Once
// to ensure thread-safety and that the client is configured only once.
// It returns the client and any error that occurred during its initialization.
func GetClient() (*dynamodb.Client, error) {
	once.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("ap-south-1"), // Hardcoded region; consider making this configurable.
		)
		if err != nil {
			// If config loading fails, store the error and stop initialization.
			initErr = fmt.Errorf("failed to load AWS config: %w", err)
			return
		}

		// Create a new DynamoDB client from the loaded configuration.
		client = dynamodb.NewFromConfig(cfg)
	})
	// Return the initialized client and any error that occurred during its setup.
	return client, initErr
}

// GetInitializedClient provides direct access to the initialized DynamoDB client.
// This function assumes that the `init()` function has successfully run,
// meaning the client has been initialized without errors.
// If the client is nil (which should ideally not happen after init() success),
// it indicates a severe issue and will panic.
func GetInitializedClient() *dynamodb.Client {
	if client == nil {
		// This panic serves as a safeguard. If the init() function panicked earlier,
		// the program would have already exited. If GetClient() was never called
		// (e.g., if the init() function was removed or failed silently),
		// this would catch the nil client.
		panic("DynamoDB client not initialized. Ensure 'database' package is imported or GetClient() is called.")
	}
	return client
}
