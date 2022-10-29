CREATE USER rest_api_user PASSWORD 'supersecret!';

CREATE SCHEMA rest_api;
GRANT ALL PRIVILEGES ON SCHEMA rest_api TO rest_api_user;
