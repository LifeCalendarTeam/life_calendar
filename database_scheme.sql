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

CREATE TABLE "activities_and_emotions" (
  "type_id" integer REFERENCES types_of_activities_and_emotions NOT NULL,
  "day_id" integer REFERENCES days ON DELETE CASCADE NOT NULL,
  "proportion" integer NOT NULL
);
