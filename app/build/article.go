package build

import (
	"Meink/app/modle"
	"html/template"
	"path/filepath"
)

func RenderArticles(tpl template.Template, articles modle.MSort, publicPath string) {
	defer WaitGroup.Done()
	//文章数
	articleCount := len(articles)
	for i, _ := range articles {
		//实现上一篇 ，下一篇文章
		currentArticle := articles[i].(modle.Article)
		var renderArticle = modle.RenderArticle{currentArticle, nil, nil}
		if i >= 1 {
			article := articles[i-1].(modle.Article)
			renderArticle.Prev = &article
		}
		if i <= articleCount-2 {
			article := articles[i+1].(modle.Article)
			renderArticle.Next = &article
		}
		outPath := filepath.Join(publicPath, currentArticle.Link)
		WaitGroup.Add(1)
		go RenderPage(tpl, renderArticle, outPath)
	}
}
