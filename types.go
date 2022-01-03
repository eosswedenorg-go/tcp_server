
package tcp_server

import (
	"net"
)

type Client struct {
	conn   net.Conn
	Addr   net.Addr
	Server *server
}

type server struct {
	address 		string
	onConnect		func(c *Client)
	onDisconnect	func(c *Client, err error)
	onMessage		func(c *Client, message string)
}
