
package tcp_server

import (
    "net"
    "sync"

    "testing"
)

var server_addr = "localhost:9989"

func TestPingPong(t *testing.T) {

    var pingMsg string
    var wg sync.WaitGroup

    wg.Add(3)

    server := New(server_addr)

    server.OnConnect(func (c *Client) {
        wg.Done()
    })

    server.OnMessage(func (c *Client, message string) {
        pingMsg = message
        err := c.WriteString("pong\n")
        if err != nil {
            t.Fatal("Failed to write pong message")
        }
        wg.Done()
    })

    server.OnDisconnect(func (c *Client, err error) {
        wg.Done()
    })

    go server.Listen()

    conn, err := net.Dial("tcp", server_addr)
    if err != nil {
		t.Fatal("Failed to connect to test server")
	}

    _, err = conn.Write([]byte("ping\n"))
	if err != nil {
		t.Fatal("Failed to send ping message.")
	}

    pongMsg := make([]byte, 5)
    _, err = conn.Read(pongMsg)
    if err != nil {
		t.Fatal("Failed to read pong message from server")
	}

    conn.Close()

    wg.Wait()

    if string(pongMsg) != "pong\n" {
        t.Fatal("Client did not receive the pong message")
    }

    if pingMsg != "ping\n" {
        t.Fatal("Server not receive the ping message")
    }
}
