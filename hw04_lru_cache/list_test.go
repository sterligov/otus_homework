package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

func TestFront(t *testing.T) {
	l := createListWithTwoValues(1, 2)
	require.Equal(t, 1, l.Front().Value)
}

func TestBack(t *testing.T) {
	l := createListWithTwoValues(1, 2)
	require.Equal(t, 2, l.Back().Value)
}

func TestPushBack(t *testing.T) {
	l := &list{}

	l.PushBack(1)
	require.Equal(t, 1, l.Head.Value)
	require.Equal(t, 1, l.Tail.Value)

	l.PushBack(2)
	require.Equal(t, 2, l.Tail.Value)
	require.Equal(t, 1, l.Head.Value)
}

func TestPushFront(t *testing.T) {
	l := &list{}

	l.PushFront(1)
	require.Equal(t, 1, l.Head.Value)
	require.Equal(t, 1, l.Tail.Value)

	l.PushFront(2)
	require.Equal(t, 2, l.Head.Value)
	require.Equal(t, 1, l.Tail.Value)
}

func TestRemove(t *testing.T) {
	l := createListWithTwoValues(1, 2)

	l.Remove(l.Tail)
	require.Nil(t, l.Head.Next)
	require.Same(t, l.Head, l.Tail)

	node := l.Head

	l.Remove(l.Head)
	require.Nil(t, l.Head)
	require.Nil(t, l.Tail)

	t.Run("remove removed node", func(t *testing.T) {
		l.Remove(node)
		require.Nil(t, l.Head)
		require.Nil(t, l.Tail)
	})
}

func TestMoveToFront(t *testing.T) {
	t.Run("tail node move to front", func(t *testing.T) {
		l := createListWithTwoValues(1, 2)
		l.MoveToFront(l.Tail)

		require.Equal(t, 2, l.Head.Value)
		require.Equal(t, 1, l.Tail.Value)
	})

	t.Run("nil node move to front", func(t *testing.T) {
		l := createListWithTwoValues(1, 2)
		l.MoveToFront(nil)

		require.Equal(t, 1, l.Head.Value)
		require.Equal(t, 2, l.Tail.Value)
	})
}

func createListWithTwoValues(headVal, tailVal int) *list {
	l := &list{}

	l.Head = &listItem{
		Value: headVal,
	}
	l.Tail = &listItem{
		Value: tailVal,
		Prev:  l.Head,
	}
	l.Head.Next = l.Tail

	return l
}
