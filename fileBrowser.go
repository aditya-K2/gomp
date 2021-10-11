package main

import (
	"strings"
	"fmt"
)

type FileNode struct{
	children [] FileNode
	path string
	parent *FileNode
}

func (f *FileNode) addChildren(path string){
	f.children = append(f.children, FileNode{children: make([]FileNode, 0), path: path, parent: f})
}

func (f *FileNode) addChildNode(m FileNode){
	m.parent = f
	f.children = append(f.children, m)
}

func generateDirectoryTree(path [] string) *FileNode{
	var head *FileNode = new(FileNode)
	var head1 *FileNode = head
	for i := range(path){
		sepPaths := strings.Split(path[i], "/")
		for j := range(sepPaths){
			if(len(head.children) == 0){
				head.addChildren(sepPaths[j])
				head = &(head.children[len(head.children) - 1])
			} else {
				var headIsChanged = false
				for k := range(head.children){
					if head.children[k].path == sepPaths[j] {
						head = &(head.children[k])
						headIsChanged = true
						break
					}
				}
				if(!headIsChanged){
					head.addChildren(sepPaths[j])
					head = &(head.children[len(head.children) - 1])
				}
			}
		}
		head = head1
	}
	fmt.Println(head)
	return head
}
