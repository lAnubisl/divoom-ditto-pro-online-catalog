# `GET https://ai.divoom-gz.com/ai/comfyui/history/{promptId}/{userId}`

[All endpoints](../../../../../endpoints.md)

Host role: **absolute URL**.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `promptId` | path | string | "Not established" | Not verified | Recovered key | Supplied from the generation response data.promptId. |
| `userId` | path | integer rendered as text | "Not established" | Not verified | Recovered key | Supplied from the account session; authorization checks untested. |

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.data` | object | Recovered key |
| `$.data.filename` | string | Recovered key |
| `$.data.subfolder` | string | Recovered key |
| `$.data.progress` | integer | Recovered key |
| `$.data.queuePosition` | integer | Recovered key |
| `$.data.status` | string | Recovered key |

## Behavior and interpretation

No request body is used. status="error" indicates generation failure. Nonempty filename and subfolder identify the image to download.
