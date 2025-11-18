package validator

import (
	"fmt"
	"picstagsbot/pkg/constants"
	"regexp"
	"strings"
	"unicode/utf8"
)

var tagRegex = regexp.MustCompile(constants.TagPattern)

func ValidateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}

	if utf8.RuneCountInString(tag) < constants.MinTagLength {
		return fmt.Errorf("tag is too short: minimum %d characters", constants.MinTagLength)
	}

	if utf8.RuneCountInString(tag) > constants.MaxTagLen {
		return fmt.Errorf("tag is too long: maximum %d characters", constants.MaxTagLen)
	}

	if !tagRegex.MatchString(tag) {
		return fmt.Errorf("tag contains invalid characters: only letters, numbers, underscore and hyphen are allowed")
	}

	return nil
}

func ValidateTags(tags []string) error {
	if len(tags) > constants.MaxTagsPerPhoto {
		return fmt.Errorf("too many tags: maximum %d tags per photo", constants.MaxTagsPerPhoto)
	}

	for i, tag := range tags {
		if err := ValidateTag(tag); err != nil {
			return fmt.Errorf("invalid tag at position %d: %w", i+1, err)
		}
	}

	return nil
}

func ValidateDescription(description string) error {
	if description == "" {
		return fmt.Errorf("description cannot be empty")
	}

	if utf8.RuneCountInString(description) > constants.MaxDescriptionLen {
		return fmt.Errorf("description is too long: maximum %d characters", constants.MaxDescriptionLen)
	}

	return nil
}

func ValidateAndParseTags(description string) ([]string, error) {
	if err := ValidateDescription(description); err != nil {
		return nil, err
	}

	tags := strings.Fields(description)
	if err := ValidateTags(tags); err != nil {
		return nil, err
	}

	return tags, nil
}

func ValidateFileSize(size int64) error {
	if size <= 0 {
		return fmt.Errorf("file size must be positive")
	}

	if size > constants.MaxFileSize {
		return fmt.Errorf("file is too large: maximum %d bytes", constants.MaxFileSize)
	}

	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	length := utf8.RuneCountInString(username)
	if length < constants.MinUsernameLength {
		return fmt.Errorf("username is too short: minimum %d characters", constants.MinUsernameLength)
	}

	if length > constants.MaxUsernameLength {
		return fmt.Errorf("username is too long: maximum %d characters", constants.MaxUsernameLength)
	}

	return nil
}

func SanitizeString(s string) string {
	s = strings.Map(func(r rune) rune {
		if r == 0 || (r < 32 && r != '\n' && r != '\r' && r != '\t') {
			return -1
		}
		return r
	}, s)
	return strings.TrimSpace(s)
}
