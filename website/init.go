package website

import (
	"webconsole/utils"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"webconsole/server"
)

var (
	Conf, err = utils.Get_Conf()
)

func Add_HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	server.HandleFunc(method, pattern, handler)
}

func Add_Handle(method, pattern string, handler http.Handler) {
	server.Handle(method, pattern, handler)
}

type Result struct {
	Ok   bool        `json:"ok" xml:"ok"`
	Msg  string      `json:"msg" xml:"msg"`
	Data interface{} `json:"data" xml:"data"`
}

type tmplPath struct {
	name   string
	path   string
	suffix string
}

var (
	templatesFuncMap template.FuncMap
	templates        *template.Template
	tmplPaths        map[string]tmplPath
)

func templatePathWalk(p string, f os.FileInfo, suffix string, err error) error {
	if f == nil {
		return err
	} else if f.IsDir() {
		return nil
	} else if (f.Mode() & os.ModeSymlink) > 0 {
		return nil
	} else {
		if f.Size() > 0 {
			sx := filepath.Ext(p)
			if strings.EqualFold(strings.ToLower(sx), strings.ToLower(suffix)) {
				p1 := strings.TrimSuffix(p, sx)
				tmplPaths[p1] = tmplPath{path: p, name: f.Name(), suffix: sx}
			}
		}
	}
	return err
}

func Init_Templates(dirName string, suffix string) error {

	templatesFuncMap = make(template.FuncMap)
	tmplPaths = make(map[string]tmplPath)

	err := filepath.Walk(strings.TrimRight(dirName, "/"), func(p string, f os.FileInfo, err error) error {
		return templatePathWalk(p, f, suffix, err)
	})

	if nil != err {
		return err
	}

	for k, v := range tmplPaths {
		if strings.EqualFold(strings.ToLower(v.suffix), strings.ToLower(suffix)) {
			tk := k[len(dirName)+1 : len(k)]
			htmlStr, err := ioutil.ReadFile(v.path)
			if err != nil {
				return err
			}
			htmlTxt := string(htmlStr)
			if len(htmlTxt) != 0 {
				var t *template.Template
				if templates == nil {
					templates = template.New(tk).Funcs(templatesFuncMap)
				}
				if tk == templates.Name() {
					t = templates
				} else {
					t = templates.New(tk).Funcs(templatesFuncMap)
				}
				t = t.Delims("<abc%", "%>")
				_, err = t.Parse(htmlTxt)
				if nil != err {
					return err
				}
				templates = template.Must(templates, err)
			}
		}
	}
	return nil
}

type Context struct {
	r       *http.Request
	w       http.ResponseWriter
	Session *server.Session
	v       map[string]interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	ctx := &Context{
		r:       r,
		w:       w,
		Session: server.WebSession,
	}
	err := ctx.parseForm()
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return ctx
}

func (c *Context) SetContentType(t string) {
	c.w.Header().Set("Content-Type", t)
}

func (c *Context) OutHtml(tpl string, obj interface{}) {
	c.SetContentType("text/html; charset=UTF-8")
	rb := bytes.NewBufferString("")
	err := templates.ExecuteTemplate(rb, tpl, obj)
	if nil != err {
		http.Error(c.w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		utils.Gzip_Html(rb, c.w, c.r)
		return
	}
}

func (c *Context) OutJson(obj interface{}) {
	c.SetContentType("application/json; charset=UTF-8")
	b, err := json.Marshal(obj)
	if nil != err {
		http.Error(c.w, err.Error(), http.StatusInternalServerError)
	} else {
		if Conf.Web.EnableJSONP {
			jsonpParam := c.GetFormValue(Conf.Web.JSONPParam)
			ret := jsonpParam + "(" + string(b) + ")"
			utils.Gzip_Binary([]byte(ret), c.w, c.r)
		} else {
			utils.Gzip_Binary(b, c.w, c.r)
		}
	}
	return
}

func (c *Context) OutXML(obj interface{}) {
	c.SetContentType("text/xml; charset=UTF-8")
	b, err := xml.Marshal(obj)
	if nil != err {
		http.Error(c.w, err.Error(), http.StatusInternalServerError)
	} else {
		utils.Gzip_Binary(b, c.w, c.r)
	}
	return
}

func (c *Context) GetFormValue(key string) string {
	fv := c.v[key]
	if nil != fv {
		return strings.TrimSpace(utils.ToStr(fv.([]string)[0]))
	} else {
		return ""
	}
}

func (c *Context) GetFormValues(v string) []string {
	fv := c.v[v]
	if nil != fv {
		return fv.([]string)
	} else {
		return nil
	}
}

func (c *Context) parseForm() error {
	err := c.r.ParseForm()
	if nil != err {
		return err
	}
	paramMap := make(map[string]interface{})
	s := c.r.Form
	for k, v := range s {
		if nil != paramMap[k] {
			paramArr := make([]interface{}, 0, 0)
			paramArr = append(paramArr, paramMap[k])
			paramArr = append(paramArr, v)
		} else {
			paramMap[k] = v
		}

	}
	c.v = paramMap
	return nil
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.w, c.r, url, http.StatusMovedPermanently)
	return
}

func (c *Context) BasicAuth(s string) {
	c.w.Header().Set("WWW-Authenticate", "Basic realm="+s)
	c.w.WriteHeader(http.StatusUnauthorized)
	return
}

func (c *Context) GetJsonByte() []byte {
	ct := c.r.Header.Get("Content-Type")
	if utils.StringUtils(ct).ContainsBool("application/json") {
		b, err := ioutil.ReadAll(c.r.Body)
		if nil != err {
			http.Error(c.w, err.Error(), http.StatusInternalServerError)
			return nil
		} else {
			return b
		}
	} else {
		return nil
	}
}

func (c *Context) IsLogin() bool {
	if nil != c.Session.Get("is_login") && c.Session.Get("is_login").(bool) {
		return true
	} else {
		return false
	}
}

func (c *Context) GetSessionUser() interface{} {
	if c.IsLogin() {
		su := c.Session.Get("user_info")
		if nil != su {
			return su
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func Run() {
	server.Run()
}

func init() {
	if nil != err {
		utils.Log_Fatal(err.Error())
	}
	utils.Load_Mime(utils.MimePath)
	if err := Init_Templates(utils.Get_Project_Dir()+utils.PathSeparator+strings.TrimLeft(Conf.Web.TemplateDir, "/"), Conf.Web.TemplateSuffix); nil != err {
		utils.Log_Fatal(err)
	}
	server.DefaultServeMux.AddStaticDir(Conf.Web.StaticDir)
}
