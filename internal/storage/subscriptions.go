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
		return fmt.Errorf("Err on UpsertEventSubscription, event=%v, user=%v, subscribedAt=%v, err=%w", eventID, userID, subscribedAt, err)
	}
	return nil
}

func (db *DB) DeleteEventSubscription(eventID, userID, subscribedAt string) error {
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
		return fmt.Errorf("DeleteEventSubscription error, event=%v, user=%v, subscribedAt=%v, err=%w", eventID, userID, subscribedAt, err)
	}
	return nil
}

func (db *DB) HasEventSubscription(eventID, userID string) (bool, error) {
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
		return false, fmt.Errorf("HasEventSubscription error, event=%v, user=%v, err=%w", eventID, userID, err)
	}

	return true, nil
}

func (db *DB) ListEventSubscriptions(eventID string) ([]EventSubscription, error) {
	return nil, nil
}
