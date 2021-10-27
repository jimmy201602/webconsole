# 欢迎大家踊跃报名，利用业余时间来拓展此项目，增加更多的功能服务大家。

## 常规部署
```
1、下载 go1.17.1.linux-amd64.tar.gz
2、tar zxvf go1.17.1.linux-amd64.tar.gz -C /usr/local
3、cd /opt && git clone https://github.com/jimmy201602/webconsole.git && cd webconsole
4、cd /opt/webconsole
5、make build
6、设置开机自动启动
  cp /opt/webconsole/bin/webconsole  /etc/init.d/   && chmod   755 /etc/init.d/webconsole 
  chkconfig   --add  webconsole  &&  chkconfig webconsole   on  && service webconsole   start
 

```

[下载 Golang](https://storage.googleapis.com/golang/go1.17.1.linux-amd64.tar.gz)


## 容器部署
```
1、下载并安装 Docker
2、cd /opt && git clone https://github.com/jimmy201602/webconsole.git && cd webconsole
3、cd /opt/webconsole
4、make build
5、docker build -t webconsole:latest .
6、docker run -d -p 8080:8080 --restart=always --name webconsole webconsole:latest
```

[安装 Docker](https://docs.docker.com/engine/installation/)


## 更新日志

2018.12.28

    更新：

        1、更新核心js库
        2、合并 pczchen 提交的分支，支持 Docker 容器访问
        3、常规修复

2017.07.31

    更新：

        1、更新 xTerm.js

2017.07.20

    修证：

        1、修证UTF-8字符集展示BUG;
        2、修证大文件展示BUG。
       
     本次BUG由 @AEGQ 修证和贡献代码

2017.04.19

    更新：
        1、删除 build.sh
        2、添加部署说明

2017.04.15

    修证：
        1、修证一些页面的展示BUG。

    更新：
        1、更新xTerm库。
        2、添加外部引用例子。

2017.04.11

    更新：
        1、添加对 JSONP 的支持，解决 JQuery 跨域请求的需求。
        2、配置文件新增 JSONP 的支持。

感谢 @朱小四(weichat:juechengke) 提出需求。

2017.03.27

    修证：
        1、修证粘贴字串的BUG。

    更新：
        1、更新 xTerm 库。

2017.03.14

    修证:
    
        1、修证 wss 协议适配 BUG

    更新：
        1、删除根目录 golang 程序包。如果运行脚本，请自行下载对应 golang 包（go1.8.linux-amd64.tar.gz）放置根目录。

2017.02.27

    修证:
    
        1、修改部分BUG
        2、更新 Golang 到 1.8 版本
        3、修证 Golang 1.8 版本中 url.Parse 解析地址 BUG
        4、替换 term.js 为 xterm.js
        5、更新对CJK（中文、日文、韩文）字符集和IME的支持
        6、更友好地支持 Linux 下主流程序，如：VIM、Tmux等
            

2016.05.25

    修证:
    
        1、修证 log 目录不存在而引起的无法启动程序的 BUG
        2、Dockerfile 无法 Build 项目的 BUG 

2016.05.21

    修证:
    
        1、SSH Session 退出异常 BUG

2016.05.17

    新增:
        添加跨域白名单的支持。(conf/conf.json     cors_white_list:"127.0.0.1,www.xxx.com");

    感谢 @玩蝴蝶的法师 提出的功能建议。

2016.05.13

    修证:
        1、修证 Ubuntu 下无法使用 VIM 的 BUG；

    替换:
        2、更换 Gorilla 的 Websocket 库。

    感谢 @玩蝴蝶的法师 提出的建议和BUG。

2016.03.05

    新增:
        1、增加 Dockerfile

2016.03.04 发布 v1.1

    新增:
        1、增加命令行的启动、停止、状态查看功能. 如: ./apibox start/stop/status

2016.03.03

    修证:
        1、修改独立模式时登陆JS验证问题.

    新增:
        1、增加后台运行模式(conf/conf.json.  daemon:true/false);
        2、增加程序运行时的PID文件(log/apibox.pid);

## 大概的数据流向：
```
    浏览器--》WebSocket--》SSH--》Linux OS
```
    

## 程序包结构：

```
├── bin
│   └── apibox
├── conf
│   ├── ssl_cert.crt (默认不存在)
│   ├── ssl_cert.key (默认不存在)
│   ├── conf.json
│   └── mime.types
├── log
├── pkg
├── static
│   ├── images
│   └── scripts
└── template


运行环境要求：

1、Linux Kernel 3.x/x86_64 及更高版本
2、建议 Linux 发行版 CentOS 7.0+ / Ubuntu 14.04+
3、启用支持 SSL/TLS 模式访问时，需要生成对应的 SSL 证书文件且放置到 "conf" 文件夹下并配置 "conf.json" 文件
4、客户端要求使用 IE9、Chrome 40、Firefox 38、Safari 9 或更高版本的浏览器访问
5、服务器端需要对防火墙开启对应的外部访问端口。具体需要开放的端口请参考 "conf" 文件夹下的 "conf.json" 文件中的端口部分的配置
6、本程序只能对 Unix/Linux 类的操作系统，且支持 SSH 协议的 OS 进行远程操作
7、基于Go1.6+，原生支持 http2 


一、部署
1、将程序解压或下载至任一目录，运行 "build.sh" 编译可执行文件
2、然后运行 "bin" 文件夹下的 "apibox" 文件即可。如:"./apibox start/stop"
3、配置文件在 "conf" 文件夹下，核心配置文件为 "conf.json"，部署的时候需要添加跨域白名单来支持其他机器的访问
4、后台运行可以配置 "conf" 文件夹下的 "conf.json" 文件,将 "daemon" 项配置为 "true" 
5、运行时日志文件存放在 "log" 文件夹下，以当天时期命名
6、也可以配置程序以 Nginx 的 fcgi 模式运行,以 Nginx 做为访问入口
```


二、使用
1、程序部署完成后，直接通过浏览器访问即可。如: http(s)://ip:port。

2、外部系统引用步骤：
    
    1）、以 GET 或 POST 的方式，提交 "vm_addr" 参数至 "http(s)://ip:port/console/chksshdaddr" ,成功后可获取到加密后的 "en_addr" 信息。注：vm_addr 格式为: "ip:port" ，若不携带端口，默认端口为 "22" 。

        获取到的结果为 JSON 格式(注:以下数据为测试数据)：

            成功：
                {
                    "ok": true,
                    "msg": "",
                    "data": { "en_addr": "0b-nDgcazQKTmUw4oBLfxott", "sshd_addr": "192.168.220.173:22" }
                }

            失败：
                { "ok": false, "msg": "Unable to resolve host address.", "data": null }

    2）、成功获取到加密的 en_addr 信息后，以 GET 或 POST 方式访问  "http(s)://ip:port/console/login/'en_addr'" 即可。注: "en_addr" 是通过第 1）步操作获取的数据。

## JQuery Demo:

        
### 第一种方式（需要二次登陆，同一个域的情况，同样可以使用跨域的方式访问）：


```javascript
var protocol = (location.protocol === "https:") ? "https://" : "http://";
var addr = protocol + location.hostname + ((location.port) ? (":" + location.port) : "")

$.post(addr+"/console/chksshdaddr?rnd=" + Math.random(), {
    "vm_addr": "192.168.220.177:22"
}, function(data) {
    var json = data;
    if (typeof(data) != "object") {
        json = $.parseJSON(data);
    }
    if (json.ok) {
        location.href = addr + json.data.sshd_addr;
    }
});
```



### 第二种方式（直接输入远端的主机地址、用户名、密码，然后直接登陆。跨域的情况）：

```html
<button class="btn btn-primary" onclick="testDemo();">Test</button>

<script type="text/javascript">
    //注意，如果远程主机的访问地址是以IP形式出现的，可以忽略此步骤，直接调用 login 方法。具体登陆地址端口请根据实际情况更改.
    var testDemo = function() {
        $.ajax({
            url: "http://a.com:8081/console/chksshdaddr?rnd=" + Math.random(),
            method: "POST",
            data: {
                "vm_addr": "172.16.18.223:22",
            },
            dataType: "jsonp",
            cache: false,
        }).done(function(data) {
            var json = data;
            if (typeof(data) != "object") {
                json = $.parseJSON(data);
            }
            if (json.ok) {
                console.log("en_addr:", json.data.sshd_addr);
                login(json.data.sshd_addr, "shibingli", "cloud123456");
            } else {
                alert("登陆失败，请确认您的主机信息。");
            }
        }).fail(function() {
            alert("未知失败，请联系管理员。");
        });
    };


    //注意，如果远程主机的访问地址是以域名或主机形式出现的，可以执行上面步骤后，再调用本方法。具体登陆地址端口请根据实际情况更改.
    var login = function(enVMAddr, username, password) {
        $.ajax({
            url: "http://a.com:8081/console/login?rnd=" + Math.random(),
            method: "POST",
            data: {
                "vm_addr": enVMAddr,
                "user_name": username,
                "user_pwd": password
            },
            dataType: "jsonp",
            cache: false,
        }).done(function(data) {
            var json = data;
            if (typeof(data) != "object") {
                json = $.parseJSON(data);
            }
            if (json.ok) {
                location.href = "http://a.com:8081" + json.data;
            } else {
                alert("登陆失败，请确认您的登陆信息。");
            }
        }).fail(function() {
            alert("未知失败，请联系管理员。");
        });
    };
</script>
```
