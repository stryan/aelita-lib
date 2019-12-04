package aelitalib

import (
	"log"
	"net/textproto"
	"strings"
	"strconv"
)

const PROTOV = "0.2"

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
	c.PrintfLine("NEW aelita %v ",PROTOV)
	c.EndRequest(id)
	c.StartResponse(id)
	resp, err := c.ReadLine()
	c.EndResponse(id)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	if resp == "OK aelita "+PROTOV {
		log.Print("Connection established")
		return &Client{c, host, port}
	} else {
		log.Print("Protocol version mismatch")
		return nil
	}
}

func (c *Client) Send(cmd string) uint {
	log.Printf("Sending '%s'",cmd)
	id := c.cn.Next()
	c.cn.StartRequest(id)
	err := c.cn.PrintfLine("CMD "+cmd)
	c.cn.EndRequest(id)
	if err != nil {
		log.Fatal("Could not send request")
	}
	return id
}

func (c *Client) Receive(id uint) string {
	c.cn.StartResponse(id)
	resp,err := c.cn.ReadLine()
	if err != nil {
		log.Fatal("Could not receive response")
	}
	resp_s := strings.Fields(resp)
	if len(resp_s) < 1 || len(resp_s) > 3 {
		log.Fatal("Bad response from server")
	}
	switch resp_s[0] {
	case "ERR":
		//return server error as is
		return resp
	case "CMD":
		//TODO implement server commands
		return "Response not implemented"
	case "END":
		//Server wants to end connection
		c.cn.EndResponse(id)
		c.cn.Close()
		return "Server closed our connection"
	case "DAT":
		if len(resp_s) != 2 {
			//DAT command but bad header
			log.Fatal("Want to receive data, but data header malformed")
		}
		numData,err := strconv.Atoi(resp_s[1])
		if err != nil {
			log.Fatal("Data header with no line count")
		}
		result := make([]string,numData)
		for i:= 0; i < numData; i++ {
			result[i],err = c.cn.ReadLine()
			if err != nil {
				log.Fatal("Could not receive response")
			}
		}
		return strings.Join(result,"\n")
	default:
		return "Server responded with an invalid header"
	}
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
	if resp == "END aelita " + PROTOV {
		log.Print("Connection closed cleanly")
	} else {
		log.Print("Connection closed poorly")
	}
	c.cn.EndResponse(id)
	c.cn.Close()
}
