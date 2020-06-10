package tests_common

import (
	"sync"
	"testing"
)

func TestPostgreSQLContainer(t *testing.T) {
	var wg sync.WaitGroup

	_, pgCloser, _, _ := PostgreSQLContainer(&wg)
	defer pgCloser()

	wg.Wait()
}

func TestClickhouseContainer(t *testing.T) {
	var wg sync.WaitGroup

	_, olapCloser := ClickhouseContainer(&wg)
	defer olapCloser()

	wg.Wait()
}