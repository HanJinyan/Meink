package tool

import (
	"Meink/app/parse"
	"Meink/app/system"
	"Meink/app/util"
	"bufio"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// 创建文章 或页面
func NewArticle(c *cli.Context) {
	var articleNewTemplate, blogTitle, fileName, articleType string
	sourcePath := parse.SiteConfig().Build.Source + "/"
	confpath := "/app/conf/post.conf"
	confPath := system.RunPath() + strings.Replace(confpath, "-/", sourcePath, -1)
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
			util.MFatal("新建文章缺少名称")
		}
	}
	//	默认为post
	if args.Len() > 1 {
		articleType = args.Slice()[1] // page
		pageTpl, err := ioutil.ReadFile(confPath)
		if err != nil {
			util.MFatal("Page模板错误" + err.Error())
		}
		articleNewTemplate = string(pageTpl)
	} else {
		postTpl, err := ioutil.ReadFile(confPath)
		if err != nil {
			util.MFatal("Post模板错误" + err.Error())
		}
		articleType = "post"
		articleNewTemplate = string(postTpl)
	}

	dateString := date.Format("2006-01-02 15:04:05")
	folderName := time.Now().Format("2006-01-02")
	folderPath := filepath.Join(sourcePath, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(folderPath, os.ModePerm) //0777也可以os.ModePerm
		os.Chmod(folderPath, os.ModePerm)
	}
	fileName = blogTitle + ".md"
	filePath := sourcePath + string(folderName) + "/" + fileName

	file, err := os.Create(filePath)
	if err != nil {
		util.MFatal(err)
	}

	postTemplate, err := template.New("post").Parse(articleNewTemplate)
	if err != nil {
		util.MFatal(err)
	}
	data := map[string]string{
		"Title":      blogTitle,
		"DateString": dateString,
		"Author":     author,
		"Draft":      draft,
		"Top":        top,
		"Type":       articleType,
		"Hide":       hide,
		"Preview":    preview,
		"Cover":      cover,
		"Tags":       "",
	}
	fileWrite := bufio.NewWriter(file)
	err = postTemplate.Execute(fileWrite, data)
	if err != nil {
		util.MFatal(err)
	}
	err = fileWrite.Flush()
	if err != nil {
		util.MFatal(err)
	}
}
