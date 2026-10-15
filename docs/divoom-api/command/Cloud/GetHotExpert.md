# `POST /Cloud/GetHotExpert`

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
| `$.ExpertList` | array | Key casing unverified |
| `$.ExpertList[].CountryISOCode` | string | Key casing unverified |
| `$.ExpertList[].ExpertHeadFileID` | string | Key casing unverified |
| `$.ExpertList[].ExpertUserId` | integer | Key casing unverified |
| `$.ExpertList[].Level` | integer | Key casing unverified |
| `$.ExpertList[].NickName` | string | Key casing unverified |
| `$.ExpertList[].PixelAmbId` | string | Key casing unverified |
| `$.ExpertList[].PixelAmbName` | string | Key casing unverified |
