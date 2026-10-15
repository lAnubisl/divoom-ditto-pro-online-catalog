# `POST /WhiteNoise/Get`

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
| `$.EndStatus` | integer | Key casing unverified |
| `$.OnOff` | integer | Key casing unverified |
| `$.Time` | integer | Key casing unverified |
| `$.Volume` | array | Key casing unverified |
| `$.Volume[]` | integer | Element type recovered |
