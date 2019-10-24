package rethinkdb

type Query interface {
	CheckIsPaused() (bool, error)
}

func (db *Client) CheckIsPaused() (bool, error) {
	err := db.Connect()
	if err != nil {
		return false, err
	}
	defer func() {
		_ = db.Disconnect()
	}()

	var state struct {
		id          string
		shouldPause bool
	}

	err = db.Get("system", "state", state)
	if err != nil {
		return false, err
	} else {
		return state.shouldPause, nil
	}
}
