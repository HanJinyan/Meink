package server

import (
	"Meink/app/httprouter"
	"Meink/app/parse"
	"Meink/app/system"
	"Meink/app/util"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

func Server() {
	runPath := system.RunPath()
	siteConfig := parse.SiteConfig()
	Port := strconv.Itoa(siteConfig.App.Port)
	router := httprouter.New()
	router.GET("/websocket", Websocket)
	router.NotFound = http.FileServer(http.Dir(filepath.Join(runPath, siteConfig.Build.Output)))
	uri := "http://localhost:" + Port + "/"
	util.MLog("Service running " + uri)
	log.Fatal(http.ListenAndServe(":"+Port, router))
}
