# `POST /Channel/GetClockCommentList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `ClockId` | body | integer | 0 | Not verified | Recovered key |  |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `MessageId` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Type` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.CommentList` | array | Key casing unverified |
| `$.CommentList[].AtList` | array | Recovered key |
| `$.CommentList[].AtList[].AtNickName` | string | Recovered key |
| `$.CommentList[].AtList[].AtUserId` | string | Recovered key |
| `$.CommentList[].AuthorLike` | integer | Recovered key |
| `$.CommentList[].Comment` | string | Recovered key |
| `$.CommentList[].CommentChildList` | array | Recovered key |
| `$.CommentList[].CommentChildList[].AckUserId` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].AckUserNickName` | string | Recovered key |
| `$.CommentList[].CommentChildList[].AtList` | array | Recovered key |
| `$.CommentList[].CommentChildList[].AtList[].AtNickName` | string | Recovered key |
| `$.CommentList[].CommentChildList[].AtList[].AtUserId` | string | Recovered key |
| `$.CommentList[].CommentChildList[].AuthorLike` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].Comment` | string | Recovered key |
| `$.CommentList[].CommentChildList[].CommentId` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].CountryISOCode` | string | Recovered key |
| `$.CommentList[].CommentChildList[].Date` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].HeadId` | string | Recovered key |
| `$.CommentList[].CommentChildList[].index` | integer | Key casing unverified |
| `$.CommentList[].CommentChildList[].IsLike` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].Level` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].LikeCnt` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].NickName` | string | Recovered key |
| `$.CommentList[].CommentChildList[].PixelAmbId` | string | Recovered key |
| `$.CommentList[].CommentChildList[].PixelAmbName` | string | Recovered key |
| `$.CommentList[].CommentChildList[].RegionId` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].Relation` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].showMore` | boolean | Key casing unverified |
| `$.CommentList[].CommentChildList[].UserId` | integer | Recovered key |
| `$.CommentList[].CommentChildList[].showNum` | integer | Key casing unverified |
| `$.CommentList[].CommentChildList[].allNum` | integer | Key casing unverified |
| `$.CommentList[].CommentChildList[].showLine` | boolean | Key casing unverified |
| `$.CommentList[].CommentId` | integer | Recovered key |
| `$.CommentList[].CountryISOCode` | string | Recovered key |
| `$.CommentList[].Date` | integer | Recovered key |
| `$.CommentList[].HeadId` | string | Recovered key |
| `$.CommentList[].index` | integer | Key casing unverified |
| `$.CommentList[].IsLike` | integer | Recovered key |
| `$.CommentList[].Level` | integer | Recovered key |
| `$.CommentList[].LikeCnt` | integer | Recovered key |
| `$.CommentList[].NickName` | string | Recovered key |
| `$.CommentList[].PixelAmbId` | string | Recovered key |
| `$.CommentList[].PixelAmbName` | string | Recovered key |
| `$.CommentList[].RegionId` | integer | Recovered key |
| `$.CommentList[].Relation` | integer | Recovered key |
| `$.CommentList[].showLine` | boolean | Key casing unverified |
| `$.CommentList[].showMore` | boolean | Key casing unverified |
| `$.CommentList[].showNum` | integer | Key casing unverified |
| `$.CommentList[].UserId` | integer | Recovered key |
| `$.CommentListNum` | integer | Key casing unverified |
| `$.CurListNum` | integer | Key casing unverified |
| `$.TotalCommentCnt` | integer | Recovered key |
