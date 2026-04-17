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
```txt
events
- event_id TEXT PRIMARY KEY
- guild_id TEXT NOT NULL
- channel_id TEXT
- name TEXT NOT NULL
- entity_type INTEGER NOT NULL
- status INTEGER NOT NULL
- scheduled_start_time TEXT NOT NULL
- scheduled_end_time TEXT
- tracking_open_time TEXT
- tracking_close_time TEXT
- created_at TEXT NOT NULL
- updated_at TEXT NOT NULL

subscriptions
- event_id TEXT NOT NULL
- user_id TEXT NOT NULL
- subscribed_at TEXT NOT NULL
- unsubscribed_at TEXT
- PRIMARY KEY (event_id, user_id)

voice_sessions
- id INTEGER PRIMARY KEY
- event_id TEXT NOT NULL
- user_id TEXT NOT NULL
- channel_id TEXT NOT NULL
- joined_at TEXT NOT NULL
- left_at TEXT
- source TEXT NOT NULL   -- gateway | reconcile

attendance_rollups
- event_id TEXT NOT NULL
- user_id TEXT NOT NULL
- total_seconds INTEGER NOT NULL
- first_joined_at TEXT
- last_left_at TEXT
- was_subscribed INTEGER NOT NULL
- attended INTEGER NOT NULL
- qualified_attendance INTEGER NOT NULL
- no_show INTEGER NOT NULL
- walk_in INTEGER NOT NULL
- PRIMARY KEY (event_id, user_id)
```

## Resources

Examples:

- [Slash Commands](https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go)

