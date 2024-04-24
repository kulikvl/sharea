package utils

import (
	"fmt"
	"time"
)

func FormatDate(dateTime time.Time) string {
	now := time.Now()
	hoursDiff := now.Sub(dateTime).Hours()

	if hoursDiff < 24 {
		return "Today"
	} else if hoursDiff < 48 {
		return "Yesterday"
	} else {
		return dateTime.Format("Jan 02, 2006")
	}
}

func FormatSize(size int64) string {
	const (
		B   = 1
		KiB = 1024 * B
		MiB = 1024 * KiB
		GiB = 1024 * MiB
		TiB = 1024 * GiB
	)

	switch {
	case size < KiB:
		return fmt.Sprintf("%dB", size)
	case size < MiB:
		return fmt.Sprintf("%.1fKiB", float64(size)/KiB)
	case size < GiB:
		return fmt.Sprintf("%.1fMiB", float64(size)/MiB)
	case size < TiB:
		return fmt.Sprintf("%.1fGiB", float64(size)/GiB)
	default:
		return fmt.Sprintf("%.1fTiB", float64(size)/TiB)
	}
}
