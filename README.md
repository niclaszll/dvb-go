# DVB Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/niclaszll/dvb-go.svg)](https://pkg.go.dev/github.com/niclaszll/dvb-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/niclaszll/dvb-go)](https://goreportcard.com/report/github.com/niclaszll/dvb-go)
[![GitHub release](https://img.shields.io/github/release/niclaszll/dvb-go.svg)](https://github.com/niclaszll/dvb-go/releases)
[![Release Workflow](https://github.com/niclaszll/dvb-go/actions/workflows/release.yml/badge.svg)](https://github.com/niclaszll/dvb-go/actions)

Unofficial Go client for the Dresden Transport (DVB) API. This library provides easy access to real-time public transport information for Dresden, Germany.

## Installation

```bash
go get github.com/niclaszll/dvb-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/niclaszll/dvb-go"
)

func main() {
    client := dvb.NewClient(dvb.Config{})

    ctx := context.Background()

    // Monitor departures from Dresden Hauptbahnhof
    options := &dvb.MonitorStopParams{
        StopId: "33000028",
        Limit:  &[]int{10}[0],
    }

    response, err := client.MonitorStop(ctx, options)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Stop: %s\n", response.Name)
    // ...
}
```

## Available Endpoints

### `client.MonitorStop(ctx, options)`

Get real-time departures and arrivals for a specific stop.

### `client.GetRoute(ctx, options)`

Find routes between two locations with journey planning.

### `client.GetLines(ctx, options)`

Get all available transport lines serving a stop.

### `client.GetPoint(ctx, options)`

Search for stops and locations by name or query.

## Examples

See the `example/` directory for some basic usage examples.
