# `POST /Lottery/GetPrizeInfo`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `CountryISOCode` | body | string | null | Not verified | Recovered key |  |
| `Langue` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.PrizeList` | array | Key casing unverified |
| `$.PrizeList[].PrizeName` | string | Key casing unverified |
| `$.PrizeList[].PrizeSmallImageId` | string | Key casing unverified |
| `$.Title1` | string | Key casing unverified |
| `$.Title2` | string | Key casing unverified |
| `$.Title3` | string | Key casing unverified |
