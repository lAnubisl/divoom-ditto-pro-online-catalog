# `POST /Mall/Buy`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Address1` | body | string | null | Not verified | Recovered key |  |
| `Address2` | body | string | null | Not verified | Recovered key |  |
| `Country` | body | string | null | Not verified | Recovered key |  |
| `Email` | body | string | null | Not verified | Recovered key |  |
| `MallId` | body | integer | 0 | Not verified | Recovered key |  |
| `Name` | body | string | null | Not verified | Recovered key |  |
| `Phone` | body | string | null | Not verified | Recovered key |  |
| `Province` | body | string | null | Not verified | Recovered key |  |
| `Remark` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
