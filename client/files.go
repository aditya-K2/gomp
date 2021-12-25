package client

import (
	"fmt"
	"strings"
)

type FileNode struct {
	Children     []FileNode
	Path         string
	Parent       *FileNode
	AbsolutePath string
}

// Source Interface For Fuzzy Searching.
type FileNodes []FileNode

func (f FileNodes) String(i int) string {
	if len(f[i].Children) == 0 {
		_s, err := CONN.ListAllInfo(f[i].AbsolutePath)
		if err != nil {
			NotificationServer.Send(fmt.Sprintf("Could Not Get Information About the Node %s", f[i].Path))
			return f[i].Path
		}
		return _s[0]["Title"]
	}
	return f[i].Path
}

func (f FileNodes) Len() int {
	return len(f)
}

func (f *FileNode) AddChildren(path string) {
	if f.Path != "" {
		f.Children = append(f.Children, FileNode{Children: make([]FileNode, 0), Path: path, Parent: f, AbsolutePath: f.AbsolutePath + "/" + path})
	} else {
		f.Children = append(f.Children, FileNode{Children: make([]FileNode, 0), Path: path, Parent: f, AbsolutePath: f.AbsolutePath + path})
	}
}

func (f *FileNode) AddChildNode(m FileNode) {
	m.Parent = f
	f.Children = append(f.Children, m)
}

func GenerateDirectoryTree(path []string) *FileNode {
	var head *FileNode = new(FileNode)
	var head1 *FileNode = head
	for i := range path {
		sepPaths := strings.Split(path[i], "/")
		for j := range sepPaths {
			if len(head.Children) == 0 {
				head.AddChildren(sepPaths[j])
				head = &(head.Children[len(head.Children)-1])
			} else {
				var headIsChanged = false
				for k := range head.Children {
					if head.Children[k].Path == sepPaths[j] {
						head = &(head.Children[k])
						headIsChanged = true
						break
					}
				}
				if !headIsChanged {
					head.AddChildren(sepPaths[j])
					head = &(head.Children[len(head.Children)-1])
				}
			}
		}
		head = head1
	}
	return head
}

func (f FileNode) Print(count int) {
	if len(f.Children) == 0 {
		return
	} else {
		for i := range f.Children {
			for j := 0; j < count; j++ {
				fmt.Print("---")
			}
			fmt.Println(f.Children[i].AbsolutePath)
			f.Children[i].Print(count + 1)
		}
	}
}
