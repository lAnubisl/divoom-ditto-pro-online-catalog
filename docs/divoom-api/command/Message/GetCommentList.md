# `POST /Message/GetCommentList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `CountryISOCode` | body | string | null | Not verified | Recovered key |  |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.LikeList` | array | Key casing unverified |
| `$.LikeList[].AckUserId` | integer | Recovered key |
| `$.LikeList[].ClockId` | integer | Recovered key |
| `$.LikeList[].Comment` | string | Recovered key |
| `$.LikeList[].CommentId` | integer | Recovered key |
| `$.LikeList[].CommentType` | integer | Recovered key |
| `$.LikeList[].CountryISOCode` | string | Recovered key |
| `$.LikeList[].Date` | integer | Recovered key |
| `$.LikeList[].ForumId` | integer | Recovered key |
| `$.LikeList[].GalleryFileId` | string | Recovered key |
| `$.LikeList[].GalleryId` | integer | Recovered key |
| `$.LikeList[].GalleryName` | string | Recovered key |
| `$.LikeList[].HeadId` | string | Recovered key |
| `$.LikeList[].IsFollow` | integer | Recovered key |
| `$.LikeList[].Level` | integer | Recovered key |
| `$.LikeList[].LikeType` | integer | Recovered key |
| `$.LikeList[].NickName` | string | Recovered key |
| `$.LikeList[].PixelAmbId` | string | Recovered key |
| `$.LikeList[].PixelAmbName` | string | Recovered key |
| `$.LikeList[].RegionId` | integer | Recovered key |
| `$.LikeList[].Relation` | integer | Recovered key |
| `$.LikeList[].UserId` | integer | Recovered key |
| `$.UnReadCnt` | integer | Key casing unverified |
