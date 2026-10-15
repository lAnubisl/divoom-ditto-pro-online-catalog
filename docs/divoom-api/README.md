# Divoom hosted HTTP API

This reference documents hosted content, account and server-service HTTP
endpoints. Findings were recorded on 14 October 2026. Requiredness, key casing
and response types that have not been verified are marked explicitly.

- [Endpoint index](endpoints.md): 238 separate documents grouped by server role,
  each with request parameters, response structures, examples and verification findings.
- [Machine-readable endpoint contracts](endpoints.json): the same 235 POST and three GET routes, with links to their documents.
- [Catalog probe evidence](catalog-probes.json): recorded requests and responses supporting verified findings.
- [Artwork formats](divoom-artwork-formats.md): encryption and decoding of downloaded artwork.

## Hosts

| Role | International | China |
| --- | --- | --- |
| command | `https://appin.divoom-gz.com` | `https://appchina.divoom-gz.com` |
| file downloads | `https://fin.divoom-gz.com` | `https://fchina.divoom-gz.com` |
| messaging | `https://rongin.divoom-gz.com` | `https://rong.divoom-gz.com` |
| upload | `https://appusa.divoom-gz.com` | `https://appchina.divoom-gz.com` |

Append a relative endpoint path to the stated host. For example,
`POST https://appin.divoom-gz.com/Discover/GetAlbumListV3`.
Full endpoint URLs override regional routing. AI routes use
`https://ai.divoom-gz.com`.

## Requests and authentication

JSON POST requests use `Content-Type: application/json; charset=utf-8`.
GET routes use path or query parameters and have no request body.
Multipart requests use POST with `Content-Type: multipart/form-data` and a
`json` form part containing the request object, plus endpoint-specific binary
parts. Binary part names are not established for every upload route.

`Token` and `UserId` identify an authenticated session. Supply them as JSON
numbers when required by an account operation. They were omitted from the
successful catalog probes. Session requirements on other routes are unverified.
Some requests include `Command` matching the endpoint path; the tested catalog
routes accepted omission. Do not assume omission works on every endpoint.

Use the [`POST /UserLogin` contract](command/UserLogin.md) for login request details. `ReturnCode=11` indicates token mismatch and requires a fresh session.

## Required and optional fields

Required in an omission probe means removing that key from a successful request
caused an API error. Optional means the request succeeded without that key;
omitting a filter can still change the returned content. Neither result proves
all accepted values, combinations or regional behavior.

Not verified means the request field was recovered but server acceptance of
omission has not been checked. Default values describe request construction,
not server defaults. Recovered response fields describe the expected contract;
only observed catalog fields have confirmed wire casing and JSON types.

## API status codes

HTTP 200 is not enough to establish success. Known status values include:

| ReturnCode | Meaning |
| --- | --- |
| 0 | Success |
| 1 | General failure |
| 2 | Email already registered |
| 3 | Required information incomplete; also observed for missing category pagination |
| 4 | Login user does not exist |
| 5 | Incorrect login password |
| 8 / 9 | Query / processing conditions do not match |
| 10 | Command mismatch |
| 11 | Token mismatch |
| 12 / 13 | Empty request / invalid JSON |
| 16 | Duplicate uploaded gallery file |
| 17 / 18 | Invalid phone format / verification code |
| 19 | Uploads too frequent |
| 22 | Review required |
| 23 / 24 / 25 | Blocked / upload restricted / account banned |
| 9999999 | Account locked |

These are recovered API interpretations. Server messages and endpoint-specific
uses can differ; the omission evidence reports actual returned numbers.


This reference is a recovered contract, not an official server specification or a guarantee of the complete hosted route set.
