# Discord Attendance Tracker

A simple bot used to track the attendance of events.  

This bot will track Discord events and log who was in attendance during the
event. It will also compare the list of attendees with those who signed up for
the event.

After the event, it will produce a report with this information.  



## Project Structure
Basic project structure is as follows:
```bash
├── README.md
├── cmd
│   └── bot
└── internal
    ├── attendance
    ├── config
    ├── discord
    ├── reconcile
    ├── report
    └── storage
```

## Database Schema
There will be 4 tables in the database with the following schema (using SQLite):

1. `events`
2. `event_subscriptions`
3. `voice_sessions`
4. `event_attendance`

### `events`

Stores metadata about Discord scheduled events and tracking windows used by the application.

- **event_id** (`TEXT`, PK):
  Unique identifier for the event (Discord event ID).

- **guild_id** (`TEXT`, NOT NULL):
  ID of the Discord guild (server) the event belongs to.

- **channel_id** (`TEXT`, nullable):
  ID of the voice/stage channel associated with the event.
  May be `NULL` for non-channel-based events.

- **name** (`TEXT`, NOT NULL):
  Name/title of the event.

- **entity_type** (`INTEGER`, NOT NULL):
  Type of event (e.g., voice, stage, external).
  Mirrors Discord's event type values.

- **scheduled_start_time** (`TEXT`, NOT NULL):
  Scheduled start time of the event (ISO 8601 format recommended).

- **scheduled_end_time** (`TEXT`, nullable):
  Scheduled end time of the event (if defined).

- **tracking_open_time** (`TEXT`, nullable):
  Time when attendance tracking begins (may be before scheduled start).

- **tracking_close_time** (`TEXT`, nullable):
  Time when attendance tracking ends (may be after scheduled end).

- **created_at** (`TEXT`, NOT NULL):
  Timestamp when this record was created in the local database.

- **updated_at** (`TEXT`, NOT NULL):
  Timestamp when this record was last updated in the local database.

---

### `event_subscriptions`

Tracks which users marked themselves as “interested” in an event.

- **event_id** (`TEXT`, PK, NOT NULL):
  ID of the associated event.

- **user_id** (`TEXT`, PK, NOT NULL):
  ID of the user who subscribed to the event.

- **subscribed_at** (`TEXT`, NOT NULL):
  Timestamp when the user subscribed to the event.

---

### `voice_sessions`

Stores raw voice channel participation data for users during events.

Each row represents a **single join/leave session**.

- **id** (`INTEGER`, PK, AUTOINCREMENT):
  Unique ID for the session record.

- **event_id** (`TEXT`, NOT NULL):
  ID of the associated event.

- **user_id** (`TEXT`, NOT NULL):
  ID of the user participating in the session.

- **channel_id** (`TEXT`, NOT NULL):
  ID of the voice channel joined.

- **joined_at** (`TEXT`, NOT NULL):
  Timestamp when the user joined the channel.

- **left_at** (`TEXT`, nullable):
  Timestamp when the user left the channel.
  `NULL` indicates the session is still active.

---

### `event_attendance`

Stores computed attendance summaries for each user per event.

This table is derived from `voice_sessions` and `event_subscriptions`.

- **event_id** (`TEXT`, PK, NOT NULL):
  ID of the associated event.

- **user_id** (`TEXT`, PK, NOT NULL):
  ID of the user.

- **total_seconds** (`INTEGER`, NOT NULL):
  Total time (in seconds) the user spent in the event's channel.

- **first_joined_at** (`TEXT`, nullable)
  Timestamp of the user's first join during the event.

- **last_left_at** (`TEXT`, nullable)
  Timestamp of the user's last leave during the event.

- **was_subscribed** (`INTEGER`, NOT NULL)
  Whether the user marked themselves as interested in the event.
    - `1` = yes
    - `0` = no

- **attended** (`INTEGER`, NOT NULL)
  Whether the user joined the event channel at least once.
    - `1` = yes
    - `0` = no

- **no_show** (`INTEGER`, NOT NULL)
  User subscribed but never joined the event channel.
    - `1` = yes
    - `0` = no

- **walk_in** (`INTEGER`, NOT NULL)
  User joined the event channel but was not subscribed.
    - `1` = yes
    - `0` = no


## Resources

Examples:

- [Slash Commands](https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go)

