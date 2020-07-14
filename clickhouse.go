package tests_common

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"os"
	"path"
	"sync"
	"time"
)

type ClickHouseAddresses map[string]string

func ClickhouseContainer(wg *sync.WaitGroup) (ClickHouseAddresses, func()) {
	wd, _ := os.Getwd()
	testdir := path.Base(wd)
	olapAddrHttp, _, _ := GetFreeLocalAddr()
	olapAddrTcp, _, _ := GetFreeLocalAddr()
	olapContainer := Container{
		Address:      "docker.io/yandex/clickhouse-server",
		Image:        "yandex/clickhouse-server",
		Name:         fmt.Sprintf("%s-clickhouse", testdir),
		Environments: []string{},
		Ports:        map[string]string{olapAddrHttp: "8123/tcp", olapAddrTcp: "9000/tcp"},
	}
	olapContainer.Run()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			r, _, _ := gorequest.New().Get(fmt.Sprintf("http://%s/?debug=1", olapAddrHttp)).End()
			if r == nil {
				time.Sleep(time.Microsecond)
				continue
			}
			break
		}
	}()

	return ClickHouseAddresses{
		"http": olapAddrHttp,
		"tcp":  olapAddrTcp,
	}, olapContainer.Close
}
