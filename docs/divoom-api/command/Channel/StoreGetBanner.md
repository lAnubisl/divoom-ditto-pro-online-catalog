# `POST /Channel/StoreGetBanner`

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

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.BannerList` | array | Recovered key |
| `$.BannerList[].BannerImageId` | string | Recovered key |
| `$.BannerList[].BannerName` | string | Recovered key |
| `$.BannerList[].ClockList` | array | Recovered key |
| `$.BannerList[].ClockList[].AddFlag` | integer | Recovered key |
| `$.BannerList[].ClockList[].ChannelType` | integer | Recovered key |
| `$.BannerList[].ClockList[].check` | boolean | Key casing unverified |
| `$.BannerList[].ClockList[].ClockExPlain` | string | Recovered key |
| `$.BannerList[].ClockList[].ClockExPlainPicId` | string | Recovered key |
| `$.BannerList[].ClockList[].ClockId` | integer | Recovered key |
| `$.BannerList[].ClockList[].ClockName` | string | Recovered key |
| `$.BannerList[].ClockList[].ClockType` | integer | Recovered key |
| `$.BannerList[].ClockList[].ConfigFlag` | integer | Recovered key |
| `$.BannerList[].ClockList[].ImagePixelId` | string | Recovered key |
| `$.BannerList[].ClockList[].isAdd` | boolean | Key casing unverified |
| `$.BannerList[].ClockList[].LcdIndependence` | integer | Recovered key |
| `$.BannerList[].ClockList[].PreviewImageId` | string | Recovered key |
| `$.BannerList[].ForumId` | integer | Recovered key |
