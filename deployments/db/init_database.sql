SELECT 'CREATE DATABASE banks'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'banks')\gexec