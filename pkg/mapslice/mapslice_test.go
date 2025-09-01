package mapslice

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap2Slice(t *testing.T) {
	for _, ts := range []struct {
		name string
		m    map[int]string
		st   []int
		sv   []string
	}{
		{
			name: "happy path - map to slice",
			m:    map[int]string{1: "one", 2: "two"},
			st:   []int{1, 2},
			sv:   []string{"one", "two"},
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			st, sv := Map2Slice(ts.m)
			sort.Ints(st)
			sort.Strings(sv)
			assert.Equal(t, ts.st, st)
			assert.Equal(t, ts.sv, sv)
		})
	}
}

func TestSlice2Map(t *testing.T) {
	for _, ts := range []struct {
		name string
		s    []int
		m    map[int]struct{}
	}{
		{
			name: "happy path - map to slice",
			s:    []int{1, 2},
			m:    map[int]struct{}{1: {}, 2: {}},
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			st := Slice2Map(ts.s)
			assert.Equal(t, ts.m, st)
		})
	}
}

func TestMapKey2Slice(t *testing.T) {
	for _, ts := range []struct {
		name string
		m    map[int]struct{}
		sk   []int
	}{
		{
			name: "happy path - map key to slice",
			m:    map[int]struct{}{1: {}, 2: {}},
			sk:   []int{1, 2},
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			st := MapKey2Slice(ts.m)
			sort.Ints(st)
			assert.Equal(t, ts.sk, st)
		})
	}
}

func TestMapValue2Slice(t *testing.T) {
	for _, ts := range []struct {
		name string
		m    map[int]string
		sv   []string
	}{
		{
			name: "happy path - map value to slice",
			m:    map[int]string{1: "one", 2: "two"},
			sv:   []string{"one", "two"},
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			sv := MapValue2Slice(ts.m)
			sort.Strings(sv)
			assert.Equal(t, ts.sv, sv)
		})
	}
}
