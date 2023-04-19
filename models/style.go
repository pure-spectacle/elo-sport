package models

type Style struct {
	StyleId   int    `json:"styleId" db:"style_id"`
	Name      string `json:"name" db:"style_name"`
}

func GetStyle() Style {
	var style Style
	return style
}

func GetStyles() []Style {
	var styles []Style
	return styles
}

func CreateStyle() Style {
	var style Style
	return style
}

func UpdateStyle() Style {
	var style Style
	return style
}

func DeleteStyle() Style {
	var style Style
	return style
}
