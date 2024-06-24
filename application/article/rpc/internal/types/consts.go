package types

const (
	SortPublishTime = iota
	SortLikeCount
)

const (
	DefaultPageSize = 20
	DefaultLimit    = 200

	DefaultSortLikeCursor = 1 << 30
)

const (
	ArticleStatusPending    int32 = iota // 待审核
	ArticleStatusNotPass                 // 审核不通过
	ArticleStatusVisible                 // 可见
	ArticleStatusUserDelete              // 用户删除
)
