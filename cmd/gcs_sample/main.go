package main

import (
	"a-project-backend/pkg/gcs"
	"context"
	"os"
)

func main() {
	data, err := os.ReadFile("./assets/sample.webp")
	if err != nil {
		return
	}
	gcsSvc, err := gcs.NewGCS()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	if err = gcsSvc.Upload(ctx, "sample.webp", data); err != nil {
		panic(err)
	}
}
