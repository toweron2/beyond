package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type Follow struct {
	ID             int64 `gorm:"primary_key"`
	UserID         int64
	FollowedUserID int64
	FollowStatus   int32
	CreateTime     time.Time
	UpdateTime     time.Time
}

func (m *Follow) TableName() string {
	return "follow"
}

type FollowModel struct {
	db *gorm.DB
}

func NewFollowModel(db *gorm.DB) *FollowModel {
	return &FollowModel{db: db}
}

func (m *FollowModel) Insert(ctx context.Context, data *Follow) error {
	return m.db.WithContext(ctx).Create(data).Error
}

func (m *FollowModel) Update(ctx context.Context, data *Follow) error {
	return m.db.WithContext(ctx).Save(data).Error
}

func (m *FollowModel) UpdateFields(ctx context.Context, id int64, values map[string]any) error {
	return m.db.WithContext(ctx).Model(&Follow{}).Where("id = ?", id).Error
}

func (m *FollowModel) FindByUserIDAndFollowedUserID(ctx context.Context, userId, followedUserId int64) (*Follow, error) {
	var result Follow
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND followed_user_id = ?", userId, followedUserId).
		Order("id desc").
		First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, err
}

func (m *FollowModel) FindByUserId(ctx context.Context, userId int64, limit int) ([]*Follow, error) {
	var result []*Follow
	err := m.db.WithContext(ctx).
		Where("user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Limit(limit).
		Find(&result).Error
	return result, err
}

func (m *FollowModel) FindByFollowedUserIds(ctx context.Context, userId int64, followedUserids []int64) ([]*Follow, error) {
	var result []*Follow
	err := m.db.WithContext(ctx).
		Where("user_id = ?", userId).
		Where("followed_user_id in (?)", followedUserids).
		Find(&result).Error
	return result, err
}

func (m *FollowModel) FindByFollowedUserId(ctx context.Context, userId int64, limit int) ([]*Follow, error) {
	var result []*Follow
	err := m.db.WithContext(ctx).
		Where("followed_user_id = ? AND follow_status = ?", userId, 1).
		Order("id desc").
		Find(&result).Error
	return result, err
}
