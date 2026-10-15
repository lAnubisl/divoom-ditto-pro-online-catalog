# `POST https://ai.divoom-gz.com/ai/image2image/modelIdList`

[All endpoints](../../../endpoints.md)

Host role: **absolute URL**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `Language` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.data` | array | Key casing unverified |
| `$.data[].imageId` | string | Key casing unverified |
| `$.data[].modelId` | integer | Key casing unverified |
| `$.data[].modelIdList` | array | Key casing unverified |
| `$.data[].modelIdList[]` | integer | Element type recovered |
| `$.data[].name` | string | Key casing unverified |
