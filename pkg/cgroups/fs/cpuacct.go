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
		fmt.Errorf("Unable to read %s cgroup param: %s", path, err)
		return paramData, err
	}
	cpuPath := filepath.Join(path, "cpuacct.stat")
	f, err := os.Open(cpuPath)
	defer f.Close()
	if err != nil {
		return paramData, err
	}
	sc := bufio.NewScanner(f)
	cpuTotal := 0.0
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		t := fields[0]
		v, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			fmt.Errorf("Error parsing cpu stats: %s", err)
			continue
		}
		// set the raw data in map
		paramData[t] = float64(v)
		cpuTotal += v
	}
	// calculate percentage from jiffies
	// get sys uptime
	uf, err := os.Open("/proc/uptime")
	defer uf.Close()
	if err != nil {
		fmt.Errorf("Unable to open /proc/uptime")
		return paramData, err
	}
	uptimeData, _ := ioutil.ReadAll(uf)
	uptimeFields := strings.Fields(string(uptimeData))
	uptimeFloat, err := strconv.ParseFloat(uptimeFields[0], 64)
	uptime := int(uptimeFloat)
	if err != nil {
		fmt.Errorf("Error parsing cpu stats: %s", err)
		return paramData, err
	}
	// find starttime of process
	pidProcsPath := filepath.Join(path, "cgroup.procs")
	pf, _ := os.Open(pidProcsPath)
	defer pf.Close()
	pr := bufio.NewReader(pf)
	l, _, _ := pr.ReadLine()
	starttime, _ := strconv.Atoi(string(l))
	// get total elapsed seconds since proc start
	seconds := uptime - (starttime / 100)
	// finally calc percentage
	cpuPercentage := 100.0 * ((cpuTotal / 100.0) / float64(seconds))
	paramData["percentage"] = cpuPercentage
	return paramData, nil
}
