package build

import (
	"Meink/app/modle"
	"Meink/app/parse"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
)

func RenderIndexPage(rootPath, publicPath string, articles modle.MSort, tagName string, tpl template.Template) {
	defer WaitGroup.Done()
	siteConfig := parse.SiteConfig()
	pagePath := filepath.Join(publicPath, rootPath)
	os.MkdirAll(pagePath, 0777)
	// 分页
	limit := siteConfig.Site.Limit
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
			"Site":     siteConfig.Site,
			"Develop":  "",
			"Minfo":    "",
			"Page":     i + 1,
			"Total":    page,
			"Prev":     prev,
			"Next":     next,
			"App":      siteConfig.App,
			"TagName":  tagName,
			"TagCount": len(articles),
		}
		WaitGroup.Add(1)
		go RenderPage(tpl, data, outPath)
	}
}
