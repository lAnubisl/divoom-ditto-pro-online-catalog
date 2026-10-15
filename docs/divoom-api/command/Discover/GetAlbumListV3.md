# `POST /Discover/GetAlbumListV3`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Verified on 14 October 2026.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Required in catalog omission probe (ReturnCode=1) | Recovered key | Inclusive end position. For a 20-position page use StartNum=(page-1)*20+1 and EndNum=page*20. |
| `FileSize` | body | integer | 127 | Optional in catalog omission probe | Recovered key | Size selector: 1 selects 16×16, 2 selects 32×32, 4 selects 64×64; 127 is the broad selector, not a pixel width. |
| `FileSort` | body | integer | 0 | Optional in catalog omission probe | Recovered key | Sort selector; complete value semantics are not established. |
| `StartNum` | body | integer | 0 | Required in catalog omission probe (ReturnCode=1) | Recovered key | Start position; use 1 for the first page. |

Session fields are described in [the guide](../../README.md).

## Request example

```http
POST https://appin.divoom-gz.com/Discover/GetAlbumListV3
Content-Type: application/json
```

```json
{
  "StartNum": 1,
  "EndNum": 20,
  "FileSize": 1,
  "FileSort": 0
}
```

## Response structure

Observed catalog sample.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$` | object | Observed response |
| `$.AlbumList` | array | Observed response |
| `$.AlbumList[]` | object | Observed response |
| `$.AlbumList[].AlbumBigImageId` | string | Observed response |
| `$.AlbumList[].AlbumId` | number | Observed response |
| `$.AlbumList[].AlbumImageId` | string | Observed response |
| `$.AlbumList[].AlbumName` | string | Observed response |
| `$.AlbumList[].GalleryCnt` | number | Observed response |
| `$.AlbumList[].GalleryList` | array | Observed response |
| `$.AlbumList[].GalleryList[]` | object | Observed response |
| `$.AlbumList[].GalleryList[].CommentCnt` | number | Observed response |
| `$.AlbumList[].GalleryList[].FileId` | string | Observed response |
| `$.AlbumList[].GalleryList[].GalleryId` | number | Observed response |
| `$.AlbumList[].GalleryList[].IsLike` | number | Observed response |
| `$.AlbumList[].GalleryList[].LikeCnt` | number | Observed response |
| `$.AlbumList[].LikeCnt` | number | Observed response |
| `$.ReturnCode` | number | Observed response |
| `$.ReturnMessage` | string | Observed response |
| `$.expire` | number | Observed response |

## Behavior and interpretation

Use increasing inclusive position ranges. The live probes checked only the first two positions and do not establish an end-of-list rule. A successful empty collection differs from an API error. A short page alone does not establish completion.
Error responses omit the successful collection; check ReturnCode before reading it.
FileSize changes artwork selection elsewhere, but its effect on album-name selection is not established by first-page comparisons.

## Verification

Observed at `2026-10-14T16:09:37.437593+00:00`. Read-only requests used no Token, UserId, Command or device context. All returned HTTP 200, including API errors. Successful baseline requests each returned two entries and ReturnCode=0. Removing one key at a time does not establish all accepted values or combinations.

Full request bodies and recorded response shapes: [catalog-probes.json](../../catalog-probes.json).

| Case | HTTP status | ReturnCode | Collection length |
| --- | --- | --- | --- |
| baseline | 200 | 0 | 2 |
| empty body object | 200 | 1 | absent |
| omit StartNum | 200 | 1 | absent |
| omit EndNum | 200 | 1 | absent |
| omit FileSize | 200 | 0 | 2 |
| omit FileSort | 200 | 0 | 2 |

## Error response structure

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.ReturnCode` | number | Observed response |
| `$.ReturnMessage` | string | Observed response |
