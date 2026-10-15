# `POST /Mall/GetListV2`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Langue` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.MallList` | array | Key casing unverified |
| `$.MallList[].ImageId` | string | Key casing unverified |
| `$.MallList[].MallId` | integer | Key casing unverified |
| `$.MallList[].Name` | string | Key casing unverified |
| `$.MallList[].Processing` | integer | Key casing unverified |
| `$.MallList[].Score` | integer | Key casing unverified |
| `$.MallList[].Text1` | string | Key casing unverified |
| `$.MallList[].Text2` | string | Key casing unverified |
| `$.MallList[].Text3` | string | Key casing unverified |
| `$.MallList[].type` | integer | Key casing unverified |
| `$.MemorialList` | array | Key casing unverified |
| `$.MemorialList[].ImageFileId` | string | Key casing unverified |
| `$.MemorialList[].LcdImageArray` | array | Key casing unverified |
| `$.MemorialList[].LcdImageArray[]` | string | Element type recovered |
| `$.MemorialList[].MemorialDay` | integer | Key casing unverified |
| `$.MemorialList[].MemorialId` | integer | Key casing unverified |
| `$.MemorialList[].MemorialMoon` | integer | Key casing unverified |
| `$.MemorialList[].MemorialName` | string | Key casing unverified |
| `$.MemorialList[].MemorialTime` | integer | Key casing unverified |
| `$.MemorialList[].RecordFileId` | string | Key casing unverified |
