package utils

import (
	"os"
)

const (
	Agent_prefix  = "Agent."
	Server_prefix = "RPC_Server."
	PathSeparator = string(os.PathSeparator)
	DevNull       = os.DevNull
	OS_MacOS      = "darwin"
	OS_FreeBSD    = "freebsd"
	OS_Linux      = "linux"
	OS_Windows    = "windows"
	Arch_AMD64    = "amd64"
	Arch_386      = "386"
	ChanLen       = 500

	CONNECT = "CONNECT"
	DELETE  = "DELETE"
	GET     = "GET"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
	PATCH   = "PATCH"
	POST    = "POST"
	PUT     = "PUT"
	TRACE   = "TRACE"

	Zip_gZip       = "gzip"
	Zip_Deflate    = "deflate"
	GZip_Mini_Size = 1024

	Log_file_suffix = ".log"
)

var (
	LogDir   = Get_Project_Dir() + PathSeparator + "log"
	ConfDir  = Get_Project_Dir() + PathSeparator + "conf"
	PidPath  = LogDir + PathSeparator + "utils.pid"
	MimePath = ConfDir + PathSeparator + "mime.types"
	ConfPath = ConfDir + PathSeparator + "conf.json"
)

type (
	WS_Conf struct {
		Addr           string `json:"addr" xml:"addr"`
		TlsAddr        string `json:"tls_addr" xml:"tls_addr"`
		TlsUrl         string `json:"tls_url" xml:"tls_url"`
		SessionTimeOut int    `json:"session_timeout" xml:"session_timeout"`
		StaticDir      string `json:"static_dir" xml:"static_dir"`
		TemplateDir    string `json:"template_dir" xml:"template_dir"`
		TemplateSuffix string `json:"template_suffix" xml:"template_suffix"`
		EnableJSONP    bool   `json:"enable_jsonp" xml:"enable_jsonp"`
		JSONPParam     string `json:"jsonp_param" xml:"jsonp_param"`
		EnableTLS      bool   `json:"enable_tls" xml:"enable_tls"`
		CorsWhiteList  string `json:"cors_white_list" xml:"cors_white_list"`
		EnableFcgi     bool   `json:"enable_fcgi" xml:"enable_fcgi"`
		TlsCert        string `json:"tls_cert,omitempty" xml:"tls_cert,omitempty"`
		TlsKey         string `json:"tls_key,omitempty" xml:"tls_key,omitempty"`
		Daemon         bool   `json:"daemon,omitempty" xml:"daemon,omitempty"`
                UrlTimeout     string `json:"url_timeout,omitempty" xml:"url_timeout,omitempty"`
                LogLevel       string `json:"log_level,omitempty" xml:"log_level,omitempty"`
	}

	DB_Conf struct {
		Url          string `json:"url" xml:"url"`
		MaxOpenConns int    `json:"max_open_conns" xml:"max_open_conns"`
		MaxIdleConns int    `json:"max_idle_conns" xml:"max_idle_conns"`
	}

	Conf struct {
		Web WS_Conf `json:"web,omitempty" xml:"web,omitempty"`
		DB  DB_Conf `json:"db,omitempty" xml:"db,omitempty"`
	}

	HttpClient struct{}
)
