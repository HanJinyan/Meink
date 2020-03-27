package app

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

const (
	Indent        = "	" //tab "按一次TAB键"
	Post_Template = `title: {{.Title}}
date: {{.DateString}}
author: {{.Author}}
cover: {{.Cover}}
draft: {{.Draft}}
top: {{.Top}}
type: {{.Type}}
hide: {{.Hide}}
tags: {{.Tags}}
preview: {{.Preview}}
---------------------
Markdown正文
`
	Page_Template = `title: {{.Title}}
author: {{.Author}}
type: {{.Type}}
---------------------
Markdown正文
`
)

// 创建文章 或页面
func NewArticle(c *cli.Context) {
	var articleNewTemplate, blogTitle, fileName, articletType string
	path := globalConfig.Build.Source + "/"

	//默认配置
	draft := "false"
	top := "false"
	hide := "false"
	author := "me"
	date := time.Now()

	// 空字符串值
	preview := ""
	cover := ""

	args := c.Args()
	if c.Args().Len() > 0 {
		blogTitle = args.Slice()[0]
	}
	if blogTitle == "" {
		if c.String("title") != "" {
			blogTitle = c.String("title")
		} else {
			MFatal("新建文章缺少名称")
		}
	}
	//	默认为post
	if args.Len() > 1 {
		articletType = args.Slice()[1] //page
		fmt.Println(articletType)
		articleNewTemplate = Page_Template
	} else {
		articletType = "post"
		articleNewTemplate = Post_Template
	}

	dateString := date.Format(Date_Fromat)
	folderName := time.Now().Format("2006-01-02")
	folderPath := filepath.Join(path, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(folderPath, os.ModePerm) //0777也可以os.ModePerm
		os.Chmod(folderPath, os.ModePerm)
	}
	fileName = blogTitle + ".md"
	filePath := path + string(folderName) + "/" + fileName

	file, err := os.Create(filePath)
	if err != nil {
		MFatal(err)
	}

	postTemplate, err := template.New("post").Parse(articleNewTemplate)
	if err != nil {
		MFatal(err)
	}
	data := map[string]string{
		"Title":      blogTitle,
		"DateString": dateString,
		"Author":     author,
		"Draft":      draft,
		"Top":        top,
		"Type":       articletType,
		"Hide":       hide,
		"Preview":    preview,
		"Cover":      cover,
		"Tags":       "",
	}
	fileWrite := bufio.NewWriter(file)
	err = postTemplate.Execute(fileWrite, data)
	if err != nil {
		MFatal(err)
	}
	err = fileWrite.Flush()
	if err != nil {
		MFatal(err)
	}
}
