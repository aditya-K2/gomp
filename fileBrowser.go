package main

import (
	"fmt"
	"strings"
)

type FileNode struct {
	children     []FileNode
	path         string
	parent       *FileNode
	absolutePath string
}

func (f *FileNode) addChildren(path string) {
	if f.path != "" {
		f.children = append(f.children, FileNode{children: make([]FileNode, 0), path: path, parent: f, absolutePath: f.absolutePath + "/" + path})
	} else {
		f.children = append(f.children, FileNode{children: make([]FileNode, 0), path: path, parent: f, absolutePath: f.absolutePath + path})
	}
}

func (f *FileNode) addChildNode(m FileNode) {
	m.parent = f
	f.children = append(f.children, m)
}

func generateDirectoryTree(path []string) *FileNode {
	var head *FileNode = new(FileNode)
	var head1 *FileNode = head
	for i := range path {
		sepPaths := strings.Split(path[i], "/")
		for j := range sepPaths {
			if len(head.children) == 0 {
				head.addChildren(sepPaths[j])
				head = &(head.children[len(head.children)-1])
			} else {
				var headIsChanged = false
				for k := range head.children {
					if head.children[k].path == sepPaths[j] {
						head = &(head.children[k])
						headIsChanged = true
						break
					}
				}
				if !headIsChanged {
					head.addChildren(sepPaths[j])
					head = &(head.children[len(head.children)-1])
				}
			}
		}
		head = head1
	}
	return head
}

func (f FileNode) Print(count int) {
	if len(f.children) == 0 {
		return
	} else {
		for i := range f.children {
			for j := 0; j < count; j++ {
				fmt.Print("---")
			}
			fmt.Println(f.children[i].absolutePath)
			f.children[i].Print(count + 1)
		}
	}
}
