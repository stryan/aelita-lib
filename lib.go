package aelitalib

import (
	"log"
	"net/textproto"
)

const PROTO = "aelita 0.1"

type Client struct {
	cn   *textproto.Conn
	host string
	port string
}

func Connect(host string, port string) *Client {
	c, err := textproto.Dial("tcp4", host+":"+port)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Print("Sending header to aelita")
	id := c.Next()
	c.StartRequest(id)
	c.PrintfLine(PROTO)
	c.EndRequest(id)
	c.StartResponse(id)
	resp, err := c.ReadLine()
	c.EndResponse(id)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if resp == "OK "+PROTO {
		log.Print("Connection established")
		return &Client{c, host, port}
	} else {
		log.Print("Protocol version mismatch")
		return nil
	}
}

func (c *Client) Send(cmd string) string {
	log.Printf("Sending '%s'",cmd)
	id := c.cn.Next()
	c.cn.StartRequest(id)
	c.cn.PrintfLine(cmd)
	c.cn.EndRequest(id)
	c.cn.StartResponse(id)
	resp,err := c.cn.ReadLine()
	if err != nil {
		log.Fatal("Could not send command")
		resp = "Failed to send command"
	}
	c.cn.EndResponse(id)
	return resp
}

func (c *Client) Disconnect() {
	log.Print("Closing connection")
	id := c.cn.Next()
	c.cn.StartRequest(id)
	c.cn.PrintfLine("close")
	c.cn.EndRequest(id)

	c.cn.StartResponse(id)
	resp, err := c.cn.ReadLine()
	if err != nil {
		log.Fatal("Error closing connection")
		return
	}
	if resp == "END" {
		log.Print("Connection closed cleanly")
	} else {
		log.Print("Connection closed poorly")
	}
	c.cn.EndResponse(id)
	c.cn.Close()
}
