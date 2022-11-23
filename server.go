package tcp_server

import (
	"net"
	"sync"
)

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

func New(address string) *Server {
	server := &Server{
		address: address,
	}

	server.OnConnect(func(c *Client) {})
	server.OnMessage(func(c *Client, message string) {})
	server.OnDisconnect(func(c *Client, err error) {})

	return server
}

// Called when a client connects
func (s *Server) OnConnect(callback func(c *Client)) {
	s.onConnect = callback
}

// Called the server gets a message from a client.
func (s *Server) OnMessage(callback func(c *Client, message string)) {
	s.onMessage = callback
}

// Called when a connection is closed.
func (s *Server) OnDisconnect(callback func(c *Client, err error)) {
	s.onDisconnect = callback
}

func (s *Server) Connect() error {
	l, err := net.Listen("tcp", s.address)
	if err == nil {
		s.listener = l
	}
	return err
}

func (s *Server) IsStarted() bool {
	return s.listener != nil
}

func (s *Server) Close() error {
	if !s.IsStarted() {
		return nil
	}

	// set running to false
	s.running = false

	// Close the listener.
	err := s.listener.Close()

	// Wait for go routines to exit.
	s.wg.Wait()

	// Cleanup
	s.listener = nil
	return err
}

func (s *Server) listenerLoop() {
	defer s.wg.Done()

	for {
		// This blocks until an client is accepted
		// or s.listener.Close() is called.
		conn, err := s.listener.Accept()
		if err != nil {
			break
		}

		c := &Client{
			conn:   conn,
			Server: s,
			Addr:   conn.RemoteAddr(),
		}
		s.wg.Add(1)
		s.onConnect(c)
		go c.read()
	}
}

func (s *Server) Listen() error {
	if !s.IsStarted() {
		err := s.Connect()
		if err != nil {
			return err
		}
	}

	s.running = true
	s.wg.Add(1)
	go s.listenerLoop()
	return nil
}
