-- First drop databases if they exist
DROP DATABASE IF EXISTS kong;
DROP DATABASE IF EXISTS konga_db;
DROP USER IF EXISTS kong;

-- Create user
CREATE USER kong WITH PASSWORD 'kong';

-- Create databases
CREATE DATABASE kong OWNER kong;
CREATE DATABASE konga_db OWNER kong;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE kong TO kong;
GRANT ALL PRIVILEGES ON DATABASE konga_db TO kong;