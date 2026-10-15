# `POST /Tag/SearchTagMoreV2`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `FileSize` | body | integer | 127 | Not verified | Recovered key |  |
| `FileSort` | body | integer | 0 | Not verified | Recovered key |  |
| `KeywordsEn` | body | string | null | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `TagKey` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.TagList` | array | Recovered key |
| `$.TagList[].FollowCnt` | integer | Recovered key |
| `$.TagList[].GalleryCnt` | integer | Recovered key |
| `$.TagList[].GalleryList` | array | Recovered key |
| `$.TagList[].GalleryList[].AIFlag` | integer | Key casing unverified |
| `$.TagList[].GalleryList[].CommentCnt` | integer | Recovered key |
| `$.TagList[].GalleryList[].FileId` | string | Recovered key |
| `$.TagList[].GalleryList[].GalleryId` | integer | Recovered key |
| `$.TagList[].GalleryList[].IsLike` | integer | Recovered key |
| `$.TagList[].GalleryList[].LikeCnt` | integer | Recovered key |
| `$.TagList[].IsFollow` | integer | Recovered key |
| `$.TagList[].TagName` | string | Recovered key |
