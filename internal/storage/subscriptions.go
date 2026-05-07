package storage

import (
	"database/sql"
	"fmt"
)

type EventSubscription struct {
	EventID      string `json:"event_id"`
	UserID       string `json:"user_id"`
	SubscribedAt string `json:"subscribed_at"`
}

func (db *DB) UpsertEventSubscription(eventID, userID, subscribedAt string) error {
	// Update event subscription, or insert a new one if it isn't there
	// Respond to GUILD_SCHEDULED_EVENT_USER_ADD with this
	const query = `
INSERT INTO event_subscriptions (
    event_id,
    user_id,
    subscribed_at
) VALUES (?, ?, ?)
ON CONFLICT(event_id, user_id) DO UPDATE 
SET subscribed_at = excluded.subscribed_at;
`
	_, err := db.conn.Exec(
		query,
		eventID,
		userID,
		subscribedAt,
	)
	if err != nil {
		return fmt.Errorf("Error on UpsertEventSubscription, event=%v, user=%v, subscribedAt=%v, err=%w", eventID, userID, subscribedAt, err)
	}
	return nil
}

func (db *DB) DeleteEventSubscription(eventID, userID, subscribedAt string) error {
	// Delete event subscription entry
	// Responds to GUILD_SCHEDULED_EVENT_USER_REMOVE
	const query = `
    DELETE FROM event_subscriptions
    WHERE event_id = ?
        AND user_id = ?;
`
	_, err := db.conn.Exec(
		query,
		eventID,
		userID,
		subscribedAt,
	)
	if err != nil {
		return fmt.Errorf("Error on DeleteEventSubscription, event=%v, user=%v, subscribedAt=%v, err=%w", eventID, userID, subscribedAt, err)
	}
	return nil
}

func (db *DB) HasEventSubscription(eventID, userID string) (bool, error) {
	// Returns bool saying if given user has subscription to given event
	// Use when determining attendance
	const query = `
SELECT 1
FROM event_subscriptions
WHERE event_id = ?
    AND user_id = ?
LIMIT 1;
`
	var exists int // assigned with Scan
	err := db.conn.QueryRow(query, eventID, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("Error on HasEventSubscription, event=%v, user=%v, err=%w", eventID, userID, err)
	}

	return true, nil
}

func (db *DB) ListEventSubscriptions(eventID string) ([]EventSubscription, error) {
	// List all users' event subscriptions for given event
	// Use for attendance report
	const query = `
SELECT
    event_id,
    user_id,
    subscribed_at
FROM event_subscriptions
WHERE event_id = ?
ORDER BY subscribed_at;
`

	rows, err := db.conn.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("Error on ListEventSubscriptions, event=%v, err=%w", eventID, err)
	}
	defer rows.Close()

	var subscriptions []EventSubscription

	for rows.Next() {
		var sub EventSubscription

		err := rows.Scan(
			&sub.EventID,
			&sub.UserID,
			&sub.SubscribedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("Error on Scan event subscription, event=%v, err=%w", eventID, err)
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on Row iteration (ListEventSubscriptions), event=%v, err=%w", eventID, err)
	}

	return subscriptions, nil
}
