package rsh

type Client struct {
	transport Transport
}

func (c *Client) RunApp(name string) error {
	return nil
}
