# `POST /App/GetOtherUserList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `TAG` | body | string | "Not established" | Not verified | Key casing unverified |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.TAG` | string | Key casing unverified |
| `$.UserList` | array | Key casing unverified |
| `$.UserList[].UserId` | integer | Key casing unverified |
| `$.UserList[].UserName` | string | Key casing unverified |
| `$.UserList[].UserNickName` | string | Key casing unverified |
