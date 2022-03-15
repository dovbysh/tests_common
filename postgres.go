package tests_common

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/go-pg/pg/v9"
)

const (
	PgPoolSize     = 150
	PgMinIdleConns = 10
)

type PgOptions struct {
	Addr         string
	User         string
	Password     string
	Database     string
	PoolSize     int
	MinIdleConns int
}

func PostgreSQLContainer(wg *sync.WaitGroup) (*PgOptions, func(), int, *Container) {
	wd, _ := os.Getwd()
	testdir := path.Base(wd)
	o := &PgOptions{
		User:         "root",
		Password:     "root",
		Database:     "db",
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
		db := pg.Connect(&pg.Options{
			Addr:         o.Addr,
			User:         o.User,
			Password:     o.Password,
			Database:     o.Database,
			PoolSize:     o.PoolSize,
			MinIdleConns: o.MinIdleConns,
		})
		defer db.Close()
		for {
			if _, e := db.WithTimeout(1 * time.Second).Exec("select 1"); e == nil {
				break
			}
		}
	}()

	return o, postgresContainer.Close, port, &postgresContainer
}
