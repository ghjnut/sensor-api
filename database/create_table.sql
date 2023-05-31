CREATE DATABASE code_challenge;

\c code_challenge

CREATE TABLE logs 
  (
    event_id serial PRIMARY KEY,
    event_date TIMESTAMPTZ NOT NULL,
    device_id VARCHAR(6) NOT NULL,
    temp_farenheit INT
  );
