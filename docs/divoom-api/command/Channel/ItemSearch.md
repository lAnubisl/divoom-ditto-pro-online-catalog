# `POST /Channel/ItemSearch`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `ClockId` | body | integer | 0 | Not verified | Recovered key |  |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `ItemFlag` | body | string | null | Not verified | Recovered key |  |
| `ItemId` | body | string | null | Not verified | Recovered key |  |
| `Key` | body | string | null | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.SearchList` | array | Key casing unverified |
| `$.SearchList[].ItemValue` | string | Recovered key |
| `$.SearchList[].ItemValueName` | string | Recovered key |
