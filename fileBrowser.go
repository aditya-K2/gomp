package main

type FileNode struct{
	children [] FileNode
	path string
	parent *FileNode
}

func (f *FileNode) addChildren(path string){
	f.children = append(f.children, FileNode{children: make([]FileNode, 0), path: path, parent: f})
}
