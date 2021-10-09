package main

import (
	"strings"
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

func mapGenerator(fileList [] string) map[string][]string {
	var a map[string][]string = make(map[string][]string)
	for i := range(fileList){
		a[strings.Split(fileList[i] , "/")[0]] = append(a[strings.Split(fileList[i] , "/")[0]], strings.Split(fileList[i] , "/")[1])
	}
	return a
}

func generateDirectoryTree(fileList [] string) *FileNode{
	var head *FileNode = new(FileNode)
	for i:=0; i<len(fileList) ; i++ {
		separatedPaths := strings.Split(fileList [i] , "/")
		var currentNode *FileNode = new(FileNode)
		tempNode := currentNode
		for  j := range(separatedPaths) {
			if j != (len(separatedPaths)-1) {
				tempNode.addChildren(separatedPaths[j])
				tempNode = &tempNode.children[0]
			} else{
				tempNode.addChildren(separatedPaths[j])
			}
		}
		head.addChildNode(*currentNode)
	}
	return head;
}
