# `GET https://ai.divoom-gz.com/ai/comfyui/image/view`

[All endpoints](../../../../endpoints.md)

Host role: **absolute URL**.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `filename` | query | string | "Not established" | Not verified | Recovered key | Supplied from result polling. |
| `subfolder` | query | string | "Not established" | Not verified | Recovered key | Supplied from result polling. |

## Response structure

Recovered contract; not live-verified.

Binary file bytes; no JSON response envelope.

## Behavior and interpretation

No request body is used. The response is image bytes, with no JSON envelope.
