package tests_common

import (
	"bytes"
	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestFtpContainer(t *testing.T) {
	var wg sync.WaitGroup
	opt, destructor, port, container := FtpContainer(&wg)
	wg.Wait()
	defer destructor()
	t.Log(opt, port)
	t.Log(container.Inspection)

	c, err := ftp.Dial(opt.Addr, ftp.DialWithTimeout(5*time.Second), ftp.DialWithDisabledEPSV(true))
	if err != nil {
		t.Fatal(err)
	}

	err = c.Login(opt.User.Login, opt.User.Password)
	if err != nil {
		t.Fatal(err)
	}

	cd, err := c.CurrentDir()
	assert.NoError(t, err)
	t.Log(cd)
	l, err := c.NameList(cd)
	assert.NoError(t, err)
	t.Log(l)

	data := bytes.NewBufferString("Hello World")
	err = c.Stor("test-file.txt", data)
	assert.NoError(t, err)
	l, err = c.NameList(cd)
	assert.NoError(t, err)
	t.Log(l)
	assert.Equal(t, []string{"/ftp/ftplogin/test-file.txt"}, l)

	if err := c.Quit(); err != nil {
		t.Fatal(err)
	}

}
