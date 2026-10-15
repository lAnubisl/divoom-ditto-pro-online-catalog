# `POST /Forum/GetForumUrl`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `CountryISOCode` | body | string | null | Not verified | Recovered key |  |
| `ForumId` | body | integer | 0 | Not verified | Recovered key |  |
| `ForumType` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.CommentCnt` | integer | Key casing unverified |
| `$.CountryISOCode` | string | Key casing unverified |
| `$.Date` | integer | Key casing unverified |
| `$.ForumId` | integer | Key casing unverified |
| `$.HeadId` | string | Key casing unverified |
| `$.ImageId` | string | Key casing unverified |
| `$.IsFollow` | integer | Key casing unverified |
| `$.IsLike` | integer | Key casing unverified |
| `$.Level` | integer | Key casing unverified |
| `$.LikeCnt` | integer | Key casing unverified |
| `$.LinkUrlNoComment` | string | Key casing unverified |
| `$.LinkUrlWithComment` | string | Key casing unverified |
| `$.NickName` | string | Key casing unverified |
| `$.PixelAmbId` | string | Key casing unverified |
| `$.PixelAmbName` | string | Key casing unverified |
| `$.Title` | string | Key casing unverified |
| `$.UserId` | integer | Key casing unverified |
| `$.WatchCnt` | integer | Key casing unverified |
