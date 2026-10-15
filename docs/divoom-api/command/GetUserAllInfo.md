# `POST /GetUserAllInfo`

[All endpoints](../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `CountryISOCode` | body | string | null | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `NotUseIpFlag` | body | integer | 0 | Not verified | Recovered key |  |
| `TimeZone` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.AreaFlag` | integer | Key casing unverified |
| `$.BackgroundId` | string | Key casing unverified |
| `$.Birthday` | string | Key casing unverified |
| `$.BuddyEmail` | string | Key casing unverified |
| `$.BuddyFlag` | integer | Key casing unverified |
| `$.BuddyHeadId` | string | Key casing unverified |
| `$.BuddyNickname` | string | Key casing unverified |
| `$.BuddyRename` | string | Key casing unverified |
| `$.BuddyUserId` | integer | Key casing unverified |
| `$.BuddyUserSign` | string | Key casing unverified |
| `$.CountryISOCode` | string | Key casing unverified |
| `$.DeviceFlag` | integer | Key casing unverified |
| `$.Email` | string | Key casing unverified |
| `$.HeadId` | string | Key casing unverified |
| `$.ManagerFlag` | boolean | Key casing unverified |
| `$.Nickname` | string | Key casing unverified |
| `$.RegionId` | integer | Key casing unverified |
| `$.Sex` | string | Key casing unverified |
| `$.TestFlag` | integer | Key casing unverified |
| `$.ThirdBindType` | integer | Key casing unverified |
| `$.TimeStamp` | integer | Key casing unverified |
| `$.Token` | integer | Key casing unverified |
| `$.UserId` | integer | Key casing unverified |
| `$.UserSign` | string | Key casing unverified |
| `$.UserToken` | string | Key casing unverified |
| `$.UserType` | integer | Key casing unverified |
| `$.VendorFlag` | boolean | Key casing unverified |
| `$.galleryHadMatch` | integer | Key casing unverified |
| `$.galleryMatchName` | string | Key casing unverified |
