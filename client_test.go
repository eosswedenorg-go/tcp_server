package tcp_server

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient_read(t *testing.T) {
	var msg_read string

	server := New(":9000")
	server.OnMessage(func(c *Client, message string) {
		// Append message to msg_read.
		msg_read += message
	})

	err := server.Listen()
	assert.NoError(t, err)

	// Connect to server
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatal("Failed to connect to test server")
	}

	// Write some data.
	_, err = conn.Write([]byte{'m', 'o', 'r', 'e', '\n'})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 100)
	_, err = conn.Write([]byte{'b', 'e', 'e', 'f', '\n'})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 100)
	conn.Close()

	server.Close()

	assert.Equal(t, "more\nbeef\n", msg_read)
}

func TestClient_readBuffered(t *testing.T) {
	var msg_read string

	server := New(":9000")
	server.OnMessage(func(c *Client, message string) {
		// Append message to msg_read.
		msg_read += message
	})

	err := server.Listen()
	assert.NoError(t, err)

	// Connect to server
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		t.Fatal("Failed to connect to test server")
	}

	// Write some data.
	_, err = conn.Write([]byte{'H', 'e'})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 100)
	_, err = conn.Write([]byte{'l', 'l'})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 100)
	_, err = conn.Write([]byte{'o', '\n'})
	assert.NoError(t, err)
	time.Sleep(time.Microsecond * 100)
	conn.Close()

	server.Close()

	assert.Equal(t, "Hello\n", msg_read)
}

func TestClient_Write(t *testing.T) {
	buf := make([]byte, 3)
	expected := []byte{0x1, 0x2, 0x3}
	rcon, wcon := net.Pipe()

	client := Client{
		conn:   wcon,
		Server: New(":9000"),
	}

	go func() {
		fmt.Println("Writing to client")
		n, err := client.Write(expected)
		assert.NoError(t, err)
		assert.Equal(t, len(expected), n)
		fmt.Println("Write Done")
	}()

	time.Sleep(time.Microsecond * 100)
	fmt.Println("Reading from client")
	n, err := rcon.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(expected), n)
	assert.Equal(t, expected, buf)
	rcon.Close()
	client.Close()
}

func TestClient_WriteError(t *testing.T) {
	_, wcon := net.Pipe()

	client := Client{
		conn:   wcon,
		Server: New(":9000"),
	}

	err := client.Close()
	assert.NoError(t, err)

	// Attempt to write to a closed connection
	n, err := client.Write([]byte{0x0, 0x1, 0x2})
	assert.Error(t, err)
	assert.Equal(t, 0, n)
}

func TestClient_Conn(t *testing.T) {
	con1, con2 := net.Pipe()

	c := Client{
		conn: con1,
	}

	assert.Equal(t, con1, c.Conn())
	assert.NotEqual(t, con2, c.Conn())
}

func TestClient_Close(t *testing.T) {
	con1, con2 := net.Pipe()

	con2.Close()

	cOpen := Client{
		conn: con1,
	}

	cClosed := Client{
		conn: con2,
	}

	assert.NoError(t, cOpen.Close())
	assert.NoError(t, cOpen.Close())
	assert.NoError(t, cClosed.Close())
}
