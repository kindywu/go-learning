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

	if len(c.map_list) == 0 {
		c.head = node
		c.tail = node
	} else {
		node.next = c.head
		c.head.prev = node
		c.head = node
	}

	c.map_list[k] = node

	if len(c.map_list) > c.capacity {
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

	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if node == c.tail {
		c.tail = node.prev
	}

	node.prev = nil
	node.next = c.head
	c.head.prev = node
	c.head = node
}
