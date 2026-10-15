# `POST /Channel/Get5LcdClockList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `ClockId` | body | integer | 0 | Not verified | Recovered key |  |
| `PageIndex` | body | integer | 0 | Not verified | Recovered key |  |
| `ParentClockId` | body | integer | 0 | Not verified | Recovered key |  |
| `LcdIndependence` | body | integer | "Not established" | Not verified | Recovered key |  |
| `LcdIndex` | body | integer | "Not established" | Not verified | Recovered key |  |
| `Language` | body | string | "Not established" | Not verified | Recovered key |  |
| `ParentItemId` | body | string | "" | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.GroupList` | array | Recovered key |
| `$.GroupList[].ClockList` | array | Recovered key |
| `$.GroupList[].ClockList[].AddFlag` | integer | Recovered key |
| `$.GroupList[].ClockList[].ChannelType` | integer | Recovered key |
| `$.GroupList[].ClockList[].check` | boolean | Key casing unverified |
| `$.GroupList[].ClockList[].ClockExPlain` | string | Recovered key |
| `$.GroupList[].ClockList[].ClockExPlainPicId` | string | Recovered key |
| `$.GroupList[].ClockList[].ClockId` | integer | Recovered key |
| `$.GroupList[].ClockList[].ClockName` | string | Recovered key |
| `$.GroupList[].ClockList[].ClockType` | integer | Recovered key |
| `$.GroupList[].ClockList[].ConfigFlag` | integer | Recovered key |
| `$.GroupList[].ClockList[].ImagePixelId` | string | Recovered key |
| `$.GroupList[].ClockList[].isAdd` | boolean | Key casing unverified |
| `$.GroupList[].ClockList[].LcdIndependence` | integer | Recovered key |
| `$.GroupList[].ClockList[].PreviewImageId` | string | Recovered key |
| `$.GroupList[].ConfigFlag` | integer | Recovered key |
| `$.GroupList[].GroupId` | integer | Recovered key |
| `$.GroupList[].GroupName` | string | Recovered key |
