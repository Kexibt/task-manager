#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "postgres" --dbname "postgres" <<-EOSQL
	CREATE DATABASE task_manager;
	\c task_manager
	CREATE TYPE status_enum AS ENUM ('todo', 'in progress', 'done');
	CREATE TABLE public.users (
		login VARCHAR(100) PRIMARY KEY NOT NULL, 
		hashPass TEXT NOT NULL
		);
	CREATE TABLE public.tasks (
		id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY, 
		title VARCHAR(100) NOT NULL, 
		description TEXT, 
		status status_enum NOT NULL, 
		user_id VARCHAR(100) NOT NULL, 
		update_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_user
      		FOREIGN KEY(user_id)
	  		REFERENCES public.users(login)
	  );
	INSERT INTO public.users ("login", "hashpass") VALUES ('test', '1234');
EOSQL