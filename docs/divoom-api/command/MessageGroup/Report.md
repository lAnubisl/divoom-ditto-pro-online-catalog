# `POST /MessageGroup/Report`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `SentTime` | body | integer | 0 | Not verified | Recovered key |  |
| `MessageImageUrl` | body | string | "" | Not verified | Recovered key |  |
| `MessagePixelFileId` | body | string | "" | Not verified | Recovered key |  |
| `ReportText` | body | string | "" | Not verified | Recovered key |  |
| `MessageText` | body | string | "" | Not verified | Recovered key |  |
| `MessageVideoUrl` | body | string | "" | Not verified | Recovered key |  |
| `TargetUserId` | body | integer | 0 | Not verified | Recovered key |  |
| `MessageUID` | body | string | "" | Not verified | Recovered key |  |
| `TargetId` | body | string | "" | Not verified | Recovered key |  |
| `BusChannel` | body | string | "" | Not verified | Recovered key |  |
| `ReportType` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
