
package tcp_server

import (
    "net"
    "sync"
    "sync/atomic"
    "time"

    "testing"
    "github.com/stretchr/testify/assert"
)

var server_addr = "127.0.0.1:9989"

func TestConnectOK(t *testing.T) {

    server := New(server_addr)
    err := server.Connect()
    assert.NoError(t, err)
    server.Close()
}

func TestConnectFail(t *testing.T) {

    // Start a blocking server on the same port.
    blocking_server, err := net.Listen("tcp", server_addr)
    assert.NoError(t, err)
    defer blocking_server.Close()

    server := New(server_addr)

    err = server.Connect()
    defer server.Close()
    assert.Error(t, err)
    assert.Equal(t, err.Error(), "listen tcp " + server_addr + ": bind: address already in use")
}

func TestClose(t *testing.T) {

    server := New(server_addr)
    err := server.Connect()

    assert.NoError(t, err)
    assert.True(t, server.IsStarted())

    err = server.Close()

    assert.NoError(t, err)
    assert.False(t, server.IsStarted())
}


func TestIsStartedTrue(t *testing.T) {

    server := New(server_addr)
    err := server.Connect()
    defer server.Close()

    assert.NoError(t, err)
    assert.True(t, server.IsStarted())
}

func TestIsStartedFalse(t *testing.T) {

    server := New(server_addr)
    assert.False(t, server.IsStarted())
}

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

    server.Listen()
    defer server.Close()

    // Add some wait so the server has time to start before we connect.
    time.Sleep(time.Second)

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

func TestGracefullClientShutdown(t *testing.T) {

    server := New(server_addr)
    var count int32 = 0;
    var max int32 = 10;

    var wg sync.WaitGroup

    wg.Add(int(max))

    server.Listen()

    // Add some wait so the server has time to start before we connect.
    time.Sleep(time.Second)

    // Spawn a bunch of clients.
    for i := int32(0); i < max; i++ {
        conn, err := net.Dial("tcp", server_addr)
        if err != nil {
            t.Fatal("Failed to connect to test server")
        }

        // Spawn a go routine that should block until server.Close() is called.
        go func() {
            b := make([]byte, 1)
            conn.Read(b)
            atomic.AddInt32(&count, 1)
            wg.Done()
        }()
    }

    err := server.Close()
    assert.NoError(t, err)

    wg.Wait()

    assert.Equal(t, count, max, "not all clients where gracefully disconnected, expected: %d but got %d", max, count)
}