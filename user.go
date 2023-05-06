package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		userAddr,
		userAddr,
		make(chan string),
		conn,
		server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()

	this.server.BroadCast(this, "已下线")
}

func (this *User) sendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {

	if msg == "who" {
		this.server.mapLock.Lock()
		for name := range this.server.OnlineMap {
			msg := "[" + name + "]" + "在线...\n"
			this.sendMsg(msg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {

		newName := msg[7:]

		if _, ok := this.server.OnlineMap[newName]; ok == true {
			this.sendMsg("该用户名已存在")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.sendMsg("您已经修改成功\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]

		if remoteName == "" {
			this.sendMsg("消息格式不正确，请使用to|username|msg格式发送私聊消息")
			return
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]

		if !ok {
			this.sendMsg("该用户不存在")
			return
		}

		remoteMsg := strings.Split(msg, "|")[2]

		if remoteMsg == "" {
			this.sendMsg("消息不能为空")
			return
		}

		remoteUser.sendMsg("[" + this.Name + "]" + "告诉你:" + remoteMsg + "\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听当前User channel方法，一旦有消息，就直接发送给看
func (this *User) ListenMessage() {
	for msg := range this.C {
		_, err := this.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println(err)
		}
	}
}
