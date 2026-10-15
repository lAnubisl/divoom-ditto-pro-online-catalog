# `POST /GetCategoryFileList`

[All endpoints](../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Category` | body | integer | 0 | Not verified | Recovered key |  |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `FileSize` | body | integer | 0 | Not verified | Recovered key |  |
| `FileType` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.CurListNum` | integer | Key casing unverified |
| `$.FileList` | array | Key casing unverified |
| `$.FileList[].FileId` | string | Key casing unverified |
| `$.FileList[].FileName` | string | Key casing unverified |
| `$.FileList[].FileType` | integer | Key casing unverified |
| `$.FileList[].IsLike` | boolean | Key casing unverified |
| `$.FileList[].LikeCnt` | integer | Key casing unverified |
