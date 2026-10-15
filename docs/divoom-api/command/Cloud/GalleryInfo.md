# `POST /Cloud/GalleryInfo`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `GalleryId` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.AIFlag` | integer | Key casing unverified |
| `$.AtList` | array | Key casing unverified |
| `$.AtList[].AtNickName` | string | Recovered key |
| `$.AtList[].AtUserId` | string | Recovered key |
| `$.CheckConfirm` | integer | Key casing unverified |
| `$.CommentCnt` | integer | Key casing unverified |
| `$.CommentUTC` | integer | Key casing unverified |
| `$.Content` | string | Key casing unverified |
| `$.CopyrightFlag` | integer | Key casing unverified |
| `$.CountryISOCode` | string | Key casing unverified |
| `$.CurListNum` | integer | Key casing unverified |
| `$.Date` | integer | Key casing unverified |
| `$.FileId` | string | Key casing unverified |
| `$.FileName` | string | Key casing unverified |
| `$.FileType` | integer | Key casing unverified |
| `$.HideFlag` | integer | Key casing unverified |
| `$.IsDel` | integer | Key casing unverified |
| `$.IsFollow` | integer | Key casing unverified |
| `$.IsLike` | integer | Key casing unverified |
| `$.IsNewComment` | integer | Key casing unverified |
| `$.LayerFileId` | string | Key casing unverified |
| `$.Level` | integer | Key casing unverified |
| `$.LikeCnt` | integer | Key casing unverified |
| `$.LikeUTC` | integer | Key casing unverified |
| `$.MusicFileId` | string | Key casing unverified |
| `$.OriginalFileId` | string | Key casing unverified |
| `$.OriginalGalleryId` | integer | Key casing unverified |
| `$.PixelAmbId` | string | Key casing unverified |
| `$.PixelAmbName` | string | Key casing unverified |
| `$.PrivateFlag` | integer | Key casing unverified |
| `$.RegionId` | integer | Key casing unverified |
| `$.ShareCnt` | integer | Key casing unverified |
| `$.UserHeaderId` | string | Key casing unverified |
| `$.UserId` | integer | Key casing unverified |
| `$.UserName` | string | Key casing unverified |
| `$.WatchCnt` | integer | Key casing unverified |
| `$.fileTagArray` | array | Key casing unverified |
| `$.fileTagArray[]` | string | Element type recovered |
