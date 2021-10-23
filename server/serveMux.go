package server

import (
	"webconsole/utils"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"sync"
)

type ServeMux struct {
	mu         sync.RWMutex
	m          map[*regexp.Regexp]muxEntry
	staticDir  []string
	hosts      bool
	serverName string
}

type muxEntry struct {
	explicit bool
	h        http.Handler
	pattern  string
	method   string
	params   map[int]string
}

func NewServeMux() *ServeMux {
	s := &ServeMux{
		m:          make(map[*regexp.Regexp]muxEntry),
		staticDir:  make([]string, 0, 0),
		serverName: "ApiBoxServer_v1.1",
	}
	return s
}

var (
	DefaultServeMux = NewServeMux()
	WebSession      *Session
)

func pathMatch(pattern *regexp.Regexp, path string) bool {
	if len(path) == 0 {
		return false
	}
	if pattern.MatchString(path) {
		return true
	} else {
		return false
	}
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

func (mux *ServeMux) match(r *http.Request, path string) (h http.Handler, pattern string) {
	paths := strings.Split(path, "/")
	values := r.URL.Query()
	for k, v := range mux.m {
		methods := strings.Split(v.method, ",")
		for _, mv := range methods {
			if strings.EqualFold(r.Method, mv) {
				if !pathMatch(k, path) {
					continue
				} else {
					if len(v.params) != 0 {
						for j, h := range v.params {
							values.Add(h, paths[j])
						}
					}
					r.URL.RawQuery = url.Values(values).Encode()
					h = v.h
					pattern = path
				}
			}
		}
	}
	return
}

func (mux *ServeMux) AddStaticDir(dir string) *ServeMux {
	mux.staticDir = append(mux.staticDir, utils.PathSeparator+strings.Trim(dir, "/"))
	return mux
}

func (mux *ServeMux) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	mux.Handle(method, pattern, handlerFunc(handler))
}

func (mux *ServeMux) Handle(method, pattern string, handler http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	params, regex, err := PathRegex(pattern)
	if nil != err {
		panic(err)
	}
	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if mux.m[regex].explicit {
		panic("http: multiple registrations for " + pattern)
	}

	if _, ok := mux.m[regex]; ok {
		panic("Conflicting Routes:" + pattern)
	} else {
		mux.m[regex] = muxEntry{explicit: true, h: handler, pattern: pattern, method: method, params: params}
		if pattern[0] != '/' {
			mux.hosts = true
		}
	}
}

func (mux *ServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			_, pattern = mux.handler(r, p)
			url := *r.URL
			url.Path = p
			return RedirectHandler(url.String(), http.StatusMovedPermanently), pattern
		}
	}
	return mux.handler(r, r.URL.Path)
}

func (mux *ServeMux) handler(r *http.Request, path string) (h http.Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	if mux.hosts {
		h, pattern = mux.match(r, r.Host+path)
	}
	if h == nil {
		h, pattern = mux.match(r, path)
	}
	if h == nil {
		h, pattern = http.NotFoundHandler(), ""
	}
	return
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	WebSession = RegSession(w, r, ABC_Conf.Web.SessionTimeOut, ABC_Conf.Web.EnableTLS)
	rPath := cleanPath(r.URL.Path)
	for _, s := range mux.staticDir {
		if strings.HasPrefix(rPath, s) {
			file_path := s + rPath[len(s):]
			file_path = strings.Trim(file_path, "/")
			file_path = utils.Get_Project_Dir() + utils.PathSeparator + file_path
			utils.Gzip_File(file_path, w, r)
			return
		}
	}
	w.Header().Add("Server", mux.serverName)
	h, _ := mux.Handler(r)
	h.ServeHTTP(w, r)
}

type handlerFunc func(http.ResponseWriter, *http.Request)

func (f handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func Handle(method, pattern string, handler http.Handler) {
	DefaultServeMux.Handle(method, pattern, handler)
}

func HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultServeMux.HandleFunc(method, pattern, handler)
}

type redirectHandler struct {
	url  string
	code int
}

func (rh *redirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	urlStr := rh.url + r.URL.String()
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Cache", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Location", urlStr)
	w.WriteHeader(rh.code)
	return
}

func RedirectHandler(url string, code int) http.Handler {
	return &redirectHandler{url, code}
}

func PathRegex(pattern string) (params map[int]string, regex *regexp.Regexp, err error) {
	parts := strings.Split(pattern, "/")
	params = make(map[int]string, 0)
	for i, part := range parts {
		utils.Log_Debug("Func path[n]:", part)
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[i] = part[1:]
			parts[i] = expr
		}
	}
	regex, err = regexp.Compile("^" + strings.Join(parts, "/") + "$")
	return
}
