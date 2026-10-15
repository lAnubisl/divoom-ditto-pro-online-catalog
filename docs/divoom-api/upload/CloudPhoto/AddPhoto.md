# `POST /CloudPhoto/AddPhoto`

[All endpoints](../../endpoints.md)

Host role: **upload**.

Request Content-Type: `multipart/form-data`.

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
| `OSSName` | body | string | null | Not verified | Recovered key |  |
| `OSSServer` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).
For multipart requests, the `json` form part holds these fields. Binary part names are endpoint-specific and may remain unverified.

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.FileId` | string | Key casing unverified |
