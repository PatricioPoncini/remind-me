package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	RemindCommand = "/r"
	StartCommand  = "/start"
)

func ParseMessage(message string) (text string, notifyTime time.Time, duration time.Duration, err error) {
	re := regexp.MustCompile(`/r\s+'(.*?)'\s+in\s+'(.*?)'`)
	matches := re.FindStringSubmatch(message)

	if len(matches) < 3 {
		return "", time.Time{}, 0, fmt.Errorf("invalid message format. no matches")
	}

	text = strings.TrimSpace(matches[1])
	durationStr := strings.TrimSpace(matches[2])

	if strings.HasSuffix(durationStr, "h") {
		hours, err := strconv.Atoi(strings.TrimSuffix(durationStr, "h"))
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("error parsing hours: %w", err)
		}
		duration = time.Duration(hours) * time.Hour
	} else if strings.HasSuffix(durationStr, "m") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(durationStr, "m"))
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("error parsing minutes: %w", err)
		}
		duration = time.Duration(minutes) * time.Minute
	} else if strings.HasSuffix(durationStr, "s") {
		seconds, err := strconv.Atoi(strings.TrimSuffix(durationStr, "s"))
		if err != nil {
			return "", time.Time{}, 0, fmt.Errorf("error parsing seconds: %w", err)
		}
		duration = time.Duration(seconds) * time.Second
	} else {
		return "", time.Time{}, 0, fmt.Errorf("invalid duration format")
	}

	notifyTime = time.Now().Add(duration)
	return text, notifyTime, duration, nil
}
