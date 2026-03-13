package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"context"
	"go-wp-tester/internal"
)

func main() {
	var url, configPath string
	var workers, duration int

	flag.StringVar(&url, "url", "", "Target WordPress URL (e.g., https://example.com)")
	flag.StringVar(&configPath, "config", "config/config.yaml", "Path to config.yaml")
	flag.IntVar(&workers, "workers", 100, "Number of concurrent workers")
	flag.IntVar(&duration, "duration", 600, "Duration of test in seconds")
	flag.Parse()

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		fmt.Println("ERROR: URL must start with http:// or https://")
		os.Exit(1)
	}

	cfg, err := internal.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Starting stress test on %s for %ds with %d workers\n", url, duration, workers)
	tester := internal.NewWordPressAPITester(cfg)
	ctx := context.Background()
	tester.Flood(ctx, url, workers, time.Duration(duration)*time.Second)
	fmt.Println("Test finished.")
} 