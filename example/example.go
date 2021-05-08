package main

import (
	"fmt"

	"github.com/nikolaynag/go-hitstree/hitstree"
)

func main() {
	ht := hitstree.NewHitsTree()
	ht.HitPath("/")
	ht.HitPath("/a")
	ht.HitPath("/b/foo")
	ht.HitPath("/b/bar")
	ht.HitPath("/b/baz")
	for i := 0; i < 1000; i++ {
		path := fmt.Sprintf("/a/%d/x", i)
		ht.HitPath(path)
	}
	ht.HitPath("/a/0000/y")
	ht.HitPath("/a/1234")
	ht.HitPath("/a/5678")
	fmt.Println(ht)
}
