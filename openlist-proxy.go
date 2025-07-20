package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"

	"github.com/OpenListTeam/OpenList/v4/pkg/sign"
)

type Link struct {
	Url    string      `json:"url"`
	Header http.Header `json:"header"`
}

type LinkResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Link   `json:"data"`
}

var (
	port              int
	https             bool
	help              bool
	showVersion       bool
	disableSign       bool
	certFile, keyFile string
	address, token    string
	s                 sign.Sign
	version           string = "dev"
)

func init() {
	flag.IntVar(&port, "port", 5243, "the proxy port.")
	flag.BoolVar(&https, "https", false, "use https protocol.")
	flag.BoolVar(&help, "help", false, "show help")
	flag.BoolVar(&showVersion, "version", false, "show version and exit")
	flag.BoolVar(&disableSign, "disable-sign", false, "disable signature verification")
	flag.StringVar(&certFile, "cert", "server.crt", "cert file")
	flag.StringVar(&keyFile, "key", "server.key", "key file")
	flag.StringVar(&address, "address", "", "openlist address")
	flag.StringVar(&token, "token", "", "openlist token")
	flag.Parse()

	s = sign.NewHMACSign([]byte(token))
}

var HttpClient = &http.Client{}

type Json map[string]interface{}

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func errorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("content-type", "text/json")
	res, _ := json.Marshal(Result{Code: code, Msg: msg})
	w.WriteHeader(200)
	_, _ = w.Write(res)
}

func downHandle(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	
	// If signature verification is not disabled, perform signature verification
	if !disableSign {
		sign := r.URL.Query().Get("sign")
		err := s.Verify(filePath, sign)
		if err != nil {
			errorResponse(w, 401, err.Error())
			return
		}
	}
	
	data := Json{
		"path": filePath,
	}
	dataByte, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/fs/link", address), bytes.NewBuffer(dataByte))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	res, err := HttpClient.Do(req)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
	defer func() {
		_ = res.Body.Close()
	}()
	dataByte, err = io.ReadAll(res.Body)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
	var resp LinkResp
	err = json.Unmarshal(dataByte, &resp)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
	if resp.Code != 200 {
		errorResponse(w, resp.Code, resp.Message)
		return
	}
	if !strings.HasPrefix(resp.Data.Url, "http") {
		resp.Data.Url = "http:" + resp.Data.Url
	}
	fmt.Println("proxy:", resp.Data.Url)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
	req2, _ := http.NewRequest(r.Method, resp.Data.Url, nil)
	maps.Copy(req2.Header, r.Header)
	maps.Copy(req2.Header, resp.Data.Header)
	res2, err := HttpClient.Do(req2)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
	defer func() {
		_ = res2.Body.Close()
	}()
	res2.Header.Del("Access-Control-Allow-Origin")
	res2.Header.Del("set-cookie")
	maps.Copy(w.Header(), res2.Header)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "range")
	w.WriteHeader(res2.StatusCode)
	_, err = io.Copy(w, res2.Body)
	if err != nil {
		errorResponse(w, 500, err.Error())
		return
	}
}

func main() {
	if help {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Println("Version:", version)
		return
	}

	fmt.Printf("OpenList-Proxy - %s\n", version)
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("listen and serve: %s\n", addr)

	srv := http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(downHandle),
	}

	if !https {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("failed to start: %s\n", err.Error())
		}
	} else {
		if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
			fmt.Printf("failed to start: %s\n", err.Error())
		}
	}
}
