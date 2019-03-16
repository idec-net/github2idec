package main

import (
	"context"
)

func main() {
	ctx := context.Background()
	config := loadConfig(filePath)

	client := config.NewClient(ctx)
	client.Run()
}
