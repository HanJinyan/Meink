package parse

import (
	"Meink/app/modle"
	"strings"
)

func RepldceRootFlag(content string) string {
	var siteConfig modle.SiteConfig
	return strings.Replace(content, "-/", siteConfig.Site.Root+"/", -1)
}
