package modle

import (
	"html/template"
	"time"
)

type SiteConfig struct {
	Site    Site
	Build   Build
	Authors map[string]Author
	App     App
}
type App struct {
	Name    string
	Version string
	Port    int
	Author  string
	Email   string
	Git     string
}
type Site struct {
	Title     string
	Subtitle  string
	Limit     int
	Theme     string
	Language  string
	Logo      string
	URL       string
	Root      string
	Copyright string
	Email     string
	Github    string
	Link      string
}
type Author struct {
	Id     string
	Name   string
	Intro  string
	Avatar string
}
type Build struct {
	Output string
	Source string
	Copy   []string
}
type Theme struct {
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
	SiteConfig
	ArticleConfig
	Time     time.Time
	MTime    time.Time
	Date     int64
	Update   int64
	Author   Author
	Category string
	Tags     []string
	Markdown string
	Preview  template.HTML
	Content  template.HTML
	Link     string
}
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
type RenderArticle struct {
	Article
	Next *Article
	Prev *Article
}
