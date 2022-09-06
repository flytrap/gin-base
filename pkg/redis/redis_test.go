package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const addr = "127.0.0.1:6379"

func TestStore(t *testing.T) {
	store := NewStore(&Config{
		Addr:      addr,
		DB:        1,
		KeyPrefix: "prefix",
	})

	defer store.Close()

	key := "test"
	err := store.Set(key, 0, 5)
	assert.Nil(t, err)

	b, err := store.Check(key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)

	b, err = store.Delete(key)
	assert.Nil(t, err)
	assert.Equal(t, true, b)
}
