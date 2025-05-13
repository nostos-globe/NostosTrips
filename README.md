# Nostos Trips Service

The **Nostos Trips Service** enables users to create and manage trips, and associate photos and videos with rich metadata, privacy controls, and location intelligence. Built in Go, it integrates tightly with services for storage, authentication, and geocoding.

---

## 🚀 Features

* Full trip lifecycle management
* Photo and video uploads per trip
* EXIF-based metadata extraction (GPS, date, altitude)
* Manual location input fallback
* Individual privacy settings (PUBLIC, PRIVATE, FRIENDS)
* MinIO-based media storage
* Presigned URLs for secure access
* Friendship-aware sharing logic
* Reverse geocoding via OpenStreetMap Nominatim
* Secrets management via HashiCorp Vault

---

## 📌 Endpoints

### 🔹 Trip Management

* **Create Trip**
  `POST /api/trips/`
  Creates a new trip.

* **Get All Trips**
  `GET /api/trips/`
  Retrieves all trips.

* **Search Trips**
  `POST /api/trips/search`
  Searches trips by query.

* **Public Trips**
  `GET /api/trips/public`
  Lists all publicly visible trips.

* **My Trips**
  `GET /api/trips/myTrips`
  Retrieves trips owned by the authenticated user.

* **Trips from Followed Users**
  `GET /api/trips/following`
  Lists trips created by users you follow.

* **Trips by User ID**
  `GET /api/trips/user/:id`
  Retrieves trips for a given user.

* **My Liked Trips**
  `GET /api/trips/myLikedTrips`
  Shows trips liked by the current user.

* **Get Trip by ID**
  `GET /api/trips/:id`
  Fetches details of a specific trip.

* **Get Trip Locations**
  `GET /api/trips/:id/locations`
  Retrieves location data tied to a trip.

* **Update Trip**
  `PUT /api/trips/update`
  Updates an existing trip.

* **Delete Trip**
  `DELETE /api/trips/delete/:id`
  Deletes a trip by ID.

### 🔹 Media Management

* **Upload Media to Trip**
  `POST /api/media/trip/:trip_id`
  Uploads media for a specific trip.

* **Get Media by ID**
  `GET /api/media/id/:media_id`
  Retrieves media metadata and details.

* **Get Media URL**
  `GET /api/media/:media_id`
  Returns a presigned URL for secure access.

* **Delete Media**
  `DELETE /api/media/:media_id`
  Deletes the specified media item.

* **Add Metadata to Media**
  `POST /api/media/:media_id/metadata`
  Attaches metadata to a media file.

* **Get Media Visibility**
  `GET /api/media/:media_id/visibility`
  Gets current visibility setting.

* **Change Media Visibility**
  `PUT /api/media/:media_id/visibility`
  Updates the visibility of a media item.

* **Get Media Location**
  `GET /api/media/:media_id/location`
  Returns location info for a given media file.

* **Get Media by Trip ID**
  `GET /api/media/trip/:trip_id`
  Retrieves all media linked to a trip.

---

## ⚙️ Installation and Configuration

### Prerequisites

* Go 1.24+
* Docker & Docker Compose
* PostgreSQL
* HashiCorp Vault
* MinIO
* Optional: OpenStreetMap-compatible geocoding access (e.g., Nominatim)

### Installation

```bash
git clone https://github.com/nostos-globe/NostosTrips.git
cd NostosTrips
go mod download
```

### Configuration

Ensure the following secrets are available in your Vault setup:

* `DATABASE_URL`
* `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`
* `JWT_SECRET` (shared with Auth Service)
* Any SMTP or geocoding credentials as needed

Vault can be accessed via token, AppRole, or Kubernetes Auth.

---

## ▶️ Running the Application

```bash
go run cmd/main.go
```

---

## 🏗️ Project Structure

```
NostosTrips/
├── cmd/                # App entry point
│   └── main.go
├── internal/
│   ├── api/            # HTTP controllers
│   ├── db/             # Database repositories
│   ├── models/         # Domain models
│   └── service/        # Core business logic
├── pkg/
│   ├── config/         # Config and secrets handling
│   └── db/             # DB initialization
├── Dockerfile
├── go.mod
└── README.md
```
