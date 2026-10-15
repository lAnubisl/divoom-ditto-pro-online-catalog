# `POST /MessageGroup/recallMessage`

[All endpoints](../../endpoints.md)

Host role: **messaging**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `busChannel` | body | string | null | Not verified | Key casing unverified |  |
| `messageUID` | body | string | null | Not verified | Key casing unverified |  |
| `sentTime` | body | integer | 0 | Not verified | Key casing unverified |  |
| `targetId` | body | string | null | Not verified | Key casing unverified |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
