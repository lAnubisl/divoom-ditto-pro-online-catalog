# `POST /ThirdLogin`

[All endpoints](../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Code` | body | string | null | Not verified | Recovered key |  |
| `IsAndroid` | body | integer | 0 | Not verified | Recovered key |  |
| `Type` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.Token` | integer | Key casing unverified |
| `$.UserId` | integer | Key casing unverified |
| `$.UserToken` | string | Key casing unverified |
