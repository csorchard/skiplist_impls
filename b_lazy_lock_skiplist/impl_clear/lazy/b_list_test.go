package lazyskiplist

import (
	"reflect"
	"skiplist/b_lazy_lock_skiplist/impl_clear/lib"
	"testing"
)

func TestPut(t *testing.T) {
	list := NewLazySkipList(lib.IntComparator)
	list.Put(1, "test", nil)
	value, _ := list.Get(1)
	if value != "test" {
		t.Errorf("Expected: %v, Got: %v", "test", value)
	}
}

func TestIterator(t *testing.T) {
	list := NewLazySkipList(lib.IntComparator)

	list.Put(1, nil, nil)
	list.Put(3, nil, nil)
	list.Put(5, nil, nil)
	list.Put(7, nil, nil)

	var slice []int
	for it := list.Begin(nil); it.Present(); it.Next() {
		slice = append(slice, it.Key().(int))
	}
	expected := []int{1, 3, 5, 7}
	if !reflect.DeepEqual(slice, expected) {
		t.Errorf("Expected: %v, Got: %v", expected, slice)
	}
}

func TestReverseIterator(t *testing.T) {
	list := NewLazySkipList(lib.IntComparator)

	list.Put(1, nil, nil)
	list.Put(3, nil, nil)
	list.Put(5, nil, nil)
	list.Put(7, nil, nil)

	var slice []int
	for it := list.End(nil); it.Present(); it.Prev() {
		slice = append(slice, it.Key().(int))
	}
	expected := []int{7, 5, 3, 1}
	if !reflect.DeepEqual(expected, slice) {
		t.Errorf("Expected: %v, Got: %v", expected, slice)
	}
}

func TestIteratorWithDeletedItem(t *testing.T) {
	list := NewLazySkipList(lib.IntComparator)

	list.Put(1, nil, nil)
	list.Put(3, nil, nil)
	list.Put(5, nil, nil)
	list.Put(7, nil, nil)

	list.Print()
	it := list.Begin(nil)

	list.Remove(1)
	list.Remove(5)

	var slice []int
	var marked []int
	for ; it.Present(); it.Next() {
		slice = append(slice, it.Key().(int))
		if it.IsMarked() {
			marked = append(marked, it.Key().(int))
		}
	}
	{
		expected := []int{1, 3, 7}
		if !reflect.DeepEqual(expected, slice) {
			t.Errorf("Expected: %v, Got: %v", expected, slice)
		}
	}
	{
		expected := []int{1}
		if !reflect.DeepEqual(marked, expected) {
			t.Errorf("Expected: %v, Got: %v", expected, marked)
		}
	}
}
