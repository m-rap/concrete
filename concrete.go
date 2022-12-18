package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

type Concrete struct {
	Nodes  map[string]*ConcreteNode
	Convos map[string]*ConcreteConvo
}

var convoLastId int = -1

const (
	ConvoDirect = 0
	ConvoGroup  = 1
)

func (c Concrete) CreateConvo(owner string, members []string, convoType int) {
	if convoType == ConvoDirect {
		for _, v := range c.Convos {
			if v.Type != ConvoDirect {
				continue
			}
			if (owner == v.Owner &&
				members[0] == v.Members[0]) ||
				(members[0] == v.Owner &&
					owner == v.Members[0]) {
				c.Nodes[owner].OnCreatedConvo(v)
				return
			}
		}
	}

	for _, v := range members {
		_, ok := c.Nodes[v]
		if !ok {
			return
		}
	}
	convoLastId++
	convoIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(convoIdBytes, uint64(convoLastId))
	convoId := base64.StdEncoding.EncodeToString(convoIdBytes)
	c.Convos[convoId] = &ConcreteConvo{
		Type:    convoType,
		Members: members,
		Owner:   owner,
	}

	pConvo := c.Convos[convoId]
	c.Nodes[owner].OnCreatedConvo(pConvo)
	for _, v := range members {
		c.Nodes[v].OnCreatedConvo(pConvo)
	}
}

func (c Concrete) Relay(convoId string, msg *ConcreteMsg) {
	convo, ok := c.Convos[convoId]
	if !ok {
		fmt.Println("server can't recognize convo")
		return
	}
	for _, member := range convo.Members {
		c.Nodes[member].Recv(convoId, msg)
	}
}

type ConcreteNode struct {
	Id     string
	Convos map[string]*ConcreteConvo
}

func (n ConcreteNode) Recv(convoId string, msg *ConcreteMsg) {
	switch msg.D1 {
	case 0:
		// respond nego
	case 1:
		convo, ok := n.Convos[convoId]
		if !ok {
			fmt.Println("not member of convo\n")
			return
		}
		msg2 := *msg
		convo.Msgs = append(convo.Msgs, &msg2)
	}
}

func (n ConcreteNode) OnCreatedConvo(convo *ConcreteConvo) {
	n.Convos[convo.Id] = &ConcreteConvo{
		Id:      convo.Id,
		Owner:   convo.Owner,
		Members: convo.Members,
	}
}

type ConcreteMsg struct {
	owner     string
	D1        int
	D2        string
	D3        bool
	D4        string
	D5        int
	Timestamp int
}

type ConcreteConvo struct {
	Id          string
	Owner       string
	Type        int
	Admins      []string
	Members     []string
	Msgs        []*ConcreteMsg
	CreatedTime int
}

func main() {
	fmt.Println("vim-go")
}
