# Clock settings

The English `/clock` page has the shared navigation and autoplay status widget.
It offers selectable preview cards for clock styles 0–15 and an RGB color picker.
The pixel layouts are transcribed from the user's numbered device photographs.
Style 14 is an empty black display (No clock). Previews show the representative
times and animation frames in the photographs, rather than the live device time.
Changing the color updates clock-colored cells immediately; decorative colors,
including greens in rainbow borders, remain fixed. Selection alone sends nothing.
Native radio controls support keyboard selection and highlight the chosen card.

**Apply and show clock** pauses autoplay, waits for any active cycle to cancel
and restore the clock, and then sends the new appearance. Resume autoplay on
the Autoplay page. Subsequent automatic clock restoration and the existing
manual clock command use the last successfully applied appearance.

Clock appearance is saved in the same external `autoplay.json` as the schedule,
using `style_id` (0–15) and `style_color` (`#RRGGBB`). Successful clock application
updates these fields while preserving the autoplay source and timing. On startup,
the client loads the saved appearance before scheduling any animation, without
sending an extra clock command. Automatic and manual clock restoration use it.
Older files without these fields receive defaults: style 0 and green `#00FF00`.
These are client defaults, not device readback. Failed device requests keep the
previous appearance and leave autoplay paused. A write failure after successful
application is shown explicitly in the UI; it does not claim the settings were saved.

## Gateway contract

Verified against `Divoom.Api/Endpoints/DisplayEndpoints.cs` and the gateway's
HTTP client guide on 12 October 2026:

```http
POST /api/display/clock
Content-Type: application/json
X-Api-Key: <backend-only key>

{"style":0,"color":"#00FF00"}
```

Style must be an integer from 0 to 15. Color must match `#RRGGBB`, with
case-insensitive hex digits. Successful response:
`{"status":"sent","mode":"clock"}`. Invalid values return HTTP 400.
Requests without a body remain supported with default appearance. The gateway
uses 24-hour format and disables weather, temperature and calendar fields.
Its response confirms the server command completed, without device readback.

UI `POST /clock/settings` requires the existing session and CSRF token and
serializes with other manual commands. The API key stays on the backend.
