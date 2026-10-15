# `POST /Ali/SendPhoto`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `ClockId` | body | integer | 0 | Not verified | Recovered key |  |
| `PageIndex` | body | integer | 0 | Not verified | Recovered key |  |
| `ParentClockId` | body | integer | 0 | Not verified | Recovered key |  |
| `LcdIndependence` | body | integer | "Not established" | Not verified | Recovered key |  |
| `LcdIndex` | body | integer | "Not established" | Not verified | Recovered key |  |
| `Language` | body | string | "Not established" | Not verified | Recovered key |  |
| `ParentItemId` | body | string | "" | Not verified | Recovered key |  |
| `FileSize` | body | integer | 0 | Not verified | Recovered key |  |
| `OSSName` | body | string | null | Not verified | Recovered key |  |
| `OSSPreviewName` | body | string | null | Not verified | Recovered key |  |
| `OSSServer` | body | integer | 0 | Not verified | Recovered key |  |
| `PhotoFlag` | body | integer | 0 | Not verified | Recovered key |  |
| `PhotoTitle` | body | string | null | Not verified | Recovered key |  |
| `PhotoTotalCnt` | body | integer | 0 | Not verified | Recovered key |  |
| `PhotoX` | body | integer | 0 | Not verified | Recovered key |  |
| `PhotoY` | body | integer | 0 | Not verified | Recovered key |  |
| `SendTime` | body | integer | 0 | Not verified | Recovered key |  |
| `TakingTime` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
