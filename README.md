<h1 align="center" style="display: flex; align-items: center; justify-content: center;">
   Whisper 
</h1>

<div align="center">
    <img alt="GitHub License" src="https://img.shields.io/github/license/IdanKoblik/whisper">
</div>

<br>

<p align="center">
   This project transforms an Android device into a fully functional SMS server.
   It provides a secure and scalable architecture for sending and managing SMS messages through a REST API and WebSocket communication.
</p>

<h1 align="center" style="display: flex; align-items: center; justify-content: center;">
   System Architecture
</h1>

![SystemDesign](assets/arch.png)

---

## Setup REST API

You can install and run the REST API in multiple ways:

**Using Docker:**

```bash
docker pull ghcr.io/idankoblik/whisper:<version>
docker run -e CONFIG_PATH=/path/to/config.yaml -e ghcr.io/idankoblik/whisper:<version>
```

**Using precompiled Linux executable:**

1. Download the binary from GitHub Releases.
2. Make it executable and run:

```bash
chmod +x whisper
CONFIG_PATH=/path/to/config.yaml ./whisper
```

**Building from source:**

1. Clone the repository and navigate to the `api` folder:

```bash
cd api
go mod tidy
make build
```

2. Run the executable with environment variables:

```bash
CONFIG_PATH=/path/to/config.yaml ./whisper
```

**Environment variables:**

* `CONFIG_PATH` → path to the configuration file

**Configuration file structure (`config.env.yaml` example):**

```yaml
addr: ""
admin_token: ""
rate_limit: 60 # optional, per minute

mongo:
  connection_string: ""
  database: ""
  collection: ""

redis:
  addr: ""
  password: ""
```

---

## API Reference

#### Register new API token

Create a new API token for authentication. Requires admin privileges. Returns the generated API token.

```http
POST /admin/register
```


**Responses:**

| HTTP Code | Description                          |
| --------- | ------------------------------------ |
| 201       | API token created successfully       |
| 400       | Failed to create token               |
| 401       | Unauthorized - Admin access required |
<br>

#### Unregister API token

Remove an API token from the database. Requires admin privileges.

```http
DELETE /admin/unregister/{token}
```


**Responses:**

| HTTP Code | Description                          |
| --------- | ------------------------------------ |
| 200       | API token removed successfully       |
| 400       | Failed to remove token               |
| 401       | Unauthorized - Admin access required |
<br>

#### Remove device

Remove a device from the authenticated user's device list

```http
DELETE /api/devices
```


**Body Parameters:**

| Parameter | Type     | Description |
| --------- | -------- | ----------- |
| `device`  | `string` |             |

**Responses:**

| HTTP Code | Description                                |
| --------- | ------------------------------------------ |
| 200       | Device removed successfully                |
| 400       | Invalid request or failed to remove device |
| 401       | Unauthorized - Invalid or missing token    |
<br>

#### Add device

Add a new device to the authenticated user's device list

```http
POST /api/devices
```


**Body Parameters:**

| Parameter | Type     | Description |
| --------- | -------- | ----------- |
| `device`  | `string` |             |

**Responses:**

| HTTP Code | Description                             |
| --------- | --------------------------------------- |
| 201       | Device added successfully               |
| 400       | Invalid request or failed to add device |
| 401       | Unauthorized - Invalid or missing token |
<br>

#### Get device

Check if a device exists and belongs to the authenticated user

```http
GET /api/devices/{id}
```


**Responses:**

| HTTP Code | Description                             |
| --------- | --------------------------------------- |
| 200       | Device found                            |
| 400       | Invalid request or validation error     |
| 401       | Unauthorized - Invalid or missing token |
| 404       | Device not found                        |
<br>

#### Send message to device

Send a message to a specific device via WebSocket. The device must be active and connected. Rate limiting may apply.

```http
POST /api/send
```


**Body Parameters:**

| Parameter | Type            | Description |
| --------- | --------------- | ----------- |
| `device`  | `string`        |             |
| `message` | `string`        |             |
| `targets` | `array[string]` |             |

**Responses:**

| HTTP Code | Description                                              |
| --------- | -------------------------------------------------------- |
| 200       | Message sent successfully                                |
| 400       | Invalid request, device not active, or invalid device id |
| 401       | Unauthorized - Invalid or missing token                  |
| 429       | Rate limit exceeded                                      |
<br>

#### Health check endpoint

Check the health status of the server and its dependencies (MongoDB and Redis)

```http
GET /health
```


**Responses:**

| HTTP Code | Description                         |
| --------- | ----------------------------------- |
| 200       | Server and dependencies are healthy |
| 503       | One or more dependencies are down   |
<br>

#### WebSocket connection endpoint

Establishes a WebSocket connection for real-time messaging. Requires authentication token. Client must send device information in JSON format.

```http
GET /ws/
```


**Body Parameters:**

| Parameter | Type     | Description |
| --------- | -------- | ----------- |
| `device`  | `string` |             |

**Responses:**

| HTTP Code | Description                                 |
| --------- | ------------------------------------------- |
| 101       | WebSocket connection upgraded successfully  |
| 400       | Invalid payload or device validation failed |
| 401       | Unauthorized - Invalid or missing token     |
<br>
