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

func connect(host string, port string) *Client {
	c, err := textproto.Dial("tcp4", host+":"+port)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	log.Print("Sending header to aelita")
	c.PrintfLine(PROTO)
	resp, err := c.ReadLine()
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

func (c *Client) send(cmd string) string {
	log.Printf("Sending '%s'",cmd)
	c.cn.PrintfLine(cmd)
	resp,err := c.cn.ReadLine()
	if err != nil {
		log.Fatal("Could not send command")
		return "Failed to send command"
	}
	return resp
}

func (c *Client) disconnect() {
	log.Print("Closing connection")
	c.cn.PrintfLine("close")
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
	c.cn.Close()
}
