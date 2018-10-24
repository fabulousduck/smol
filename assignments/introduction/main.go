package main

import "fmt"

type List struct {
	Head, Tail *Node
}

type Node struct {
	Prev *Node
	Data int
	Next *Node
}

func main() {
	list := NewList()
	//create the list
	for i := 0; i < 10; i++ {
		node := &Node{Prev: nil, Data: i, Next: nil}
		list.PushRear(node)
	}

	//print the list
	for node := list.Head; node != nil; node = node.Next {
		fmt.Println(node.Data)
	}
}

func NewList() *List {
	list := new(List)
	list.Head = nil
	list.Tail = nil
	return list

}

//pointers persist beyond the runtime of the function, so we can do this.
//lmao
func (list *List) PushRear(node *Node) {
	if list.Head == nil {
		list.Head = node
		list.Tail = list.Head
		return
	}
	list.Tail.Next = node
	node.Prev = list.Tail
	list.Tail = node
}
