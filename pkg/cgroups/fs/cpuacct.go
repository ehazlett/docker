package fs

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dotcloud/docker/pkg/cgroups"
)

type cpuacctGroup struct {
}

func (s *cpuacctGroup) Set(d *data) error {
	// we just want to join this group even though we don't set anything
	if _, err := d.join("cpuacct"); err != nil && err != cgroups.ErrNotFound {
		return err
	}
	return nil
}

func (s *cpuacctGroup) Remove(d *data) error {
	return removePath(d.path("cpuacct"))
}

func (s *cpuacctGroup) Stats(d *data) (map[string]float64, error) {
	paramData := make(map[string]float64)
	path, err := d.path("cpuacct")
	if err != nil {
		return paramData, fmt.Errorf("Unable to read %s cgroup param: %s", path, err)
	}
	cpuPath := filepath.Join(path, "cpuacct.stat")
	f, err := os.Open(cpuPath)
	if err != nil {
		return paramData, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	cpuTotal := 0.0
	for sc.Scan() {
		t, v, err := getCgroupParamKeyValue(sc.Text())
		if err != nil {
			return paramData, fmt.Errorf("Error parsing param data: %s", err)
		}
		// set the raw data in map
		paramData[t] = v
		cpuTotal += v
	}
	// calculate percentage from jiffies
	// get sys uptime
	uf, err := os.Open("/proc/uptime")
	if err != nil {
		return paramData, fmt.Errorf("Unable to open /proc/uptime")
	}
	defer uf.Close()
	uptimeData, err := ioutil.ReadAll(uf)
	if err != nil {
		return paramData, fmt.Errorf("Error reading /proc/uptime: %s", err)
	}
	uptimeFields := strings.Fields(string(uptimeData))
	uptime, err := strconv.ParseFloat(uptimeFields[0], 64)
	if err != nil {
		return paramData, fmt.Errorf("Error parsing cpu stats: %s", err)
	}
	// find starttime of process
	pidProcsPath := filepath.Join(path, "cgroup.procs")
	pf, err := os.Open(pidProcsPath)
	if err != nil {
		return paramData, fmt.Errorf("Error parsing cpu stats: %s", err)
	}
	defer pf.Close()
	pr := bufio.NewReader(pf)
	l, _, err := pr.ReadLine()
	if err != nil {
		return paramData, fmt.Errorf("Error reading param file: %s", err)
	}
	starttime, err := strconv.ParseFloat(string(l), 64)
	if err != nil {
		return paramData, fmt.Errorf("Unable to parse starttime: %s", err)
	}
	// get total elapsed seconds since proc start
	seconds := uptime - (starttime / 100)
	// finally calc percentage
	cpuPercentage := 100.0 * ((cpuTotal / 100.0) / float64(seconds))
	paramData["percentage"] = cpuPercentage
	return paramData, nil
}
