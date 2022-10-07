package client

import (
	"fmt"
	"strings"

	"github.com/fhs/gompd/v2/mpd"
)

type FileNode struct {
	Children     []FileNode
	Path         string
	Parent       *FileNode
	AbsolutePath string
	Title        string
	Artist       string
	Album        string
}

// Source Interface For Fuzzy Searching.
type FileNodes []FileNode

func (f FileNodes) String(i int) string {
	if len(f[i].Children) == 0 {
		return f[i].Title
	}
	return f[i].Path
}

func (f FileNodes) Len() int {
	return len(f)
}

func (f *FileNode) AddChildren(path string, title string, artist string, album string) {
	if f.Path != "" {
		f.Children = append(f.Children, FileNode{Children: make([]FileNode, 0), Path: path, Parent: f, AbsolutePath: f.AbsolutePath + "/" + path, Title: title, Artist: artist, Album: album})
	} else {
		f.Children = append(f.Children, FileNode{Children: make([]FileNode, 0), Path: path, Parent: f, AbsolutePath: f.AbsolutePath + path})
	}
}

func (f *FileNode) AddChildNode(m FileNode) {
	m.Parent = f
	f.Children = append(f.Children, m)
}

func GenerateDirectoryTree(path []mpd.Attrs) *FileNode {
	var head *FileNode = new(FileNode)
	var head1 *FileNode = head
	for i := range path {
		sepPaths := strings.Split(path[i]["file"], "/")
		for j := range sepPaths {
			if len(head.Children) == 0 {
				head.AddChildren(sepPaths[j], path[i]["Title"], path[i]["Artist"], path[i]["Album"])
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
					head.AddChildren(sepPaths[j], path[i]["Title"], path[i]["Artist"], path[i]["Album"])
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
