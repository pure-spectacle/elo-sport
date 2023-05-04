package models

type Follow struct {
	FollowerId  int `json:"followerId" db:"follower_id"`
	FollowingId int `json:"followingId" db:"following_id"`
}

func GetFollow() Follow {
	var follow Follow
	return follow
}

func GetFollows() []Follow {
	var follows []Follow
	return follows
}

func CreateFollow() Follow {
	var follow Follow
	return follow
}

func UpdateFollow() Follow {
	var follow Follow
	return follow
}

func DeleteFollow() Follow {
	var follow Follow
	return follow
}
