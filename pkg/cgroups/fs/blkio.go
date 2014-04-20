package fs

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dotcloud/docker/pkg/cgroups"
)

type blkioGroup struct {
}

func (s *blkioGroup) Set(d *data) error {
	// we just want to join this group even though we don't set anything
	if _, err := d.join("blkio"); err != nil && err != cgroups.ErrNotFound {
		return err
	}
	return nil
}

func (s *blkioGroup) Remove(d *data) error {
	return removePath(d.path("blkio"))
}

func (s *blkioGroup) Stats(d *data) (map[string]float64, error) {
	paramData := make(map[string]float64)
	path, _ := d.path("blkio")
	params := []string{
		"blkio.sectors",
		"blkio.io_service_bytes",
		"blkio.io_serviced",
		"blkio.io_queued",
	}
	for _, param := range params {
		paramPath := filepath.Join(path, param)
		f, err := os.Open(paramPath)
		if err != nil {
			return paramData, err
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			fields := strings.Fields(sc.Text())
			v, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Errorf("Error parsing %s stats: %s", param, err)
				continue
			}
			paramName := strings.Split(param, ".")[1]
			paramData[paramName] = float64(v)
		}
	}
	return paramData, nil
}
