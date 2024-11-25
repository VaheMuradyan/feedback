package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateFeedbackAndUpdateRating(c *gin.Context) {
    log.Println("Starting CreateFeedbackAndUpdateRating handler")

    var feedbackResponse feedbackResponse
    if err := c.ShouldBindJSON(&feedbackResponse); err != nil {
        log.Printf("Failed to bind JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    log.Println("JSON binding successful:", feedbackResponse)

    // Convert Rating
    rating, err := strconv.Atoi(feedbackResponse.Rating)
    if err != nil {
        log.Printf("Invalid Rating: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Rating must be a valid integer"})
        return
    }

    // Check if AdminRating exists or create one
    var adminRating AdminRating
    err = DB.Where("admin_id = ?", feedbackResponse.AdminID).First(&adminRating).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            adminIDUint, _ := strconv.ParseUint(feedbackResponse.AdminID, 10, 32)
            adminRating = AdminRating{
                AdminID: uint(adminIDUint),
                Raiting: rating,
            }
            if err := DB.Create(&adminRating).Error; err != nil {
                log.Printf("Failed to create admin rating: %v\n", err)
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin rating"})
                return
            }
            log.Println("AdminRating created successfully")
        } else {
            log.Printf("Failed to fetch admin rating: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admin rating"})
            return
        }
    }

    // Create the feedback and link it to the AdminRating
    feedback := feedbeck{
        UserID:         feedbackResponse.UserID,
        AdminID:        feedbackResponse.AdminID,
        Feedback:       feedbackResponse.Feedback,
        Rating:         rating,
        CreatedAt:      time.Now(),
        AdminRatingId:  adminRating.ID, // Link to the correct AdminRating ID
    }

    log.Println("Saving feedback to the database...")
    if err := DB.Create(&feedback).Error; err != nil {
        log.Printf("Failed to create feedback: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create feedback"})
        return
    }
    log.Println("Feedback saved successfully")

    c.JSON(http.StatusCreated, gin.H{
        "message":     "Feedback created and admin rating updated successfully",
        "feedback":    feedback,
        "adminRating": adminRating,
    })
}


type feedbackResponse struct {
	UserID   string `json:"user_id" binding:"required"`
	AdminID  string `json:"admin_id" binding:"required"`
	Feedback string `json:"feedback" binding:"required"`
	Rating   string `json:"rating" binding:"required"` // Accept rating as a string
}

type feedbeck struct {
	ID            uint      `gorm:"primaryKey"`
	UserID        string    `gorm:"not null"`
	AdminID       string    `gorm:"not null"`
	Feedback      string    `gorm:"not null"`
	Rating        int       `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	AdminRatingId uint
	AdminRating   AdminRating `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type AdminRating struct {
	gorm.Model
	AdminID uint
	Raiting int
}

///Users/user/Documents/Programs/karevor/feedback/feedback.go