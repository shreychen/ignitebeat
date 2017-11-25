package beater

import (
	"io/ioutil"
	"net"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"
)

// Get my IP
func GetMyIP() (ip string, err error) {
	if conn, err := net.Dial("udp", "8.8.8.8:80"); err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)

		ip = localAddr.IP.String()
	}
	return ip, err
}

// Open URL
func OpenURL(u string) (body []byte, err error) {

	resp, err := http.Get(u)
	defer resp.Body.Close()

	if err != nil {
		logp.Info(err.Error())
		return body, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logp.Info(err.Error())
		return body, err
	}

	return body, nil
}
