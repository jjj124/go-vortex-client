package client

import (
	"context"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse/pathfs"
	"log"
	"testing"
	"time"
)

type RootNode struct {
	fs.Inode
	Name string
}

type PidNode struct {
	fs.Inode
	Name string
	pathfs.FileSystem
}

func (r *RootNode) OnAdd(ctx context.Context) {
	pid := &PidNode{}
	ch := r.NewInode(ctx, pid, fs.StableAttr{Ino: 2})
	r.AddChild("pid", ch, false)
}

func TestFuse(t *testing.T) {
	root := &RootNode{}
	opts := &fs.Options{}
	server, err := fs.Mount("/Users/caocg/adapter/xxx", root, opts)
	if err != nil {
		log.Fatalf("Failed to mount filesystem: %v", err)
	}
	// 启动文件系统服务
	server.Serve()
	//server.Wait()
	server.Unmount()

	time.Sleep(time.Hour * 1)

}
