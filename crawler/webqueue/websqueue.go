package webqueue

import "log"

type node struct {
	url  string
	next *node
}

type WebQueue struct {
	head   *node
	tail   *node
	memory map[string]bool
}

func (wq *WebQueue) Enqueue(url string) {
	n := &node{url: url}
	if wq.memory[url] {
		return
	}
	wq.memory[url] = true
	if wq.IsEmpty() {
		wq.head = n
		wq.tail = n
	} else {
		if wq.tail == nil {
			log.Println("tail is nil")
		}

		if n == nil {
			log.Println("n is nil")
		}

		wq.tail.next = n
		wq.tail = n
	}
}

func (wq *WebQueue) Dequeue() string {
	if wq.IsEmpty() {
		return ""
	}
	url := wq.head.url
	wq.head = wq.head.next
	if wq.head == nil {
		wq.tail = nil
	}
	return url
}

func (wq *WebQueue) IsEmpty() bool {
	return wq.head == nil
}

func NewWebQueue() *WebQueue {
	mem := make(map[string]bool)
	return &WebQueue{memory: mem}
}
