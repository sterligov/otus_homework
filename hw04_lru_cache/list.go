package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(interface{}) *listItem
	PushBack(interface{}) *listItem
	Remove(*listItem)
	MoveToFront(*listItem)
}

type listItem struct {
	Value interface{}
	Next  *listItem
	Prev  *listItem
}

type list struct {
	Capacity int
	Head     *listItem
	Tail     *listItem
}

func (l *list) Len() int {
	return l.Capacity
}

func (l *list) MoveToFront(item *listItem) {
	if item == nil {
		return
	}

	l.Remove(item)

	l.Capacity++
	item.Prev = nil
	item.Next = l.Head

	if l.Head != nil {
		l.Head.Prev = item
	} else {
		l.Tail = item
	}
	l.Head = item
}

func (l *list) Front() *listItem {
	return l.Head
}

func (l *list) Back() *listItem {
	return l.Tail
}

func (l *list) PushBack(v interface{}) *listItem {
	node := &listItem{
		Value: v,
		Prev:  l.Tail,
	}

	if l.Tail != nil {
		l.Tail.Next = node
	} else {
		l.Head = node
	}

	l.Tail = node
	l.Capacity++

	return l.Tail
}

func (l *list) PushFront(v interface{}) *listItem {
	node := &listItem{
		Value: v,
		Next:  l.Head,
	}

	if l.Head != nil {
		l.Head.Prev = node
	} else {
		l.Tail = node
	}

	l.Head = node
	l.Capacity++

	return l.Head
}

func (l *list) Remove(item *listItem) {
	if item == nil {
		return
	}

	l.Capacity--

	if item.Prev != nil {
		item.Prev.Next = item.Next
	} else {
		l.Head = item.Next
	}

	if item.Next != nil {
		item.Next.Prev = item.Prev
	} else {
		l.Tail = item.Prev
	}
}

func NewList() List {
	return &list{}
}
