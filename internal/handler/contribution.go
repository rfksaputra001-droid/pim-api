package handler

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"pim-api-go/internal/middleware"
	"pim-api-go/internal/service"
	"pim-api-go/pkg/cloudinary"
)

type ContributionHandler struct{ svc *service.ContributionService }

func NewContributionHandler(svc *service.ContributionService) *ContributionHandler {
	return &ContributionHandler{svc}
}

func (h *ContributionHandler) Dashboard(c *gin.Context) {
	result, err := h.svc.Dashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data dashboard"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ContributionHandler) Leaderboard(c *gin.Context) {
	entries, err := h.svc.Leaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat leaderboard"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": entries})
}

func (h *ContributionHandler) ListPublic(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	year, _ := strconv.Atoi(c.Query("year"))
	month, _ := strconv.Atoi(c.Query("month"))

	result, err := h.svc.ListPublic(year, month, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data iuran"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ContributionHandler) Submit(c *gin.Context) {
	file, header, err := c.Request.FormFile("buktiTransfer")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bukti transfer wajib diupload"})
		return
	}
	defer file.Close()

	rawMIME := header.Header.Get("Content-Type")
	mimeType := strings.SplitN(rawMIME, ";", 2)[0]
	mimeType = strings.TrimSpace(mimeType)
	if !cloudinary.AllowedMIMEs[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Gunakan JPG, PNG, atau PDF."})
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membaca file"})
		return
	}

	nominal, _ := strconv.ParseInt(c.PostForm("nominal"), 10, 64)

	input := service.SubmitInput{
		Nama:        strings.TrimSpace(c.PostForm("nama")),
		NoTelepon:   strings.TrimSpace(c.PostForm("noTelepon")),
		Nominal:     nominal,
		Catatan:     strings.TrimSpace(c.PostForm("catatan")),
		IsAnonymous: c.PostForm("isAnonymous") == "true",
		FileData:    data,
		MimeType:    mimeType,
	}

	uploadFn := func(d []byte, folder string) (string, error) {
		return cloudinary.Upload(context.Background(), d, folder)
	}

	if err := h.svc.Submit(input, uploadFn); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Terima kasih! Pembayaran Anda sedang diverifikasi oleh admin. 🙏"})
}

func (h *ContributionHandler) ListAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	result, err := h.svc.ListAdmin(
		c.Query("status"),
		c.DefaultQuery("sortBy", "createdAt"),
		c.DefaultQuery("sortDir", "desc"),
		page,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data iuran"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ContributionHandler) Verify(c *gin.Context) {
	claims := c.MustGet("admin").(*middleware.Claims)
	waLink, err := h.svc.Verify(c.Param("id"), claims.Username)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Iuran tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Iuran berhasil diverifikasi", "waLink": waLink})
}

func (h *ContributionHandler) Reject(c *gin.Context) {
	var body struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&body)
	if strings.TrimSpace(body.Reason) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Alasan penolakan wajib diisi"})
		return
	}

	waLink, err := h.svc.Reject(c.Param("id"), strings.TrimSpace(body.Reason))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "Iuran tidak ditemukan" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Iuran ditolak", "waLink": waLink})
}

func (h *ContributionHandler) ExportCSV(c *gin.Context) {
	rows, err := h.svc.ExportCSV()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengekspor data"})
		return
	}

	escCSV := func(v string) string {
		if len(v) > 0 && strings.ContainsRune("=+-@\t\r", rune(v[0])) {
			v = "'" + v
		}
		return `"` + strings.ReplaceAll(v, `"`, `""`) + `"`
	}

	var buf strings.Builder
	for _, row := range rows {
		cols := make([]string, len(row))
		for i, v := range row {
			cols[i] = escCSV(v)
		}
		buf.WriteString(strings.Join(cols, ",") + "\n")
	}

	filename := "laporan-iuran-" + time.Now().Format("2006-01-02") + ".csv"
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="`+filename+`"`)
	c.String(http.StatusOK, "\xef\xbb\xbf"+buf.String())
}
