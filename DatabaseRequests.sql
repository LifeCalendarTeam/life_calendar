CREATE FUNCTION get_user_information 
	(IN user_id int, OUT first_name VARCHAR, OUT second_name VARCHAR, OUT city VARCHAR, 
     OUT birthday DATE, OUT avatar VARCHAR, OUT password_hash VARCHAR) AS
	$$ SELECT users.first_name, users.second_name, users.city, users.birthday, users.avatar, users.password_hash 
    FROM users WHERE users.id = user_id; $$ 
    LANGUAGE SQL;
    
CREATE FUNCTION get_day_information 
	(IN day_id int, OUT user_id INT, OUT date DATE) AS
	$$ SELECT days.user_id, days.date
    FROM days WHERE days.id = day_id; $$ 
    LANGUAGE SQL;

CREATE FUNCTION get_user_days (IN input_user_id INT) RETURNS TABLE(day_id INT, date DATE) AS
	$$ SELECT days.id, days.date FROM days
    WHERE days.user_id = input_user_id; $$
    LANGUAGE SQL;

CREATE FUNCTION is_it_activity_or_emotion (in type_id INT, OUT activity_or_emotion activity_or_emotion) AS
	$$ SELECT activity_or_emotion FROM types_of_activities_and_emotions
    WHERE types_of_activities_and_emotions.id = type_id; $$
    LANGUAGE SQL;
    
CREATE FUNCTION get_day_activities (IN input_day_id INT) RETURNS TABLE(type_id INT, proportion INT) AS
	$$ SELECT activities_and_emotions.type_id, activities_and_emotions.proportion FROM activities_and_emotions
    WHERE activities_and_emotions.day_id = input_day_id AND is_it_activity_or_emotion(activities_and_emotions.type_id) = 'activity'; $$
    LANGUAGE SQL;


