CREATE TYPE activity_or_emotion AS ENUM ('activity', 'emotion');

CREATE TABLE "users" (
  "id" serial PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "second_name" varchar NOT NULL,
  "city" varchar,
  "birthday" date,
  "avatar" varchar,
  "password_hash" varchar NOT NULL
);

CREATE TABLE "days" (
  "id" serial PRIMARY KEY,
  "user_id" integer REFERENCES users NOT NULL,
  "date" date NOT NULL
);

CREATE TABLE "types_of_activities_and_emotions" (
  "id" serial PRIMARY KEY,
  "user_id" integer REFERENCES users NOT NULL,
  "name" varchar NOT NULL,
  "color" varchar NOT NULL,
  "is_everyday" bool NOT NULL,
  "activity_or_emotion" activity_or_emotion NOT NULL
);

CREATE FUNCTION does_activity_or_emotion_belong_to_user_of_the_day(type_id integer, day_id integer) RETURNS bool
    STABLE
AS
$$
SELECT days.user_id = types_of_activities_and_emotions.user_id
FROM days,
     types_of_activities_and_emotions
WHERE days.id = day_id
  AND types_of_activities_and_emotions.id = type_id
$$ LANGUAGE sql;

CREATE TABLE "activities_and_emotions" (
  "type_id" integer REFERENCES types_of_activities_and_emotions NOT NULL,
  "day_id" integer REFERENCES days NOT NULL,
  "proportion" integer NOT NULL,
  CHECK (does_activity_or_emotion_belong_to_user_of_the_day(type_id,day_id))
);
