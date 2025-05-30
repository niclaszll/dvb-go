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

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Run all examples
	monitorStopExample(ctx, client)
	getLinesExample(ctx, client)
	getRouteExample(ctx, client)
	getPointExample(ctx, client)

	fmt.Println("=== All examples completed ===")
}

// Example 1: Monitor a specific stop (Dresden Hauptbahnhof)
func monitorStopExample(ctx context.Context, client *dvb.Client) {
	fmt.Println("=== Dresden Transport Monitor Example ===")
	fmt.Println()

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
	} else {
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
}

// Example 2: Get available lines for a stop
func getLinesExample(ctx context.Context, client *dvb.Client) {
	fmt.Println("\n=== Available Lines Example ===")
	fmt.Println()

	stopID := "33000028"
	linesOptions := &dvb.GetLinesOptions{
		StopId: stopID,
	}

	fmt.Printf("Fetching available lines for stop ID: %s\n", stopID)
	fmt.Println("---")

	linesResponse, err := client.GetLines(ctx, linesOptions)
	if err != nil {
		log.Printf("Error fetching lines: %v", err)
		return
	}

	fmt.Printf("Status: %s - %s\n", linesResponse.Status.Code, linesResponse.Status.Message)
	fmt.Printf("Expiration Time: %s\n", linesResponse.ExpirationTime)
	fmt.Println()

	if len(linesResponse.Lines) == 0 {
		fmt.Println("No lines found.")
	} else {
		fmt.Printf("Found %d available lines:\n", len(linesResponse.Lines))
		fmt.Println("---")

		for i, line := range linesResponse.Lines {
			fmt.Printf("%d. Line %s\n", i+1, line.Name)
			if len(line.Directions) > 0 {
				fmt.Printf("   Directions: ")
				for j, dir := range line.Directions {
					if j > 0 {
						fmt.Printf(", ")
					}
					fmt.Printf("%s", dir.Name)
				}
				fmt.Printf("\n")
			}
			if line.Diva.Number != "" {
				fmt.Printf("   Number: %s\n", line.Diva.Number)
			}
			fmt.Println()
		}
	}
}

// Example 3: Get route between two stops
func getRouteExample(ctx context.Context, client *dvb.Client) {
	fmt.Println("\n=== Route Planning Example ===")
	fmt.Println()

	origin := "33000742"
	destination := "33000037"

	routeOptions := &dvb.GetRouteOptions{
		Origin:      origin,
		Destination: destination,
	}

	fmt.Printf("Finding route from '%s' to '%s'\n", origin, destination)
	fmt.Println("---")

	routeResponse, err := client.GetRoute(ctx, routeOptions)
	if err != nil {
		log.Printf("Error fetching route: %v", err)
		return
	}

	fmt.Printf("Status: %s - %s\n", routeResponse.Status.Code, routeResponse.Status.Message)
	fmt.Printf("Session ID: %s\n", routeResponse.SessionId)
	fmt.Println()

	if len(routeResponse.Routes) == 0 {
		fmt.Println("No routes found.")
	} else {
		fmt.Printf("Found %d route(s):\n", len(routeResponse.Routes))
		fmt.Println("---")

		for i, route := range routeResponse.Routes {
			fmt.Printf("%d. Route %d (Duration: %d min, Interchanges: %d)\n",
				i+1, route.RouteId, route.Duration, route.Interchanges)
			fmt.Printf("   Price: %s (Day ticket: %s)\n", route.Price, route.PriceDayTicket)
			fmt.Printf("   Fare zones: %s\n", route.FareZoneNames)

			if len(route.PartialRoutes) > 0 {
				fmt.Printf("   Route segments:\n")
				for j, partial := range route.PartialRoutes {
					fmt.Printf("     %d. %s", j+1, partial.Mot.Type)
					if partial.Mot.Name != nil && *partial.Mot.Name != "" {
						fmt.Printf(" %s", *partial.Mot.Name)
					}
					if partial.Mot.Direction != nil && *partial.Mot.Direction != "" {
						fmt.Printf(" → %s", *partial.Mot.Direction)
					}
					fmt.Printf(" (%d min)\n", partial.Duration)
				}
			}
			fmt.Println()
		}
	}
}

// Example 4: Get point information
func getPointExample(ctx context.Context, client *dvb.Client) {
	fmt.Println("\n=== Point Search Example ===")
	fmt.Println()

	query := "Dresden Hauptbahnhof"
	limit := 1

	pointOptions := &dvb.GetPointOptions{
		Query: query,
		Limit: &limit,
	}

	fmt.Printf("Searching for points matching: %s\n", query)
	fmt.Println("---")

	pointResponse, err := client.GetPoint(ctx, pointOptions)
	if err != nil {
		log.Printf("Error fetching point information: %v", err)
		return
	}

	fmt.Printf("Status: %s - %s\n", pointResponse.Status.Code, pointResponse.Status.Message)
	fmt.Printf("Point Status: %s\n", pointResponse.PointStatus)
	fmt.Printf("Expiration Time: %s\n", pointResponse.ExpirationTime)
	fmt.Println()

	if len(pointResponse.Points) == 0 {
		fmt.Println("No points found.")
	} else {
		fmt.Printf("Found %d point(s):\n", len(pointResponse.Points))
		fmt.Println("---")

		for i, point := range pointResponse.Points {
			fmt.Printf("%d. %s\n", i+1, point)
		}
	}
	fmt.Println()
}
