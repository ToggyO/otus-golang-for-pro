package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.head == nil {
		l.tail = item
	} else {
		l.head.Prev = item
		item.Next = l.head
	}

	l.head = item
	l.len++

	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.tail == nil {
		l.head = item
	} else {
		item.Prev = l.tail
		l.tail.Next = item
	}

	l.tail = item
	l.len++

	return item
}

func (l *list) Remove(item *ListItem) {
	if l.len > 0 {
		if item == l.head {
			l.head = l.head.Next
		}

		if item == l.tail {
			l.tail = l.tail.Prev
		}
	}
	l.safelyAssignNeighborNodes(item)
	l.len--
}

func (l *list) MoveToFront(item *ListItem) {
	if item == l.head {
		return
	}

	// Remove item from its position
	l.safelyAssignNeighborNodes(item)

	// Update tail
	if item == l.tail {
		l.tail = item.Prev
	}

	// Set item as head
	item.Prev = nil
	item.Next = l.head
	l.head.Prev = item
	l.head = item
}

func (l *list) safelyAssignNeighborNodes(item *ListItem) {
	if item.Prev != nil {
		item.Prev.Next = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	}
}
