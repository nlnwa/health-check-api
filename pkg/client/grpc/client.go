package grpc

type Options struct {
	VeidemannApiUrl string
}

type Client struct {
	veidemannApiUrl string
}

func New(options Options) Client {
	return Client{
		veidemannApiUrl: options.VeidemannApiUrl,
	}
}
