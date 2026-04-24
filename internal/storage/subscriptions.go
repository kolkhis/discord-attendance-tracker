package storage

type EventSubscription struct {
	EventID      string `json:"event_id"`
	UserID       string `json:"user_id"`
	SubscribedAt string `json:"subscribed_at"`
}

func (db *DB) UpsertEventSubscription(eventID, userID, subscribedAt string) error {
	return nil
}

func (db *DB) DeleteEventSubscription(eventID, userID, subscribedAt string) error {
	return nil
}

func (db *DB) HasEventSubscription(eventID, userID string) (bool, error) {
	return false, nil
}

func (db *DB) ListEventSubscriptions(eventID string) ([]EventSubscription, error) {
	return nil, nil
}
