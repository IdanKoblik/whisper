<h1 align="center" style="display: flex; align-items: center; justify-content: center;">
   Whisper 
</h1>

<div align="center">
    <img alt="GitHub License" src="https://img.shields.io/github/license/IdanKoblik/whisper">
    <img alt="Sonar" src="https://sonarcloud.io/api/project_badges/measure?project=IdanKoblik_whisper&metric=security_rating&token=80a1e7fdd01c76f58138be77745de3116285aa66">
</div>

<br>

<p align="center">
   This project transforms an Android device into a fully functional SMS server.
   It provides a secure and scalable architecture for sending and managing SMS messages through a REST API and WebSocket communication.
</p>

<h1 align="center" style="display: flex; align-items: center; justify-content: center;">
   System architecture
</h1>

![SystemDesign](https://raw.githubusercontent.com/IdanKoblik/assets/refs/heads/main/whisper.png)

## API Reference

#### Ping the server

Returns "pong"

```http
GET /ping
```


**Responses:**

| HTTP Code | Description |
| --------- | ----------- |
| 200       | pong        |
<br>

#### Register a new user

Requires X-Admin-Token header and RawUser JSON body

```http
POST /register
```


**Header Parameters:**

| Parameter       | Type     | Description               |
| --------------- | -------- | ------------------------- |
| `X-Admin-Token` | `string` | **Required**. Admin Token |

**Body Parameters:**

| Parameter     | Type            | Description |
| ------------- | --------------- | ----------- |
| `owner`       | `string`        |             |
| `subject`     | `string`        |             |
| `subscribers` | `array[string]` |             |

**Responses:**

| HTTP Code | Description   |
| --------- | ------------- |
| 200       | JWT Token     |
| 400       | Invalid input |
| 401       | Unauthorized  |
<br>
