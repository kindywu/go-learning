package algorithm

type Node[K comparable, V any] struct {
	key   K
	value V
	prev  *Node[K, V]
	next  *Node[K, V]
}

type LuaCache[K comparable, V any] struct {
	capacity int
	head     *Node[K, V]
	tail     *Node[K, V]
	map_list map[K]*Node[K, V]
}

func newLuaCache[K comparable, V any](capacity int) LuaCache[K, V] {
	return LuaCache[K, V]{
		capacity: capacity,
		head:     nil,
		tail:     nil,
		map_list: make(map[K]*Node[K, V]),
	}
}

func (c *LuaCache[K, V]) put(k K, v V) *Node[K, V] {
	if node, exists := c.map_list[k]; exists {
		node.value = v
		c.moveToHead(node)
		return nil
	}

	node := &Node[K, V]{key: k, value: v}
	c.map_list[k] = node
	c.addNodeToHead(node)

	if len(c.map_list) > c.capacity {
		evicted := c.removeTailNode()
		return evicted
	}

	return nil
}

func (c *LuaCache[K, V]) get(k K) (V, bool) {
	if node, exists := c.map_list[k]; exists {
		c.moveToHead(node)
		return node.value, exists
	}
	var defaultValue V
	return defaultValue, false
}

func (c *LuaCache[K, V]) remove(k K) bool {
	if node, exists := c.map_list[k]; exists {
		c.removeNode(node)
		delete(c.map_list, k)
		return true
	}
	return false
}

func (c *LuaCache[K, V]) info() (int, int) {
	return len(c.map_list), c.capacity
}

func (c *LuaCache[K, V]) moveToHead(node *Node[K, V]) {
	if c.head == node {
		return
	}
	c.removeNode(node)
	c.addNodeToHead(node)
}

func (c *LuaCache[K, V]) addNodeToHead(node *Node[K, V]) {
	node.next = c.head
	node.prev = nil

	if c.head != nil {
		c.head.prev = node
	}

	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func (c *LuaCache[K, V]) removeNode(node *Node[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		c.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		c.tail = node.prev
	}
}

func (c *LuaCache[K, V]) removeTailNode() *Node[K, V] {
	if c.tail == nil {
		return nil
	}

	evicted := c.tail
	if c.tail.prev != nil {
		c.tail = c.tail.prev
		c.tail.next = nil
	} else {
		c.head = nil
		c.tail = nil
	}

	delete(c.map_list, evicted.key)
	return evicted
}
