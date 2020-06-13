-- initially revoke every grants from the public schema: (this should be normal in every production environments)
REVOKE ALL ON SCHEMA PUBLIC FROM PUBLIC;

CREATE SCHEMA db_stage;
CREATE SCHEMA db_aggregation;

COMMENT ON SCHEMA db_stage IS 'working schema that keeps (most) volatile data, that is needed for loading date into the main and mart schema tables';
COMMENT ON SCHEMA db_aggregation IS 'smaller aggregation schema that keeps detailed smaller datasets for faster reporting';

-- comes from the bootstrap;
CREATE ROLE db_viewer;
CREATE ROLE db_agg_viewer;
COMMENT ON ROLE db_viewer IS 'role that will be assigned to all users who can read data from result tables in the main and mart schemas';
COMMENT ON ROLE db_agg_viewer IS 'role that will be assigned to all users who can read data only from mart schemas';

GRANT USAGE ON SCHEMA db_stage TO db_viewer, db_agg_viewer;
GRANT USAGE ON SCHEMA db_aggregation TO db_viewer, db_agg_viewer;

-- //@UNDO
-- SQL to undo the change goes here.
REVOKE USAGE ON SCHEMA db_aggregation FROM db_viewer, db_agg_viewer;
REVOKE USAGE ON SCHEMA db_stage FROM db_viewer, db_agg_viewer;

DROP ROLE db_viewer;
DROP ROLE db_agg_viewer;

DROP SCHEMA db_aggregation;
DROP SCHEMA db_stage;
