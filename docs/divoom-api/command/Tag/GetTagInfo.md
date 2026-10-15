# `POST /Tag/GetTagInfo`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `TagName` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.FollowCnt` | integer | Key casing unverified |
| `$.GalleryCnt` | integer | Key casing unverified |
| `$.IsFollow` | integer | Key casing unverified |
| `$.UserCnt` | integer | Key casing unverified |
| `$.UserList` | array | Key casing unverified |
| `$.UserList[].HeadId` | string | Key casing unverified |
| `$.UserList[].UserId` | integer | Key casing unverified |
