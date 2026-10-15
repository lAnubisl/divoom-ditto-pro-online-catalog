# `POST /User/GetKidsMode`

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
| `$.DeviceId` | integer | Recovered key |
| `$.ForbidAI` | integer | Recovered key |
| `$.ForbidComment` | integer | Recovered key |
| `$.ForbidGroupChat` | integer | Recovered key |
| `$.ForbidLogout` | integer | Recovered key |
| `$.ForbidPrivateChat` | integer | Recovered key |
| `$.ForbidSearch` | integer | Recovered key |
| `$.KidsMode` | integer | Recovered key |
| `$.KidsPassword` | integer | Recovered key |
| `$.OnlyLikeSort` | integer | Recovered key |
