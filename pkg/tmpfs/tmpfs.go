package tmpfs

import (
	"fmt"
	"strings"

	"github.com/docker/docker/daemon/graphdriver"
	"github.com/inaccel/daemon/pkg/runtime"
	"github.com/moby/sys/mount"
)

type Options struct {
	Size     string
	NrBlocks string
	NrInodes string
	Mode     string
	GID      string
	UID      string
	Huge     string
	Mpol     string
}

func IsMounted(path string) bool {
	return graphdriver.NewFsChecker(graphdriver.FsMagicTmpFs).IsMounted(path)
}

func Mount(path string, options Options) error {
	if !IsMounted(path) {
		var o []string
		if options.Size != "" {
			o = append(o, fmt.Sprintf("size=%s", options.Size))
		}
		if options.NrBlocks != "" {
			o = append(o, fmt.Sprintf("nr_blocks=%s", options.NrBlocks))
		}
		if options.NrInodes != "" {
			o = append(o, fmt.Sprintf("nr_inodes=%s", options.NrInodes))
		}
		if options.Mode != "" {
			o = append(o, fmt.Sprintf("mode=%s", options.Mode))
		}
		if options.GID != "" {
			if runtime.LinuxVersionCode >= runtime.KernelVersion(2, 5, 7) {
				o = append(o, fmt.Sprintf("gid=%s", options.GID))
			}
		}
		if options.UID != "" {
			if runtime.LinuxVersionCode >= runtime.KernelVersion(2, 5, 7) {
				o = append(o, fmt.Sprintf("uid=%s", options.UID))
			}
		}
		if options.Huge != "" {
			if runtime.LinuxVersionCode >= runtime.KernelVersion(4, 7, 0) {
				o = append(o, fmt.Sprintf("huge=%s", options.Huge))
			}
		}
		if options.Mpol != "" {
			if runtime.LinuxVersionCode >= runtime.KernelVersion(2, 6, 15) {
				o = append(o, fmt.Sprintf("mpol=%s", options.Mpol))
			}
		}

		return mount.Mount("tmpfs", path, "tmpfs", strings.Join(o, ","))
	}
	return nil
}

func Unmount(path string) error {
	if IsMounted(path) {
		return mount.Unmount(path)
	}
	return nil
}
