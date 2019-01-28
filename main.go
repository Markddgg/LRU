package main

import (
	"container/list"
	"fmt"
)

func main() {
	cache := NewLRUCache(5)

	arr := []int{5, 3, 5, 2, 1, 4, 5, 4, 3, 7, 8, 8, 9, 3, 1, 8, 7, 3, 9, 7}

	for _, v := range arr {
		cache.Set(v, v)
		fmt.Printf("%d \t", v)

		e := cache.List.Front()
		for i := 0; i < cache.List.Len(); i++ {
			fmt.Printf("%v ", e.Value.(*CacheNode).Value)
			e = e.Next()
		}

		fmt.Printf("\n")
	}
}

type CacheNode struct {
	Key, Value interface{}
}

func (node *CacheNode) NewCacheNode(k, v interface{}) *CacheNode {
	return &CacheNode{k, v}
}

type LRUCache struct {
	Cap     int
	List    *list.List
	HashMap map[interface{}]*list.Element
}

func NewLRUCache(len int) *LRUCache {
	return &LRUCache{len, list.New(), make(map[interface{}]*list.Element)}
}

func (cache LRUCache) Set(k, v interface{}) error {

	if cache.List == nil {
		return fmt.Errorf("uninitialized cache")
	}

	// 1 如果命中缓存，被命中的缓存移动到list头部(最近访问)
	if e, ok := cache.HashMap[k]; ok {
		e.Value.(*CacheNode).Value = v
		cache.List.MoveToFront(e)

		return nil
	}

	// 2 如果缓存未命中
	// 2.1 生产新CacheNode放入List头部及HashMap
	e := cache.List.PushFront(&CacheNode{k, v})
	cache.HashMap[k] = e

	// 2.2 如果超出缓存cap,淘汰最近最少访问页面
	if cache.List.Len() > cache.Cap {
		e := cache.List.Back()
		if e == nil {
			return nil
		}

		node := e.Value.(*CacheNode)
		fmt.Printf("remove from cache, k: %v, v: %v \n", node.Key, node.Value)
		delete(cache.HashMap, node.Key)
		cache.List.Remove(e)
	}

	return nil
}

func (cache LRUCache) Get(k interface{}) (interface{}, int, error) {

	if cache.List == nil {
		return nil, 0, fmt.Errorf("uninitialized cache")
	}

	if e, ok := cache.HashMap[k]; ok {
		cache.List.MoveToFront(e)
		return e.Value.(*CacheNode).Value, 1, nil
	}

	return nil, 0, nil
}

func (cache *LRUCache) Remove(k interface{}) error {

	if cache.List == nil {
		return fmt.Errorf("uninitialized cache")
	}

	if e, ok := cache.HashMap[k]; ok {
		delete(cache.HashMap, k)
		cache.List.Remove(e)
	}

	return nil
}
