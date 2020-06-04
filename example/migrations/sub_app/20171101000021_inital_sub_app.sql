CREATE SCHEMA data_loader;

COMMENT ON SCHEMA data_loader IS 'working schema data_loader service';

CREATE ROLE data_loader INHERIT LOGIN PASSWORD 'data_loader';

COMMENT ON ROLE data_loader IS 'user for the data load service jobs';

GRANT ALL PRIVILEGES ON SCHEMA data_loader TO data_loader;

-- //@UNDO
-- SQL to undo the change goes here.

DROP SCHEMA data_loader;

DROP ROLE data_loader;

