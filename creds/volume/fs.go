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
}

var _ = (fs.NodeOpener)((*credFile)(nil))

var _ = (fs.NodeGetattrer)((*credFile)(nil))

func (cf *credFile) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	out.Size = uint64(len(cf.data))
	return 0
}

func (cf *credFile) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	return nil, fuse.FOPEN_KEEP_CACHE, fs.OK
}

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

//var _ = (fs.NodeFsyncer)((*credRoot)(nil))

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
		ch := p.NewPersistentInode(ctx, &credFile{data: data}, fs.StableAttr{})
		p.AddChild(base, ch, true)
	}
}

func (cr *credRoot) Flush(ctx context.Context, f fs.FileHandle) syscall.Errno {
	return 0
}
