package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	novugo "github.com/novuhq/novu-go"
	"github.com/novuhq/novu-go/models/components"
	"github.com/pick-cee/events-api/internal/config"
	"github.com/pick-cee/events-api/internal/database"
	"github.com/pick-cee/events-api/internal/middleware"
	"github.com/pick-cee/events-api/internal/models"
)


type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
	}
}

// Request/Response DTOs
type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// sign up
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user is an existing user
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
    return
	}

	// create user 
	user := models.User{
		Name: req.Name,
		Email: req.Email,
		Password: req.Password,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
    return
	}

	// generate a token
	token, err := h.generateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		User: UserResponse{
			ID: user.ID,
			Email: user.Email,
			Name: user.Name,
		},
		Token: token,
	}

	sendWelcomeEmail(response.User.Email)


	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check if user exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// check password
	if !existingUser.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// generate token
	token, err := h.generateToken(existingUser.ID, existingUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	response := AuthResponse{
		User: UserResponse{
			ID: existingUser.ID,
			Email: existingUser.Email,
			Name: existingUser.Name,
		},
		Token: token,
	}

	c.JSON(http.StatusOK, response)

}


func (h *AuthHandler) generateToken(userId uint, email string) (string, error) {
	claims := middleware.Claims{
		UserID: userId,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}

func sendWelcomeEmail(email string) {
	// Run in a goroutine to avoid blocking the response
	go func() {
		log.Println("Sending welcome email to:", email)
		ctx := context.Background()
		secretKey := config.GetEnv("NOVU_SECRET_KEY", "")

		s := novugo.New(novugo.WithSecurity(secretKey))

		_, err := s.Trigger(ctx, components.TriggerEventRequestDto{
			WorkflowID: "manimo-welcome-email",
			Payload: map[string]any{
				"style": ":root{color-scheme:light dark;supported-color-schemes:light dark;}body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Oxygen,Ubuntu,Cantarell,sans-serif;line-height:1.6;color:#333;max-width:600px;margin:0 auto;padding:20px;background-color:#f9f9f9;}.email-container{background-color:#fff;border-radius:12px;box-shadow:0 2px 10px rgba(0,0,0,0.05);}.email-body{padding:30px 40px 40px;}.logo-container{text-align:center;background-color:#2592FF;height:60px;border-top-left-radius:12px;border-top-right-radius:12px;padding:10px 0;display:flex;align-items:center;justify-content:center;margin-bottom:30px;}.logo{width:130px;margin:0 auto;object-fit:contain;}.logo-light{display:block;margin:0 auto;width:100px;}.logo-dark{display:none;margin:0 auto;width:100px;}h1{text-align:center;font-size:24px;margin-bottom:20px;font-weight:600;}h2{font-size:20px;margin-top:0;margin-bottom:20px;font-weight:600;}.welcome-message{text-align:center;margin-bottom:30px;font-size:16px;color:#555;}.button-container{text-align:center;margin:30px 0;}.button{background-color:#2592FF;color:white;padding:12px 24px;border-radius:50px;text-decoration:none;font-weight:500;display:inline-block;font-size:16px;}.features{margin:40px 0;}.feature-table{width:100%;border-collapse:collapse;margin-bottom:20px;}.feature-icon-cell{width:60px;vertical-align:top;}.feature-icon{background-color:#E6F3FF;width:40px;height:40px;border-radius:50%;text-align:center;font-size:18px;color:#2592FF;display:inline-block;line-height:40px;}.feature-text-cell{vertical-align:top;}.feature-title{font-weight:600;margin-bottom:2px;}.feature-description{color:#777;font-size:14px;}.divider{border:none;height:1px;background-color:#eee;margin:30px 0;}.footer{text-align:center;color:#777;font-size:14px;}.social-links{text-align:center;margin:20px 0;}.social-link{width:30px;height:30px;background-color:#f2f2f2;border-radius:50%;display:inline-block;text-align:center;line-height:30px;font-size:14px;color:#555;text-decoration:none;margin:0 10px;}.account-info{background-color:#EBF6FF;border-radius:8px;padding:16px;margin:20px 0;font-size:14px;color:#555;}.invite-message{text-align:center;margin-bottom:30px;font-size:16px;color:#555;}.team-info{background-color:#EBF6FF;border-radius:8px;padding:20px;margin:25px 0;text-align:center;}.team-name{font-size:18px;font-weight:600;color:#2592FF;margin-bottom:5px;}.invited-by{font-size:14px;color:#666;}.expiry-notice{text-align:center;font-size:14px;color:#777;margin-top:15px;}.role-info{font-size:15px;margin-top:10px;padding:12px;background-color:#F5F9FC;border-radius:6px;}.role-label{font-weight:600;margin-bottom:3px;}.success-message{text-align:center;margin-bottom:30px;font-size:16px;color:#555;}.success-icon{text-align:center;margin:20px 0;font-size:48px;}.confirmation-message{text-align:center;margin-bottom:30px;font-size:16px;color:#555;}.subscription-details{background-color:#EBF6FF;border-radius:8px;padding:20px;margin:25px 0;}.details-header{font-size:18px;font-weight:600;color:#2592FF;margin-bottom:15px;}.details-item{margin-bottom:10px;}.details-label{font-weight:500;display:inline-block;width:150px;}.what-to-expect{margin:30px 0;}.expect-table{width:100%;border-collapse:collapse;margin-bottom:20px;}.expect-icon-cell{width:60px;vertical-align:top;}.expect-icon{background-color:#E6F3FF;width:40px;height:40px;border-radius:50%;text-align:center;font-size:18px;color:#2592FF;display:inline-block;line-height:40px;}.expect-text-cell{vertical-align:top;}.expect-title{font-weight:600;margin-bottom:2px;}.expect-description{color:#777;font-size:14px;}.preference-link{text-align:center;margin:25px 0;font-size:14px;}.unsubscribe{font-size:12px;color:#999;margin-top:15px;}.info-box{background-color:#EBF6FF;border-radius:8px;padding:20px;margin:25px 0;text-align:center;}.info-box p{margin:10px 0;}.info-header{color:#2592FF;font-weight:600;margin-bottom:8px;}.info-header-icon{display:inline-block;margin-right:8px;}.info-content{color:#555;font-size:14px;line-height:1.5;}.feedback-section{margin:30px 0;text-align:center;}.feedback-title{font-weight:600;margin-bottom:15px;}.feedback-options{text-align:center;margin:20px 0;}.feedback-option{background-color:#f5f5f5;border-radius:50px;padding:8px 16px;font-size:14px;color:#555;text-decoration:none;border:1px solid #eee;display:inline-block;margin:5px;}.resubscribe{text-align:center;margin:25px 0;color:#777;font-size:14px;}.next-steps{margin:30px 0;}.step-table{width:100%;border-collapse:collapse;margin-bottom:20px;}.step-number-cell{width:40px;vertical-align:top;padding-top:3px;}.step-number{background-color:#2592FF;color:white;width:24px;height:24px;border-radius:50%;text-align:center;font-size:14px;font-weight:600;line-height:24px;display:inline-block;}.step-text-cell{vertical-align:top;}.step-title{font-weight:600;margin-bottom:2px;}.step-description{color:#777;font-size:14px;}.verification-code{text-align:center; margin:30px auto;width:100%;}.verification-code-table{width:auto;margin:0 auto;border-collapse:separate;border-spacing:8px 0;}.code-digit{width:48px; height:48px; background-color:#f7f7f7; border-radius:8px; font-size:24px; font-weight:600;color:#333;text-align:center;vertical-align:middle;border:1px solid #e0e0e0;padding:0;}.warning{text-align:center;color:#777;margin-bottom:30px;font-size:15px;font-weight:500;}.ip-info{display:inline-block;background-color:#f2f2f2;border-radius:50px;padding:8px 18px;font-size:14px;margin-top:12px;font-weight:500;}.message{text-align:center;font-size:16px;color:#555;margin-bottom:30px;}.member-card{background-color:#EBF6FF;border-radius:8px;padding:20px;margin:25px 0;}.member-table{width:100%;border-collapse:collapse;}.member-avatar-cell{width:80px;vertical-align:top;}.member-avatar{width:60px;height:60px;border-radius:50%;background-color:#D6EDFF;text-align:center;font-size:24px;color:#2592FF;line-height:60px;}.member-info-cell{vertical-align:top;}.member-name{font-size:18px;font-weight:600;}.member-role{color:#2592FF;font-weight:500;margin:5px 0;}.member-email{color:#777;font-size:14px;}.team-stats{background-color:#F7F9FB;border-radius:8px;padding:16px;margin:25px 0;text-align:center;}.stat-value{font-size:24px;font-weight:600;color:#2592FF;margin-bottom:5px;}.stat-label{font-size:14px;color:#777;}.notification-message{text-align:center;margin-bottom:30px;font-size:16px;color:#555;}@media (prefers-color-scheme:dark){body{background-color:#121212;color:#f0f0f0;}.email-container{background-color:#1e1e1e;box-shadow:0 2px 10px rgba(0,0,0,0.3);}.logo-light{display:none!important;}.logo-dark{display:block!important;}h1, h2{color:#ffffff;}.welcome-message, .confirmation-message, .invite-message, .success-message, .message, .notification-message{color:#d0d0d0;}.button{background-color:#4dabff;}.feature-icon, .expect-icon{background-color:#1e2730;color:#4dabff;}.feature-title, .expect-title, .step-title, .member-name{color:#ffffff;}.feature-description, .expect-description, .step-description, .member-email{color:#bcbcbc;}.divider{background-color:#333;}.footer{color:#bcbcbc;}.social-link{background-color:#2a2a2a;color:#d0d0d0;border:1px solid #3d3d3d;}.account-info, .team-info, .subscription-details, .info-box, .member-card{background-color:#1e2730;border:1px solid #2a3744;color:#d0d0d0;}.team-name, .details-header, .info-header, .member-role{color:#4dabff;}.invited-by, .expiry-notice, .preference-link, .resubscribe, .stat-label{color:#bcbcbc;}.role-info{background-color:#252e38;color:#d0d0d0;border:1px solid #2a3744;}.code-digit{background-color:#2a2a2a;color:#ffffff;border:1px solid #3d3d3d;}.warning{color:#bcbcbc;}.ip-info{background-color:#2a2a2a;color:#d0d0d0;border:1px solid #3d3d3d;}.step-number{background-color:#4dabff;}.feedback-option{background-color:#2a2a2a;color:#d0d0d0;border:1px solid #3d3d3d;}.unsubscribe{color:#888;}.team-stats{background-color:#252e38;border:1px solid #2a3744;}.stat-value{color:#4dabff;}.member-avatar{background-color:#2a3744;color:#4dabff;}a{color:#4dabff;}}@media only screen and (max-width:768px){body{padding:15px;}.email-container{padding:30px;}}@media only screen and (max-width:480px){body{padding:10px;}.email-container{padding:25px 15px;border-radius:10px;}h1{font-size:22px;}h2{font-size:18px;}.button{padding:10px 20px;font-size:15px;}.feature-icon, .expect-icon{width:36px;height:36px;line-height:36px;font-size:16px;}.feature-icon-cell, .expect-icon-cell{width:50px;}.confirmation-message, .invite-message, .success-message, .message, .warning, .notification-message{font-size:15px;}.team-info, .subscription-details, .info-box, .member-card{padding:15px;}.details-label{width:120px;}.code-digit{width:36px;height:36px;font-size:18px;}.verification-code-table{border-spacing:6px 0;}.feedback-option{padding:6px 14px;font-size:13px;margin:3px;}.role-info{padding:10px;}.member-avatar-cell{width:70px;}.member-avatar{width:50px;height:50px;font-size:20px;line-height:50px;}.member-name{font-size:16px;}.team-stats{padding:12px;}.stat-value{font-size:20px;}}@media only screen and (max-width:320px){.email-container{padding:20px 12px;}.feature-icon, .expect-icon{width:32px;height:32px;line-height:32px;font-size:14px;}.feature-icon-cell, .expect-icon-cell{width:45px;}.details-label{width:100%;display:block;margin-bottom:2px;}.code-digit{width:30px;height:30px;font-size:16px;}.verification-code-table{border-spacing:4px 0;}.feedback-option{display:block;margin:8px auto;width:80%;}.step-number{width:22px;height:22px;line-height:22px;font-size:12px;}.step-number-cell{width:35px;}.member-avatar-cell{width:60px;}.member-avatar{width:45px;height:45px;font-size:18px;line-height:45px;}}",
				"image": "https://manimo.s3.us-east-1.amazonaws.com/new-logo-full+2.png",
			},
			To: components.CreateToSubscriberPayloadDto(components.SubscriberPayloadDto{
				Email: &email,
				SubscriberID: email,
			}),
		},
		nil)

		if err != nil {
			log.Println("Error sending welcome email:", err)
		}
	}()
}