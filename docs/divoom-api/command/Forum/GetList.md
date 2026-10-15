# `POST /Forum/GetList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `RegionId` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Tag` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.ForumList` | array | Key casing unverified |
| `$.ForumList[].CommentCnt` | integer | Key casing unverified |
| `$.ForumList[].CountryISOCode` | string | Key casing unverified |
| `$.ForumList[].Date` | integer | Key casing unverified |
| `$.ForumList[].ForumId` | integer | Key casing unverified |
| `$.ForumList[].HeadId` | string | Key casing unverified |
| `$.ForumList[].ImageId` | string | Key casing unverified |
| `$.ForumList[].IsFollow` | integer | Key casing unverified |
| `$.ForumList[].IsLike` | integer | Key casing unverified |
| `$.ForumList[].Level` | integer | Key casing unverified |
| `$.ForumList[].LikeCnt` | integer | Key casing unverified |
| `$.ForumList[].LinkUrlNoComment` | string | Key casing unverified |
| `$.ForumList[].LinkUrlWithComment` | string | Key casing unverified |
| `$.ForumList[].NickName` | string | Key casing unverified |
| `$.ForumList[].PixelAmbId` | string | Key casing unverified |
| `$.ForumList[].PixelAmbName` | string | Key casing unverified |
| `$.ForumList[].Title` | string | Key casing unverified |
| `$.ForumList[].UserId` | integer | Key casing unverified |
| `$.ForumList[].WatchCnt` | integer | Key casing unverified |
| `$.HeaderNum` | integer | Key casing unverified |
| `$.TotalNum` | integer | Key casing unverified |
