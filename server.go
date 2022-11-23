package tcp_server

import (
	"net"
)

func New(address string) *server {
	server := &server{
		address: address,
	}

	server.OnConnect(func(c *Client) {})
	server.OnMessage(func(c *Client, message string) {})
	server.OnDisconnect(func(c *Client, err error) {})

	return server
}

// Called when a client connects
func (s *server) OnConnect(callback func(c *Client)) {
	s.onConnect = callback
}

// Called the server gets a message from a client.
func (s *server) OnMessage(callback func(c *Client, message string)) {
	s.onMessage = callback
}

// Called when a connection is closed.
func (s *server) OnDisconnect(callback func(c *Client, err error)) {
	s.onDisconnect = callback
}

func (s *server) Connect() error {
	l, err := net.Listen("tcp", s.address)
	if err == nil {
		s.listener = l
	}
	return err
}

func (s *server) IsStarted() bool {
	return s.listener != nil
}

func (s *server) Close() error {
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

func (s *server) listenerLoop() {
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

func (s *server) Listen() error {
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
