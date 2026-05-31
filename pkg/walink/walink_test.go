package walink_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"pim-api-go/pkg/walink"
)

func TestVerifiedLink(t *testing.T) {
	link := walink.Verified("081234567890", "Budi", 50000, 1500000)
	assert.True(t, strings.HasPrefix(link, "https://wa.me/6281234567890"))
	assert.Contains(t, link, "?text=")
	assert.Contains(t, link, "Budi")
}

func TestRejectedLink(t *testing.T) {
	link := walink.Rejected("081234567890", "Siti", 50000, "Bukti tidak jelas")
	assert.True(t, strings.HasPrefix(link, "https://wa.me/6281234567890"))
	assert.Contains(t, link, "Siti")
}

func TestPhoneNormalization(t *testing.T) {
	link1 := walink.Verified("081234567890", "A", 1000, 1000)
	assert.Contains(t, link1, "wa.me/6281234567890")

	link2 := walink.Verified("6281234567890", "A", 1000, 1000)
	assert.Contains(t, link2, "wa.me/6281234567890")
}
