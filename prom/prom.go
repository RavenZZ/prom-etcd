package prom

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/coreos/etcd/client"
)

var (
	targetFile = flag.String("target-file", "tgroups.json", "the file that contains the target groups")
)

func Run() {
	var servicesPrefix = fmt.Sprintf("/services/%s/", ZTOHostName)
	var (
		srvs = services{}
	)

	ctx := context.TODO()
	flag.Parse()
	cfg := client.Config{Endpoints: []string{"http://10.9.13.6:2379"}}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)
	if err != nil {
		if err == context.Canceled {
			// ctx is canceled by another routine
		} else if err == context.DeadlineExceeded {
			// ctx is attached with a deadline and it exceeded
		} else if cerr, ok := err.(*client.ClusterError); ok {
			logger.Errorf(cerr.Error())
			// process (cerr.Errors)
		} else {
			// bad cluster endpoints, which are not etcd servers
		}
	}

	//// Retrieve the subtree of the /services path.
	resp, err := kapi.Get(ctx, servicesPrefix, &client.GetOptions{Recursive: true})

	if err != nil {
		logger.Fatalln(err)
	} else {
		srvs.handle(resp.Node, srvs.update)
		srvs.persist()
	}

	watcher := kapi.Watcher(servicesPrefix, &client.WatcherOptions{Recursive: true, AfterIndex: 0})
	// Start recursively watching for updates.
	for {

		res, err := watcher.Next(ctx)
		if err != nil {
			logger.Errorln(err)
		}
		handler := srvs.update
		if res.Action == "delete" || res.Action == "expire" {
			handler = srvs.delete
		}
		srvs.handle(res.Node, handler)
		srvs.persist()
	}
}
