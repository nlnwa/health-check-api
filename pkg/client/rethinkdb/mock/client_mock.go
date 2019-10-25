package mock

type Client struct {
	isPaused bool
}

func New() *Client {
	return &Client{}
}

func NewPausedClient() *Client {
	return &Client{true}
}

func (db Client) CheckIsPaused() (bool, error) {
	return true, nil
}
