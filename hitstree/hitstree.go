package hitstree

import (
	"fmt"
	"path"
	"sort"
	"strings"
)

var (
	// MaxChildrenCnt define maximum number of children before using template
	MaxChildrenCnt = 100
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
	Tags     map[string]bool
	Children map[string]*HitsTree
}

// NewHitsTree returns initialized empty HitsTree
func NewHitsTree() *HitsTree {
	return &HitsTree{
		Hits:     0,
		Tags:     map[string]bool{},
		Children: map[string]*HitsTree{},
	}
}

func mergeTags(dst, src map[string]bool) {
	for k, v := range src {
		if v {
			dst[k] = true
		}
	}
}

// Merge stores all hits and sub-paths of tm HitsTree into t HitsTree
func (t *HitsTree) Merge(tm *HitsTree) {
	left := t
	right := tm
	left.Hits += right.Hits
	mergeTags(left.Tags, right.Tags)
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
	t.AddHits(components, 1, nil)
}

// AddHits appends given number of tagged hits to path represented by given set of components
func (t *HitsTree) AddHits(components []string, numHits int64, tags map[string]bool) {
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
		if len(current.Children) >= MaxChildrenCnt {
			next = current.MergeChildren()
		} else {
			next = NewHitsTree()
			current.Children[component] = next
		}
		current = next
	}
	current.Hits += numHits
	mergeTags(current.Tags, tags)
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

// AddHitsToPath adds one hit to path represented by string p
func (t *HitsTree) AddHitsToPath(p string, numHits int64, tags map[string]bool) {
	trimmedPath := strings.Trim(path.Clean(p), Delimiter)
	var pathComponents []string
	if trimmedPath != "" {
		pathComponents = strings.Split(trimmedPath, Delimiter)
	} else {
		pathComponents = []string{}
	}
	t.AddHits(pathComponents, numHits, tags)
}

// outputToMap updates hitsByPath map by adding number of hits for each
// path
func (t *HitsTree) outputToMaps(parentPath string, hitsByPath map[string]int64, tagsByPath map[string]map[string]bool) {
	currPath := strings.TrimRight(parentPath, Delimiter)
	if currPath == "" {
		currPath = Delimiter
	}
	if t.Hits > 0 && hitsByPath != nil {
		hitsByPath[currPath] += t.Hits
	}
	if t.Tags != nil && tagsByPath != nil && len(t.Tags) > 0 {
		tags := tagsByPath[currPath]
		if tags == nil {
			tags = make(map[string]bool)
			tagsByPath[currPath] = tags
		}
		mergeTags(tags, t.Tags)
	}
	for key, child := range t.Children {
		child.outputToMaps(path.Join(parentPath, key), hitsByPath, tagsByPath)
	}
}

// HitsMap return map where path is key and hits count is value
func (t *HitsTree) HitsMap() (hitsMap map[string]int64) {
	hitsMap = make(map[string]int64)
	t.outputToMaps(Delimiter, hitsMap, nil)
	return
}

// TagsMap return map where path is key and combined tags set is value
func (t *HitsTree) TagsMap() (tagsMap map[string]map[string]bool) {
	tagsMap = make(map[string]map[string]bool)
	t.outputToMaps(Delimiter, nil, tagsMap)
	return
}

// String returns string representation of t
func (t *HitsTree) String() (result string) {
	hitsByPath := t.HitsMap()
	tagsByPath := t.TagsMap()
	pathList := []string{}
	for p := range hitsByPath {
		pathList = append(pathList, p)
	}
	sort.Strings(pathList)
	lines := make([]string, len(pathList))
	for i, p := range pathList {
		tags := make([]string, 0, len(tagsByPath[p]))
		for tag := range tagsByPath[p] {
			tags = append(tags, tag)
		}
		lines[i] = fmt.Sprintf("%d\t%s\t%s", hitsByPath[p], p, strings.Join(tags, ", "))
	}
	return strings.Join(lines, "\n")
}
