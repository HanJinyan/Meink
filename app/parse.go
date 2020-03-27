package app

import (
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//所有config的入口 ，同时是一种结构体的嵌入

type GlobalConfig struct {
	Site    SiteConfig
	Build   BuildConfig
	Authors map[string]AuthorConfig //Author 是一个多层的结构
	I18n    map[string]string       //I18n 是一个多层的结构
	Develop bool
}
type SiteConfig struct {
	Title     string
	Subtitle  string
	Limit     int
	Theme     string
	Language  string
	Logo      string
	URL       string
	Root      string
	Version   string
	Copyright string
	Link      string
}
type AuthorConfig struct {
	Id     string //作者ID
	Name   string //作者名称
	Intro  string //昵称
	Avatar string //头像
}
type BuildConfig struct {
	Output  string
	Source  string
	Port    string
	Copy    []string
	Publish string
}
type ThemeConfig struct {
	Copy     []string
	Language map[string]map[string]string
}
type ArticleConfig struct {
	Title      string
	Date       string
	Update     string
	Author     string
	Tags       []string
	Categories []string //分类
	Cover      string
	Draft      bool
	Preview    template.HTML
	Top        bool
	Type       string
	Hide       bool
}
type Article struct {
	GlobalConfig
	ArticleConfig
	Time       time.Time
	MTime      time.Time
	Date       int64
	Update     int64
	PageDate   string
	PageUpdate string
	Author     AuthorConfig
	Category   string
	Tags       []string
	Markdown   string
	Preview    template.HTML
	Content    template.HTML
	Link       string
}

var globalConfig *GlobalConfig
var rootPath string

//解析所有config的总入口 ，传递给命令行使用
func ParseGlobalConfigForWrap(develop bool) {
	//获取当前的执行文件所在的目录 ,类是于 "pwd"
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//如果有错误即获取不到当前目录，打印错误并退出
		MFatal("获取当前的执行文件所在的目录:" + err.Error())
	}
	//传递给ParseGlobalConfig这个函数开始解析
	globalConfig = ParseGlobalConfig(filepath.Join(rootPath, "config.yml"), develop)
	//处理 config.yml 不存在，或者为空的情况，打印错误并退出
	if globalConfig == nil {
		MFatal("没有找到 config.yml 这个文件 ，或者文件为空:" + err.Error())
	}
}
func ParseGlobalConfig(configPath string, develop bool) *GlobalConfig {
	var globalConfig *GlobalConfig
	//读取config.yml 文件
	//上一步判断过config.yml不存在，或者为空的情况，这里可以不处理 ，返回空
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil
	}
	//解码 .yml文件
	if err = yaml.Unmarshal(data, &globalConfig); err != nil {
		MFatal("解码config.yml 文件出错:" + err.Error())
	}
	//配置一些默认信息
	globalConfig.Develop = develop
	if develop {
		globalConfig.Site.Root = ""
	}
	globalConfig.Site.Logo = strings.Replace(globalConfig.Site.Logo, "-/", globalConfig.Site.Logo, 1)
	//判断URL是否有 “/”，如果没有加上
	// 即 www.jinyan.ink123 修改为 www.jinyan.ink/123 ,防止URL写错
	if globalConfig.Site.URL != "" && strings.HasSuffix(globalConfig.Site.URL, "/") {
		globalConfig.Site.URL = strings.TrimSuffix(globalConfig.Site.URL, "/")
	}
	///默认输出目录
	if globalConfig.Build.Output == "" {
		globalConfig.Build.Output = "public"
	}
	if globalConfig.Build.Source == "" {
		globalConfig.Build.Output = "source"
	}
	//	解析Theme的config.yml
	themeConfig := ParseThemeConfig(filepath.Join(rootPath, globalConfig.Site.Theme, "config.yml"))
	for _, copyItem := range themeConfig.Copy {
		//把Theme的config.yml里面要复制的内容添加到 Build.Copy 里面，一起复制
		globalConfig.Build.Copy = append(globalConfig.Build.Copy, filepath.Join(globalConfig.Site.Theme, copyItem))
	}
	//解析Theme.yml 里面的language语言配置
	globalConfig.I18n = make(map[string]string)
	for item, languageItem := range themeConfig.Language {
		globalConfig.I18n[item] = languageItem[globalConfig.Site.Language]
	}
	return globalConfig
}

// 解析theme的config.yml
func ParseThemeConfig(configPath string) *ThemeConfig {
	var themeConfig *ThemeConfig
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		MFatal("theme里面没有找到config.yml文件:" + err.Error())
	}
	//解码theme的config.yml
	if err := yaml.Unmarshal(data, &themeConfig); err != nil {
		MFatal("theme里面的onfig.yml文件存在错误:" + err.Error())

	}
	return themeConfig
}

const (
	//文章中config 与 内容的分割线
	Config_Split = "---"
)

