package modle

type MSort []interface{}

// 排序
func (v MSort) Len() int      { return len(v) }
func (v MSort) Swap(i, j int) { v[i], v[j] = v[j], v[i] } //实现交换
func (v MSort) Less(i, j int) bool { //
	switch v[i].(type) {
	case ArticleInfo:
		return v[i].(ArticleInfo).DetailDate > v[j].(ArticleInfo).DetailDate
	//文章时间，是否置顶作排序
	case Article:
		article1 := v[i].(Article)
		article2 := v[j].(Article)
		if article1.Top && !article2.Top {
			return true
		} else if !article1.Top && article2.Top {
			return false
		} else {
			return article1.Date > article2.Date
		}
	//Archive页面的时间排序
	case Archive:
		return v[i].(Archive).Year > v[j].(Archive).Year
	//Tge 页面tge出现次数排序
	case Tag:
		if v[i].(Tag).Count == v[j].(Tag).Count {
			return v[i].(Tag).Name > v[j].(Tag).Name
		}
		return v[i].(Tag).Count > v[j].(Tag).Count
	}
	return false
}
