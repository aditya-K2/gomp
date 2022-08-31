package globals

import (
	"github.com/aditya-K2/fuzzy"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/render"
	"github.com/aditya-K2/gomp/ui"
	"github.com/fhs/gompd/mpd"
)

var (
	Conn               *mpd.Client
	Notify             *notify.NotificationServer
	Renderer           *render.Renderer
	Ui                 *ui.Application
	DirTree            *client.FileNode
	SearchContentSlice []interface{}
	Matches            fuzzy.Matches
)
