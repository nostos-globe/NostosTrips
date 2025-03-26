package controller

import (
    "main/internal/models"
    "main/internal/service"
    "net/http"
    "strconv"
	"time"
    "github.com/gin-gonic/gin"
)

type MediaController struct {
    MediaService *service.MediaService
    AuthClient   *service.AuthClient
}

func (c *MediaController) UploadMedia(ctx *gin.Context) {
    tripID, err := strconv.ParseInt(ctx.Param("trip_id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid trip ID"})
        return
    }

    tokenCookie, err := ctx.Cookie("auth_token")
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
        return
    }

    userID, err := c.AuthClient.GetUserID(tokenCookie)
    
    if err != nil || userID == 0 {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
        return
    }

    file, header, err := ctx.Request.FormFile("media")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "no file provided"})
        return
    }

    visibility := models.VisibilityEnum(ctx.PostForm("visibility"))
    if visibility == "" {
        visibility = models.Public
    }

    objectName, err := c.MediaService.UploadMedia(int64(userID), file, header, visibility)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload media"})
        return
    }

    // Create media record with trip association
    media := models.Media{
        TripID:      tripID,
        UserID:      int64(userID),
        LocationID:  3,  // Will be handled as NULL by GORM
        Type:        "image",
        FilePath:    objectName,
        Visibility:  visibility,
        UploadDate:  time.Now(),
        CaptureDate: time.Now(),
    }

    err = c.MediaService.SaveMedia(&media)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save media metadata"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "media uploaded successfully",
        "path": objectName,
    })
}

func (c *MediaController) GetMediaURL(ctx *gin.Context) {
    mediaID, err := strconv.ParseInt(ctx.Param("media_id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
        return
    }

    tokenCookie, err := ctx.Cookie("auth_token")
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
        return
    }

    userID, err := c.AuthClient.GetUserID(tokenCookie)
    if err != nil || userID == 0 {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
        return
    }

    url, err := c.MediaService.GetMediaURL(mediaID, int64(userID))
    if err != nil {
        ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"url": url})
}

func (c *MediaController) DeleteMedia(ctx *gin.Context) {
    tripID := ctx.Param("trip_id")
    mediaID, err := strconv.ParseInt(ctx.Param("media_id"), 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid media ID"})
        return
    }

    tokenCookie, err := ctx.Cookie("auth_token")
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no token found"})
        return
    }

    userID, err := c.AuthClient.GetUserID(tokenCookie)
    if err != nil || userID == 0 {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "failed to authenticate user"})
        return
    }

    // Delete specific media instead of all trip media
    err = c.MediaService.DeleteMedia(mediaID, tripID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete media"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "media deleted successfully"})
}