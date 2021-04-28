package handler

import (
	"bytes"
	_ "embed" // go:embed requires import of "embed"
	"encoding/json"
	"errors"
	"github.com/ammario/ipisp"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

var (
	InvalidRequest = errors.New("you sent an invalid request")
)

//go:embed template/result.gohtml
var ASNIPTemplate string

//go:embed template/error.gohtml
var ErrorTemplate string

func WriteError(w http.ResponseWriter, t *template.Template, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	err = t.Execute(w, struct {
		Error error
	}{Error: err})
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	return
}

// Handler for routing.
func Handler(w http.ResponseWriter, r *http.Request) {

	tmplError, err := template.New("error").Parse(ErrorTemplate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, tmplError, err)
		return
	}

	splittedURL := bytes.Split(body, []byte("="))

	if len(splittedURL) < 2 {
		WriteError(w, tmplError, InvalidRequest)
		return
	}

	if len(splittedURL[1]) == 0 {
		WriteError(w, tmplError, InvalidRequest)
		return
	}

	client, err := ipisp.NewDNSClient()
	if err != nil {
		WriteError(w, tmplError, err)
		return
	}
	defer client.Close()

	tmpl, err := template.New("result").Parse(ASNIPTemplate)
	if err != nil {
		WriteError(w, tmplError, err)
		return
	}

	// Detecting what we need to do: IP or ASN queue
	if bytes.ContainsRune(splittedURL[1], '.') {
		resp, err := client.LookupIP(net.ParseIP(string(splittedURL[1])))
		if err != nil {
			WriteError(w, tmplError, err)
			return
		}

		w.Header().Set("Content-type", "text/html; charset=UTF-8")
		JSON, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			WriteError(w, tmplError, err)
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
			WriteError(w, tmplError, err)
			return
		}

	} else {
		if bytes.ContainsAny(splittedURL[1], "ASas") {
			splittedURL[1] = splittedURL[1][2:]
		}
		asnNum, err := strconv.Atoi(string(splittedURL[1]))
		if err != nil {
			WriteError(w, tmplError, err)
			return
		}
		resp, err := client.LookupASN(ipisp.ASN(asnNum))
		if err != nil {
			WriteError(w, tmplError, err)
			return
		}

		w.Header().Set("Content-type", "text/html; charset=UTF-8")
		JSON, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			WriteError(w, tmplError, err)
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
			WriteError(w, tmplError, err)
			return
		}

	}
}
