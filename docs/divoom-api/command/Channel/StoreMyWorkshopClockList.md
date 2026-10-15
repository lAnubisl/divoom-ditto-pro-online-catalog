# `POST /Channel/StoreMyWorkshopClockList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | "Not established" | Not verified | Recovered key |  |
| `CountryISOCode` | body | string | "Not established" | Not verified | Recovered key |  |
| `Flag` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.ClockList` | array | Recovered key |
| `$.ClockList[].AddFlag` | integer | Recovered key |
| `$.ClockList[].ChannelType` | integer | Recovered key |
| `$.ClockList[].check` | boolean | Key casing unverified |
| `$.ClockList[].ClockExPlain` | string | Recovered key |
| `$.ClockList[].ClockExPlainPicId` | string | Recovered key |
| `$.ClockList[].ClockId` | integer | Recovered key |
| `$.ClockList[].ClockName` | string | Recovered key |
| `$.ClockList[].ClockType` | integer | Recovered key |
| `$.ClockList[].ConfigFlag` | integer | Recovered key |
| `$.ClockList[].ImagePixelId` | string | Recovered key |
| `$.ClockList[].isAdd` | boolean | Key casing unverified |
| `$.ClockList[].LcdIndependence` | integer | Recovered key |
| `$.ClockList[].PreviewImageId` | string | Recovered key |
