package main

import (
	"fmt"

	"github.com/nikolaynag/go-hitstree/hitstree"
)

func main() {
	ht := hitstree.NewHitsTree()
	hitstree.MaxChildrenCnt = 10
	ht.HitPath("/")
	ht.HitPath("/users")
	ht.HitPath("/content/foo")
	ht.HitPath("/content/bar")
	ht.AddHitsToPath("/content/baz", 10, map[string]bool{"tag1": true})
	ht.HitPath("/content/baz")
	ht.AddHitsToPath("/content/baz", 1, map[string]bool{"tag2": true})
	for i := 0; i < 200; i++ {
		if i%10 != 0 {
			path := fmt.Sprintf("/users/%d/posts", i)
			ht.HitPath(path)
		} else {
			for j := 0; j < 20; j++ {
				path := fmt.Sprintf("/users/%d/posts/%d", i, j)
				ht.HitPath(path)
			}
		}
	}
	ht.HitPath("/users/0000/posts")
	ht.HitPath("/users/1234")
	ht.HitPath("/users/5678")
	fmt.Println(ht)
}
