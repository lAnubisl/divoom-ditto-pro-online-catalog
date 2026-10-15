# `POST /Discover/GetAlbumImageListV3`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Verified on 14 October 2026.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Classify` | body | integer | 0 | Optional in catalog omission probe | Recovered key | Category identifier; also accepted by album-content requests. |
| `EndNum` | body | integer | 0 | Required in catalog omission probe (ReturnCode=1) | Recovered key | Inclusive end position. For a 20-position page use StartNum=(page-1)*20+1 and EndNum=page*20. |
| `FileSort` | body | integer | 0 | Optional in catalog omission probe | Recovered key | Sort selector; complete value semantics are not established. |
| `RefreshIndex` | body | integer | 0 | Optional in catalog omission probe | Recovered key | Refresh state, normally 0; broader semantics are unverified. |
| `StartNum` | body | integer | 0 | Required in catalog omission probe (ReturnCode=1) | Recovered key | Start position; use 1 for the first page. |
| `FileType` | body | integer | 5 | Optional in catalog omission probe | Recovered key | Picture 0, animation 1, multiple pictures 2, multiple animations 3, LED 4, all 5, sand 6. |
| `FileSize` | body | integer | 127 | Optional in catalog omission probe | Recovered key | Size selector: 1 selects 16×16, 2 selects 32×32, 4 selects 64×64; 127 is the broad selector, not a pixel width. |
| `Version` | body | integer | 19 | Optional in catalog omission probe | Recovered key | Gallery request version; use 19 to reproduce the recovered requests. |
| `AlbumId` | body | integer | 0 | Required in catalog omission probe (ReturnCode=1) | Recovered key | Album identifier returned in AlbumList. |

Session fields are described in [the guide](../../README.md).

## Request example

```http
POST https://appin.divoom-gz.com/Discover/GetAlbumImageListV3
Content-Type: application/json
```

```json
{
  "AlbumId": 34,
  "StartNum": 1,
  "EndNum": 20,
  "Classify": 0,
  "FileSize": 1,
  "FileType": 1,
  "FileSort": 0,
  "Version": 19,
  "RefreshIndex": 0
}
```

## Response structure

Observed catalog sample.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$` | object | Observed response |
| `$.CurListNum` | number | Observed response |
| `$.FileList` | array | Observed response |
| `$.FileList[]` | object | Observed response |
| `$.FileList[].AIFlag` | number | Observed response |
| `$.FileList[].AtList` | array | Observed response |
| `$.FileList[].CheckConfirm` | number | Observed response |
| `$.FileList[].Classify` | number | Observed response |
| `$.FileList[].CommentCnt` | number | Observed response |
| `$.FileList[].CommentUTC` | number | Observed response |
| `$.FileList[].Content` | string | Observed response |
| `$.FileList[].CopyrightFlag` | number | Observed response |
| `$.FileList[].CountryISOCode` | string | Observed response |
| `$.FileList[].Date` | number | Observed response |
| `$.FileList[].FileId` | string | Observed response |
| `$.FileList[].FileName` | string | Observed response |
| `$.FileList[].FileSize` | number | Observed response |
| `$.FileList[].FileTagArray` | array | Observed response |
| `$.FileList[].FileTagArray[]` | string | Observed response |
| `$.FileList[].FileType` | number | Observed response |
| `$.FileList[].FileURL` | string | Observed response |
| `$.FileList[].FillGameIsFinish` | number | Observed response |
| `$.FileList[].FillGameScore` | number | Observed response |
| `$.FileList[].GalleryId` | number | Observed response |
| `$.FileList[].HideFlag` | number | Observed response |
| `$.FileList[].IsAddNew` | number | Observed response |
| `$.FileList[].IsAddRecommend` | number | Observed response |
| `$.FileList[].IsDel` | number | Observed response |
| `$.FileList[].IsFollow` | number | Observed response |
| `$.FileList[].IsLike` | number | Observed response |
| `$.FileList[].LayerFileId` | string | Observed response |
| `$.FileList[].Level` | number | Observed response |
| `$.FileList[].LikeCnt` | number | Observed response |
| `$.FileList[].LikeUTC` | number | Observed response |
| `$.FileList[].MusicFileId` | string | Observed response |
| `$.FileList[].OriginalGalleryId` | number | Observed response |
| `$.FileList[].PixelAmbId` | string | Observed response |
| `$.FileList[].PixelAmbName` | string | Observed response |
| `$.FileList[].PrivateFlag` | number | Observed response |
| `$.FileList[].RegionId` | string | Observed response |
| `$.FileList[].ShareCnt` | number | Observed response |
| `$.FileList[].UserHeaderId` | string | Observed response |
| `$.FileList[].UserId` | number | Observed response |
| `$.FileList[].UserName` | string | Observed response |
| `$.FileList[].WatchCnt` | number | Observed response |
| `$.FileListNum` | number | Observed response |
| `$.ReturnCode` | number | Observed response |
| `$.ReturnMessage` | string | Observed response |
| `$.expire` | number | Observed response |

## Behavior and interpretation

Use increasing inclusive position ranges. The live probes checked only the first two positions and do not establish an end-of-list rule. A successful empty collection differs from an API error. A short page alone does not establish completion.
Error responses omit the successful collection; check ReturnCode before reading it.
Observed file-list responses use capitalized keys such as Content, FileTagArray, ShareCnt, UserHeaderId and UserName. RegionId was a string. List counters, expire and PacketFlag have unverified semantics. Empty arrays do not establish their element types; observed fields are not guaranteed in every response.

## Verification

Observed at `2026-10-14T16:09:37.437593+00:00`. Read-only requests used no Token, UserId, Command or device context. All returned HTTP 200, including API errors. Successful baseline requests each returned two entries and ReturnCode=0. Removing one key at a time does not establish all accepted values or combinations.

Full request bodies and recorded response shapes: [catalog-probes.json](../../catalog-probes.json).

| Case | HTTP status | ReturnCode | Collection length |
| --- | --- | --- | --- |
| baseline | 200 | 0 | 2 |
| empty body object | 200 | 1 | absent |
| omit AlbumId | 200 | 1 | absent |
| omit StartNum | 200 | 1 | absent |
| omit EndNum | 200 | 1 | absent |
| omit Classify | 200 | 0 | 2 |
| omit FileSize | 200 | 0 | 2 |
| omit FileType | 200 | 0 | 2 |
| omit FileSort | 200 | 0 | 2 |
| omit Version | 200 | 0 | 2 |
| omit RefreshIndex | 200 | 0 | 2 |

## Error response structure

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.ReturnCode` | number | Observed response |
| `$.ReturnMessage` | string | Observed response |
