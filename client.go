package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn

	flag int
}

func NewClient(ServerIp string, ServerPort int) *Client {
	client := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))

	if err != nil {
		fmt.Println("net.Dial error: ", err)
		return nil
	}

	client.conn = conn
	return client
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.修改名称")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if 0 <= flag && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>请输入合法范围内的数字<<<")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>请输入用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}

	return true
}

func (client *Client) PublicChat() {

	var chatMsg string

	fmt.Println("请输入聊天内容, exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"

			_, err := client.conn.Write([]byte(sendMsg))

			if err != nil {
				fmt.Println("conn Write err: ", err)
				break
			}

			chatMsg = ""
			fmt.Println(">>>请输入聊天内容, exit退出")
			fmt.Scanln(&chatMsg)
		}
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err: ", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>请输入聊天对象[用户名], exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>请输入消息内容, exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"

				_, err := client.conn.Write([]byte(sendMsg))

				if err != nil {
					fmt.Println("conn Write err: ", err)
					break
				}

				chatMsg = ""
				fmt.Println(">>>请输入聊天内容, exit退出")
				fmt.Scanln(&chatMsg)
			}
		}

		client.SelectUsers()
		fmt.Println(">>>请输入聊天对象[用户名], exit退出:")
		fmt.Scanln(&remoteName)
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

var ServerIp string
var ServerPort int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "默认ip(127.0.0.1)")
	flag.IntVar(&ServerPort, "port", 8888, "默认端口(8888)")
}

func main() {
	flag.Parse()

	client := NewClient(ServerIp, ServerPort)
	if client == nil {
		fmt.Println(">>> 链接服务器失败...")
		return
	}

	go client.DealResponse()
	fmt.Println(">>> 链接服务器成功...")

	client.Run()
}
