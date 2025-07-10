package models

import "time"

// ChatErrorResponse holds any errors that occurred during a chat request
type ChatErrorResponse struct {
	Error Error `json:"error"` // Error details
}


type Error struct { 
	ResponseCode      int    `json:"response_code"` 
	Message           string `json:"message"` 
	Detail            string `json:"detail"` 
	ExternalReference string `json:"ext_ref"`
	Date              string `json:"date"`
}

func NewError() *Error {
	return &Error{
		Date: time.Now().Format("02-01-2006"),
	}
}

