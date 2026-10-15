# `POST /Medal/GetList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Langue` | body | string | null | Not verified | Recovered key |  |
| `TargetUserId` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.MedalList` | array | Key casing unverified |
| `$.MedalList[].ActionType` | integer | Key casing unverified |
| `$.MedalList[].EndTime` | string | Key casing unverified |
| `$.MedalList[].Explain` | string | Key casing unverified |
| `$.MedalList[].InValidImageId` | string | Key casing unverified |
| `$.MedalList[].IsValid` | integer | Key casing unverified |
| `$.MedalList[].MedalId` | integer | Key casing unverified |
| `$.MedalList[].Name` | string | Key casing unverified |
| `$.MedalList[].StarTime` | string | Key casing unverified |
| `$.MedalList[].SubTitle` | string | Key casing unverified |
| `$.MedalList[].ValidImageId` | string | Key casing unverified |
| `$.MedalList[].ValidTime` | integer | Key casing unverified |
| `$.MedalValidCnt` | integer | Key casing unverified |
