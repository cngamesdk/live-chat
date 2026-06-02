package v1

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cngamesdk/live-chat/api/global"
	"github.com/cngamesdk/live-chat/api/logger"
	"github.com/cngamesdk/live-chat/api/model"
	"github.com/gin-gonic/gin"
)

type UploadApi struct{}

// Upload godoc
// @Summary Upload image or video
// @Tags H5-Upload
// @Param file formData file true "file to upload"
// @Success 200 {object} model.Response
// @Router /api/v1/upload [post]
func (a *UploadApi) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		model.FailWithMessage("file required", c)
		return
	}

	cfg := global.GVA_CONFIG.Chat
	if file.Size > cfg.MaxUploadSize {
		model.FailWithMessage(fmt.Sprintf("file too large, max %d bytes", cfg.MaxUploadSize), c)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	mimeType := file.Header.Get("Content-Type")
	allowed := false
	for _, t := range cfg.AllowedUploadTypes {
		if t == mimeType {
			allowed = true
			break
		}
	}
	if !allowed {
		model.FailWithMessage("unsupported file type", c)
		return
	}

	fileType := "image"
	if strings.HasPrefix(mimeType, "video/") {
		fileType = "video"
	}

	uploadDir := global.GVA_CONFIG.Upload.Path
	fileName := fmt.Sprintf("%s/%d%s", fileType, time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, fileName)

	dir := filepath.Dir(savePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		model.FailWithMessage("upload failed", c)
		return
	}

	src, err := file.Open()
	if err != nil {
		model.FailWithMessage("upload failed", c)
		return
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		model.FailWithMessage("upload failed", c)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		model.FailWithMessage("upload failed", c)
		return
	}

	userID := c.GetString("user_id")
	upload := model.Upload{
		UserID:   userID,
		FileName: file.Filename,
		FileURL:  "/uploads/" + fileName,
		FileType: fileType,
		FileSize: file.Size,
		MimeType: mimeType,
	}

	if err := global.GVA_DB.Create(&upload).Error; err != nil {
		logger.Error(c.Request.Context(), "failed to save upload record")
	}

	model.OkWithData(gin.H{
		"url":       upload.FileURL,
		"file_type": fileType,
		"file_name": file.Filename,
		"file_size": file.Size,
	}, c)
}
