package handler

import (
	"bytes"
	"fmt"
	"github.com/ammario/ipisp"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

// Handler for routing
func Handler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	splittedURL := bytes.Split(body, []byte("="))

	if len(splittedURL) < 2 {
		_, err := w.Write([]byte("bruh"))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
		return
	}

	if len(splittedURL[1]) == 0 {
		_, err := w.Write([]byte("bruh"))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
		}
		return
	}

	client, err := ipisp.NewDNSClient()
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer client.Close()

	// Detecting what we need to do: IP or ASN queue
	if bytes.ContainsRune(splittedURL[1], '.') {
		resp, err := client.LookupIP(net.ParseIP(string(splittedURL[1])))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, err = w.Write([]byte(fmt.Sprintf("Resolved IP: %+v\n", resp)))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	} else {
		asnNum, err := strconv.Atoi(string(splittedURL[1]))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		resp, err := client.LookupASN(ipisp.ASN(asnNum))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		_, err = w.Write([]byte(fmt.Sprintf("Resolved ASN: %+v\n", resp)))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}
	}
}
