package disks

import (
	"syscall"
)



type DiskStats struct {
  Total     uint64
  Used      uint64
  Free      uint64
}

func GetDiskUsage() (*DiskStats, error) {
  total, free, used, err := getDiskUsage()
  if err != nil {
    return nil, err
  }
  return &DiskStats{
    Total: total,
    Used: used,
    Free: free,
  },nil
}

func getDiskUsage() (total, free, used uint64, err error) {
  path := "/"
	var stat syscall.Statfs_t
	err = syscall.Statfs(path, &stat)
	if err != nil {
		return 0, 0, 0, err
	}

	total = stat.Blocks * uint64(stat.Bsize)
	free = stat.Bfree * uint64(stat.Bsize)
	used = total - free
	return total, free, used, nil
}




