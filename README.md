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
   System architecture
</h1>

![SystemDesign](assets/arch.png)

## API Reference

#### Register a new API user

Allows an admin to create a new API user and receive an API token

```http
POST /api/admin/register
```


**Header Parameters:**

| Parameter       | Type     | Description               |
| --------------- | -------- | ------------------------- |
| `X-Admin-Token` | `string` | **Required**. Admin token |

**Responses:**

| HTTP Code | Description                       |
| --------- | --------------------------------- |
| 200       | API token for the new user        |
| 400       | Bad Request                       |
| 401       | Unauthorized: Invalid admin token |
<br>

#### Unregister a user

Allows an admin to delete a user by API token

```http
DELETE /api/admin/unregister/{ApiToken}
```


**Header Parameters:**

| Parameter       | Type     | Description               |
| --------------- | -------- | ------------------------- |
| `X-Admin-Token` | `string` | **Required**. Admin token |

**Responses:**

| HTTP Code | Description                       |
| --------- | --------------------------------- |
| 200       | Deleted {ApiToken}                |
| 400       | Bad Request                       |
| 401       | Unauthorized: Invalid admin token |
<br>

#### Send a message

Sends a message through the API using the user's token

```http
POST /api/message
```


**Header Parameters:**

| Parameter     | Type     | Description                  |
| ------------- | -------- | ---------------------------- |
| `X-Api-Token` | `string` | **Required**. User API token |

**Body Parameters:**

| Parameter     | Type            | Description |
| ------------- | --------------- | ----------- |
| `device_id`   | `string`        |             |
| `message`     | `string`        |             |
| `subscribers` | `array[string]` |             |

**Responses:**

| HTTP Code | Description         |
| --------- | ------------------- |
| 200       | Message sent        |
| 400       | Bad Request         |
| 401       | Unauthorized        |
| 429       | Rate limit exceeded |
<br>

#### Ping the server

Simple health check endpoint

```http
GET /api/ping
```


**Responses:**

| HTTP Code | Description |
| --------- | ----------- |
| 200       | pong        |
<br>
