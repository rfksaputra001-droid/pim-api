package walink

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var nonDigit = regexp.MustCompile(`\D`)

func normalizePhone(phone string) string {
	digits := nonDigit.ReplaceAllString(phone, "")
	if strings.HasPrefix(digits, "0") {
		digits = "62" + digits[1:]
	}
	return digits
}

func formatRupiah(amount int64) string {
	s := fmt.Sprintf("%d", amount)
	n := len(s)
	var b strings.Builder
	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}
	return "Rp" + b.String()
}

func Verified(phone, name string, amount, balance int64) string {
	msg := fmt.Sprintf(
		"Terima kasih %s! Iuran Anda sebesar %s telah dikonfirmasi oleh admin. Saldo kas RT 22/06 kini %s. 🙏",
		name, formatRupiah(amount), formatRupiah(balance),
	)
	return "https://wa.me/" + normalizePhone(phone) + "?text=" + url.QueryEscape(msg)
}

func Rejected(phone, name string, amount int64, reason string) string {
	msg := fmt.Sprintf(
		"Mohon maaf %s, iuran Anda sebesar %s belum dapat diverifikasi.\n\nAlasan: %s\n\nSilakan hubungi admin untuk informasi lebih lanjut.",
		name, formatRupiah(amount), reason,
	)
	return "https://wa.me/" + normalizePhone(phone) + "?text=" + url.QueryEscape(msg)
}
