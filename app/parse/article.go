package parse

import (
	"Meink/app/modle"
	"Meink/app/util"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

func ArticleConfig(markdownPath string) (articleConfig *modle.ArticleConfig, content string) {
	var (
		configStr    = ""
		Config_Split = "---"
	)
	data, err := ioutil.ReadFile(markdownPath)
	if err != nil {
		util.MFatal("读取文章出错 :" + err.Error())

	}
	contentStr := string(data)
	contentStr = RepldceRootFlag(contentStr)
	markdownStr := strings.SplitN(contentStr, Config_Split, 2)
	contentLen := len(markdownStr)
	if contentLen > 0 {
		configStr = markdownStr[0]
	}
	if contentLen > 1 {
		content = markdownStr[1]
	}
	if err := yaml.Unmarshal([]byte(configStr), &articleConfig); err != nil {
		util.MError("解析文章配置出错" + err.Error())
		return nil, ""
	}
	if articleConfig == nil {
		return nil, ""
	}
	if articleConfig.Type == "" {
		articleConfig.Type = "post"
	}
	previewArry := strings.SplitN(content, "", 2)
	if len(articleConfig.Preview) <= 0 && len(previewArry) > 1 {
		articleConfig.Preview = Markdown(previewArry[0])
		content = strings.Replace(content, "", "", 1)
	} else {
		articleConfig.Preview = Markdown(string(articleConfig.Preview))
	}
	return articleConfig, content
}
func Article(markdownPath string) *modle.Article {
	var (
		article modle.Article
		tags    = map[string]bool{}
	)
	articleConfig, content := ArticleConfig(markdownPath)
	if articleConfig == nil {
		util.MWarn("文章配置信息错误: %s" + markdownPath)
		return nil
	}
	article.Hide = articleConfig.Hide
	article.Type = articleConfig.Type
	article.Preview = articleConfig.Preview
	article.Markdown = content
	article.Content = Markdown(content)
	if articleConfig.Date != "" {
		article.Time = Date(articleConfig.Date) //time.timw
		article.Date = article.Time.Unix()      //int
	}
	if articleConfig.Update != "" {
		article.MTime = Date(articleConfig.Update)
		article.Update = article.MTime.Unix()

	}
	article.Title = articleConfig.Title
	article.Draft = articleConfig.Draft
	article.Top = articleConfig.Top
	if author, ok := SiteConfig().Authors[articleConfig.Author]; ok {
		author.Id = articleConfig.Author
		author.Avatar = RepldceRootFlag(author.Avatar)
		article.Author = author
	}
	if len(articleConfig.Categories) > 0 {
		article.Category = article.Categories[0]

	} else {
		article.Category = "No_Classification"
	}
	article.Tags = articleConfig.Tags
	for _, tag := range articleConfig.Tags {
		tags[tag] = true
	}
	for _, category := range articleConfig.Categories {
		if _, ok := tags[category]; !ok {
			article.Tags = append(article.Tags, category)
		}
	}
	fileName := strings.TrimSuffix(strings.ToLower(filepath.Base(markdownPath)), ".md")
	link := fileName + ".html"
	//另一种模式，以文章时间作为访问路径
	if article.Type == "post" {
		dataPrefix := article.Time.Format("2006-01-02-")
		if strings.HasPrefix(fileName, dataPrefix) {
			fileName = fileName[len(dataPrefix):]
		}
		if SiteConfig().Site.Link != "" {
			linkMap := map[string]string{
				"{year}":     article.Time.Format("2006"),
				"{month}":    article.Time.Format("01"),
				"{day}":      article.Time.Format("02"),
				"{hour}":     article.Time.Format("15"),
				"{minute}":   article.Time.Format("04"),
				"{second}":   article.Time.Format("05"),
				"{category}": article.Category,
				"{title}":    fileName,
			}
			link = SiteConfig().Site.Link
			for key, val := range linkMap {
				link = strings.Replace(link, key, val, -1)
			}
		}
	}

	article.Link = link
	article.SiteConfig = *siteConfig
	return &article
}
