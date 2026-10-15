# `POST /Dialog/GetInfo`

[All endpoints](../../endpoints.md)

Host role: **command**.

Request Content-Type: `application/json`.

Anonymous access: Not verified.

## Request parameters

| Parameter | Location | JSON / text type | Request default | Required / optional | Key accuracy | Purpose |
| --- | --- | --- | --- | --- | --- | --- |
| `AnnLastIndex` | body | integer | 0 | Not verified | Recovered key |  |
| `isAndroid` | body | integer | 0 | Not verified | Recovered key |  |
| `Language` | body | string | null | Not verified | Recovered key |  |
| `RegionId` | body | integer | 0 | Not verified | Recovered key |  |

Session fields are described in [the guide](../../README.md).

## Response structure

Recovered contract; not live-verified.

| JSON path | JSON type | Key accuracy |
| --- | --- | --- |
| `$.DeviceId` | integer | Key casing unverified |
| `$.ReturnCode` | integer | Key casing unverified |
| `$.ReturnMessage` | string | Key casing unverified |
| `$.AnnImageId` | string | Key casing unverified |
| `$.AnnLastIndex` | integer | Key casing unverified |
| `$.AnnLinkUrl` | string | Key casing unverified |
| `$.AppDescribe` | string | Key casing unverified |
| `$.AppVersion` | integer | Key casing unverified |
| `$.DiscountImageId` | string | Key casing unverified |
| `$.LotteryCnt` | integer | Key casing unverified |
| `$.LotteryImageFileId` | string | Key casing unverified |
| `$.MedalList` | array | Key casing unverified |
| `$.MedalList[].ActionType` | integer | Key casing unverified |
| `$.MedalList[].EndTime` | string | Key casing unverified |
| `$.MedalList[].Explain` | string | Key casing unverified |
| `$.MedalList[].InValidImageId` | string | Key casing unverified |
| `$.MedalList[].IsValid` | integer | Key casing unverified |
| `$.MedalList[].MedalId` | integer | Key casing unverified |
| `$.MedalList[].Name` | string | Key casing unverified |
| `$.MedalList[].StarTime` | string | Key casing unverified |
| `$.MedalList[].SubTitle` | string | Key casing unverified |
| `$.MedalList[].ValidImageId` | string | Key casing unverified |
| `$.MedalList[].ValidTime` | integer | Key casing unverified |
| `$.MedalValidCnt` | integer | Key casing unverified |
| `$.activity` | object | Key casing unverified |
| `$.activity.f11345k` | integer | Key casing unverified |
| `$.activity.f11346l` | object | Key casing unverified |
| `$.activity.f11347m` | object | Key casing unverified |
| `$.activity.f11347m.btn_neg` | object | Key casing unverified |
| `$.activity.f11347m.btn_pos` | object | Key casing unverified |
| `$.activity.f11347m.context` | object | Key casing unverified |
| `$.activity.f11347m.dialog` | object | Key casing unverified |
| `$.activity.f11347m.display` | object | Key casing unverified |
| `$.activity.f11347m.editTextHint` | string | Key casing unverified |
| `$.activity.f11347m.fix_text` | object | Key casing unverified |
| `$.activity.f11347m.img_line` | object | Key casing unverified |
| `$.activity.f11347m.img_loading` | object | Key casing unverified |
| `$.activity.f11347m.itemList` | array | Key casing unverified |
| `$.activity.f11347m.lLayout_bg` | object | Key casing unverified |
| `$.activity.f11347m.lLayout_content` | object | Key casing unverified |
| `$.activity.f11347m.loading_layout` | object | Key casing unverified |
| `$.activity.f11347m.thread` | object | Key casing unverified |
| `$.activity.f11347m.txt_edit` | object | Key casing unverified |
| `$.activity.f11347m.txt_edit.DefMaxBytesLen` | integer | Key casing unverified |
| `$.activity.f11347m.txt_edit.DefMaxStrlen` | integer | Key casing unverified |
| `$.activity.f11347m.txt_edit.fixStr` | string | Key casing unverified |
| `$.activity.f11347m.txt_edit.hideVisual` | boolean | Key casing unverified |
| `$.activity.f11347m.txt_edit.maxBytesLen` | integer | Key casing unverified |
| `$.activity.f11347m.txt_edit.maxStrLen` | integer | Key casing unverified |
| `$.activity.f11347m.txt_edit.showToast` | boolean | Key casing unverified |
| `$.activity.f11347m.txt_edit.utf8BytesMode` | boolean | Key casing unverified |
| `$.activity.f11347m.txt_edit_layout` | object | Key casing unverified |
| `$.activity.f11347m.txt_loading` | object | Key casing unverified |
| `$.activity.f11347m.txt_loading_cancel` | object | Key casing unverified |
| `$.activity.f11347m.txt_msg` | object | Key casing unverified |
| `$.activity.f11347m.txt_title` | object | Key casing unverified |
| `$.activity.f11347m.view_loading_cancel` | object | Key casing unverified |
| `$.activity.f11347m.showTitle` | boolean | Key casing unverified |
| `$.activity.f11347m.showMsg` | boolean | Key casing unverified |
| `$.activity.f11347m.showPosBtn` | boolean | Key casing unverified |
| `$.activity.f11347m.showNegBtn` | boolean | Key casing unverified |
| `$.activity.f11347m.showEdit` | boolean | Key casing unverified |
| `$.activity.f11347m.showItem` | boolean | Key casing unverified |
| `$.activity.f11347m.showLoading` | boolean | Key casing unverified |
| `$.activity.f11347m.loadingListener` | object | Key casing unverified |
| `$.activity.f11347m.loadingTimeoutTime` | integer | Key casing unverified |
| `$.activity.f11347m.loadingTimeout` | object | Key casing unverified |
| `$.activity.f11347m.backgroundDark` | boolean | Key casing unverified |
| `$.activity.f11347m.TAG` | string | Key casing unverified |
| `$.activity.f11347m.view` | object | Key casing unverified |
| `$.activity.f11347m.itemHeight` | integer | Key casing unverified |
| `$.activity.f11347m.itemTextSize` | integer | Key casing unverified |
| `$.activity.f11348n` | object | Key casing unverified |
| `$.activity.f11351q` | object | Key casing unverified |
| `$.activity.f11338d` | string | Key casing unverified |
| `$.activity.f11339e` | array | Key casing unverified |
| `$.activity.f11340f` | integer | Key casing unverified |
| `$.activity.f11341g` | integer | Key casing unverified |
| `$.activity.f11342h` | integer | Key casing unverified |
| `$.activity.f11343i` | integer | Key casing unverified |
| `$.activity.f11344j` | integer | Key casing unverified |
| `$.activity.f11349o` | boolean | Key casing unverified |
| `$.activity.f11350p` | boolean | Key casing unverified |
