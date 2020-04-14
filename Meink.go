package main

import (
	"Meink/app"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	appName        = "Meink"
	appUsage       = "轻量级静态博客构建工具"
	appVersion     = "1.0.1(bate)"
	appAuthor      = "HanJinyan"
	appAuthorEmail = "jinyanink@outlook.com"
)

/*
函数（功能）入口，也是命令入口 ---> Meink run ,Meink init ....
*/
func main() {
	app := &cli.App{
		Name:    appName,
		Usage:   appUsage,
		Version: appVersion,
		Authors: []*cli.Author{
			{Name: appAuthor, Email: appAuthorEmail},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "运行博客",
				Action: func(c *cli.Context) error {
					app.ParseGlobalConfigForWrap(true)
					app.Build()
					app.DynamicMonitoringFile()
					app.Serve()
					return nil

				},
			},
			{
				Name:  "init",
				Usage: "清空public文件夹(用于调试)",
				Action: func(c *cli.Context) error {
					app.ParseGlobalConfigForWrap(false)
					app.CleanPublic()
					return nil

				},
			},
			{
				Name:  "release",
				Usage: "发布版打包到release文件夹",
				Action: func(c *cli.Context) error {
					app.ParseGlobalConfigForWrap(false)
					app.Release()
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "创建一篇新文章 --> 创建文章用 new article_name  创建页面用 new page_name page ",
				Action: func(c *cli.Context) error {
					app.ParseGlobalConfigForWrap(false)
					app.NewArticle(c)
					app.Sync()
					app.Serve()
					return nil
				},
			},
			{
				Name:  "sync",
				Usage: "同步文章到服务器",
				Action: func(c *cli.Context) error {
					app.ParseGlobalConfigForWrap(false)
					app.Sync()
					app.Serve()
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
