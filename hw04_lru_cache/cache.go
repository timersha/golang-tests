package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mx       sync.RWMutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type CacheItem struct {
	k Key
	v interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lruC *lruCache) cleanOutdatedElements() {
	if lruC.capacity < lruC.queue.Len() {
		delete(lruC.items, lruC.queue.Back().Value.(*CacheItem).k)
		lruC.queue.Remove(lruC.queue.Back())
	}
}

func (lruC *lruCache) Set(key Key, value interface{}) bool {
	lruC.mx.Lock()
	defer lruC.mx.Unlock()

	v, ok := lruC.items[key]
	if ok {
		v.Value = &CacheItem{key, value}
		lruC.queue.MoveToFront(v)
	} else {
		v = lruC.queue.PushFront(value)
		v.Value = &CacheItem{key, value}
		lruC.items[key] = v
		lruC.cleanOutdatedElements()
	}
	return ok
}

func (lruC *lruCache) Get(key Key) (interface{}, bool) {
	lruC.mx.RLock()
	defer lruC.mx.RUnlock()

	v, ok := lruC.items[key]
	if ok {
		lruC.queue.MoveToFront(v)
		return (v.Value).(*CacheItem).v, true
	}
	return nil, false
}

func (lruC *lruCache) Clear() {
	l := lruC.queue.Len()
	for range l {
		delete(lruC.items, lruC.queue.Front().Value.(*CacheItem).k)
		lruC.queue.Remove(lruC.queue.Front())
	}
}
