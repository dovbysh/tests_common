package tests_common

import (
	"sync"
	"testing"
)

func TestPostgreSQLContainer(t *testing.T) {
	var wg sync.WaitGroup
	_, pgCloser, _, _ := PostgreSQLContainer(&wg)
	defer pgCloser()

	_, olapCloser := ClickhouseContainer(&wg)
	defer olapCloser()

	wg.Wait()
}
