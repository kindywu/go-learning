package main

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	id, err := inner_shorten("zIatZC", "www.abc.com")

	assert.Empty(t, id)
	assert.IsType(t, &pq.Error{}, err)
	var pqErr *pq.Error
	assert.ErrorAs(t, err, &pqErr)
	assert.Equal(t, pqErr.Code, pq.ErrorCode("23505"))
}