//解析文章的配置信息
func ParseArticleConfig(markdownPath string) (articleConfig *ArticleConfig, content string) {
	var configStr string
	data, err := ioutil.ReadFile(markdownPath)
	if err != nil {
		MFatal("文章的配置信息读取:" + err.Error())

	}
	//内容数据流
	contentStr := string(data)
	contentStr = RepldceRootFlag(contentStr)
	//处理文章配置信息和内容 ，以"---"为i分割线
	markdownStr := strings.SplitN(contentStr, Config_Split, 2)
	//分割线以上为文章配置信息 ，以下为文章内容
	contentLen := len(markdownStr)
	if contentLen > 0 {
		configStr = markdownStr[0]
	}
	if contentLen > 1 {
		content = markdownStr[1]
	}

	//解析文章的配置信息
	if err := yaml.Unmarshal([]byte(configStr), &articleConfig); err != nil {
		MError(err.Error())
		return nil, ""
	}
	if articleConfig == nil {
		return nil, ""
	}
	//默认文章为post类型
	if articleConfig.Type == "" {
		articleConfig.Type = "post"
	}
	//文章的preview
	previewArry := strings.SplitN(content, "", 2)
	if len(articleConfig.Preview) <= 0 && len(previewArry) > 1 {
		articleConfig.Preview = ParseMarkdown(previewArry[0])
		content = strings.Replace(content, "", "", 1)
	} else {
		articleConfig.Preview = ParseMarkdown(string(articleConfig.Preview))
	}
	return articleConfig, content
}

//
func RepldceRootFlag(content string) string {
	return strings.Replace(content, "-/", globalConfig.Site.Root+"/", -1)
}
func ParseMarkdown(markdown string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(markdown)))
}

//解析文章
func ParseArticle(markdownPath string) *Article {
	var article Article
	articleConfig, content := ParseArticleConfig(markdownPath)
	if articleConfig == nil {
		MWarn("文章配置信息错误: %s" + markdownPath)
		return nil
	}
	//绑定信息，后面会使用到
	article.Hide = articleConfig.Hide
	article.Type = articleConfig.Type
	article.Preview = articleConfig.Preview
	article.Markdown = content
	article.Content = ParseMarkdown(content)
	if articleConfig.Date != "" {
		article.Time = ParseDate(articleConfig.Date) //time.timw
		article.Date = article.Time.Unix()           //int
		article.PageDate = time.Unix(article.Date, 0).Format("2006-01-02")
	}
	if articleConfig.Update != "" {
		article.MTime = ParseDate(articleConfig.Update)
		article.Update = article.MTime.Unix()
		article.PageUpdate = time.Unix(article.Update, 0).Format("2006-01-02")

	}
	article.Title = articleConfig.Title
	article.Draft = articleConfig.Draft
	article.Top = articleConfig.Top
	if author, ok := globalConfig.Authors[articleConfig.Author]; ok {
		author.Id = articleConfig.Author
		author.Avatar = RepldceRootFlag(author.Avatar)
		article.Author = author
	}
	//以tag信息给文章分类
	if len(articleConfig.Categories) > 0 {
		article.Category = article.Categories[0]
	} else {
		article.Category = "misc"
	}
	tags := map[string]bool{}
	article.Tags = articleConfig.Tags
	for _, tag := range articleConfig.Tags {
		tags[tag] = true
	}
	for _, category := range articleConfig.Categories {
		if _, ok := tags[category]; !ok {
			article.Tags = append(article.Tags, category)
		}
	}
	if articleConfig.Cover != "" {
		article.Cover = articleConfig.Cover

	}
	//以文章名字作为访问路径
	fileName := strings.TrimSuffix(strings.ToLower(filepath.Base(markdownPath)), ".md")
	link := fileName + ".html"
	//另一种模式，以文章时间作为访问路径
	if article.Type == "post" {
		dataPrefix := article.Time.Format("2006-01-02-")
		if strings.HasPrefix(fileName, dataPrefix) {
			fileName = fileName[len(dataPrefix):]
		}

		if globalConfig.Site.Link != "" {
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
			link = globalConfig.Site.Link
			for key, val := range linkMap {
				link = strings.Replace(link, key, val, -1)
			}
		}
	}
	article.Link = link
	article.GlobalConfig = *globalConfig
	return &article
}

const (
	//基本的时间格式化
	Date_Fromat               = "2006-01-02 15:04:05"
	Date_Format_With_Timezone = "2006-01-02 15:04:05 -0700"
)

//格式化时间
func ParseDate(dataStr string) time.Time {
	date, err := time.Parse(fmt.Sprintf(Date_Format_With_Timezone), dataStr)
	if err != nil {
		date, err = time.ParseInLocation(fmt.Sprintf(Date_Fromat), dataStr, time.Now().Location())
		if err != nil {
			MError("文章时间解析错误:" + err.Error())
		}
	}
	return date
}
