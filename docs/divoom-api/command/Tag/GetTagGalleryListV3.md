# `POST /Tag/GetTagGalleryListV3`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Classify` | body | integer | 0 | Not verified | Recovered key |  |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `FileSort` | body | integer | 0 | Not verified | Recovered key |  |
| `RefreshIndex` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `FileType` | body | integer | 5 | Not verified | Recovered key |  |
| `FileSize` | body | integer | 127 | Not verified | Recovered key |  |
| `Version` | body | integer | 19 | Not verified | Recovered key |  |
| `Mode` | body | integer | 0 | Not verified | Recovered key |  |
| `TagName` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.FileList` | array | Key casing unverified |
| `$.FileList[].AIFlag` | integer | Key casing unverified |
| `$.FileList[].AtList` | array | Key casing unverified |
| `$.FileList[].AtList[].AtNickName` | string | Recovered key |
| `$.FileList[].AtList[].AtUserId` | string | Recovered key |
| `$.FileList[].CheckConfirm` | integer | Key casing unverified |
| `$.FileList[].Classify` | integer | Key casing unverified |
| `$.FileList[].CommentCnt` | integer | Key casing unverified |
| `$.FileList[].CommentUTC` | integer | Key casing unverified |
| `$.FileList[].CopyrightFlag` | integer | Key casing unverified |
| `$.FileList[].CountryISOCode` | string | Key casing unverified |
| `$.FileList[].Date` | integer | Key casing unverified |
| `$.FileList[].FileId` | string | Key casing unverified |
| `$.FileList[].FileName` | string | Key casing unverified |
| `$.FileList[].FileSize` | integer | Key casing unverified |
| `$.FileList[].FileType` | integer | Key casing unverified |
| `$.FileList[].FillGameIsFinish` | integer | Key casing unverified |
| `$.FileList[].FillGameScore` | integer | Key casing unverified |
| `$.FileList[].GalleryId` | integer | Key casing unverified |
| `$.FileList[].GoodLevel` | integer | Key casing unverified |
| `$.FileList[].HideFlag` | integer | Key casing unverified |
| `$.FileList[].IsAddNew` | integer | Key casing unverified |
| `$.FileList[].IsAddRecommend` | integer | Key casing unverified |
| `$.FileList[].IsDel` | integer | Key casing unverified |
| `$.FileList[].IsFollow` | integer | Key casing unverified |
| `$.FileList[].IsLike` | integer | Key casing unverified |
| `$.FileList[].IsNewComment` | integer | Key casing unverified |
| `$.FileList[].LayerFileId` | string | Key casing unverified |
| `$.FileList[].Level` | integer | Key casing unverified |
| `$.FileList[].LikeCnt` | integer | Key casing unverified |
| `$.FileList[].LikeUTC` | integer | Key casing unverified |
| `$.FileList[].MusicFileId` | string | Key casing unverified |
| `$.FileList[].OriginalGalleryId` | integer | Key casing unverified |
| `$.FileList[].PixelAmbId` | string | Key casing unverified |
| `$.FileList[].PixelAmbName` | string | Key casing unverified |
| `$.FileList[].PrivateFlag` | integer | Key casing unverified |
| `$.FileList[].RegionId` | integer | Key casing unverified |
| `$.FileList[].ReportInfo` | string | Key casing unverified |
| `$.FileList[].ReportUser` | integer | Key casing unverified |
| `$.FileList[].UserId` | integer | Key casing unverified |
| `$.FileList[].WatchCnt` | integer | Key casing unverified |
| `$.FileList[].content` | string | Key casing unverified |
| `$.FileList[].fileTagArray` | array | Key casing unverified |
| `$.FileList[].fileTagArray[]` | string | Element type recovered |
| `$.FileList[].shareCnt` | integer | Key casing unverified |
| `$.FileList[].userHeaderId` | string | Key casing unverified |
| `$.FileList[].userName` | string | Key casing unverified |
