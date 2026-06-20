package lib

import "errors"

var (
	ErrVideoUnavailable    = errors.New("video is unavailable")
	ErrTranscriptNotFound  = errors.New("no transcript found for this video")
	ErrTranslationDisabled = errors.New("translation is disabled for this transcript")
	ErrInvalidLanguage     = errors.New("invalid language code")
)
