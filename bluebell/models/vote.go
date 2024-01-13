package models

type VoteForPost struct {
	*ParamVoteForPost
	UserID int64
}
