package main

import (
	"Meink/app/build"
	"Meink/app/parse"
	"Meink/app/server"
	"Meink/app/tool"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	appInfo := parse.SiteConfig().App
	app := &cli.App{
		Name:    appInfo.Name,
		Version: appInfo.Version,
		Authors: []*cli.Author{
			{Name: appInfo.Author, Email: appInfo.Email},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "运行博客",
				Action: func(c *cli.Context) error {
					parse.I18n()
					build.Build()
					tool.Copy()
					tool.DynamicMonitoringFile()
					server.Server()
					return nil

				},
			},
			{
				Name:  "init",
				Usage: "清空Public文件夹",
				Action: func(c *cli.Context) error {
					tool.CleanPublic()
					return nil

				},
			},
			{
				Name:  "new",
				Usage: "创建一篇新文章 --> 创建文章用 new article_name  创建页面用 new page_name page ",
				Action: func(c *cli.Context) error {
					tool.NewArticle(c)
					return nil
				},
			},
		},
	}
	app.Run(os.Args)
}
