package app

import (
	"github.com/gorilla/feeds"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type RenderArticle struct {
	Article
	Next *Article
	Prev *Article
}

var WaitGroup sync.WaitGroup //并发
//按数据呈现html文件
func RenderPage(tpl template.Template, tplData interface{}, outPath string) {
	outFile, err := os.Create(outPath)
	if err != nil {
		MFatal(err.Error())
	}
	defer func() {
		outFile.Close()
	}()
	defer WaitGroup.Done()
	// 模板渲染
	err = tpl.Execute(outFile, tplData)
	if err != nil {
		MFatal(err.Error())
	}
}

//单个文章
func RenderArticles(tpl template.Template, articles MSort) {
	defer WaitGroup.Done()
	//文章数
	articleCount := len(articles)
	for i, _ := range articles {
		//实现上一篇 ，下一篇文章
		currentArticle := articles[i].(Article)
		var renderArticle = RenderArticle{currentArticle, nil, nil}
		if i >= 1 {
			article := articles[i-1].(Article)
			renderArticle.Prev = &article
		}
		if i <= articleCount-2 {
			article := articles[i+1].(Article)
			renderArticle.Next = &article
		}
		outPath := filepath.Join(publicPath, currentArticle.Link)
		WaitGroup.Add(1)
		go RenderPage(tpl, renderArticle, outPath)
	}
}
func RenderIndexPage(rootPath string, articles MSort, tagName string) {
	defer WaitGroup.Done()
	pagePath := filepath.Join(publicPath, rootPath)
	os.MkdirAll(pagePath, 0777)
	// 分页
	limit := globalConfig.Site.Limit
	total := len(articles)
	page := total / limit
	rest := total % limit
	if rest != 0 {
		page++
	}
	if total < limit {
		page = 1
	}
	for i := 0; i < page; i++ {
		var prev = filepath.Join(rootPath, "page"+strconv.Itoa(i)+".html")
		var next = filepath.Join(rootPath, "page"+strconv.Itoa(i+2)+".html")
		outPath := filepath.Join(pagePath, "index.html")
		if i != 0 {
			fileName := "page" + strconv.Itoa(i+1) + ".html"
			outPath = filepath.Join(pagePath, fileName)
		} else {
			prev = ""
		}
		if i == 1 {
			prev = filepath.Join(rootPath, "index.html")
		}
		first := i * limit
		count := first + limit
		if i == page-1 {
			if rest != 0 {
				count = first + rest
			}
			next = ""
		}
		//	绑定数据
		var data = map[string]interface{}{
			"Articles": articles[first:count],
			"Site":     globalConfig.Site,
			"Develop":  globalConfig.Develop,
			"Page":     i + 1,
			"Total":    page,
			"Prev":     prev,
			"Next":     next,
			"TagName":  tagName,
			"TagCount": len(articles),
		}
		WaitGroup.Add(1)
		go RenderPage(pageTpl, data, outPath)
	}
}
//生成RSS页面
func GenerateRss(articles MSort)  {
	defer WaitGroup.Done()
	var feedArticle MSort
	if len(articles) < globalConfig.Site.Limit{
		feedArticle = articles
	}else {
		feedArticle = articles[0:globalConfig.Site.Limit]
	}
	if globalConfig.Site.URL != ""{
		feed :=&feeds.Feed{
			Title:       globalConfig.Site.Title,
			Link:        &feeds.Link{Href:globalConfig.Site.URL},
			Description: globalConfig.Site.Subtitle,
			Author:      &feeds.Author{globalConfig.Site.Title,""},
			Created:     time.Time{},
		}
		feed.Items = make([]*feeds.Item , 0)
		for _ , item :=range feedArticle{
			artilce := item.(Article)
			feed.Items = append(feed.Items , &feeds.Item{
				Title:       artilce.Title,
				Link:        &feeds.Link{Href:globalConfig.Site.URL},
				Author:      &feeds.Author{artilce.Author.Name,""},
				Description: string(artilce.Preview),
				Updated:     artilce.MTime,
				Created:     artilce.Time,
			})
		}
		if rss , err := feed.ToAtom() ; err == nil {
			err := ioutil.WriteFile(filepath.Join(publicPath , "rss.xml"),[]byte(rss) ,0644)
		if err !=nil {
			MFatal(err.Error())
		}
		}else {
			MFatal(err.Error())
		}
	}
}