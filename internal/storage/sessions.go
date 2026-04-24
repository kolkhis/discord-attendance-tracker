package storage

import (
	"database/sql"
	"fmt"
)

type VoiceSession struct {
	ID        int64          `json:"id"`
	EventID   string         `json:"event_id"`
	UserID    string         `json:"user_id"`
	ChannelID string         `json:"channel_id"`
	JoinedAt  string         `json:"joined_at"`
	LeftAt    sql.NullString `json:"left_at"`
}

func (db *DB) StartVoiceSession(eventID, userID, channelID, joinedAt string) error {
	const query = `
INSERT INTO voice_sessions (
    event_id,
    user_id,
    channel_id,
    joined_at
) VALUES (?, ?, ?, ?);
`
	_, err := db.conn.Exec(
		query,
		eventID,
		userID,
		channelID,
		joinedAt,
	)

	if err != nil {
		return fmt.Errorf("Start Voice Session error, event=%v user=%v: %v", eventID, userID, err)
	}

	return nil
}

func (db *DB) EndVoiceSession(eventID, userID, leftAt string) error {
	const query = `
UPDATE voice_sessions
SET left_at = ?
WHERE event_id = ?
    AND user_id = ?
    AND left_at IS NULL;
`
	result, err := db.conn.Exec(query, leftAt, eventID, userID)
	if err != nil {
		return fmt.Errorf("End voice session error: event=%v, user=%v, err=%v\n", eventID, userID, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected error in EndVoiceSession, event=%v, user=%v, err=%v\n", eventID, userID, err)
	}

	if rows == 0 {
		return fmt.Errorf("No open voice session found. event=%v, user=%v\n", eventID, userID)
	}
	return nil
}

func (db *DB) GetOpenVoiceSession(eventID, userID string) (*VoiceSession, error) {
	const query = `
SELECT
    id,
    event_id,
    user_id,
    channel_id,
    joined_at,
    left_at
FROM voice_sessions
WHERE event_id = ?
    AND user_id = ?
    AND left_at IS NULL
ORDER BY joined_at DESC
LIMIT 1;
`
	var s VoiceSession
	err := db.conn.QueryRow(query, eventID, userID).Scan(
		&s.ID,
		&s.EventID,
		&s.UserID,
		&s.ChannelID,
		&s.JoinedAt,
		&s.LeftAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No open voice sessions
		}
		return nil, fmt.Errorf("Error fetching open voice sessions: event=%v, user=%v, err=%v\n", eventID, userID, err)
	}
	return &s, nil
}

// Helper fn to check if event/user has open voice sessions for the given event
func (db *DB) HasOpenVoiceSession(eventID, userID string) (bool, error) {
	session, err := db.GetOpenVoiceSession(eventID, userID)
	if err != nil {
		return false, err
	}
	return session != nil, nil
}
