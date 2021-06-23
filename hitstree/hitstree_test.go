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
		"1\t/\t",
		"2\t/1/2/3\t",
		"100\t/1/2/3/4\t",
		"1\t/test\t",
	}, "\n")
	assert.Equal(t, expectedString, ht.String())
}

func TestMerge(t *testing.T) {
	hitstree.Delimiter = "/"
	a := hitstree.NewHitsTree()
	b := hitstree.NewHitsTree()
	a.HitPath("/a/1")
	a.AddHitsToPath("/common", 1, map[string]bool{"tag1": true})
	a.HitPath("/common")
	a.AddHitsToPath("/common/a", 1, map[string]bool{"tag0": true})
	b.HitPath("/b/1")
	b.AddHitsToPath("/common", 1, map[string]bool{"tag2": true})
	b.HitPath("/common/b")
	expectedHitsMap := map[string]int64{
		"/a/1":      1,
		"/b/1":      1,
		"/common":   3,
		"/common/a": 1,
		"/common/b": 1,
	}
	expectedTagsMap := map[string]map[string]bool{
		"/common":   map[string]bool{"tag1": true, "tag2": true},
		"/common/a": map[string]bool{"tag0": true},
	}
	a.Merge(b)
	assert.Equal(t, expectedHitsMap, a.HitsMap())
	assert.Equal(t, expectedTagsMap, a.TagsMap())
}

func TestMergeChildren(t *testing.T) {
	ht := hitstree.NewHitsTree()
	hitstree.MaxChildrenCnt = 3
	hitstree.Placeholder = "{}"
	hitstree.Delimiter = "/"
	ht.HitPath("/00/")
	ht.HitPath("/01/foo")
	ht.AddHitsToPath("/02/bar", 1, map[string]bool{"bar": true})
	ht.HitPath("/03/bar")
	ht.HitPath("/04/baz")
	ht.AddHitsToPath("/05/bar", 2, map[string]bool{"05": true})
	ht.HitPath("/")
	expectedHits := map[string]int64{
		"/":       1,
		"/{}":     1,
		"/{}/foo": 1,
		"/{}/bar": 4,
		"/{}/baz": 1,
	}
	expectedTags := map[string]map[string]bool{
		"/{}/bar": map[string]bool{"bar": true, "05": true},
	}
	assert.Equal(t, expectedHits, ht.HitsMap())
	assert.Equal(t, expectedTags, ht.TagsMap())
}

func TestZeroMaxChildren(t *testing.T) {
	ht := hitstree.NewHitsTree()
	hitstree.MaxChildrenCnt = 0
	hitstree.Placeholder = "{}"
	hitstree.Delimiter = "/"
	ht.HitPath("/0")
	ht.HitPath("/0/1")
	ht.HitPath("/0/1")
	ht.HitPath("/0/2")
	expected := map[string]int64{
		"/0":   1,
		"/0/1": 2,
		"/0/2": 1,
	}
	assert.Equal(t, expected, ht.HitsMap())
}
