package types

const ArticleBizID = "article"
const (
	ArticleStatusPending    int32 = iota // 待审核
	ArticleStatusNotPass                 // 审核不通过
	ArticleStatusVisible                 // 可见
	ArticleStatusUserDelete              // 用户删除
)
