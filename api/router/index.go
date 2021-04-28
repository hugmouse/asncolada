package handler

import (
	"bytes"
	_ "embed" // go:embed requires import of "embed"
	"github.com/ammario/ipisp"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

//go:embed template/result.gohtml
var cumtemplate string

// Handler for routing.
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

	tmpl, err := template.New("result").Parse(cumtemplate)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	// Detecting what we need to do: IP or ASN queue
	if bytes.ContainsRune(splittedURL[1], '.') {
		resp, err := client.LookupIP(net.ParseIP(string(splittedURL[1])))
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-type", "text/html; charset=UTF-8")
		err = tmpl.Execute(w, struct {
			IPorASN   string
			RawOutput *ipisp.Response
		}{
			IPorASN:   resp.IP.String(),
			RawOutput: resp,
		})

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

		w.Header().Set("Content-type", "text/html; charset=UTF-8")
		err = tmpl.Execute(w, struct {
			IPorASN   string
			RawOutput *ipisp.Response
		}{
			IPorASN:   resp.IP.String(),
			RawOutput: resp,
		})
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

	}
}
