
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

func (s *server) Listen() error {

	sock, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer sock.Close()

	for {
		conn, _ := sock.Accept()
		c := &Client{
			conn:   conn,
			Server: s,
			Addr: conn.RemoteAddr(),
		}
		s.onConnect(c)
		go c.read()
	}

	return nil
}
