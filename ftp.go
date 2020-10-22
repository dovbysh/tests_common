package tests_common

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type FtpUser struct {
	Login    string
	Password string
}

func (u FtpUser) String() string {
	return fmt.Sprintf("%s|%s", u.Login, u.Password)
}

type FtpOptions struct {
	User FtpUser
	Addr string
}

func FtpContainer(wg *sync.WaitGroup) (*FtpOptions, func(), int, *Container) {
	wd, _ := os.Getwd()
	testdir := path.Base(wd)
	addr, _, port := GetFreeLocalAddr()
	o := &FtpOptions{
		User: FtpUser{
			Login:    "ftplogin",
			Password: "ftppass",
		},
		Addr: addr,
	}
	container := Container{
		Address: "docker.io/delfer/alpine-ftp-server",
		Image:   "delfer/alpine-ftp-server",
		Name:    fmt.Sprintf("%s-ftp", testdir),
		Environments: []string{
			"USERS=" + o.User.String(),
		},
		Ports: map[string]string{addr: "21/tcp"},
	}

	container.Run()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			c, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
			if err != nil {
				continue
			}

			err = c.Login(o.User.Login, o.User.Password)
			if err != nil {
				log.Fatal(err)
			}

			if err := c.Quit(); err != nil {
				continue
			}

			break
		}
	}()

	return o, container.Close, port, &container
}
