# `POST /Hot/GetHotFiles32`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `DeviceType` | body | integer | 0 | Not verified | Recovered key |  |
| `IsTest` | body | boolean | false | Not verified | Recovered key |  |
| `Type` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.VendorList` | array | Key casing unverified |
| `$.VendorList[].FileList` | array | Key casing unverified |
| `$.VendorList[].FileList[].FileId` | string | Key casing unverified |
| `$.VendorList[].FileList[].Sha1` | string | Key casing unverified |
| `$.VendorList[].FileList[].Version` | integer | Key casing unverified |
| `$.VendorList[].VendorId` | integer | Key casing unverified |
