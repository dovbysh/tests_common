package tests_common

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/parnurzeal/gorequest"
	"os"
	"path"
	"sync"
	"time"
)

const (
	PgPoolSize     = 150
	PgMinIdleConns = 10
)

type ClickHouseAddresses map[string]string

func PostgreSQLContainer(wg *sync.WaitGroup) (*pg.Options, func(), int, *Container) {
	wd, _ := os.Getwd()
	testdir := path.Base(wd)
	o := &pg.Options{
		User:         "root",
		Password:     "root",
		Database:     "db",
		TLSConfig:    nil,
		PoolSize:     PgPoolSize,
		MinIdleConns: PgMinIdleConns,
	}
	dbAddr, _, port := GetFreeLocalAddr()
	o.Addr = dbAddr
	postgresContainer := Container{
		Address: "docker.io/library/postgres",
		Image:   "postgres",
		Name:    fmt.Sprintf("%s-postgres", testdir),
		Environments: []string{
			"POSTGRES_USER=" + o.User,
			"POSTGRES_PASSWORD=" + o.Password,
			"POSTGRES_DB=" + o.Database,
		},
		Ports: map[string]string{dbAddr: "5432/tcp"},
	}

	postgresContainer.Run()
	wg.Add(1)
	go func() {
		defer wg.Done()
		db := pg.Connect(o)
		defer db.Close()
		for {
			if _, e := db.WithTimeout(1 * time.Second).Exec("select 1"); e == nil {
				break
			}
		}
	}()

	return o, postgresContainer.Close, port, &postgresContainer
}

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
