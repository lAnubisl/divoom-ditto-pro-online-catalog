# `POST /AidSleep/Add`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `AudioType` | body | integer | 0 | Not verified | Recovered key |  |
| `FileId` | body | string | null | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `Name` | body | string | null | Not verified | Recovered key |  |
| `SleepId` | body | integer | 0 | Not verified | Recovered key |  |
| `Type` | body | integer | 0 | Not verified | Recovered key |  |
| `VideoType` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
