# `POST /Tag/GetUserList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `TagName` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.CurListNum` | integer | Recovered key |
| `$.FollowListNum` | integer | Recovered key |
| `$.UserList` | array | Recovered key |
| `$.UserList[].CountryISOCode` | string | Key casing unverified |
| `$.UserList[].FansCnt` | integer | Key casing unverified |
| `$.UserList[].HeadId` | string | Key casing unverified |
| `$.UserList[].IsFollow` | integer | Key casing unverified |
| `$.UserList[].Level` | integer | Key casing unverified |
| `$.UserList[].NickName` | string | Key casing unverified |
| `$.UserList[].PixelAmbId` | string | Key casing unverified |
| `$.UserList[].PixelAmbName` | string | Key casing unverified |
| `$.UserList[].RegionId` | integer | Key casing unverified |
| `$.UserList[].Score` | integer | Key casing unverified |
| `$.UserList[].SelectFlag` | integer | Key casing unverified |
| `$.UserList[].UserId` | integer | Key casing unverified |
| `$.UserList[].isAdd` | boolean | Key casing unverified |
| `$.UserList[].select` | boolean | Key casing unverified |
