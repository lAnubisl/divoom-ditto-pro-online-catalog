# `POST /User/ClickAdvert`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `AdvertName` | body | string | null | Not verified | Recovered key |  |
| `CountryISOCode` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.CommentCnt` | integer | Key casing unverified |
| `$.ForumId` | integer | Key casing unverified |
| `$.KeyWork` | string | Key casing unverified |
| `$.LikeCnt` | integer | Key casing unverified |
| `$.ShareCnt` | integer | Key casing unverified |
