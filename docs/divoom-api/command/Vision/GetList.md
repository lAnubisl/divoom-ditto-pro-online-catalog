# `POST /Vision/GetList`

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
| `ClassifyId` | body | integer | 0 | Not verified | Recovered key |  |
| `Flag` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.VisionList` | array | Key casing unverified |
| `$.VisionList[].AddFlag` | integer | Key casing unverified |
| `$.VisionList[].ImageFileId` | string | Key casing unverified |
| `$.VisionList[].ImageSmallFileId` | string | Key casing unverified |
| `$.VisionList[].VideoFileId` | string | Key casing unverified |
| `$.VisionList[].VisionId` | integer | Key casing unverified |
