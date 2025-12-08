# Nextbike Collector

- [Nextbike Collector](#nextbike-collector)
  - [Features](#features)
  - [Data Model](#data-model)
    - [Supported Data Entities](#supported-data-entities)
      - [Countries](#countries)
      - [Cities](#cities)
      - [Places (Stations/Hubs)](#places-stationshubs)
      - [Bikes](#bikes)
  - [Quick Start](#quick-start)
    - [Prerequisites](#prerequisites)
    - [Configuration](#configuration)
    - [Running with Docker](#running-with-docker)
    - [Development Setup](#development-setup)
  - [Architecture](#architecture)
    - [Data Flow](#data-flow)
    - [Processing Pipeline](#processing-pipeline)
      - [Stage 1: JSON Collection](#stage-1-json-collection)
      - [Stage 2: Parquet Compaction](#stage-2-parquet-compaction)
    - [Directory Structure](#directory-structure)
    - [Monitoring and Logs](#monitoring-and-logs)

A specialized data collector for Nextbike bike-sharing systems. This tool automatically collects, processes, and stores mobility data from Nextbike's live API in both JSON and efficient Parquet formats for analytics and research purposes.

## Features

- **Nextbike API Integration**: Direct integration with Nextbike's live JSON API for real-time data collection
- **Automated Scheduling**: Configurable cron-based data collection with independent scraping and compaction schedules
- **Efficient Data Processing**: Converts JSON data to columnar Parquet format with ZSTD compression for optimal storage and querying
- **Hierarchical Data Structure**: Processes nested Nextbike data (Countries → Cities → Places → Bikes) into flat, analytics-ready tables
- **Remote Storage**: Automatically transfers processed data to Samba/SMB shares for centralized storage
- **Hot Configuration Reload**: Updates collection schedules and targets without service restart
- **Production Ready**: Dockerized with comprehensive structured logging and error handling

## Data Model

The Nextbike API provides hierarchical bike-sharing data that is processed into four main entities:

### Supported Data Entities

#### Countries

- System-level information including geographic bounds, pricing, and operational details
- Metadata like hotlines, domains, currencies, and mobile app store links
- Aggregate statistics (total bikes, available bikes, booked bikes)

#### Cities

- City-level bike-sharing operations within countries
- Geographic boundaries, zoom levels, and operational settings
- City-specific bike type distributions and availability metrics
- Website links and refresh rate configurations

#### Places (Stations/Hubs)

- Individual bike stations or virtual hubs
- Real-time availability (bikes, racks, maintenance status)
- Location coordinates and addressing information
- Station-specific bike inventories and rack configurations

#### Bikes

- Individual bike tracking with unique identifiers
- Bike type classification and technical specifications
- Battery levels for electric bikes (pedelecs)
- Lock types and operational status

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24+ (for development)

### Configuration

1. Copy the example configuration:

   ```bash
   cp example.values.yaml values.yaml
   ```

2. Edit `values.yaml` with your collection preferences:

   ```yaml
   url: https://maps.nextbike.net/maps/nextbike-live.json
   output_root_path_json: /app/output-json
   output_root_path_parquet: /app/output-parquet
   
   targets:
     - name: "Germany"
       filter: "?countries=DE"                    # Collect only German data
       scrape_cron_expression: "* * * * *"        # Every minute
       compact_cron_expression: "0 12 * * *"      # Daily at noon
     - name: "Global"
       filter: ""                                 # Collect all countries
       scrape_cron_expression: "*/5 * * * *"      # Every 5 minutes  
       compact_cron_expression: "0 13 * * *"      # Daily at 1 PM
   
   # Optional: Samba/SMB remote storage
   move_parquet_files_to_samba_share_cron_expression: '0 8,11,14,17 * * *'
   SMB_address: '192.168.1.100:445'
   SMB_username: 'your-username'
   SMB_password: 'your-password'
   SMB_workgroup: 'WORKGROUP'
   SMB_path_to_mount: 'data-share'
   SMB_path_inside_mount: 'nextbike-data'
   ```

### Running with Docker

1. Build the image:

   ```bash
   docker build -t nextbike-collector:latest .
   ```

2. Run with Docker Compose:

   ```bash
   docker-compose up -d
   ```

### Development Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/PhD-Kerger/nextbike-collector.git
   cd nextbike-collector
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run locally:

   ```bash
   go run main.go [path/to/values.yaml]
   ```

## Architecture

### Data Flow

```text
Nextbike API → JSON Collection → Daily Compaction → Parquet Storage → Optional SMB Transfer
```

1. **Collection**: Scheduled HTTP requests to Nextbike's live API endpoint
2. **JSON Storage**: Raw API responses stored with timestamp-based organization
3. **Compaction**: Daily processing of JSON files into compressed Parquet format
4. **Transfer**: Optional automated upload to remote Samba shares

### Processing Pipeline

The system performs a two-stage data processing pipeline:

#### Stage 1: JSON Collection

- Fetches data from `https://maps.nextbike.net/maps/nextbike-live.json`
- Applies optional country/region filters
- Stores raw JSON responses organized by date and timestamp
- Preserves complete API response structure for data lineage

#### Stage 2: Parquet Compaction

- Processes previous day's JSON files (excludes current day to avoid incomplete data)
- Flattens hierarchical JSON structure into analytics-ready tables
- Sorts data by entity IDs and timestamps for optimal query performance
- Applies ZSTD compression for efficient storage
- Automatically cleans up source JSON files after successful processing

### Directory Structure

```text
output-json/
├── Germany/
│   ├── 2024-10-30/
│   │   ├── 1698667200.json    # Unix timestamp
│   │   ├── 1698667260.json
│   │   └── ...
│   └── 2024-10-31/
└── Global/
    └── 2024-10-30/
        └── 1698667500.json

output-parquet/
├── Germany/
│   ├── 2024-10-30/
│   │   ├── Countries.parquet   # Flattened country data
│   │   ├── Cities.parquet      # Flattened city data  
│   │   ├── Places.parquet      # Flattened station data
│   │   └── Bikes.parquet       # Flattened bike data
│   └── 2024-10-31/
└── Global/
    └── 2024-10-30/
        ├── Countries.parquet
        ├── Cities.parquet
        ├── Places.parquet
        └── Bikes.parquet
```

### Monitoring and Logs

The application provides structured logging with the following information:

- Configuration changes and hot reloads
- Data collection status and API response validation
- Parquet compaction progress and performance metrics
- Remote storage transfer status and error handling
- File system operations and cleanup activities

Logs are written to both console and the `./log` directory with rotation.
