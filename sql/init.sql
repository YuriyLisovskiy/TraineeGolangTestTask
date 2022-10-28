CREATE USER rest_api_user PASSWORD 'supersecret!';

CREATE SCHEMA rest_api_schema;
GRANT ALL PRIVILEGES ON SCHEMA rest_api_schema TO rest_api_user;
