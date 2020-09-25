package build

import (
	"Meink/app/modle"
	"Meink/app/parse"
	"Meink/app/system"
	"Meink/app/util"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var WaitGroup sync.WaitGroup

func Build() {
	startTime := time.Now()
	siteConfig := parse.SiteConfig()
	runPath := system.RunPath()
	var (
		themePath, publicPath, sourcePath        string
		articleTpl, indexTpl, archiveTpl, tagTpl template.Template
		articles                                 = make(modle.MSort, 0)
		pages                                    = make(modle.MSort, 0)
		tagMap                                   = make(map[string]modle.MSort)
		archiveMap                               = make(map[string]modle.MSort)
	)
	themePath = filepath.Join(runPath, siteConfig.Site.Theme)
	publicPath = filepath.Join(runPath, siteConfig.Build.Output)
	sourcePath = filepath.Join(runPath, siteConfig.Build.Source)
	//编译模板文件
	articleTpl = CompileTpl(filepath.Join(themePath, "article.html"), MergeTpl(themePath), "article")
	indexTpl = CompileTpl(filepath.Join(themePath, "index.html"), MergeTpl(themePath), "index")
	archiveTpl = CompileTpl(filepath.Join(themePath, "archive.html"), MergeTpl(themePath), "archive")
	tagTpl = CompileTpl(filepath.Join(themePath, "tag.html"), MergeTpl(themePath), "tag")
	//清除public文件夹中的部分模板文件，
	publicPath = filepath.Join(runPath, siteConfig.Build.Output)
	cleanPatterns := []string{"tag", "*.html", "msic", "*.xml", "*.css", "*.js"}
	for _, pattern := range cleanPatterns {
		files, _ := filepath.Glob(filepath.Join(publicPath, pattern))
		for _, path := range files {
			os.RemoveAll(path)
		}
	}
	//查找所有.md文件
	util.SymWalk(sourcePath, func(path string, info os.FileInfo, err error) error {
		fileExt := strings.ToLower(filepath.Ext(path))
		if fileExt == ".md" {
			//解析 markdown 文件
			article := parse.Article(path)
			//排除为空 ,草稿，隐藏 的文章
			if article == nil || article.Draft || article.Hide == true {
				return nil
			}
			util.MLog("Building: " + article.Title)
			//创建生成的html文件的目录
			derectory := filepath.Dir(article.Link)
			err := os.MkdirAll(filepath.Join(publicPath, derectory), 0777)
			if err != nil {
				util.MFatal(err.Error())
			}
			//文章属性为页面，就创建一个新的 HTML文件
			if article.Type == "page" {
				pages = append(pages, *article)
				return nil
			}
			//把符合条件的文章加入渲染队列
			articles = append(articles, *article)

			//添加 tag
			for _, tag := range article.Tags {
				if _, ok := tagMap[tag]; !ok {
					tagMap[tag] = make(modle.MSort, 0)
				}
				tagMap[tag] = append(tagMap[tag], *article)
			}
			dateYear := article.Time.Format("2006")
			if _, ok := archiveMap[dateYear]; !ok {
				archiveMap[dateYear] = make(modle.MSort, 0)
			}
			articleInfo := modle.ArticleInfo{
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
		util.MFatal("必须要有一篇文章")
	}
	//首页
	sort.Sort(articles)
	WaitGroup.Add(1)
	go RenderArticles(articleTpl, articles, publicPath)
	//文章
	WaitGroup.Add(1)
	go RenderArticles(articleTpl, pages, publicPath)
	//首页
	WaitGroup.Add(1)
	go RenderIndexPage("", publicPath, articles, "", indexTpl)

	for tagName, articles := range tagMap {
		sort.Sort(articles)
		WaitGroup.Add(1)
		go RenderIndexPage(filepath.Join("tag", tagName), publicPath, articles, tagName, indexTpl)
	}

	archives := make(modle.MSort, 0)
	for year, articleInfo := range archiveMap {
		sort.Sort(articleInfo)
		archives = append(archives, modle.Archive{
			Year:     year,
			Articles: articleInfo,
		})
	}
	sort.Sort(archives)
	WaitGroup.Add(1)
	go RenderPage(archiveTpl, map[string]interface{}{
		"Total":   len(articles),
		"Archive": archives,
		"Site":    siteConfig.Site,
		"App":     siteConfig.App,
	}, filepath.Join(publicPath, "archive.html"))

	//生成Teg 页面
	tags := make(modle.MSort, 0)
	for tagName, tagArticles := range tagMap {
		articleInfo := make(modle.MSort, 0)
		for _, article := range tagArticles {
			articleValue := article.(modle.Article)
			articleInfo = append(articleInfo, modle.ArticleInfo{
				DetailDate: articleValue.Date,
				Date:       articleValue.Time.Format("2006-01-02"),
				Title:      articleValue.Title,
				Link:       articleValue.Link,
				Top:        articleValue.Top,
			})
		}
		sort.Sort(articleInfo)
		tags = append(tags, modle.Tag{
			Name:     tagName,
			Count:    len(tagArticles),
			Articles: articleInfo,
		})
	}
	sort.Sort(modle.MSort(tags))
	WaitGroup.Add(1)
	go RenderPage(tagTpl, map[string]interface{}{
		"Total": len(articles),
		"Tag":   tags,
		"Site":  siteConfig.Site,
		"App":   siteConfig.App,
	}, filepath.Join(publicPath, "tag.html"))
	WaitGroup.Wait()
	endTime := time.Now()
	usedTime := endTime.Sub(startTime)
	fmt.Printf("\nBuild completed at public folder (%v)\n", usedTime)
}
