package mock

type MockClient struct {
}

func NewMockClient() MockClient {
	return MockClient{}
}

func (c MockClient) GetRunningJobs() (bool, error) {
	return true, nil
}
