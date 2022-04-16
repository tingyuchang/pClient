DROP DATABASE IF EXISTS go;
CREATE DATABASE go;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Userdata;
\c go;
CREATE TABLE Users (
                       ID serial,
                       Username VARCHAR(100) PRIMARY KEY
);
CREATE TABLE Userdata (
                          UserId int NOT NULL,
                          name varchar(100),
                          Surname varchar(100),
                          Description varchar(200)
)