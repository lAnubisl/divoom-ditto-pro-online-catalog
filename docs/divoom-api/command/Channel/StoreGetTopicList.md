# `POST /Channel/StoreGetTopicList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `EndNum` | body | integer | 0 | Not verified | Recovered key |  |
| `StartNum` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | "Not established" | Not verified | Recovered key |  |
| `CountryISOCode` | body | string | "Not established" | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.TopicList` | array | Recovered key |
| `$.TopicList[].TopicDesc` | string | Recovered key |
| `$.TopicList[].TopicId` | integer | Recovered key |
| `$.TopicList[].TopicImageId` | string | Recovered key |
| `$.TopicList[].TopicName` | string | Recovered key |
