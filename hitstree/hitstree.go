package hitstree

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

var (
	// MaxChildCnt define maximum number of children before using template
	MaxChildCnt = 100
	// Placeholder is used as template path component
	Placeholder = "{}"
	// Delimiter is used to split path components in text representation
	Delimiter = "/"
)

// HitsTree stores number of hits for each member of hierarchical path set,
// like set of HTTP request paths for example. The key feature of HitsTree
// is an ability to detect templates at any level of hierarchy. When number of
// children at some level exceeds MaxChildCnt, all children are merged into
// one template child with Placeholder as path component.
type HitsTree struct {
	Hits     int64
	Children map[string]*HitsTree
}

// NewHitsTree returns initialized empty HitsTree
func NewHitsTree() *HitsTree {
	return &HitsTree{
		Hits:     0,
		Children: map[string]*HitsTree{},
	}
}

// Merge stores all hits and sub-paths of tm HitsTree into t HitsTree
func (t *HitsTree) Merge(tm *HitsTree) {
	left := t
	right := tm
	left.Hits += right.Hits
	for key, rightChild := range right.Children {
		leftChild, ok := left.Children[key]
		if ok {
			leftChild.Merge(rightChild)
		} else {
			left.Children[key] = rightChild
		}
	}
}

// MergeChildren merges all 1st-level children of tree into one template
// child
func (t *HitsTree) MergeChildren() (mergedChild *HitsTree) {
	for key, child := range t.Children {
		if mergedChild == nil {
			mergedChild = child
			delete(t.Children, key)
			continue
		}
		mergedChild.Merge(child)
		delete(t.Children, key)
	}
	if mergedChild != nil {
		t.Children[Placeholder] = mergedChild
	}
	return
}

// Hit stores one hit to path represented by given set of components
func (t *HitsTree) Hit(components []string) {
	current := t
	for _, component := range components {
		next, ok := current.Children[Placeholder]
		if ok {
			current = next
			continue
		}
		next, ok = current.Children[component]
		if ok {
			current = next
			continue
		}
		if len(current.Children) >= MaxChildCnt {
			next = current.MergeChildren()
		} else {
			next = &HitsTree{
				Hits:     0,
				Children: map[string]*HitsTree{},
			}
			current.Children[component] = next
		}
		current = next
	}
	current.Hits++
}

// HitPath adds one hit to path represented by string p
func (t *HitsTree) HitPath(p string) {
	trimmedPath := strings.Trim(path.Clean(p), Delimiter)
	var pathComponents []string
	if trimmedPath != "" {
		pathComponents = strings.Split(trimmedPath, Delimiter)
	} else {
		pathComponents = []string{}
	}
	t.Hit(pathComponents)
}

// outputToMap updates hitsByPath map by adding number of hits for each
// path
func (t *HitsTree) outputToMap(parentPath string, hitsByPath map[string]int64) {
	if t.Hits > 0 {
		currPath := strings.TrimRight(parentPath, Delimiter)
		if currPath == "" {
			currPath = Delimiter
		}
		hitsByPath[currPath] += t.Hits
	}
	for key, child := range t.Children {
		child.outputToMap(path.Join(parentPath, key), hitsByPath)
	}
}

// HitsMap return map where path is key and hits count is value
func (t *HitsTree) HitsMap() (hitsMap map[string]int64) {
	hitsMap = make(map[string]int64)
	t.outputToMap(Delimiter, hitsMap)
	return
}

// String returns string representation of t
func (t *HitsTree) String() (result string) {
	hitsByPath := t.HitsMap()
	pathList := []string{}
	for p := range hitsByPath {
		pathList = append(pathList, p)
	}
	sort.Strings(pathList)
	lines := make([]string, len(pathList))
	for i, p := range pathList {
		lines[i] = fmt.Sprintf("%d\t%s", hitsByPath[p], p)
	}
	return strings.Join(lines, "\n")
}
