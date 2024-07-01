package types

const (
	SortPublishTime int32 = iota + 1
	SortLikeCount
)

const (
	DefaultPageSize = 20
	DefaultLimit    = 200

	DefaultSortLikeCursor = 1 << 30
)
