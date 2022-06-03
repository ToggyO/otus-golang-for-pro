package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem

	lock sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	listItem, ok := lru.items[key]
	if ok {
		listItem.Value.(*cacheItem).value = value
		lru.queue.MoveToFront(listItem)
		return true
	}

	cacheElem := &cacheItem{key: key, value: value}

	listItem = lru.queue.PushFront(cacheElem)
	if lru.queue.Len() > lru.capacity {
		itemToRemove := lru.queue.Back()
		lru.queue.Remove(itemToRemove)
		delete(lru.items, itemToRemove.Value.(*cacheItem).key)
	}

	lru.items[key] = listItem
	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.lock.RLock()
	defer lru.lock.RUnlock()

	listItem, ok := lru.items[key]
	if !ok {
		return nil, false
	}

	lru.queue.MoveToFront(listItem)
	return listItem.Value.(*cacheItem).value, true
}

func (lru *lruCache) Clear() {
	lru.lock.Lock()
	defer lru.lock.Unlock()

	// O(1) =))
	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
}
