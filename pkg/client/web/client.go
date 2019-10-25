package web

type Options struct {
	VeidemannDashboardUrl string
}

type Client struct {
	veidemannDashboardUrl string
}

func New(options Options) Client {
	return Client{
		veidemannDashboardUrl: options.VeidemannDashboardUrl,
	}
}
