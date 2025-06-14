# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go web client for P2PQuake, a Japanese earthquake monitoring system. The application displays earthquake information, tsunami warnings, and user-reported earthquake sensing data from MongoDB collections.

## Architecture

- **Backend**: Go HTTP server with MongoDB integration
- **Frontend**: Server-side rendered HTML templates with TailwindCSS
- **Data Sources**: Two MongoDB collections - `whole` (general events) and `jma` (Japan Meteorological Agency data)
- **Event Types**: 
  - Code 551/552/556: Official earthquake/tsunami/EEW data
  - Code 9611: User-reported earthquake sensing data

### Key Components

- `main.go`: HTTP server setup with MongoDB connection
- `handler/`: HTTP request handlers for index and item pages
- `model/`: Data conversion from MongoDB BSON to Go structs
- `renderer/`: HTML template rendering engine
- `template/`: HTML templates for different event types
- `static/`: CSS and image assets

### Data Processing

The earthquake model includes complex logic for:
- Scale prioritization (震度5弱以上と推定 gets lower priority than actual 震度5弱)
- Point deduplication by city/municipality
- Prefecture-based grouping
- Time formatting for Japanese locale

## Development Commands

### Build and Run
```bash
go run main.go
```

### CSS Generation
```bash
./generate.sh
# or manually:
npx tailwindcss -i template/input.css -o static/main.css
```

### Docker Build
```bash
docker build -t web-client .
```

## Environment Variables

Required for runtime:
- `MONGODB_URL`: MongoDB connection string
- `DATABASE`: MongoDB database name  
- `COLLECTION`: MongoDB collection name for general events
- `GTM_CONTAINER_ID`: Google Tag Manager container ID (optional)

## Template System

Templates use Go's html/template with custom functions:
- `date`: Current timestamp in Japanese format
- `gtag`: Google Tag Manager container ID

Templates are organized by event type (earthquake.html, tsunami.html, eew.html, userquake.html) with a shared layout.html.