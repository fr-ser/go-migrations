-- check that the schema exists
SELECT 1 / COUNT(*) FROM information_schema.schemata WHERE schema_name = 'data_loader';

-- check for the user exists
SELECT 1/ COUNT(*) FROM pg_roles WHERE rolname='data_loader';
