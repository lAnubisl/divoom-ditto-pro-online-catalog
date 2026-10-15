# `POST https://ai.divoom-gz.com/ai/text2image/generate`

[All endpoints](../../../endpoints.md)

Host role: **absolute URL**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `AiEncryption` | body | string | null | Nonempty valid signature required in live probes | Verified key | Lowercase hex HMAC-MD5 of the decimal `UserId`; see below. |
| `JetBraking` | body | integer | 0 | Not verified | Recovered key |  |
| `Text` | body | string | null | Omission not tested | Verified key | Image prompt. |
| `Size` | body | integer | 128 | Omission not tested | Verified key | 128 produced a 128x128 image; 16 submitted a job but completion was not observed. App offers 128 and 256. |
| `Weight` | body | integer | 100 | Not verified | Recovered key |  |
| `Seed` | body | integer | 100 | Not verified | Recovered key |  |
| `Fast` | body | integer | 0 | Not verified | Recovered key |  |

The successful
request included numeric `Token` and `UserId` from `/UserLogin`; their omission
was not tested. `Command` was omitted successfully.

## Signing `AiEncryption`

The Android 3.8.40 app's `AiPixelModel.l` calls `O2.a.a` to sign the decimal
user ID. That helper uses HMAC-MD5 with the UTF-8 key
`DivoomBluetoothDevice<>?`, then encodes the digest as 32 lowercase hex characters:

```python
import hashlib
import hmac

ai_encryption = hmac.new(
    b"DivoomBluetoothDevice<>?",
    str(user_id).encode("utf-8"),
    hashlib.md5,
).hexdigest()
```

Use your own authenticated session's user ID. Plain MD5 of the user ID does
not work. The HMAC method was confirmed by a successful live request.

## Response structure

Live generation and result retrieval verified on 14 October 2026.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Verified key |
| `$.ReturnMessage` | string | Verified key |
| `$.data` | object or null | Verified key |
| `$.data.promptId` | string | Verified key |
| `$.data.status` | null in submission response | Verified key; non-null type unverified |
| `$.data.seed` | null in submission response | Verified key; non-null type unverified |
| `$.data.errorMessage` | null in submission response | Verified key; non-null type unverified |

`DeviceId` was not present in the observed responses.

## Verified request and result workflow

The successful request used this prompt:

> Pixel art of a cute orange cat sitting on a crescent moon, dark blue starry background, vibrant colors, simple clean shapes

It included `AiEncryption`, `Token`, `UserId`, `JetBraking=0`, `Size=128`,
`Weight=100`, `Seed=100`, and `Fast=0`. No credentials or account identifiers
are included in the examples below; placeholder values are illustrative.

Submission returned HTTP 200 with this structure:

```json
{
  "ReturnCode": 0,
  "ReturnMessage": "",
  "data": {
    "promptId": "<job-id>",
    "status": null,
    "seed": null,
    "errorMessage": null
  }
}
```

1. Keep `data.promptId` and the session's `UserId`.
2. Poll [`GET /ai/comfyui/history/{promptId}/{userId}`](../comfyui/history/%7BpromptId%7D/%7BuserId%7D.md). The successful test polled at five-second intervals.
3. When polling returns a nonempty `filename`, download using
   [`GET /ai/comfyui/image/view`](../comfyui/image/view.md) with URL-encoded
   `filename` and `subfolder` query parameters.

The first poll returned `ReturnCode=1`, `status="sending"`,
`queuePosition=-1`, and null `filename`, `progress`, and `subfolder`. Its
message was `获取历史任务详情失败: null` (failed to retrieve task details: null).
The next poll, about five seconds later, returned:

```json
{
  "ReturnCode": 0,
  "ReturnMessage": "",
  "data": {
    "filename": "<generated-file>.webp",
    "status": "success",
    "queuePosition": -1,
    "progress": 100,
    "subfolder": "<returned-subfolder>"
  }
}
```

The downloaded file was a 128x128 still WebP image. Preserve the returned
file format instead of assuming PNG. The transient polling error above did
not prevent this job from completing; it should not alone be treated as a
terminal generation failure. Use bounded polling and report persistent
errors or explicit `status="error"`.

## Observed failures and size limitations

HTTP 200 does not establish API success. Both signing failures returned
`ReturnCode=1` with `data=null`:

| Signature supplied | `ReturnMessage` | Interpretation |
| --- | --- | --- |
| `AiEncryption=null` | `生成图片失败: 用户ID或md5加密用户ID为空` | User ID or its MD5 signature is empty. |
| Plain MD5 of `UserId` | `生成图片失败: 用户ID与md5加密用户ID不匹配` | User ID and signature do not match. |

The server messages say MD5, but the app and successful probe establish
**HMAC-MD5** as the signing algorithm.

A second request with the same prompt and settings except `Size=16` returned
`ReturnCode=0` and a job ID. Repeated polling continued to return the
`sending`/task-details-error response, with no image filename. Polling was
stopped at the user's request; no terminal failure or completed 16x16 image
was observed. Job acceptance therefore does not confirm native 16x16 support.

The app's `AiPixelFragment` offers only sizes 128 and 256, and displays a
size warning for smaller devices. Size 256 was not live-tested here. A
practical route for a 16x16 display is to generate at a verified size and
resize the result locally.

No animation-generation endpoint or frame-count/duration parameter was found
in the recovered AI workflow. The verified result was a still image; native
AI animation support remains unestablished.
