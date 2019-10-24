package mock

type Client struct {
}

func New() Client {
	return Client{}
}

func (c Client) GetHarvesterPodNames() ([]string, error) {
	return []string{}, nil
}
