package upstream

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Echo echoes the HTTP request details:
// method, path, host, headers, cookies and body
type Echo struct{}

func (e Echo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := &bytes.Buffer{}

	e.basicRequestDetails(buf, r)
	e.queries(buf, r)
	e.headers(buf, r)
	e.cookies(buf, r)
	err := e.body(buf, r)
	if err != nil {
		http.Error(w, "couldn't read body: "+err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}

	_, err = w.Write(buf.Bytes())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "couldn't write response: "+err.Error(), http.StatusInternalServerError)
	}
}

func (Echo) basicRequestDetails(sw io.StringWriter, r *http.Request) {
	_, _ = sw.WriteString("Method: " + r.Method + "\n")
	if scheme := r.URL.Scheme; scheme != "" {
		_, _ = sw.WriteString("Scheme: " + r.URL.Scheme + "\n")
	}

	if host := r.URL.Host; host != "" {
		_, _ = sw.WriteString("Host: " + r.URL.Host + "\n")
	}

	_, _ = sw.WriteString("Path: " + r.URL.Path + "\n")
}

func (Echo) queries(sw io.StringWriter, r *http.Request) {
	if queries := r.URL.Query(); len(queries) > 0 {
		_, _ = sw.WriteString("\n" + "Query values:" + "\n")
		for key := range queries {
			_, _ = sw.WriteString(fmt.Sprintf("  %s: %s"+"\n", key, queries.Get(key)))
		}
	}
}

func (Echo) headers(sw io.StringWriter, r *http.Request) {
	if headers := r.Header; len(headers) > 0 {
		_, _ = sw.WriteString("\n" + "Headers:" + "\n")
		for key := range headers {
			_, _ = sw.WriteString(fmt.Sprintf("  %s: %s"+"\n", key, headers.Get(key)))
		}
	}
}

func (Echo) cookies(sw io.StringWriter, r *http.Request) {
	if cookies := r.Cookies(); len(cookies) > 0 {
		_, _ = sw.WriteString("\n" + "Cookies:" + "\n")
		for _, cookie := range cookies {
			if cookie == nil {
				continue
			}

			_, _ = sw.WriteString(fmt.Sprintf("  - %s: %s"+"\n", cookie.Name, cookie.Value))
			_, _ = sw.WriteString(fmt.Sprintf("    Domain: %s"+"\n", cookie.Domain))
			_, _ = sw.WriteString(fmt.Sprintf("    Path: %s"+"\n", cookie.Path))
			_, _ = sw.WriteString(fmt.Sprintf("    Expires: %s"+"\n", cookie.Expires))
			_, _ = sw.WriteString(fmt.Sprintf("    HTTP only: %v"+"\n", cookie.HttpOnly))
			_, _ = sw.WriteString(fmt.Sprintf("    Secure: %v"+"\n", cookie.Secure))
		}
	}
}

func (Echo) body(buf *bytes.Buffer, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(body) > 0 {
		buf.WriteString("\n" + "Body:" + "\n")
		buf.Write(body)
	}

	return nil
}
