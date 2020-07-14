package tests_common

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"os"
	"path"
	"sync"
	"time"
)

type NatsStreamingOptions struct {
	ClusterId string
	Url       string
}

func NatsStreamingContainer(wg *sync.WaitGroup) (*NatsStreamingOptions, func(), int, *Container) {
	wd, _ := os.Getwd()
	testdir := path.Base(wd)
	o := &NatsStreamingOptions{
		ClusterId: fmt.Sprintf("cluster%x", time.Now().UnixNano()),
	}
	addr, _, port := GetFreeLocalAddr()
	o.Url = fmt.Sprintf("nats://%s", addr)
	c := Container{
		Address: "docker.io/library/nats-streaming",
		Image:   "nats-streaming",
		Name:    fmt.Sprintf("%s-nats-streaming", testdir),
		Cmd:     []string{"--cluster_id", o.ClusterId},
		Ports:   map[string]string{addr: "4222/tcp"},
	}

	c.Run()
	wg.Add(1)
	go func() {
		defer wg.Done()
		opts := []nats.Option{nats.Name("NATS Streaming test")}
		for {
			nc, err := nats.Connect(o.Url, opts...)
			if err != nil {
				continue
			}
			defer nc.Close()
			sc, err := stan.Connect(o.ClusterId, "ClientId", stan.NatsConn(nc))
			if err != nil {
				continue
			}
			defer sc.Close()

			break
		}
	}()

	return o, c.Close, port, &c
}
