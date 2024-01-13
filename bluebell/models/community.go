package models

import "time"

type Community struct {
	ID         int64     `json:"community_id" db:"community_id"`
	Name       string    `json:"community_name" db:"community_name"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}

type CommunityDetail struct {
	ID         int64     `json:"community_id" db:"community_id"`
	Name       string    `json:"community_name" db:"community_name"`
	Introduction string `json:"introduction" db:"introduction"`
}

