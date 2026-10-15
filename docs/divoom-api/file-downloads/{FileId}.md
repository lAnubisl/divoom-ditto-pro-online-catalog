# `GET /{FileId}`

[All endpoints](../endpoints.md)

Host role: **file downloads**.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `FileId` | path | string | "Not established" | Required to identify the requested file; omission behavior untested. | Recovered key | Artwork path supplied by a gallery response. |

## Response structure

Recovered contract; not live-verified.

Binary file bytes; no JSON response envelope.

## Behavior and interpretation

No request body is used. The response contains proprietary artwork bytes, not necessarily GIF data. See the artwork format reference for decoding. See [artwork formats](../divoom-artwork-formats.md).
An already absolute artwork URL can be requested directly instead of prefixing the regional file host.
