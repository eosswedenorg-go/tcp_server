
package tcp_server

import (
    "net"
    "sync"
)

// Struct that represents a client connection
type Client struct {

    // Actual network connections.
    conn    net.Conn

    // Client address
    Addr    net.Addr

    // Pointer to server object.
    Server  *server
}

// Struct that represents a network listener
type server struct {
    // Actual listener connection
    listener        net.Listener

    // Address to bind the server on
    address         string

    // WaitGroup to make all go routines exit gracefully.
    wg              sync.WaitGroup

    // Callback functions
    onConnect       func(c *Client)
    onDisconnect    func(c *Client, err error)
    onMessage       func(c *Client, message string)
}
