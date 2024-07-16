package algorithm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLuaCache(t *testing.T) {
	cache := newLuaCache[int, string](10)

	cache.put(1, "abc")
	v, exist := cache.get(1)
	len, capacity := cache.info()
	assert.Equal(t, true, exist)
	assert.Equal(t, "abc", v)
	assert.Equal(t, 1, len)
	assert.Equal(t, 10, capacity)

	cache.put(1, "efg")
	v, exist = cache.get(1)
	len, capacity = cache.info()
	assert.Equal(t, true, exist)
	assert.Equal(t, "efg", v)
	assert.Equal(t, 1, len)
	assert.Equal(t, 10, capacity)

	cache.put(2, "abc")
	v, exist = cache.get(2)
	len, capacity = cache.info()
	assert.Equal(t, true, exist)
	assert.Equal(t, "abc", v)
	assert.Equal(t, 2, len)
	assert.Equal(t, 10, capacity)
}

func TestLuaCache2(t *testing.T) {
	cache := newLuaCache[int, int](3)
	cache.put(1, 1)
	cache.put(2, 2)
	cache.put(3, 3)
	cache.put(1, 1)

	deleted_node := cache.put(4, 4)
	assert.NotNil(t, deleted_node)
	assert.Equal(t, 2, deleted_node.value)

	deleted_node = cache.put(5, 5)
	assert.NotNil(t, deleted_node)
	assert.Equal(t, 3, deleted_node.value)

	deleted_node = cache.put(6, 6)
	assert.NotNil(t, deleted_node)
	assert.Equal(t, 1, deleted_node.value)

	assert.Equal(t, false, cache.remove(1))
	assert.Equal(t, true, cache.remove(6))

}
