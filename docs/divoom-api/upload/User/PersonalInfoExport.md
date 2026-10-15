# `POST /User/PersonalInfoExport`

[All endpoints](../../endpoints.md)

Host role: **upload**.

Request Content-Type: `multipart/form-data`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `ExportEmail` | body | string | null | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).
For multipart requests, the `json` form part holds these fields. Binary part names are endpoint-specific and may remain unverified.

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
