package volume

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/tweedproject/tweed/creds"
)

type credFile struct {
	fs.Inode

	data []byte
	root *credRoot
}

func (cf *credFile) resize(sz uint64) {
	if sz > uint64(cap(cf.data)) {
		n := make([]byte, sz)
		copy(n, cf.data)
		cf.data = n
	} else {
		cf.data = cf.data[:sz]
	}
}

func (cf *credFile) getattr(out *fuse.AttrOut) {
	out.Size = uint64(len(cf.data))
}

var _ = (fs.NodeGetattrer)((*credFile)(nil))

func (cf *credFile) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	cf.getattr(out)
	return 0
}

var _ = (fs.NodeSetattrer)((*credFile)(nil))

func (cf *credFile) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if sz, ok := in.GetSize(); ok {
		cf.resize(sz)
	}
	cf.getattr(out)
	return 0
}

var _ = (fs.NodeWriter)((*credFile)(nil))

func (cf *credFile) Write(ctx context.Context, fh fs.FileHandle, buf []byte, off int64) (uint32, syscall.Errno) {
	sz := int64(len(buf))
	if off+sz > int64(len(cf.data)) {
		cf.resize(uint64(off + sz))
	}
	copy(cf.data[off:], buf)
	cf.root.data[cf.Path(nil)] = cf.data
	return uint32(sz), cf.root.sync()
}

var _ = (fs.NodeOpener)((*credFile)(nil))

func (cf *credFile) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	return nil, 0, 0
}

var _ = (fs.NodeReader)((*credFile)(nil))

func (cf *credFile) Read(ctx context.Context, f fs.FileHandle, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	end := int(off) + len(dest)
	if end > len(cf.data) {
		end = len(cf.data)
	}
	return fuse.ReadResultData(cf.data[off:end]), fs.OK
}

type credRoot struct {
	fs.Inode
	secrets creds.Secrets
	secret  string
	data    map[string][]byte
}

var _ = (fs.NodeOnAdder)((*credRoot)(nil))

func createRoot(ctx context.Context, secrets creds.Secrets, secret string) (*credRoot, error) {
	root := &credRoot{
		secrets: secrets,
		secret:  secret,
	}
	data, exists, err := root.secrets.Get(root.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial secret volume data: %s", err)
	}
	if !exists {
		return root, nil
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to encode secret volume date %s", err)
	}

	var sdata map[string][]byte
	err = json.Unmarshal(raw, &sdata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret volume date %s", err)
	}
	root.data = sdata

	return root, nil
}

func (root *credRoot) OnAdd(ctx context.Context) {
	for f, data := range root.data {
		dir, base := filepath.Split(f)
		p := &root.Inode
		for _, component := range strings.Split(dir, "/") {
			if len(component) == 0 {
				continue
			}
			ch := p.GetChild(component)
			if ch == nil {
				ch = p.NewPersistentInode(ctx, &fs.Inode{},
					fs.StableAttr{Mode: fuse.S_IFDIR})
				p.AddChild(component, ch, true)
			}

			p = ch
		}
		ch := p.NewPersistentInode(ctx, &credFile{data: data, root: root}, fs.StableAttr{})
		p.AddChild(base, ch, true)
	}
}

func (root *credRoot) sync() syscall.Errno {
	err := root.secrets.Set(root.secret, root.data)
	if err != nil {
		fmt.Printf("Encountered an error while persisting secrets volume changes: %s", err)
		return 1
	}
	return 0
}
