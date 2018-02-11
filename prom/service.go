package prom

import (
	"encoding/json"
	"regexp"
	"strings"

	etcd "github.com/coreos/etcd/client"
)

type (
	instances map[string]string
	services  map[string]instances
)

var pathPat = regexp.MustCompile(`/services/([^/]+)(?:/(\d+))?`)

func (srvs services) handle(node *etcd.Node, handler func(*etcd.Node)) {
	if node.Dir {
		for _, n := range node.Nodes {
			srvs.handle(n, handler)
		}
	} else {
		handler(node)
	}
}

func (srvs services) update(node *etcd.Node) {
	if !strings.Contains(node.Key, "metrics") {
		return
	}
	i := strings.LastIndex(node.Key, "/")
	srv := node.Key[:i]
	instanceID := node.Key[i+1:]
	insts, ok := srvs[srv]
	if !ok {
		insts = instances{}
		srvs[srv] = insts
	}
	insts[instanceID] = node.Value
}

func (srvs services) delete(node *etcd.Node) {
	if !strings.Contains(node.Key, "metrics") {
		return
	}
	i := strings.LastIndex(node.Key, "/")
	srv := node.Key[:i]
	instanceID := node.Key[i+1:]

	// Deletion of an entire service.
	if instanceID == "" {
		delete(srvs, srv)
		return
	}
	// Delete a single instance from the service.
	delete(srvs[srv], instanceID)
}

type TargetGroup struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

func (srvs services) persist() {
	var tgroups []*TargetGroup
	// Write files for current services.
	for job, instances := range srvs {
		var targets []string
		for _, addr := range instances {
			targets = append(targets, addr)
		}

		tgroups = append(tgroups, &TargetGroup{
			Targets: targets,
			Labels:  map[string]string{"job": job},
		})
	}

	content, err := json.Marshal(tgroups)
	if err != nil {
		logger.Errorln(err)
		return
	}

	f, err := create(*targetFile)
	if err != nil {
		logger.Errorln(err)
		return
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		logger.Errorln(err)
	}
}
