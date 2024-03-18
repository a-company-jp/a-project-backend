package handler

import (
	"a-project-backend/gen/gModel"
	"a-project-backend/gen/gQuery"
	"a-project-backend/pkg/config"
	"a-project-backend/pkg/gcs"
	"a-project-backend/svc/pkg/domain/model/exception"
	"a-project-backend/svc/pkg/domain/model/pkg_time"
	"a-project-backend/svc/pkg/middleware"
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/nickalie/go-webpbin"
	"github.com/openhacku-team-a/a-project-frontend/proto/golang/pb_out"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
	q  *gQuery.Query
	g  *gcs.GCS
}

const (
	GCSUserIconFolder       = "user_icon"
	GCSMilestoneImageFolder = "milestone_img"
)

func NewUser(db *gorm.DB, g *gcs.GCS) User {
	return User{
		db: db,
		q:  gQuery.Use(db),
		g:  g,
	}
}

// GetMe Meユーザーの取得
func (h User) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		uAny, exists := c.Get(middleware.AuthorizedUserField)
		if !exists {
			c.AbortWithStatusJSON(500, gin.H{"error": "userId not set"})
		}
		u, ok := uAny.(gModel.User)
		if !ok {
			c.AbortWithStatusJSON(500, gin.H{"error": "user not castable"})
		}

		tags := make([]*pb_out.Tag, len(u.Tags))
		for i, t := range u.Tags {
			tags[i] = &pb_out.Tag{
				TagId:   t.TagID,
				TagName: t.TagName,
			}
		}

		ms, err := h.q.Milestone.Where(h.q.Milestone.UserID.Eq(u.UserID)).Find()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		conf := config.Get()

		var milestones []*pb_out.Milestone
		for _, m := range ms {
			imgUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSMilestoneImageFolder, m.ImageHash)
			milestones = append(milestones, &pb_out.Milestone{
				MilestoneId: m.MilestoneID,
				UserId:      m.UserID,
				Title:       m.Title,
				Content:     m.Content,
				ImageUrl:    &imgUrl,
				BeginDate:   pkg_time.DateOnly(m.BeginDate).String(),
				FinishDate:  pkg_time.DateOnly(m.FinishDate).String(),
			})
		}

		iconUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSUserIconFolder, u.IconImageHash)

		resp := pb_out.UserInfoResponse{
			UserData: &pb_out.UserData{
				UserId:        u.UserID,
				Username:      u.Username,
				Firstname:     u.Firstname,
				Lastname:      u.Lastname,
				FirstnameKana: u.FirstnameKana,
				LastnameKana:  u.LastnameKana,
				StatusMessage: u.StatusMessage,
				Tag:           tags,
				IconUrl:       &iconUrl,
			},
			Milestones: milestones,
		}
		c.ProtoBuf(200, &resp)
	}
}

// GetUserInfo ユーザー情報の単体取得
func (h User) GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		u, err := h.q.User.Where(h.q.User.UserID.Eq(userID)).First()
		if err != nil {
			if errors.Is(err, exception.ErrNotFound) {
				c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
			} else {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}
		}

		ms, err := h.q.Milestone.Where(h.q.Milestone.UserID.Eq(userID)).Find()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		tags := make([]*pb_out.Tag, len(u.Tags))
		for i, t := range u.Tags {
			tags[i] = &pb_out.Tag{
				TagId:   t.TagID,
				TagName: t.TagName,
			}
		}

		conf := config.Get()

		var milestones []*pb_out.Milestone
		for _, m := range ms {
			imgUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSMilestoneImageFolder, m.ImageHash)
			milestones = append(milestones, &pb_out.Milestone{
				MilestoneId: m.MilestoneID,
				UserId:      m.UserID,
				Title:       m.Title,
				Content:     m.Content,
				ImageUrl:    &imgUrl,
				BeginDate:   pkg_time.DateOnly(m.BeginDate).String(),
				FinishDate:  pkg_time.DateOnly(m.FinishDate).String(),
			})
		}

		iconUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSUserIconFolder, u.IconImageHash)

		resp := pb_out.UserInfoResponse{
			UserData: &pb_out.UserData{
				UserId:        u.UserID,
				Username:      u.Username,
				Firstname:     u.Firstname,
				Lastname:      u.Lastname,
				FirstnameKana: u.FirstnameKana,
				LastnameKana:  u.LastnameKana,
				StatusMessage: u.StatusMessage,
				Tag:           tags,
				IconUrl:       &iconUrl,
			},
			Milestones: milestones,
		}
		c.ProtoBuf(200, &resp)
	}
}

