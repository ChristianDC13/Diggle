package webqueue

import "testing"

func TestWebQueue(t *testing.T) {

	t.Run("Enqueue", func(t *testing.T) {
		wq := NewWebQueue()
		wq.Enqueue("http://example.com")
		if wq.head == nil {
			t.Errorf("Value: wq.head = got %v; want %v", wq.head, "http://example.com")
		}
		if wq.tail == nil {
			t.Errorf("Value: wq.tail = got %v; want %v", wq.tail, "http://example.com")
		}
		if wq.head.url != "http://example.com" {
			t.Errorf("Value: wq.head.url = got %v; want %v", wq.head.url, "http://example.com")
		}
		if wq.tail.url != "http://example.com" {
			t.Errorf("Value: wq.tail.url = got %v; want %v", wq.tail.url, "http://example.com")
		}
	})

	t.Run("Dequeue", func(t *testing.T) {
		wq := NewWebQueue()
		wq.Enqueue("http://example.com")
		url := wq.Dequeue()
		if url != "http://example.com" {
			t.Errorf("Value: wq.Dequeue() = got %v; want %v", url, "http://example.com")
		}
		if wq.head != nil {
			t.Errorf("Value: wq.head = got %v; want %v", wq.head, nil)
		}
		if wq.tail != nil {
			t.Errorf("Value: wq.tail = got %v; want %v", wq.tail, nil)
		}
	})

	t.Run("IsEmpty", func(t *testing.T) {
		wq := NewWebQueue()
		if !wq.IsEmpty() {
			t.Errorf("Value: wq.IsEmpty() = got %v; want %v", wq.IsEmpty(), true)
		}
		wq.Enqueue("http://example.com")
		if wq.IsEmpty() {
			t.Errorf("Value: wq.IsEmpty() = got %v; want %v", wq.IsEmpty(), false)
		}
		wq.Dequeue()
		if !wq.IsEmpty() {
			t.Errorf("Value: wq.IsEmpty() = got %v; want %v", wq.IsEmpty(), true)
		}
	})

}
