package mock

type Client struct {
	hasActivity bool
}

func New() Client {
	return Client{}
}

func NewWithActivity() Client {
	return Client{true}
}

func (pc Client) GetActivity() (bool, error) {
	return pc.hasActivity, nil
}
