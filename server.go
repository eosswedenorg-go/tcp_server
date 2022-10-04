
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
    if s.IsStarted() {
        return s.exit()
    }
    return nil
}

func (s *server) exit() error {
    err := s.listener.Close()
    if err == nil {
        s.listener = nil
    }
    return err
}

func (s *server) Listen() error {

    if ! s.IsStarted() {
        err := s.Connect()
        if err != nil {
            return err
        }
    }
    defer s.exit()

    for {
        conn, err := s.listener.Accept()
        if err != nil {
            return err
        }
        c := &Client{
            conn:   conn,
            Server: s,
            Addr:   conn.RemoteAddr(),
        }
        s.onConnect(c)
        go c.read()
    }
}
