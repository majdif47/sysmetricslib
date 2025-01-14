package diskinfo

import (
	"syscall"
)

func getDiskUsage() (total, free, used uint64, err error) {
  path := "/"
	var stat syscall.Statfs_t

	// Get filesystem stats for the given path
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, 0, err
	}

	// Total space (in bytes)
	total = stat.Blocks * uint64(stat.Bsize)

	// Free space (in bytes)
	free = stat.Bfree * uint64(stat.Bsize)

	// Used space (in bytes)
	used = total - free

	return total, free, used, nil
}




