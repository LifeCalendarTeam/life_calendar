// All functions use data types from objects.go

func get_user_days(input_user_id int) []briefDay{
	days := make([]briefDay, 0)
	panicIfError(db.Select(&days, "SELECT id, date FROM days WHERE user_id=$1", input_user_id))
	return days
}


func get_days_color_proportions(input_day_id int) []proportionAndColor {
	colorsProportions := make([]proportionAndColor, 0)
	panicIfError(db.Select(&colorsProportions,
		"SELECT CAST(proportion AS FLOAT), (SELECT color FROM types_of_activities_and_emotions WHERE id = type_id) FROM activities_and_emotions WHERE day_id = $1", input_day_id))
	return colorsProportions
}


func get_user_information(input_user_id int) User {
	var user_information User
	panicIfError(db.Select(&user_information, "SELECT * FROM users WHERE id=$1", input_user_id))
	return user_information
}


