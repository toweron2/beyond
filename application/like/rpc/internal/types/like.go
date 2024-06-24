package types

type ThumbupMsg struct {
	BizId    string `json:"bizId,omitempty"`
	ObjId    int64  `json:"objId,omitempty"`
	UserId   int64  `json:"userId,omitempty"`
	LikeType int32  `json:"likeType,omitempty"`
}
