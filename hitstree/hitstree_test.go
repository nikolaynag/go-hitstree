package hitstree_test

import (
	"strings"
	"testing"

	"github.com/nikolaynag/go-hitstree/hitstree"
	"github.com/stretchr/testify/assert"
)

func TestHitPath(t *testing.T) {
	hitstree.Delimiter = "/"
	ht := hitstree.NewHitsTree()
	ht.HitPath("/")
	ht.HitPath("/1/2/3")
	for i := 0; i < 100; i++ {
		ht.HitPath("/1/2/3/4")
	}
	ht.HitPath("/1/2/3")
	ht.HitPath("/test/")
	expected := map[string]int64{
		"/":        1,
		"/1/2/3":   2,
		"/1/2/3/4": 100,
		"/test":    1,
	}
	assert.Equal(t, expected, ht.HitsMap())
	expectedString := strings.Join([]string{
		"1\t/",
		"2\t/1/2/3",
		"100\t/1/2/3/4",
		"1\t/test",
	}, "\n")
	assert.Equal(t, expectedString, ht.String())
}

func TestMerge(t *testing.T) {
	hitstree.Delimiter = "/"
	a := hitstree.NewHitsTree()
	b := hitstree.NewHitsTree()
	a.HitPath("/a/1")
	a.HitPath("/common")
	a.HitPath("/common")
	a.HitPath("/common/a")
	b.HitPath("/b/1")
	b.HitPath("/common")
	b.HitPath("/common/b")
	expected := map[string]int64{
		"/a/1":      1,
		"/b/1":      1,
		"/common":   3,
		"/common/a": 1,
		"/common/b": 1,
	}
	a.Merge(b)
	assert.Equal(t, expected, a.HitsMap())
}

func TestMergeChildren(t *testing.T) {
	ht := hitstree.NewHitsTree()
	hitstree.MaxChildrenCnt = 3
	hitstree.Placeholder = "{}"
	hitstree.Delimiter = "/"
	ht.HitPath("/00/")
	ht.HitPath("/01/foo")
	ht.HitPath("/02/bar")
	ht.HitPath("/03/bar")
	ht.HitPath("/04/baz")
	ht.HitPath("/05/bar")
	ht.HitPath("/")
	expected := map[string]int64{
		"/":       1,
		"/{}":     1,
		"/{}/foo": 1,
		"/{}/bar": 3,
		"/{}/baz": 1,
	}
	assert.Equal(t, expected, ht.HitsMap())
}
