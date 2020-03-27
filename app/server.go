package app

import "path/filepath"

func Serve() {
	service := NewHttpServer()
	service.Get("/websocket", Websocket)
	service.Get("*", Static(filepath.Join(rootPath, globalConfig.Build.Output)))
	uri := "http://localhost:" + globalConfig.Build.Port + "/"
	MLog("Open " + uri + " to  preview")
	service.Listen(":" + globalConfig.Build.Port)
}
