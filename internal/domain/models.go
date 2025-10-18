package domain

import "time"

type File struct {
	ID        string    `json:"id"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileList struct {
	Files []File `json:"files"`
	Total int    `json:"total"`
}
