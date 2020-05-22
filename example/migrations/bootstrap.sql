CREATE SCHEMA extensions_functions;
CREATE EXTENSION pg_stat_statements WITH SCHEMA extensions_functions;
CREATE EXTENSION pgcrypto WITH SCHEMA extensions_functions;
CREATE EXTENSION tablefunc WITH SCHEMA extensions_functions;

CREATE TABLE public.migrations_changelog (
  id VARCHAR(14) NOT NULL PRIMARY KEY
  , name TEXT NOT NULL
  , applied_at timestamptz NOT NULL
);
