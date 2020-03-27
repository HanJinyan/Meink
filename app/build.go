package app

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ArticleInfo struct {
	DetailDate int64
	Date       string
	Title      string
	Link       string
	Top        bool
}
type Archive struct {
	Year     string
	Articles MSort
}

type Tag struct {
	Name     string
	Count    int //计数
	Articles MSort
}

var themePath, publicPath, sourcePath string
var articleTpl, pageTpl, archiveTpl, tagTpl template.Template

//把parse 解析出来的信息绑定到 .HTML模板上
func Build() {
	//开始时间
	startTime := time.Now()
	//排序
	var (
		articles   = make(MSort, 0)
		pages      = make(MSort, 0)
		tagMap     = make(map[string]MSort)
		archiveMap = make(map[string]MSort)
	)
	//加入目录的path
	themePath = filepath.Join(rootPath, globalConfig.Site.Theme)
	publicPath = filepath.Join(rootPath, globalConfig.Build.Output)
	sourcePath = filepath.Join(rootPath, globalConfig.Build.Source)

	// _.html 插入 *.html
	var partialTpl string
	files, _ := filepath.Glob(filepath.Join(themePath, "*.html"))
	for _, path := range files {
		fileExt := strings.ToLower(filepath.Ext(path))
		baseName := strings.ToLower(filepath.Base(path))
		if fileExt == ".html" && strings.HasPrefix(baseName, "_") {
			html, err := ioutil.ReadFile(path)
			if err != nil {
				MFatal(err.Error())

			}
			tplName := strings.TrimPrefix(baseName, "_")
			tplName = strings.TrimSuffix(tplName, ".html")
			htmlStr := "{{define \"" + tplName + "\"}}" + string(html) + "{{end}}"
			partialTpl = partialTpl + htmlStr
		}

	}
	//编译模板文件
	articleTpl = CompileTpl(filepath.Join(themePath, "article.html"), partialTpl, "article")
	pageTpl = CompileTpl(filepath.Join(themePath, "page.html"), partialTpl, "page")
	archiveTpl = CompileTpl(filepath.Join(themePath, "archive.html"), partialTpl, "archive")
	tagTpl = CompileTpl(filepath.Join(themePath, "tag.html"), partialTpl, "tag")
	//清除部分模板文件
	publicPath = filepath.Join(rootPath, globalConfig.Build.Output)
	cleanPatterns := []string{"tag", "*.html", "msic"}
	for _, pattern := range cleanPatterns {
		files, _ := filepath.Glob(filepath.Join(publicPath, pattern))
		for _, path := range files {
			os.RemoveAll(path)
			MLog("Cleaning: " + path)
		}
	}
	//查找所有.md以生成文章
	SymWalk(sourcePath, func(path string, info os.FileInfo, err error) error {
		fileExt := strings.ToLower(filepath.Ext(path))
		if fileExt == ".md" {
			//解析 markdown 文件
			article := ParseArticle(path)
			//文章不为空 ，不是草稿， 不是隐藏 就渲染
			if article == nil || article.Draft || article.Hide == true {
				return nil
			}
			MLog("Building: " + article.Title)
			//创建生成的html文件
			derectory := filepath.Dir(article.Link)
			err := os.MkdirAll(filepath.Join(publicPath, derectory), 0777)
			if err != nil {
				MFatal(err.Error())
			}
			//文章属性为页面，就创建一个新的 HTML文件
			if article.Type == "page" {
				pages = append(pages, *article)
				return nil
			}
			//隐藏文章不渲染
			if article.Hide == false {
				articles = append(articles, *article)
			}
			//添加 tag
			for _, tag := range article.Tags {
				if _, ok := tagMap[tag]; !ok {
					tagMap[tag] = make(MSort, 0)
				}
				tagMap[tag] = append(tagMap[tag], *article)
			}
			dateYear := article.Time.Format("2006")
			if _, ok := archiveMap[dateYear]; !ok {
				archiveMap[dateYear] = make(MSort, 0)
			}
			articleInfo := ArticleInfo{
				DetailDate: article.Date,
				Date:       article.Time.Format("2006-01-02"),
				Title:      article.Title,
				Link:       article.Link,
				Top:        article.Top,
			}
			archiveMap[dateYear] = append(archiveMap[dateYear], articleInfo)
		}
		return nil
	})

	//逐个渲染页面
	if len(articles) == 0 {
		MFatal("必须要有一篇文章")
	}

	sort.Sort(articles)
	//生成 page 页面
	WaitGroup.Add(1)
	go RenderIndexPage("", articles, "")

	// 生成 Article 页面
	WaitGroup.Add(1)
	go RenderArticles(articleTpl, articles)

	// 在首页点击文章标签可以到另一个页面
	for tagName, articles := range tagMap {
		sort.Sort(articles)
		WaitGroup.Add(1)
		go RenderIndexPage(filepath.Join("tag", tagName), articles, tagName)
	}

	//生成 archive 页面并绑定数据
	archives := make(MSort, 0)
	for year, articleInfo := range archiveMap {
		sort.Sort(articleInfo)
		archives = append(archives, Archive{
			Year:     year,
			Articles: articleInfo,
		})
	}
	sort.Sort(archives)
	WaitGroup.Add(1)
	go RenderPage(archiveTpl, map[string]interface{}{
		"Total":   len(articles),
		"Archive": archives,
		"Site":    globalConfig.Site,
		"I18n":    globalConfig.I18n,
	}, filepath.Join(publicPath, "archive.html"))

	//	生成Teg 页面 点击tag 转到对应tag的列表页
	tags := make(MSort, 0)
	for tagName, tagArticles := range tagMap {
		articleInfo := make(MSort, 0)
		for _, article := range tagArticles {
			articleValue := article.(Article)
			articleInfo = append(articleInfo, ArticleInfo{
				DetailDate: articleValue.Date,
				Date:       articleValue.Time.Format("2006-01-02"),
				Title:      articleValue.Title,
				Link:       articleValue.Link,
				Top:        articleValue.Top,
			})
		}
		sort.Sort(articleInfo)
		tags = append(tags, Tag{
			Name:     tagName,
			Count:    len(tagArticles),
			Articles: articleInfo,
		})
	}
	sort.Sort(MSort(tags))
	WaitGroup.Add(1)
	go RenderPage(tagTpl, map[string]interface{}{
		"Total": len(articles),
		"Tag":   tags,
		"Site":  globalConfig.Site,
		"I18n":  globalConfig.I18n,
	}, filepath.Join(publicPath, "tag.html"))

	//生成 RSS页面
	WaitGroup.Add(1)
	go GenerateRss(articles)

	// 生成其他页面 如auout.html
	WaitGroup.Add(1)
	go RenderArticles(articleTpl, pages)

	Copy()
	WaitGroup.Wait()
	endTime := time.Now()
	usedTime := endTime.Sub(startTime)
	fmt.Printf("\nFinished to build in public folder (%v)\n", usedTime)

}

//连接并编译模板
func CompileTpl(tplPath string, partialTpl string, name string) template.Template {
	htmlTpl, err := ioutil.ReadFile(tplPath)
	if err != nil {
		MFatal("连接并编译模板出错")
	}
	//把模板合并
	htmlStr := string(htmlTpl) + partialTpl
	//插入 I18n数据
	funcMap := template.FuncMap{
		"i18n": func(val string) string {
			return globalConfig.I18n[val]
		},
	}
	tpl, err := template.New(name).Funcs(funcMap).Parse(htmlStr)
	if err != nil {
		MFatal(err.Error())
	}
	return *tpl
}