// GetUserInfos ユーザー情報の全件取得
func (h User) GetUserInfos() gin.HandlerFunc {
	return func(c *gin.Context) {
		// conf
		conf := config.Get()

		// ユーザー取得
		users, err := h.q.User.Find()
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		// レスポンス作成
		var resp pb_out.UserInfosResponse
		for _, user := range users {
			// Tag
			tags := make([]*pb_out.Tag, len(user.Tags))
			for i, t := range user.Tags {
				tags[i] = &pb_out.Tag{
					TagId:   t.TagID,
					TagName: t.TagName,
				}
			}

			// milestone取得
			ms, err := h.q.Milestone.Where(h.q.Milestone.UserID.Eq(user.UserID)).Find()
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}

			// response方に変換
			var milestones []*pb_out.Milestone
			for _, m := range ms {
				imgUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSMilestoneImageFolder, m.ImageHash)
				milestones = append(milestones, &pb_out.Milestone{
					MilestoneId: m.MilestoneID,
					UserId:      m.UserID,
					Title:       m.Title,
					Content:     m.Content,
					ImageUrl:    &imgUrl,
					BeginDate:   pkg_time.DateOnly(m.BeginDate).String(),
					FinishDate:  pkg_time.DateOnly(m.FinishDate).String(),
				})
			}

			iconUrl := fmt.Sprintf("https://storage.googleapis.com/%s/%s/%s", conf.Application.GCS.BucketName, GCSUserIconFolder, user.IconImageHash)

			resp.UserInfoResponses = append(resp.UserInfoResponses, &pb_out.UserInfoResponse{
				UserData: &pb_out.UserData{
					UserId:        user.UserID,
					Username:      user.Username,
					Firstname:     user.Firstname,
					Lastname:      user.Lastname,
					FirstnameKana: user.FirstnameKana,
					LastnameKana:  user.LastnameKana,
					StatusMessage: user.StatusMessage,
					Tag:           tags,
					IconUrl:       &iconUrl,
				},
				Milestones: milestones,
			})
		}
		c.ProtoBuf(200, &resp)
	}
}

// UpdateUserInfo ユーザー情報の更新
func (h User) UpdateUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if userId == "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "id should not be null"})
		}

		data, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		var userInfoUpdateRequest pb_out.UserInfoUpdateRequest
		err = proto.Unmarshal(data, &userInfoUpdateRequest)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		_, err = h.q.User.WithContext(c).Where(h.q.User.UserID.Eq(userId)).Updates(
			map[string]interface{}{
				"username":       userInfoUpdateRequest.UserData.Username,
				"firstname":      userInfoUpdateRequest.UserData.Firstname,
				"lastname":       userInfoUpdateRequest.UserData.Lastname,
				"firstname_kana": userInfoUpdateRequest.UserData.FirstnameKana,
				"lastname_kana":  userInfoUpdateRequest.UserData.LastnameKana,
				"status_message": userInfoUpdateRequest.UserData.StatusMessage,
			},
		)

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		c.AbortWithStatus(200)
	}
}

// UpdateUserIcon ユーザーアイコンの更新
func (h User) UpdateUserIcon() gin.HandlerFunc {
	return func(c *gin.Context) {
		// decode
		img, _, err := image.Decode(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		// convert to webp
		webpData, err := convertIMG2WEBP(img)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		// get userId
		userId := c.Param("user_id")
		if userId == "" {
			c.AbortWithStatusJSON(500, gin.H{"error": "id should not be null"})
		}

		// create object name
		now := time.Now()
		formatted := now.Format("20060102-150405")
		objectName := fmt.Sprintf("%s/%s-%s", GCSUserIconFolder, userId, formatted)

		// save to storage
		err = h.g.Upload(c, objectName, webpData)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		// save to user table
		_, err = h.q.User.WithContext(c).Where(h.q.User.UserID.Eq(userId)).Updates(
			map[string]interface{}{
				"image_hash": objectName,
			},
		)

		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}

		c.AbortWithStatus(200)
	}
}

func convertIMG2WEBP(img image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := webpbin.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
