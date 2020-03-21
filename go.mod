module github.com/malice-plugins/windows-defender

go 1.14

require (
	github.com/fatih/structs v1.1.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.6.2
	github.com/konsorten/go-windows-terminal-sequences v1.0.1
	github.com/mailru/easyjson v0.0.0-20180823135443-60711f1a8329
	github.com/malice-plugins/pkgs v0.0.0-20190107161315-79532f02e4f0
	github.com/mattn/go-runewidth v0.0.4
	github.com/moul/http2curl v1.0.0
	github.com/olivere/elastic v6.2.15+incompatible
	github.com/parnurzeal/gorequest v0.2.15
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.3.0
	github.com/urfave/cli v1.20.0
	golang.org/dl v0.0.0-20200319204010-bf12898a6070 // indirect
	golang.org/x/crypto v0.0.0-20190103213133-ff983b9c42bc
	golang.org/x/net v0.0.0-20190107155100-1a61f4433d85
	golang.org/x/sys v0.0.0-20190107070147-cb59ee366067
	golang.org/x/text v0.3.0
)

replace github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.4.2