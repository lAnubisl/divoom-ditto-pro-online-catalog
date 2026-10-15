# `POST /GetExpertListV4`

[All endpoints](../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `FileType` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `RefreshIndex` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `FileSize` | body | integer | 127 | Not verified | Recovered key |  |
| `Version` | body | integer | 19 | Not verified | Recovered key |  |

Session fields are described in [the guide](../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.CurListNum` | integer | Key casing unverified |
| `$.ExpertList` | array | Key casing unverified |
| `$.ExpertList[].CountryISOCode` | string | Key casing unverified |
| `$.ExpertList[].FansCnt` | integer | Key casing unverified |
| `$.ExpertList[].FileList` | array | Key casing unverified |
| `$.ExpertList[].FileList[].AIFlag` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].AtList` | array | Key casing unverified |
| `$.ExpertList[].FileList[].AtList[].AtNickName` | string | Recovered key |
| `$.ExpertList[].FileList[].AtList[].AtUserId` | string | Recovered key |
| `$.ExpertList[].FileList[].CheckConfirm` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].Classify` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].CommentCnt` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].CommentUTC` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].CopyrightFlag` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].CountryISOCode` | string | Key casing unverified |
| `$.ExpertList[].FileList[].Date` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].FileId` | string | Key casing unverified |
| `$.ExpertList[].FileList[].FileName` | string | Key casing unverified |
| `$.ExpertList[].FileList[].FileSize` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].FileType` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].FillGameIsFinish` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].FillGameScore` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].GalleryId` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].GoodLevel` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].HideFlag` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].IsAddNew` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].IsAddRecommend` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].IsDel` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].IsFollow` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].IsLike` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].IsNewComment` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].LayerFileId` | string | Key casing unverified |
| `$.ExpertList[].FileList[].Level` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].LikeCnt` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].LikeUTC` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].MusicFileId` | string | Key casing unverified |
| `$.ExpertList[].FileList[].OriginalGalleryId` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].PixelAmbId` | string | Key casing unverified |
| `$.ExpertList[].FileList[].PixelAmbName` | string | Key casing unverified |
| `$.ExpertList[].FileList[].PrivateFlag` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].RegionId` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].ReportInfo` | string | Key casing unverified |
| `$.ExpertList[].FileList[].ReportUser` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].UserId` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].WatchCnt` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].content` | string | Key casing unverified |
| `$.ExpertList[].FileList[].fileTagArray` | array | Key casing unverified |
| `$.ExpertList[].FileList[].fileTagArray[]` | string | Element type recovered |
| `$.ExpertList[].FileList[].shareCnt` | integer | Key casing unverified |
| `$.ExpertList[].FileList[].userHeaderId` | string | Key casing unverified |
| `$.ExpertList[].FileList[].userName` | string | Key casing unverified |
| `$.ExpertList[].HeadId` | string | Key casing unverified |
| `$.ExpertList[].IsFollow` | integer | Key casing unverified |
| `$.ExpertList[].Level` | integer | Key casing unverified |
| `$.ExpertList[].MedalList` | array | Key casing unverified |
| `$.ExpertList[].MedalList[].ImageId` | string | Key casing unverified |
| `$.ExpertList[].NickName` | string | Key casing unverified |
| `$.ExpertList[].PixelAmbId` | string | Key casing unverified |
| `$.ExpertList[].PixelAmbName` | string | Key casing unverified |
| `$.ExpertList[].RegionId` | integer | Key casing unverified |
| `$.ExpertList[].Score` | integer | Key casing unverified |
| `$.ExpertList[].UserId` | integer | Key casing unverified |
| `$.ExpertListNum` | integer | Key casing unverified |
