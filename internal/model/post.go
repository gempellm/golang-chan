package model

import "google.golang.org/protobuf/types/known/timestamppb"

type Post struct {
	ID            uint64
	ThreadID      uint64
	Title         string
	Message       string
	Created       *timestamppb.Timestamp
	AttachedMedia []Media
}

type Media struct {
	Data []byte
	Name string
}
