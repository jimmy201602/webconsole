package server

import (
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"
	"strings"
	"time"
	"webconsole/utils"
)

var (
	ABC_Conf, conf_err = utils.Get_Conf()
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if nil != conf_err {
		utils.Log_Fatal(conf_err.Error())
	}
	_ = ABC_Conf.Web.Addr
}

func GetPID() string {
	return utils.ToStr(os.Getpid())
}

func Run() {
	go func() {
		for {
			err := utils.WritePidFile(utils.PidPath, GetPID())
			if nil != err {
				utils.Log_Fatal(err)
				break
			}
			time.Sleep(time.Duration(1 * time.Second))
		}
	}()

	if ABC_Conf.Web.Daemon {
		ret, err := utils.Daemon(0, 0)
		if nil != err && ret == -1 {
			utils.Log_Fatal(err)
			return
		}
	}

	runMsgArr := make([]string, 0, 0)
	runMsgArr = append(runMsgArr, "==>")
	runMsgArr = append(runMsgArr, "Web server running.")
	runMsgArr = append(runMsgArr, "PID:"+GetPID()+",")
	runMsgArr = append(runMsgArr, "Addr:"+ABC_Conf.Web.Addr+".")

	if ABC_Conf.Web.EnableFcgi {
		runMsgArr = append(runMsgArr, "Fcgi:"+utils.ToStr(ABC_Conf.Web.EnableFcgi)+".")
	}

	if ABC_Conf.Web.EnableTLS {
		runMsgArr = append(runMsgArr, "SSL:"+utils.ToStr(ABC_Conf.Web.EnableTLS)+",")
		runMsgArr = append(runMsgArr, "TLS_Addr:"+utils.ToStr(ABC_Conf.Web.TlsAddr)+",")
		runMsgArr = append(runMsgArr, "TLS_Url:"+utils.ToStr(ABC_Conf.Web.TlsUrl)+".")
	}

	runMsg := strings.Join(runMsgArr, " ")

	utils.Log_Info(runMsg)

	if ABC_Conf.Web.EnableFcgi {
		listener, err := net.Listen("tcp", ABC_Conf.Web.Addr)
		if err != nil {
			utils.Log_Fatal(err.Error())
			return
		}
		err = fcgi.Serve(listener, DefaultServeMux)
		if nil != err {
			utils.Log_Fatal(err)
			return
		}
	} else {
		if ABC_Conf.Web.EnableTLS {
			go func() {
				err := http.ListenAndServe(ABC_Conf.Web.Addr, RedirectHandler(ABC_Conf.Web.TlsUrl, http.StatusMovedPermanently))
				if nil != err {
					utils.Log_Fatal(err)
					return
				}
			}()

			certFile := utils.ConfDir + utils.PathSeparator + ABC_Conf.Web.TlsCert
			keyFile := utils.ConfDir + utils.PathSeparator + ABC_Conf.Web.TlsKey

			err := http.ListenAndServeTLS(ABC_Conf.Web.TlsAddr, certFile, keyFile, DefaultServeMux)
			if nil != err {
				utils.Log_Fatal(err)
				return
			}
		} else {
			err := http.ListenAndServe(ABC_Conf.Web.Addr, DefaultServeMux)
			if nil != err {
				utils.Log_Fatal(err)
				return
			}
		}
	}
}
