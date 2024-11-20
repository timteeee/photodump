package html

import (
	"time"
)

type PhotoFeedParams struct{}

type PhotoCardParams struct {
	Src        string
	Filename   string
	UploadedAt time.Time
}
