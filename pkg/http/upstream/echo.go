package upstream

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Echo echoes the HTTP request details:
// method, path, host, headers, cookies and body
type Echo struct{}

func (Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}

	buf.WriteString("Method: " + r.Method + "\n")
	if scheme := r.URL.Scheme; scheme != "" {
		buf.WriteString("Scheme: " + r.URL.Scheme + "\n")
	}

	if host := r.URL.Host; host != "" {
		buf.WriteString("Host: " + r.URL.Host + "\n")
	}

	buf.WriteString("Path: " + r.URL.Path + "\n")

	if queries := r.URL.Query(); len(queries) > 0 {
		buf.WriteString("\n" + "Query values:" + "\n")
		for key := range queries {
			buf.WriteString(fmt.Sprintf("  %s: %s"+"\n", key, queries.Get(key)))
		}
	}

	if headers := r.Header; len(headers) > 0 {
		buf.WriteString("\n" + "Headers:" + "\n")
		for key := range headers {
			buf.WriteString(fmt.Sprintf("  %s: %s"+"\n", key, headers.Get(key)))
		}
	}

	if cookies := r.Cookies(); len(cookies) > 0 {
		buf.WriteString("\n" + "Cookies:" + "\n")
		for _, cookie := range cookies {
			if cookie == nil {
				continue
			}

			buf.WriteString(fmt.Sprintf("  - %s: %s"+"\n", cookie.Name, cookie.Value))
			buf.WriteString(fmt.Sprintf("    Domain: %s"+"\n", cookie.Domain))
			buf.WriteString(fmt.Sprintf("    Path: %s"+"\n", cookie.Path))
			buf.WriteString(fmt.Sprintf("    Expires: %s"+"\n", cookie.Expires))
			buf.WriteString(fmt.Sprintf("    HTTP only: %v"+"\n", cookie.HttpOnly))
			buf.WriteString(fmt.Sprintf("    Secure: %v"+"\n", cookie.Secure))
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "couldn't read body: "+err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if len(body) > 0 {
		buf.WriteString("\n" + "Body:" + "\n")
		buf.Write(body)
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "couldn't write response: "+err.Error(), http.StatusInternalServerError)
	}
}
