package helper

import (
	"strings"

	"github.com/joeig/go-powerdns/v3"
)

// ToPowerDNS memformat content sesuai aturan PowerDNS
func ToPowerDNS(recordType powerdns.RRType, content string) string {
	clean := strings.TrimSpace(content)

	switch recordType {
	case powerdns.RRTypeA, powerdns.RRTypeAAAA:
		// IP address tidak berubah
		return clean

	case powerdns.RRTypeNS, powerdns.RRTypeCNAME, powerdns.RRTypeMX, powerdns.RRTypeSRV, powerdns.RRTypePTR:
		// domain harus diakhiri titik

		return EnsureDot(clean)

	case powerdns.RRTypeTXT:
		// harus dibungkus tanda kutip
		if !strings.HasPrefix(clean, "\"") {
			clean = `"` + clean
		}
		if !strings.HasSuffix(clean, "\"") {
			clean = clean + `"`
		}
		return clean
	}

	return clean
}

// FromPowerDNS mengembalikan content ke bentuk normal (untuk UI/DB)
func FromPowerDNS(recordType powerdns.RRType, content string) string {
	clean := strings.TrimSpace(content)

	switch recordType {
	case powerdns.RRTypeNS, powerdns.RRTypeCNAME, powerdns.RRTypeMX, powerdns.RRTypeSRV, powerdns.RRTypePTR:
		// hapus titik di akhir jika ada
		return strings.TrimSuffix(clean, ".")

	case powerdns.RRTypeTXT:
		// hapus kutip di awal & akhir jika ada
		clean = strings.TrimPrefix(clean, `"`)
		clean = strings.TrimSuffix(clean, `"`)
		return clean
	}

	// default: biarkan apa adanya
	return clean
}

func EnsureDot(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if !strings.HasSuffix(s, ".") {
		return s + "."
	}
	return s
}

// EqualContent membandingkan dua content setelah normalisasi sesuai tipe
func EqualContent(a, b string, t powerdns.RRType) bool {
	return ToPowerDNS(t, a) == ToPowerDNS(t, b)
}
