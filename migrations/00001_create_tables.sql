-- +goose Up
CREATE TABLE "users" ( 
   "user_id" INT, 
   "name" VARCHAR, 
   "age" INT, 
   "spouse" INT 
);
CREATE UNIQUE INDEX "users_user_id" ON "users" ("user_id"); 
CREATE TABLE "activities" ( 
   "user_id" INT, 
   "date" TIMESTAMP, 
   "name" VARCHAR 
) PARTITION BY RANGE("date"); 
CREATE INDEX "activities_user_id_date" ON "activities" ("user_id", "date"); 
CREATE TABLE "activities_202011" PARTITION OF "activities" FOR VALUES FROM 
('2020-11-01'::TIMESTAMP) TO ('2020-12-01'::TIMESTAMP); 
CREATE TABLE "activities_202012" PARTITION OF "activities" FOR VALUES FROM 
('2020-12-01'::TIMESTAMP) TO ('2021-01-01'::TIMESTAMP);
-- +goose Down
DROP TABLE "users";
DROP TABLE "activities";