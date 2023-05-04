package models

type Follow struct {
	FollowerId int `json:"followerId" db:"follower_id"`
	FollowedId int `json:"followedId" db:"followed_id"`
	CreatedDate string `json:"createdDate" db:"created_dt"`
	UpdatedDate string `json:"updatedDate" db:"updated_dt"`
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
