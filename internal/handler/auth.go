package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"pim-api-go/internal/middleware"
	"pim-api-go/internal/service"
)

type AuthHandler struct{ svc *service.AuthService }

func NewAuthHandler(svc *service.AuthService) *AuthHandler { return &AuthHandler{svc} }

func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username dan password wajib diisi"})
		return
	}

	result, err := h.svc.Login(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	isProd := os.Getenv("NODE_ENV") == "production"
	sameSite := http.SameSiteStrictMode
	if isProd {
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    result.Token,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: sameSite,
		MaxAge:   8 * 60 * 60,
		Path:     "/",
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "username": result.Username})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   "",
		MaxAge:  -1,
		Path:    "/",
		Expires: time.Unix(0, 0),
	})
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims := c.MustGet("admin").(*middleware.Claims)
	username, role, err := h.svc.GetMe(claims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data admin"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": username, "role": role})
}
