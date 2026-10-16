# Hillview Assets API

A RESTful API service for managing assets and their checkouts at Hillview TV. This service handles asset tracking, user management, and checkout/check-in workflows.

## Table of Contents

- [Overview](#overview)
- [Base URL](#base-url)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
  - [Health Check](#health-check)
  - [Checkout Operations](#checkout-operations)
  - [Create Operations](#create-operations)
  - [Read Operations](#read-operations)
  - [Validation Operations](#validation-operations)
  - [Search Operations](#search-operations)
- [Data Models](#data-models)
- [Error Handling](#error-handling)
- [Development](#development)

---

## Overview

The Hillview Assets API provides endpoints to:

- Track physical assets (equipment, devices, etc.)
- Manage users who can check out assets
- Handle checkout and check-in workflows
- Validate asset and user existence
- Search and retrieve checkout history

## Base URL

```
/assets/v1.1
```

All endpoints (except healthcheck) are prefixed with this base path.

## Authentication

CORS is enabled for all origins. The following headers are allowed:

- `X-Requested-With`
- `Content-Type`
- `Origin`
- `Authorization`
- `Accept`
- `X-CSRF-Token`

---

## API Endpoints

### Health Check

#### Check Service Health

```
GET /healthcheck
```

Returns a simple health check response to verify the service is running.

**Response:**

- `200 OK` - Service is healthy

---

### Checkout Operations

#### Checkout an Asset

```
POST /assets/v1.1/checkout
```

Checks out an asset to a user. The asset must not already be checked out.

**Request Body:**

| Field      | Type              | Required | Description                               |
| ---------- | ----------------- | -------- | ----------------------------------------- |
| `user_id`  | integer           | Yes      | ID of the user checking out the asset     |
| `asset_id` | integer           | Yes      | ID of the asset to check out              |
| `duration` | string (ISO 8601) | Yes      | Expected return date/time                 |
| `notes`    | string            | No       | Optional notes about the checkout         |
| `offsite`  | boolean           | Yes      | Whether the asset is leaving the premises |

**Example Request:**

```json
{
  "user_id": 1,
  "asset_id": 42,
  "duration": "2025-01-15T17:00:00Z",
  "notes": "Taking camera for field shoot",
  "offsite": true
}
```

**Responses:**

| Status                      | Description                    |
| --------------------------- | ------------------------------ |
| `200 OK`                    | Asset successfully checked out |
| `400 Bad Request`           | Asset is already checked out   |
| `500 Internal Server Error` | Server error                   |

---

#### Check In an Asset

```
POST /assets/v1.1/checkin
```

Checks in a previously checked-out asset.

**Request Body:**

| Field      | Type    | Required | Description                       |
| ---------- | ------- | -------- | --------------------------------- |
| `asset_id` | integer | Yes      | ID of the asset to check in       |
| `notes`    | string  | No       | Optional notes about the check-in |

**Example Request:**

```json
{
  "asset_id": 42,
  "notes": "Returned in good condition"
}
```

**Responses:**

| Status                      | Description                        |
| --------------------------- | ---------------------------------- |
| `200 OK`                    | Asset successfully checked in      |
| `400 Bad Request`           | Asset is not currently checked out |
| `500 Internal Server Error` | Server error                       |

---

### Create Operations

#### Create a New Asset

```
POST /assets/v1.1/create/asset
```

Creates a new asset in the system.

**Request Body:**

| Field           | Type    | Required | Description                         |
| --------------- | ------- | -------- | ----------------------------------- |
| `name`          | string  | Yes      | Name of the asset                   |
| `identifier`    | string  | Yes      | Unique tag/identifier for the asset |
| `category`      | integer | Yes      | Category ID for the asset           |
| `image_url`     | string  | No       | URL to an image of the asset        |
| `description`   | string  | No       | Description of the asset            |
| `serial_number` | string  | No       | Serial number of the asset          |
| `manufacturer`  | string  | No       | Manufacturer name                   |
| `model`         | string  | No       | Model name/number                   |
| `notes`         | string  | No       | Additional notes about the asset    |

**Example Request:**

```json
{
  "name": "Canon EOS R5",
  "identifier": "CAM-001",
  "category": 1,
  "image_url": "https://example.com/images/canon-r5.jpg",
  "description": "Professional mirrorless camera",
  "serial_number": "SN123456789",
  "manufacturer": "Canon",
  "model": "EOS R5",
  "notes": "Purchased January 2024"
}
```

**Responses:**

| Status                      | Description                                |
| --------------------------- | ------------------------------------------ |
| `201 Created`               | Asset successfully created                 |
| `400 Bad Request`           | Missing required fields or invalid request |
| `500 Internal Server Error` | Server error                               |

---

#### Create a New User

```
POST /assets/v1.1/create/user
```

Creates a new user who can check out assets.

**Request Body:**

| Field       | Type   | Required | Description                        |
| ----------- | ------ | -------- | ---------------------------------- |
| `name`      | string | Yes      | Full name of the user              |
| `email`     | string | Yes      | Email address                      |
| `tag`       | string | Yes      | Unique identifier/tag for the user |
| `photo_url` | string | Yes      | URL to the user's profile photo    |

**Example Request:**

```json
{
  "name": "John Doe",
  "email": "john.doe@hillview.tv",
  "tag": "USR-001",
  "photo_url": "https://example.com/photos/john.jpg"
}
```

**Responses:**

| Status                      | Description               |
| --------------------------- | ------------------------- |
| `201 Created`               | User successfully created |
| `400 Bad Request`           | Missing required fields   |
| `500 Internal Server Error` | Server error              |

---

### Read Operations

#### Get Asset by ID

```
GET /assets/v1.1/read/assetByID/{id}
```

Retrieves a single asset by its database ID, including active checkout information if applicable.

**Path Parameters:**

| Parameter | Type    | Description             |
| --------- | ------- | ----------------------- |
| `id`      | integer | The asset's database ID |

**Example Response:**

```json
{
  "id": 42,
  "name": "Canon EOS R5",
  "image_url": "https://example.com/images/canon-r5.jpg",
  "identifier": "CAM-001",
  "description": "Professional mirrorless camera",
  "category": {
    "id": 1,
    "name": "Camera",
    "short_name": "CAM"
  },
  "status": {
    "id": 1,
    "name": "Available",
    "short_name": "AVL"
  },
  "metadata": {
    "serial_number": "SN123456789",
    "manufacturer": "Canon",
    "model": "EOS R5",
    "notes": "Purchased January 2024"
  },
  "active_tab": null
}
```

**Responses:**

| Status                      | Description              |
| --------------------------- | ------------------------ |
| `200 OK`                    | Returns the asset object |
| `500 Internal Server Error` | Server error             |

---

#### Get Asset by Tag

```
GET /assets/v1.1/read/assetByTag/{id}
```

Retrieves a single asset by its identifier/tag.

**Path Parameters:**

| Parameter | Type   | Description                       |
| --------- | ------ | --------------------------------- |
| `id`      | string | The asset's unique identifier/tag |

**Example:** `GET /assets/v1.1/read/assetByTag/CAM-001`

**Response:** Same as [Get Asset by ID](#get-asset-by-id)

---

#### Get User by ID

```
GET /assets/v1.1/read/userByID/{id}
```

Retrieves a single user by their database ID.

**Path Parameters:**

| Parameter | Type    | Description            |
| --------- | ------- | ---------------------- |
| `id`      | integer | The user's database ID |

**Example Response:**

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@hillview.tv",
  "identifier": "USR-001",
  "profile_image_url": "https://example.com/photos/john.jpg",
  "inserted_at": "2024-01-15T10:30:00Z"
}
```

**Responses:**

| Status                      | Description             |
| --------------------------- | ----------------------- |
| `200 OK`                    | Returns the user object |
| `500 Internal Server Error` | Server error            |

---

#### Get User by Tag

```
GET /assets/v1.1/read/userByTag/{id}
```

Retrieves a single user by their identifier/tag.

**Path Parameters:**

| Parameter | Type   | Description                      |
| --------- | ------ | -------------------------------- |
| `id`      | string | The user's unique identifier/tag |

**Example:** `GET /assets/v1.1/read/userByTag/USR-001`

**Response:** Same as [Get User by ID](#get-user-by-id)

---

#### Get Asset Checkout History

```
GET /assets/v1.1/read/assetCheckoutHistory
```

Retrieves the checkout history for an asset. Returns the last 10 checkout records.

**Query Parameters:**

| Parameter | Type   | Required | Description          |
| --------- | ------ | -------- | -------------------- |
| `id`      | string | No\*     | Asset database ID    |
| `tag`     | string | No\*     | Asset identifier/tag |

\*At least one of `id` or `tag` should be provided.

**Example Request:** `GET /assets/v1.1/read/assetCheckoutHistory?tag=CAM-001`

**Example Response:**

```json
[
  {
    "id": 100,
    "asset_id": 42,
    "checkout_status": {
      "id": 2,
      "name": "Checked In",
      "short_name": "IN"
    },
    "associated_user": 1,
    "checkout_notes": "Returned in good condition",
    "time_out": "2024-12-20T09:00:00Z",
    "time_in": "2024-12-20T17:00:00Z",
    "expected_in": "2024-12-20T18:00:00Z",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john.doe@hillview.tv",
      "identifier": "USR-001",
      "profile_image_url": "https://example.com/photos/john.jpg",
      "inserted_at": "2024-01-15T10:30:00Z"
    }
  }
]
```

**Responses:**

| Status                      | Description                       |
| --------------------------- | --------------------------------- |
| `200 OK`                    | Returns array of checkout records |
| `500 Internal Server Error` | Server error                      |

---

### Validation Operations

These endpoints validate whether a resource exists without returning the full object. Useful for quick existence checks.

#### Validate Asset by ID

```
GET /assets/v1.1/valid/assetByID/{id}
```

Checks if an asset exists by database ID.

**Path Parameters:**

| Parameter | Type    | Description             |
| --------- | ------- | ----------------------- |
| `id`      | integer | The asset's database ID |

**Responses:**

| Status                      | Description          |
| --------------------------- | -------------------- |
| `200 OK`                    | Asset exists         |
| `204 No Content`            | Asset does not exist |
| `500 Internal Server Error` | Server error         |

---

#### Validate Asset by Tag

```
GET /assets/v1.1/valid/assetByTag/{id}
```

Checks if an asset exists by identifier/tag.

**Path Parameters:**

| Parameter | Type   | Description                       |
| --------- | ------ | --------------------------------- |
| `id`      | string | The asset's unique identifier/tag |

**Responses:**

| Status                      | Description          |
| --------------------------- | -------------------- |
| `200 OK`                    | Asset exists         |
| `204 No Content`            | Asset does not exist |
| `500 Internal Server Error` | Server error         |

---

#### Validate User by ID

```
GET /assets/v1.1/valid/userByID/{id}
```

Checks if a user exists by database ID.

**Path Parameters:**

| Parameter | Type    | Description            |
| --------- | ------- | ---------------------- |
| `id`      | integer | The user's database ID |

**Responses:**

| Status                      | Description         |
| --------------------------- | ------------------- |
| `200 OK`                    | User exists         |
| `204 No Content`            | User does not exist |
| `500 Internal Server Error` | Server error        |

---

#### Validate User by Tag

```
GET /assets/v1.1/valid/userByTag/{id}
```

Checks if a user exists by identifier/tag.

**Path Parameters:**

| Parameter | Type   | Description                      |
| --------- | ------ | -------------------------------- |
| `id`      | string | The user's unique identifier/tag |

**Responses:**

| Status                      | Description         |
| --------------------------- | ------------------- |
| `200 OK`                    | User exists         |
| `204 No Content`            | User does not exist |
| `500 Internal Server Error` | Server error        |

---

### Search Operations

#### Search Checkouts

```
GET /assets/v1.1/search/checkouts
```

Search for checkouts. _(Currently a placeholder - to be implemented)_

**Responses:**

| Status   | Description                                       |
| -------- | ------------------------------------------------- |
| `200 OK` | Search results (currently returns empty response) |

---

## Data Models

### Asset

| Field         | Type          | Description                      |
| ------------- | ------------- | -------------------------------- |
| `id`          | integer       | Unique database ID               |
| `name`        | string        | Asset name                       |
| `image_url`   | string        | URL to asset image               |
| `identifier`  | string        | Unique tag/identifier            |
| `description` | string        | Asset description                |
| `category`    | GeneralNSN    | Category information             |
| `status`      | GeneralNSN    | Current status                   |
| `metadata`    | AssetMetadata | Additional asset information     |
| `active_tab`  | AssetCheckout | Current active checkout (if any) |

### AssetMetadata

| Field           | Type   | Description       |
| --------------- | ------ | ----------------- |
| `serial_number` | string | Serial number     |
| `manufacturer`  | string | Manufacturer name |
| `model`         | string | Model name/number |
| `notes`         | string | Additional notes  |

### AssetCheckout

| Field             | Type       | Description                                     |
| ----------------- | ---------- | ----------------------------------------------- |
| `id`              | integer    | Unique checkout ID                              |
| `asset_id`        | integer    | Associated asset ID                             |
| `checkout_status` | GeneralNSN | Checkout status                                 |
| `associated_user` | integer    | User ID who checked out                         |
| `checkout_notes`  | string     | Notes about the checkout                        |
| `time_out`        | datetime   | When the asset was checked out                  |
| `time_in`         | datetime   | When the asset was returned (null if still out) |
| `expected_in`     | datetime   | Expected return time                            |
| `user`            | User       | User information                                |

### User

| Field               | Type     | Description               |
| ------------------- | -------- | ------------------------- |
| `id`                | integer  | Unique database ID        |
| `name`              | string   | Full name                 |
| `email`             | string   | Email address             |
| `identifier`        | string   | Unique tag/identifier     |
| `profile_image_url` | string   | URL to profile photo      |
| `inserted_at`       | datetime | When the user was created |

### GeneralNSN

A general-purpose object for name/short_name pairs (used for statuses, categories, etc.)

| Field        | Type    | Description      |
| ------------ | ------- | ---------------- |
| `id`         | integer | Unique ID        |
| `name`       | string  | Full name        |
| `short_name` | string  | Abbreviated name |

---

## Error Handling

The API returns standard HTTP status codes:

| Status Code                 | Description                                            |
| --------------------------- | ------------------------------------------------------ |
| `200 OK`                    | Request succeeded                                      |
| `201 Created`               | Resource created successfully                          |
| `204 No Content`            | Resource not found (validation endpoints)              |
| `400 Bad Request`           | Invalid request (missing fields, business logic error) |
| `500 Internal Server Error` | Server-side error                                      |

Error responses are returned as plain text messages describing the issue.

---

## Development

### Prerequisites

- Go 1.19+
- MySQL database

### Environment Variables

The service requires environment variables to be set. Create a `.env` file with your configuration.

### Running Locally

```bash
# Load environment and run
source .env && go run .
```

### Code Formatting

```bash
gofmt -w -s .
```

### Run Tests

```bash
go test ./...
```

### Check for Issues

```bash
go vet ./...
```

---

## License

See [LICENSE](LICENSE) for details.
