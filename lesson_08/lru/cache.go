package lru

type node struct {
	key, value string
	prev, next *node
}

type lruCache struct {
	capacity int
	cache    map[string]*node
	head     *node
	tail     *node
}

type LruCache interface {
 Put(key, value string) 
 Get(key string) (string, bool)
}

func NewLruCache(capacity int) LruCache {
	c := &lruCache{
		capacity: capacity,
		cache:    make(map[string]*node),
		head:     &node{},
		tail:     &node{},
	}
	c.head.next = c.tail
	c.tail.prev = c.head
	return c
}

func (c *lruCache) Get(key string) (string, bool) {
	if node, exists := c.cache[key]; exists {
		c.moveToHead(node)
		return node.value, true
	}
	return "", false
}

func (c *lruCache) Put(key, value string) {
	if node, exists := c.cache[key]; exists {
		node.value = value
		c.moveToHead(node)
		return
	}

	newNode := &node{key: key, value: value}
	c.cache[key] = newNode
	c.addToHead(newNode)

	if len(c.cache) > c.capacity {
		tail := c.removeTail()
		delete(c.cache, tail.key)
	}
}

func (c *lruCache) addToHead(node *node) {
	node.prev = c.head
	node.next = c.head.next
	c.head.next.prev = node
	c.head.next = node
}

func (c *lruCache) removeNode(node *node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (c *lruCache) moveToHead(node *node) {
	c.removeNode(node)
	c.addToHead(node)
}

func (c *lruCache) removeTail() *node {
	lastNode := c.tail.prev
	c.removeNode(lastNode)
	return lastNode
}
