# `POST /MessageGroup/GetGroupList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

No operation-specific body fields were recovered. Session requirements are not verified.

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.ClassifyList` | array | Key casing unverified |
| `$.ClassifyList[].ClassifyName` | string | Key casing unverified |
| `$.ClassifyList[].GroupList` | array | Key casing unverified |
| `$.ClassifyList[].GroupList[].ChannelId` | string | Key casing unverified |
| `$.ClassifyList[].GroupList[].ClassifyName` | string | Key casing unverified |
| `$.ClassifyList[].GroupList[].GroupExplain` | string | Key casing unverified |
| `$.ClassifyList[].GroupList[].GroupId` | string | Key casing unverified |
| `$.ClassifyList[].GroupList[].GroupName` | string | Key casing unverified |
| `$.ClassifyList[].GroupList[].GroupPixelFileId` | string | Key casing unverified |
| `$.ClassifyList[].GroupList[].Joined` | integer | Key casing unverified |
| `$.ClassifyList[].GroupList[].MessageType` | integer | Key casing unverified |
| `$.ClassifyList[].GroupList[].PictureInterval` | integer | Key casing unverified |
| `$.ClassifyList[].GroupList[].PixelInterval` | integer | Key casing unverified |
| `$.ClassifyList[].GroupList[].TextInterval` | integer | Key casing unverified |
| `$.ClassifyList[].GroupList[].VideoInterval` | integer | Key casing unverified |
| `$.ClassifyList[].GroupList[].unread` | integer | Key casing unverified |
