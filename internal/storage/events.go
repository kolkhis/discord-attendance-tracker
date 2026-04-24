package storage

import (
	"database/sql"
	"fmt"
)

type Event struct {
	EventID            string         `json:"event_id"`
	GuildID            string         `json:"guild_id"`
	ChannelID          sql.NullString `json:"channel_id"`
	Name               string         `json:"name"`
	EntityType         string         `json:"entity_type"` // 1=STAGE_INSTANCE, 2=VOICE, 3=EXTERNAL
	ScheduledStartTime string         `json:"scheduled_start_time"`
	ScheduledEndTime   sql.NullString `json:"scheduled_end_time"`
	TrackingOpenTime   sql.NullString `json:"tracking_open_time"`
	TrackingCloseTime  sql.NullString `json:"tracking_close_time"`
	CreatedAt          string         `json:"created_at"`
	UpdatedAt          string         `json:"updated_at"`
}

const (
	EntityTypeStageInstance = 1
	EntityTypeVoice         = 2
	EntityTypeExternal      = 3
)

func (db *DB) UpsertEvent(e *Event) error {
	// Update event if it already exists, otherwise insert a new one
	const query = `
INSERT INTO events (
	event_id,
	guild_id,
	channel_id,
	name,
	entity_type,
	scheduled_start_time,
	scheduled_end_time,
	tracking_open_time,
	tracking_close_time,
	created_at,
	updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(event_id) DO UPDATE SET
	guild_id = excluded.guild_id,
	channel_id = excluded.channel_id,
	name = excluded.name,
	entity_type = excluded.entity_type,
	scheduled_start_time = excluded.scheduled_start_time,
	scheduled_end_time = excluded.scheduled_end_time,
	tracking_open_time = excluded.tracking_open_time,
	tracking_close_time = excluded.tracking_close_time,
	created_at = excluded.created_at,
	updated_at = excluded.updated_at;
`

	_, err := db.conn.Exec(query,
		e.EventID,
		e.GuildID,
		e.ChannelID,
		e.Name,
		e.EntityType,
		e.ScheduledStartTime,
		e.ScheduledEndTime,
		e.TrackingOpenTime,
		e.TrackingCloseTime,
		e.CreatedAt,
		e.UpdatedAt,
	)
	return err

}

func (db *DB) GetEvent(eventID string) (*Event, error) {
	const query = `
SELECT
	event_id,
	guild_id,
	channel_id,
	name,
	entity_type,
	scheduled_start_time,
	scheduled_end_time,
	tracking_open_time,
	tracking_close_time,
	created_at,
	updated_at
FROM events
WHERE event_ID = ?;
`
	var e Event
	err := db.conn.QueryRow(query, eventID).Scan(
		&e.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Error on get event %s: %w", eventID, err)
	}
	return &e, nil
}

func (db *DB) ListOpenTrackingEvents() ([]Event, error) {
	const query = `
SELECT 
	event_id,
	guild_id,
	channel_id,
	name,
	entity_type,
	scheduled_start_time,
	scheduled_end_time,
	tracking_open_time,
	tracking_close_time,
	created_at,
	updated_at
FROM events
WHERE tracking_open_time IS NOT NULL
	AND tracking_close_time IS NULL
ORDER_BY scheduled_start_time;
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("List open tracking events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.EventID,
			&e.GuildID,
			&e.ChannelID,
			&e.Name,
			&e.EntityType,
			&e.ScheduledStartTime,
			&e.ScheduledEndTime,
			&e.TrackingOpenTime,
			&e.TrackingCloseTime,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Scan event row: %w", err)
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error in interating event rows: %w", err)
	}
	return events, nil
}

// func (db *DB) DeleteEvent(eventID string) error         { return nil }
