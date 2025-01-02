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
	Next  *ListItem // право
	Prev  *ListItem // лево
}

type list struct {
	count int
	front *ListItem
	back  *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.count
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	newFront := &ListItem{
		Value: v,
		Next:  l.front,
		Prev:  nil,
	}
	if l.front == nil && l.back == nil {
		l.front = newFront
		l.back = newFront
	} else {
		l.front.Prev = newFront
		l.front = newFront
	}
	l.count++
	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	newBack := &ListItem{
		Value: v,
		Next:  nil,
		Prev:  l.back,
	}

	if l.front == nil && l.back == nil {
		l.front = newBack
		l.back = newBack
	} else {
		l.back.Next = newBack
		l.back = newBack
	}

	l.count++
	return l.back
}

func (l *list) Remove(i *ListItem) {
	if l.front == i {
		newFront := l.front.Prev
		l.front = newFront
		l.front.Prev = nil
		return
	}

	if l.back == i {
		newBack := l.back.Next
		l.back = newBack
		l.back.Next = nil
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	l.count--
}

func (l *list) MoveToFront(i *ListItem) {

	if l.front == i {
		return
	}

	if l.back == i {
		i.Prev.Next = nil
		l.back = i.Prev
		i.Prev = nil

		l.front.Prev = i
		i.Next = l.front
		l.front = i
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev

	l.front.Prev = i
	i.Next = l.front
	l.front = i
}
