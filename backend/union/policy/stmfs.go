package policy

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/rclone/rclone/backend/union/upstream"
	"github.com/rclone/rclone/fs"
)

func init() {
	registerPolicy("stmfs", &StMfs{})
}

// St stands for syncthing.
// Mfs stands for most free space
// Search category: same as epmfs.
// Action category: same as epmfs.
// Create category: Put the temporary files on the first drive. Pick the drive with the most free space among the others.
type StMfs struct {
	Mfs
}

const syncthingTempPrefix = ".syncthing."

// Create category policy, governing the creation of files and directories
func (p *StMfs) Create(ctx context.Context, upstreams []*upstream.Fs, path string) ([]*upstream.Fs, error) {
	if len(upstreams) < 2 {
		return nil, fs.ErrorPermissionDenied
	}

	name := filepath.Base(path)
	if strings.HasPrefix(name, syncthingTempPrefix) {
		return []*upstream.Fs{upstreams[0]}, nil
	}

	return p.Mfs.Create(ctx, upstreams[1:], path)
}