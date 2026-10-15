# Divoom HTTP endpoint index

238 endpoint documents: 235 POST routes and three GET routes, grouped by server role. Each document contains its request parameters, response structure and available examples and verification findings.

See [shared hosts and authentication](README.md), [machine-readable contracts](endpoints.json), and [artwork formats](divoom-artwork-formats.md).

The directory structure is `<role>/<endpoint path>.md`. For a full URL, the endpoint path is taken from that URL; its server remains stated in the endpoint heading.

| Role | Endpoint documents | Host selection |
| --- | --- | --- |
| [Command](command/README.md) | 220 | Regional command hosts |
| [Messaging](messaging/README.md) | 4 | Regional messaging hosts |
| [Upload](upload/README.md) | 6 | Regional upload hosts |
| [File downloads](file-downloads/README.md) | 1 | Regional file hosts |
| [Absolute URLs](absolute-url/README.md) | 7 | Explicit server URLs, currently the AI service |

For example, `POST /MessageGroup/JoinGroup` is documented in [messaging/MessageGroup/JoinGroup.md](messaging/MessageGroup/JoinGroup.md).
