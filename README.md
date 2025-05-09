# Nostos Trips Service

## Description
The Nostos Trips service manages the creation and organization of trips, allowing users to associate photos and videos directly with each trip. It offers granular control over privacy and metadata for each stored media item.

## Features

- Complete trip creation and management
- Storage and administration of photos and videos associated with trips
- Individual privacy settings (PUBLIC/PRIVATE/FRIENDS) for each media item
- Detailed metadata storage (GPS location, capture date, altitude)
- Automatic location extraction from media EXIF data
- Manual location input for media without GPS data
- MinIO integration for efficient image and video storage
- Presigned URLs for secure media access
- Friendship-based access control for media sharing
- Automatic geocoding of coordinates to city and country information

## Technologies Used

- **Language:** Go 1.24
- **Framework:** Gin
- **Database:** PostgreSQL with GORM
- **Storage:** MinIO for media files
- **Authentication:** JWT via Auth Service
- **Secrets Management:** HashiCorp Vault
- **Orchestration:** Docker
- **EXIF Data Extraction:** rwcarlsen/goexif
- **Geocoding:** OpenStreetMap Nominatim API

## Architecture

The service follows a clean architecture pattern with the following components:

- **API Controllers:** Handle HTTP requests and responses
- **Services:** Implement business logic
- **Repositories:** Handle database operations
- **Models:** Define data structures
- **Configuration:** Manage environment and secrets

## Database Schema

The service uses multiple schemas in PostgreSQL:

- `trips.trips`: Stores trip information
- `media.media`: Stores media metadata
- `locations.locations`: Stores location information
- `albums.album_trips`: Stores album-trip relationships

## Media Features

### Automatic Metadata Extraction
The service automatically extracts the following metadata from uploaded media:
- Media type (photo/video)
- GPS coordinates (latitude, longitude, altitude)
- City and country (reverse geocoded from coordinates)
- Capture date

### Manual Location Input
For media without GPS data, the service returns a special response code (202 Accepted) with a flag indicating that manual location input is required. The frontend can then prompt the user to provide location information.

### Privacy Controls
Each media item can have one of the following visibility settings:
- PUBLIC: Visible to all users
- PRIVATE: Visible only to the owner
- FRIENDS: Visible to the owner and their friends

## Security
- Authentication: Implemented using JWT tokens from the Auth Service
- Access Control: Based on user permissions and media visibility settings
- Media Access: Secured using MinIO presigned URLs with expiration times
- Secrets Management: HashiCorp Vault for secure storage of sensitive configuration

## Structure
```
NostosTrips/
├── cmd/
│   └── main.go           # Application entry point
├── internal/
│   ├── api/              # HTTP controllers
│   ├── db/               # Database repositories
│   ├── models/           # Data models
│   └── service/          # Business logic
├── pkg/
│   ├── config/           # Configuration management
│   └── db/               # Database connection
├── Dockerfile            # Container definition
├── go.mod                # Go module definition
└── README.md             # This file
```

## Installation

Clone the repository:

```bash
git clone https://github.com/nostos-globe/NostosTrips.git
cd NostosTrips
```

## Development

To run the service locally, ensure you have Docker and Docker Compose installed. Then, execute:

```bash
go mod download
go run cmd/main.go
```

