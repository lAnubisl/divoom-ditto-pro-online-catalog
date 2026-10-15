# `POST /Cloud/GetSubscribeAndFollow`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Keywords` | body | string | null | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Version` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.CurListNum` | integer | Key casing unverified |
| `$.FollowList` | array | Key casing unverified |
| `$.FollowList[].CountryISOCode` | string | Key casing unverified |
| `$.FollowList[].FansCnt` | integer | Key casing unverified |
| `$.FollowList[].HeadId` | string | Key casing unverified |
| `$.FollowList[].IsFollow` | integer | Key casing unverified |
| `$.FollowList[].Level` | integer | Key casing unverified |
| `$.FollowList[].NickName` | string | Key casing unverified |
| `$.FollowList[].PixelAmbId` | string | Key casing unverified |
| `$.FollowList[].PixelAmbName` | string | Key casing unverified |
| `$.FollowList[].RegionId` | integer | Key casing unverified |
| `$.FollowList[].Score` | integer | Key casing unverified |
| `$.FollowList[].SelectFlag` | integer | Key casing unverified |
| `$.FollowList[].UserId` | integer | Key casing unverified |
| `$.FollowList[].isAdd` | boolean | Key casing unverified |
| `$.FollowList[].select` | boolean | Key casing unverified |
| `$.FollowListNum` | integer | Key casing unverified |
