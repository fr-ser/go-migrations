-- check that the schema exists
SELECT 1 / COUNT(*) FROM information_schema.schemata WHERE schema_name = 'db_stage';
SELECT 1 / COUNT(*) FROM information_schema.schemata WHERE schema_name = 'db_aggregation';

-- check for the roles exists
SELECT 1/ COUNT(*) FROM pg_roles WHERE rolname='db_viewer';
SELECT 1/ COUNT(*) FROM pg_roles WHERE rolname='db_agg_viewer';
