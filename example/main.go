package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"dvb-go"
)

func main() {
	// Create a new DVB client with default configuration
	client := dvb.NewClient(dvb.Config{})

	// Example 1: Monitor a specific stop (Dresden Hauptbahnhof)
	fmt.Println("=== Dresden Transport Monitor Example ===")
	fmt.Println()

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Monitor Dresden Hauptbahnhof (stop ID: 33000028)
	stopID := "33000028"
	limit := 10

	options := &dvb.MonitorStopOptions{
		StopId: stopID,
		Limit:  &limit,
	}

	fmt.Printf("Fetching departures for stop ID: %s\n", stopID)
	fmt.Println("---")

	response, err := client.MonitorStop(ctx, options)
	if err != nil {
		log.Fatalf("Error fetching stop information: %v", err)
	}

	// Display the results
	fmt.Printf("Stop: %s\n", response.Name)
	fmt.Printf("Place: %s\n", response.Place)
	fmt.Printf("Status: %s - %s\n", response.Status.Code, response.Status.Message)
	fmt.Printf("Expiration Time: %s\n", response.ExpirationTime)
	fmt.Println()

	if len(response.Departures) == 0 {
		fmt.Println("No departures found.")
		return
	}

	fmt.Printf("Found %d departures:\n", len(response.Departures))
	fmt.Println("---")

	for i, departure := range response.Departures {
		fmt.Printf("%d. Line %s → %s\n", i+1, departure.LineName, departure.Direction)
		fmt.Printf("   Platform: %s (%s)\n", departure.Platform.Name, departure.Platform.Type)
		fmt.Printf("   Scheduled: %s\n", departure.ScheduledTime)
		fmt.Printf("   Real-time: %s\n", departure.RealTime)
		fmt.Printf("   State: %s\n", departure.State)
		if departure.Occupancy != "" {
			fmt.Printf("   Occupancy: %s\n", departure.Occupancy)
		}
		if len(departure.RouteChanges) > 0 {
			fmt.Printf("   Route Changes: %v\n", departure.RouteChanges)
		}
		fmt.Println()
	}
}
