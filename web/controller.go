package web

import (
	"net/http"
	"strings"

	"../assassin"
	"../poc"
	"../seeker"
	"github.com/AmyangXYZ/sweetygo"
	"github.com/gorilla/websocket"
)

func index(ctx *sweetygo.Context) {
	ctx.Render(200, "index")
}

func static(ctx *sweetygo.Context) {
	staticHandle := http.StripPrefix("/static",
		http.FileServer(http.Dir("./web/static")))
	staticHandle.ServeHTTP(ctx.Resp, ctx.Req)
}

func newAssassin(ctx *sweetygo.Context) {
	target := ctx.Param("target")
	if strings.Contains(target, ",") == false {
		a = assassin.New(target)
		ret := map[string]string{
			"target": target,
		}
		ctx.JSON(201, ret, "success")
		return
	}
	targets := strings.Split(target, ",")
	for _, t := range targets {
		ateam = append(ateam, assassin.New(t))
	}

	ret := map[string][]string{
		"targets": targets,
	}
	ctx.JSON(201, ret, "success")
}

func basicInfo(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["basicInfo"].Set(conn, a.Target)
	a.Gatherers["basicInfo"].Run()
	conn.Close()
}

func whois(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["whois"].Set(conn, a.Target)
	a.Gatherers["whois"].Run()
	conn.Close()
}

func tracert(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["tracert"].Set(conn, a.Target)
	a.Gatherers["tracert"].Run()
	conn.Close()
}

func cmsDetect(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["cms"].Set(conn, a.Target)
	a.Gatherers["cms"].Run()
	conn.Close()
}

func portScan(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["port"].Set(conn, a.Target, "tcp")
	a.Gatherers["port"].Run()
	conn.Close()
}

type dirbMsg struct {
	// Payload   string `json:"payload"`
	GortCount int `json:"gort_count"`
}

func dirBrute(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := dirbMsg{}
	conn.ReadJSON(&m)
	a.Gatherers["dirb"].Set(conn, a.Target, m.GortCount)
	a.Gatherers["dirb"].Run()
	conn.Close()
}

func crawl(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Gatherers["crawl"].Set(conn, a.Target, 4)
	a.Gatherers["crawl"].Run()
	a.FuzzableURLs = a.Gatherers["crawl"].Report().([]string)
}

func checkSQLi(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["sqli"].Set(conn, a.FuzzableURLs)
	a.Attackers["sqli"].Run()
	conn.Close()
}

func checkXSS(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	a.Attackers["xss"].Set(conn, a.FuzzableURLs)
	a.Attackers["xss"].Run()
	conn.Close()
}

type intruderMsg struct {
	Header    string `json:"header"`
	Payload   string `json:"payload"`
	GortCount int    `json:"gort_count"`
}

func intrude(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := intruderMsg{}
	conn.ReadJSON(&m)
	a.Attackers["intruder"].Set(conn, a.Target, m.Header, m.Payload, m.GortCount)
	a.Attackers["intruder"].Run()
	conn.Close()
}

type seekerMsg struct {
	Query   string `json:"query"`
	SE      string `json:"se"`
	MaxPage int    `json:"max_page"`
}

func seek(ctx *sweetygo.Context) {
	conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	m := seekerMsg{}
	conn.ReadJSON(&m)
	S := seeker.NewSeeker(m.Query, m.SE, m.MaxPage)
	S.Run(conn)
	conn.Close()
}

func getPoCs(ctx *sweetygo.Context) {
	var pocList []string
	for pocNames := range poc.PoCMap {
		pocList = append(pocList, pocNames)
	}

	ret := map[string][]string{
		"poclist": pocList,
	}
	ctx.JSON(200, ret, "success")
}

// POST -d "poc=xxx"
func setPoC(ctx *sweetygo.Context) {
	pocName := ctx.Param("poc")
	for _, aa := range ateam {
		aa.PoC = poc.PoCMap[pocName]
	}

	ret := map[string]string{
		"poc": pocName,
	}
	ctx.JSON(201, ret, "success")
}

func runPoC(ctx *sweetygo.Context) {
	// conn, _ := websocket.Upgrade(ctx.Resp, ctx.Req, ctx.Resp.Header(), 1024, 1024)
	// concurrency := 2
	// blockers := make(chan struct{}, concurrency)
	// var existedList []string

	// for _, aa := range ateam {
	// 	blockers <- struct{}{}
	// 	go func(a *assassin.Assassin, blocker chan struct{}) {
	// 		defer func() { <-blocker }()
	// 		a.POC.Run(conn)
	// 		if result := a.POC.Report().(string); result == "true" {
	// 			existedList = append(existedList, a.Target)
	// 		}
	// 	}(aa, blockers)
	// }
	// for i := 0; i < cap(blockers); i++ {
	// 	blockers <- struct{}{}
	// }

	// ret := map[string][]string{
	// 	"existed": existedList,
	// }
	// ctx.JSON(200, ret, "success")
}
