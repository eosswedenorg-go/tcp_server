
package tcp_server

import (
	"bufio"
	"net"
)

// Read client messages.
func (c *Client) read() {
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onDisconnect(c, err)
			return
		}
		c.Server.onMessage(c, message)
	}
}

// Write string to client.
func (c *Client) WriteString(message string) error {
	return c.Write([]byte(message))
}

// Write bytes to client
func (c *Client) Write(b []byte) error {
	_, err := c.conn.Write(b)
	if err != nil {
		c.conn.Close()
		c.Server.onDisconnect(c, err)
	}
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}
