# `POST /WhiteNoise/Set`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndStatus` | body | integer | 0 | Not verified | Recovered key |  |
| `OnOff` | body | integer | 0 | Not verified | Recovered key |  |
| `Time` | body | integer | 0 | Not verified | Recovered key |  |
| `Volume` | body | array | [0, 0, 0, 0, 0, 0, 0, 0] | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
