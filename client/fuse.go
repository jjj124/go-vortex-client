package client

import (
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
	"github.com/hanwen/go-fuse/v2/fuse/pathfs"
	"time"
)

/**

../{pid}/{client_id}/
					├──state
					   |———metrics
					├──recent
                       |———recv
                       |———recv
*/

type Root struct {
	AdapterClient
}

func (r *Root) String() string {
	//TODO implement me
	panic("implement me")
}

func (r *Root) SetDebug(debug bool) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Chmod(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Chown(name string, uid uint32, gid uint32, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Truncate(name string, size uint64, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Access(name string, mode uint32, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Link(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Rename(oldName string, newName string, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Rmdir(name string, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Unlink(name string, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) GetXAttr(name string, attribute string, context *fuse.Context) (data []byte, code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) ListXAttr(name string, context *fuse.Context) (attributes []string, code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	//TODO implement me
	panic("implement me")
}

func (r *Root) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	//TODO implement me
	panic("implement me")
}

func (r *Root) OnMount(nodeFs *pathfs.PathNodeFs) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) OnUnmount() {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) OpenDir(name string, context *fuse.Context) (stream []fuse.DirEntry, code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Symlink(value string, linkName string, context *fuse.Context) (code fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	//TODO implement me
	panic("implement me")
}

func (r *Root) StatFs(name string) *fuse.StatfsOut {
	//TODO implement me
	panic("implement me")
}
