package handler

import (
	"bytes"
	_ "embed" // go:embed requires import of "embed"
	"encoding/json"
	"github.com/ammario/ipisp"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

//go:embed template/result.gohtml
var ASNIPTemplate string

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

	tmpl, err := template.New("result").Parse(ASNIPTemplate)
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
		JSON, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = tmpl.Execute(w, struct {
			RawOutput *ipisp.Response
			JSON      string
		}{
			RawOutput: resp,
			JSON:      string(JSON),
		})

		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

	} else {
		if bytes.ContainsAny(splittedURL[1], "ASas") {
			splittedURL[1] = splittedURL[1][2:]
		}
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
		JSON, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		err = tmpl.Execute(w, struct {
			RawOutput *ipisp.Response
			JSON      string
		}{
			RawOutput: resp,
			JSON:      string(JSON),
		})
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

	}
}
