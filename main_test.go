package tests_common

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
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

func TestNatsStreamingContainer(t *testing.T) {
	var wg sync.WaitGroup

	o, closer, port, c := NatsStreamingContainer(&wg)
	wg.Wait()
	defer closer()

	assert.NotEmpty(t, port)
	assert.NotNil(t, c)

	opts := []nats.Option{nats.Name("NATS Streaming testing")}
	nc, err := nats.Connect(o.Url, opts...)
	assert.NoError(t, err)
	defer nc.Close()

	sc, err := stan.Connect(o.ClusterId, "ClientId", stan.NatsConn(nc))
	assert.NoError(t, err)

	defer sc.Close()

	subj := "testSubj"
	mesg := []byte("testMsg")

	wg.Add(1)
	subsc, err := sc.Subscribe(subj, func(msg *stan.Msg) {
		defer wg.Done()
		assert.Equal(t, mesg, msg.Data)
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, subsc)
	err = sc.Publish(subj, mesg)
	assert.NoError(t, err)

	wg.Wait()

	d, err := subsc.Delivered()
	assert.Equal(t, int64(1), d)
	assert.NoError(t, err)
}
