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

CREATE TABLE "types_of_activities" (
  "id" serial PRIMARY KEY,
  "user_id" integer REFERENCES users NOT NULL,
  "name" varchar NOT NULL,
  "color" varchar NOT NULL,
  "is_everyday" bool NOT NULL
);

CREATE TABLE "types_of_emotions" (
  "id" serial PRIMARY KEY,
  "user_id" integer REFERENCES users NOT NULL,
  "name" varchar NOT NULL,
  "color" varchar NOT NULL,
  "is_everyday" bool NOT NULL
);

CREATE TABLE "activities" (
  "type_id" integer REFERENCES types_of_activities NOT NULL,
  "day_id" integer REFERENCES days NOT NULL,
  "proportion" integer NOT NULL
);

CREATE TABLE "emotions" (
  "type_id" integer REFERENCES types_of_emotions NOT NULL,
  "day_id" integer REFERENCES days NOT NULL,
  "proportion" integer NOT NULL
);
