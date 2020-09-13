package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear logic", func(t *testing.T) {
		const nValues = 10
		c := NewCache(nValues)

		for i := 0; i < nValues; i++ {
			k := Key(strconv.Itoa(i))
			c.Set(k, i)
		}

		c.Clear()

		for i := 0; i < nValues; i++ {
			k := Key(strconv.Itoa(i))
			v, ok := c.Get(k)
			require.Nil(t, v)
			require.Equal(t, false, ok)
		}
	})

	t.Run("cache overflow", func(t *testing.T) {
		c := NewCache(3)

		c.Set("1", 1)
		c.Set("2", 2)
		c.Set("3", 3)
		c.Set("4", 4)

		v, ok := c.Get("1")
		require.Nil(t, v)
		require.False(t, ok)

		for i := 2; i <= 4; i++ {
			k := Key(strconv.Itoa(i))
			v, ok := c.Get(k)
			require.Equal(t, v, i)
			require.True(t, ok)
		}
	})

	t.Run("complex", func(t *testing.T) {
		const nValues = 100
		c := NewCache(nValues)

		for i := 0; i < nValues; i++ {
			k := Key(strconv.Itoa(i))
			c.Set(k, i)
		}

		for i := 0; i < nValues/4; i++ { // get first quarter of values
			k := Key(strconv.Itoa(i))
			c.Get(k)
		}

		for i := nValues / 4; i < nValues/2; i++ { // change values of second quarter
			k := Key(strconv.Itoa(i))
			c.Set(k, i+1)
		}

		const nNewValues = 10
		for i := 0; i < nNewValues; i++ { // add nNewValues values to cache
			k := Key(strconv.Itoa(nValues + i))
			c.Set(k, i)
		}

		for i := 0; i < nNewValues; i++ { // try to get nNewValues of third quarter
			k := Key(strconv.Itoa(i + nValues/2))
			v, ok := c.Get(k)
			require.Nil(t, v)
			require.False(t, ok)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
