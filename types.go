package tcp_server

import (
	"net"
	"sync"
)

// Struct that represents a client connection
type Client struct {
	// Actual network connections.
	conn net.Conn

	// Client address
	Addr net.Addr

	// Pointer to server object.
	Server *Server
}

// Struct that represents a network listener
type Server struct {
	// Actual listener connection
	listener net.Listener

	// Address to bind the server on
	address string

	// WaitGroup to make all go routines exit gracefully.
	wg sync.WaitGroup

	// Simple bool to signal that we are in listener loop.
	// Clients should never write to this (only server.Close() may do so)
	// but read the value and exit when it's false.
	running bool

	// Callback functions
	onConnect    func(c *Client)
	onDisconnect func(c *Client, err error)
	onMessage    func(c *Client, message string)
}
