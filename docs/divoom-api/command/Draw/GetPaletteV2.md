# `POST /Draw/GetPaletteV2`

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
| `AndroidFlag` | body | integer | 1 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.Index` | integer | Key casing unverified |
| `$.TAG` | string | Key casing unverified |
| `$.PaletteList` | array | Key casing unverified |
| `$.PaletteList[].ColorList` | array | Recovered key |
| `$.PaletteList[].ColorList[]` | string | Element type recovered |
| `$.PaletteList[].Name` | string | Recovered key |
| `$.PaletteList[].PaletteId` | integer | Recovered key |
