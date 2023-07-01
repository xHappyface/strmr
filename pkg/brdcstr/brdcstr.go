package brdcstr

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Brdcstr struct {
	URL    string
	Port   string
	client http.Client
}

func New(url string, port string) *Brdcstr {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	return &Brdcstr{
		URL:    url,
		Port:   port,
		client: c,
	}
}

func (b *Brdcstr) Alive() (bool, error) {
	url := fmt.Sprintf("%s:%s/-/ping", b.URL, b.Port)
	resp, err := b.client.Get(url)
	if err != nil {
		return false, errors.New("cannot check status on brdcstr: " + err.Error())
	}
	if resp.StatusCode != 200 {
		return false, errors.New("non 200 OK status from brdcstr alive check: " + strconv.Itoa(resp.StatusCode))
	}
	return true, nil
}
