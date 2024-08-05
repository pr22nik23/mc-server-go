package main

import "github.com/gofrs/uuid/v5"

type StatusResponse struct {
	Version            Version     `json:"version"`
	Players            Players     `json:"players"`
	Description        Description `json:"description"`
	Favicon            string      `json:"favion"`
	EnforcesSecureChat bool        `json:"enforcesSecureChat"`
}

type Version struct {
	Name     string `json:"name"`
	Protocol uint16 `json:"protocol"`
}

type Players struct {
	Max    uint32 `json:"max"`
	Online uint32 `json:"online"`
	Sample []Sample `json:"sample"`
}

type Sample struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}

type Description struct {
	Text string `json:"text"`
}
