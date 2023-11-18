package models

import "errors"

var ErrNoRecord = errors.New("models: no matching record found")
var ErrPublish = errors.New("drive: already published")
