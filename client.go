
package tcp_server

import (
    "bufio"
    "net"
    "time"
)

// Read client messages.
func (c *Client) read() {

    defer c.Close()
    defer c.Server.wg.Done()

    timeout := 100 * time.Millisecond

    reader := bufio.NewReader(c.conn)
    for c.Server.running {

        c.conn.SetReadDeadline(time.Now().Add(timeout))

        message, err := reader.ReadString('\n')
        if err != nil {
            c.Server.onDisconnect(c, err)
            return
        }
        c.Server.onMessage(c, message)
    }

    c.Server.onDisconnect(c, nil)
}

// Write string to client.
func (c *Client) WriteString(message string) (int, error) {
    return c.Write([]byte(message))
}

// Write bytes to client
func (c *Client) Write(b []byte) (int, error) {
    c.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
    n, err := c.conn.Write(b)
    if err != nil {
        c.conn.Close()
        c.Server.onDisconnect(c, err)
    }
    return n, err
}

func (c *Client) Conn() net.Conn {
    return c.conn
}

func (c *Client) Close() error {
    return c.conn.Close()
}
