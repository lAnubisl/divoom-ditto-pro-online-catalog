# `POST /GetSomeoneInfoV2`

[All endpoints](../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `SomeOneUserId` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.BackgroundId` | string | Key casing unverified |
| `$.BlackFlag` | integer | Key casing unverified |
| `$.CountryISOCode` | string | Key casing unverified |
| `$.FansCnt` | integer | Key casing unverified |
| `$.FollowCnt` | integer | Key casing unverified |
| `$.HeadId` | string | Key casing unverified |
| `$.IsFollow` | integer | Key casing unverified |
| `$.Level` | integer | Key casing unverified |
| `$.LikeCnt` | integer | Key casing unverified |
| `$.MedalList` | array | Key casing unverified |
| `$.MedalList[].ImageId` | string | Key casing unverified |
| `$.MessageFlag` | integer | Key casing unverified |
| `$.NickName` | string | Key casing unverified |
| `$.PixelAmbId` | string | Key casing unverified |
| `$.PixelAmbName` | string | Key casing unverified |
| `$.PixelAmbTypeId` | integer | Key casing unverified |
| `$.PixelFlag` | integer | Key casing unverified |
| `$.RegionId` | integer | Key casing unverified |
| `$.Relation` | integer | Key casing unverified |
| `$.Score` | integer | Key casing unverified |
| `$.UserNewSign` | string | Key casing unverified |
| `$.WebUrl` | string | Key casing unverified |
