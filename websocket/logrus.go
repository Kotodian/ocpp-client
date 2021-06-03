package websocket

func (c *Client) withSN(sn string) {
	c.entry = c.entry.WithField("sn", sn)
}

func (c *Client) withAddr(addr string) {
	c.entry = c.entry.WithField("addr", addr)
}
