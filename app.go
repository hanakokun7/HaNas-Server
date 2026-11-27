package main

import (
	_ "embed"
	"fmt"
	_ "image/gif"
	_ "image/png"
	"net/http"
)

//go:embed index.html
var indexHtmlContent string

//go:embed assets/index.js
var js_script string

//go:embed assets/index.css
var css_style string

//go:embed assets/icon.ico
var favicon []byte

//go:embed assets/apple-touch-icon.png
var appleTouchIcon []byte

//go:embed assets/config.js
var configJS string

func main() {
	config_js := Init()
	if config_js == "" {
		return
	}
	go gomikago_akeru()
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/me", authMiddleware(Me))
	http.HandleFunc("/file/", authMiddleware(GetFile))
	http.HandleFunc("/thumbnail/", authMiddleware(GetThumbnail))
	http.HandleFunc("/node/", authMiddleware(GetJson))
	http.HandleFunc("/upload", authMiddleware(UpFile))
	http.HandleFunc("/copy", authMiddleware(CpFile))
	http.HandleFunc("/move", authMiddleware(MvFile))
	http.HandleFunc("/rename", authMiddleware(RnFile))
	http.HandleFunc("/delete", authMiddleware(DlFile))
	http.HandleFunc("/trash", authMiddleware(get_trash_list))
	http.HandleFunc("/trash/in", authMiddleware(tDlFile))
	http.HandleFunc("/trash/restore/", authMiddleware(restore_trash))
	http.HandleFunc("/trash/info/", authMiddleware(get_trash_info))
	http.HandleFunc("/share/create", authMiddleware(CreateShare))
	http.HandleFunc("/share/delete", authMiddleware(DeleteShare))
	http.HandleFunc("/s/", GetSharedFile)
	http.HandleFunc("/delete-account", authMiddleware(DeleteAccount))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHtmlContent))
	})
	http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(js_script))
	})
	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(css_style))
	})
	http.HandleFunc("/icon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(favicon)
	})
	http.HandleFunc("/apple-touch-icon.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(appleTouchIcon)
	})
	http.HandleFunc("/config.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(config_js))
	})
	if acfg.Https_f {
		go func() {
			fmt.Printf("%s server started at :8443\n", "HaNas")
			http.ListenAndServeTLS(":8443", acfg.Cert, acfg.Key, nil)
		}()
	}
	fmt.Printf("%s server started at :80\n", "HaNas")
	http.ListenAndServe(":80", nil)
}
