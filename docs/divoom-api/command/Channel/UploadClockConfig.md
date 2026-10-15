# `POST /Channel/UploadClockConfig`

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
| `Explain` | body | string | null | Not verified | Recovered key |  |
| `ItemList` | body | array | null | Not verified | Recovered key |  |
| `Name` | body | string | null | Not verified | Recovered key |  |

Nested request fields:

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.ItemList[].ControlOtherItemShow` | array | Recovered key |
| `$.ItemList[].ControlOtherItemShow[]` | string | Element type recovered |
| `$.ItemList[].Invisible` | integer | Recovered key |
| `$.ItemList[].ItemFlagByOtherItemId` | string | Recovered key |
| `$.ItemList[].ItemId` | string | Recovered key |
| `$.ItemList[].ItemName` | string | Recovered key |
| `$.ItemList[].ItemType` | integer | Recovered key |
| `$.ItemList[].ItemValue` | string | Recovered key |
| `$.ItemList[].ItemValueName` | string | Recovered key |
| `$.ItemList[].ItemValueS` | array | Recovered key |
| `$.ItemList[].ItemValueS[]` | string | Element type recovered |
| `$.ItemList[].SelectionId` | array | Recovered key |
| `$.ItemList[].SelectionId[]` | string | Element type recovered |
| `$.ItemList[].SelectionName` | array | Recovered key |
| `$.ItemList[].SelectionName[]` | string | Element type recovered |
| `$.ItemList[].ShowByOtherItemId` | string | Recovered key |
| `$.ItemList[].ShowByOtherItemValue` | string | Recovered key |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
