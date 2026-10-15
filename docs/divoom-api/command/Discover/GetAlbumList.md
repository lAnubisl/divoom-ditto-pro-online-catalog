# `POST /Discover/GetAlbumList`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `CountryISOCode` | body | string | null | Not verified | Recovered key |  |
| `Langue` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.AlbumList` | array | Key casing unverified |
| `$.AlbumList[].AlbumBigImageId` | string | Recovered key |
| `$.AlbumList[].AlbumId` | integer | Recovered key |
| `$.AlbumList[].AlbumImageId` | string | Recovered key |
| `$.AlbumList[].AlbumName` | string | Recovered key |
