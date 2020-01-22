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
	//	mu      sync.Mutex
	secrets creds.Secrets
	secret  string
}

var _ = (fs.NodeOnAdder)((*credRoot)(nil))

//var _ = (fs.NodeFsyncer)((*credRoot)(nil))

func createRoot(ctx context.Context, secrets creds.Secrets, secret string) (*credRoot, error) {
	root := credRoot{
		secrets: secrets,
		secret:  secret,
	}
	data, exists, err := root.secrets.Get(root.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial secret volume data: %s", err)
	}
	if !exists {
		return &root, nil
	}

	raw, err := json.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to encode secret volume date %s", err)
	}

	var all map[string][]byte
	err = json.Unmarshal(raw, &all)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret volume date %s", err)
	}

	for f, data := range all {
		dir, base := filepath.Split(f)

		p := &zr.Inode
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
		root.AddChild(base, ch, true)
	}
	return root, nil
}

// func (cr *credRoot) Flush(ctx context.Context, f fs.FileHandle) syscall.Errno {
// 	return 0
// }
