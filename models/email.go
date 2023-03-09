package models

type Email struct {
	Id       string `json:"id"`
	Title    string `json:"title" `
	Body     string `json:"body" binding:"required"`
	Created  int64  `json:"created"`
	Folder   string `json:"folder" binding:"required"`
	Username string `json:"username" binding:"required"`
}
