package algorithm

type Node[V any] struct {
	value V
	prev  *Node[V]
	next  *Node[V]
}

type LuaCache[K comparable, V any] struct {
	capacity int
	head     *Node[V]
	map_list map[K]*Node[V]
}

func newLuaCache[K comparable, V any](capacity int) LuaCache[K, V] {
	var defaultValue V
	head := Node[V]{value: defaultValue, prev: nil, next: nil}
	return LuaCache[K, V]{capacity: capacity, head: &head, map_list: make(map[K]*Node[V])}
}

func (c *LuaCache[K, V]) put(k K, v V) *Node[V] {
	if node, exists := c.map_list[k]; exists {
		node.value = v
		c.moveToHead(node)
		return nil
	} else {
		node := Node[V]{value: v, prev: c.head, next: c.head.next}
		if c.head.next != nil {
			c.head.next.prev = &node
		}
		c.head.next = &node
		c.map_list[k] = &node

		if c.head.prev == nil {
			c.head.prev = &node
		}

		var deletedNode *Node[V] = nil
		if len(c.map_list) > c.capacity {
			deletedNode = c.head.prev
			c.head.prev = deletedNode.prev
			c.head.prev.next = nil
			deletedNode.prev = nil
		}
		return deletedNode
	}
}

func (c *LuaCache[K, V]) get(k K) (V, bool) {
	if node, exists := c.map_list[k]; exists {
		c.moveToHead(node)
		return node.value, exists
	} else {
		var defaultValue V
		return defaultValue, exists
	}
}

func (c *LuaCache[K, V]) info() (int, int) {
	return len(c.map_list), c.capacity
}

func (c *LuaCache[K, V]) moveToHead(node *Node[V]) {
	if node.prev != c.head {

		if node.next == nil {
			c.head.prev = node.prev
		}

		node.prev.next = node.next

		if node.next != nil {
			node.next.prev = node.prev
		}

		node.next = c.head.next
		c.head.next.prev = node

		c.head.next = node
		node.prev = c.head

	}
}
