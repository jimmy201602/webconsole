package server

import (
	"webconsole/utils"
	"container/list"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	SessionName    = "SID"
	SessionManager = &Manager{list: list.New(), timeOut: 30, sessions: make(map[string]*list.Element)}
)

type Session struct {
	lock    sync.RWMutex
	sid     string
	regTime time.Time
	value   map[interface{}]interface{}
}

type Manager struct {
	lock     sync.RWMutex
	list     *list.List
	timeOut  int
	sessions map[string]*list.Element
}

func (s *Session) Set(key, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.value[key] = value
}

func (s *Session) Get(key interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if v, ok := s.value[key]; ok {
		return v
	} else {
		return nil
	}
}

func (s *Session) Del(key interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.value[key]; ok {
		delete(s.value, key)
	}
}

func (s *Session) ID() string {
	return s.sid
}

func (m *Manager) Init(sid string, timeOut int) *Session {
	m.timeOut = timeOut
	m.lock.RLock()
	if sess, ok := m.sessions[sid]; ok {
		go m.Update(sid)
		m.lock.RUnlock()
		return sess.Value.(*Session)
	} else {
		m.lock.RUnlock()
		m.lock.Lock()
		nsess := &Session{sid: sid, regTime: time.Now(), value: make(map[interface{}]interface{})}
		sess := m.list.PushBack(nsess)
		m.sessions[sid] = sess
		m.lock.Unlock()
		return nsess
	}
	return nil
}

func (m *Manager) Update(sid string) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if sess, ok := m.sessions[sid]; ok {
		sess.Value.(*Session).regTime = time.Now()
		m.list.MoveToFront(sess)
	}
}

func (m *Manager) Destroy(sid string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if sess, ok := m.sessions[sid]; ok {
		delete(m.sessions, sid)
		m.list.Remove(sess)
	}
}

func (m *Manager) GC() {
	m.lock.RLock()
	for {
		sess := m.list.Back()
		if sess == nil {
			break
		}
		cTime := sess.Value.(*Session).regTime.Add(+time.Second * time.Duration(m.timeOut))
		if cTime.Before(time.Now()) {
			m.lock.RUnlock()
			m.lock.Lock()
			m.list.Remove(sess)
			delete(m.sessions, sess.Value.(*Session).sid)
			m.lock.Unlock()
			m.lock.RLock()
		} else {
			break
		}
	}
	m.lock.RUnlock()
	time.AfterFunc(time.Duration(m.timeOut)*time.Second, func() { SessionManager.GC() })
}

func RegSession(w http.ResponseWriter, r *http.Request, timeout int, enableTls bool) *Session {
	var session *Session
	cookie, err := r.Cookie(SessionName)
	if nil != err && err == http.ErrNoCookie {
		sid, _ := utils.StringUtils("").UUID()
		session = SessionManager.Init(sid, timeout)
		cookie = &http.Cookie{
			Name:     SessionName,
			Value:    url.QueryEscape(sid),
			Path:     "/",
			HttpOnly: true,
		}
	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session = SessionManager.Init(sid, timeout)
	}
	if enableTls {
		cookie.Secure = true
	}
	cookie.MaxAge = 0
	cookie.Expires = time.Now().Add(+time.Duration(timeout) * time.Second)
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
	go SessionManager.GC()
	return session
}
