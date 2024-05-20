CREATE DATABASE costurai;

\c costurai;

CREATE TABLE "dressmakers" (
    "id" TEXT NOT NULL,
    "name" varchar(100) NOT NULL,
    "contact" varchar(50) NOT NULL,
    "grade" decimal(2, 1) NOT NULL,
    "location" GEOGRAPHY(Point, 4326),
    "services" varchar(100) NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "dressmakers_pk" PRIMARY KEY ("id")
);

CREATE TABLE "users" (
    "id" TEXT NOT NULL,
    "name" varchar(100) NOT NULL,
    "email" varchar(100) NOT NULL,
    "password" varchar(100) NOT NULL,
    "location" GEOGRAPHY(Point, 4326),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "users_pk" PRIMARY KEY ("id")
);

CREATE TABLE "dressmakers_reviews" (
  "id" TEXT NOT NULL,
  "dressmaker_id" TEXT NOT NULL,
  "user_id" TEXT NOT NULL,
  "comment" TEXT NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT "dressmakers_reviews_pk" PRIMARY KEY ("id"),
  FOREIGN KEY ("dressmaker_id") REFERENCES "dressmakers" ("id"),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
