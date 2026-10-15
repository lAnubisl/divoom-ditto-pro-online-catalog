# `POST /PhotoFrame/GetList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Language` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.PhotoList` | array | Key casing unverified |
| `$.PhotoList[].BigImageId` | string | Key casing unverified |
| `$.PhotoList[].PhotoName` | string | Key casing unverified |
| `$.PhotoList[].PixelStartX` | integer | Key casing unverified |
| `$.PhotoList[].PixelStartY` | integer | Key casing unverified |
| `$.PhotoList[].SmallImageId` | string | Key casing unverified |
