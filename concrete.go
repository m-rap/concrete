package main

import (
	"fmt"
	"strconv"
)

type Concrete struct {
	Nodes  map[string]ConcreteNode
	Convos map[string]ConcreteConvo
}

var convoLastId int = -1

const (
	ConvoDirect = 0
	ConvoGroup  = 1
)

func (Concrete c) CreateConvo(owner ConcreteNode, members []ConcreteNode, convoType int) {
	if convoType == ConvoDirect {
		for _, v := range c.Convos {
			if v.Type != ConvoDirect {
				continue
			}
			if (owner.Id == v.Owner.Id &&
				members[0].Id == v.Members[0].Id) ||
				(members[0].Id == v.Owner.Id &&
					owner.Id == v.Members[0].Id) {
				owner.OnCreatedConvo(&v)
				return
			}
		}
	}

	for _, v := range members {
		sm, ok := c.Nodes
		if !ok {
			return
		}
	}
	convoLastId++
	convoId := strconv.Itoa(convoLastId)
	c.Convos[convoId] = ConcreteConvo{}
	pConvo := &c.Convos[convoId]
	pConvo.Type = convoType
	pConvo.Members = members
	pConvo.Owner = owner
	owner.OnCreatedConvo(pConvo)
	for _, v := range members {
		v.OnCreatedConvo(pConvo)
	}
}

func (Concrete c) Relay(convoId string, msg ConcreteMsg) {
	convo, ok := c.Convos[convoId]
	if !ok {
		fmt.Println("server can't recognize convo")
		return
	}
	for i, member := range convo.Members {
		member.Recv(convoId, msg)
	}
}

type ConcreteNode struct {
	Id     string
	Convos map[string]ConcreteConvo
}

func (ConcreteNode n) Recv(convoId string, msg ConcreteMsg) {
	switch msg.D1 {
	case 0:
		// respond nego
	case 1:
		convo, ok := n.Convos[convoId]
		if !ok {
			fmt.Println("not member of convo\n")
			return
		}
		convo.Msgs[convoId] = msg
	}
}

func (ConcreteNode n) OnCreatedConvo(convo *ConcreteConvo) {
	n.Convos[convoId] = ConcreteConvo{
		Id:      convo.Id,
		Owner:   convo.Owner,
		Members: convo.Members,
	}
}

type ConcreteMsg struct {
	owner     string
	D1        string
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
	Msgs        []ConcreteMsg
	CreatedTime int
}

func main() {
	fmt.Println("vim-go")
}
