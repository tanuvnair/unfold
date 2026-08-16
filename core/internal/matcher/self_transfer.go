package matcher

import "strings"

// LooksLikeSelfTransfer reports bank internet-banking / own-account moves
// that are not merchant autopay (e.g. Kotak "IB:..." descriptions).
func LooksLikeSelfTransfer(description string) bool {
	s := strings.ToUpper(strings.TrimSpace(description))
	if s == "" {
		return false
	}
	return strings.HasPrefix(s, "IB:") || strings.HasPrefix(s, "IB ")
}
