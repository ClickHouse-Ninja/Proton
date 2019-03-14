package server

import (
	"sort"
	"testing"
)

func TestTagSort(t *testing.T) {
	name, value := []string{"b name", "c name", "a name"}, []string{"b value", "c value", "a value"}
	sort.Sort(&tagSort{
		name:  name,
		value: value,
	})
	t.Log(name, value)
}
