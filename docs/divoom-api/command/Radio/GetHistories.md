# `POST /Radio/GetHistories`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

No operation-specific body fields were recovered. Session requirements are not verified.

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.HistoryList` | array | Key casing unverified |
| `$.HistoryList[].StationId` | string | Key casing unverified |
| `$.HistoryList[].StationLogo` | string | Key casing unverified |
| `$.HistoryList[].StationName` | string | Key casing unverified |
