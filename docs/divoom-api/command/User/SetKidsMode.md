# `POST /User/SetKidsMode`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `ForbidAI` | body | integer | 0 | Not verified | Recovered key |  |
| `ForbidComment` | body | integer | 0 | Not verified | Recovered key |  |
| `ForbidGroupChat` | body | integer | 0 | Not verified | Recovered key |  |
| `ForbidLogout` | body | integer | 0 | Not verified | Recovered key |  |
| `ForbidPrivateChat` | body | integer | 0 | Not verified | Recovered key |  |
| `ForbidSearch` | body | integer | 0 | Not verified | Recovered key |  |
| `KidsMode` | body | integer | 0 | Not verified | Recovered key |  |
| `KidsPassword` | body | integer | 0 | Not verified | Recovered key |  |
| `OnlyLikeSort` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
