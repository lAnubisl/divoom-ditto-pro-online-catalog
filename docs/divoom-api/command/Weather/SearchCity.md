# `POST /Weather/SearchCity`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `CityName` | body | string | null | Not verified | Recovered key |  |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.CityList` | array | Key casing unverified |
| `$.CityList[].CityId` | integer | Key casing unverified |
| `$.CityList[].CityName` | string | Key casing unverified |
| `$.CityList[].Country` | string | Key casing unverified |
| `$.CityList[].Lat` | number | Key casing unverified |
| `$.CityList[].Lon` | number | Key casing unverified |
