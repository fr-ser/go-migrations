# go-migrations

If you are looking for a tool do apply database migration written in GO look here
https://github.com/golang-migrate/migrate.

In contrast to the repo above, this repo is very use case specific and (therefore) not as rich in
drivers and functionality as the above mentioned.

## Migration Layout

An example of a migration structure can be found in the example folder `./example`.

Folder names starting with an underscore (`/_`) generally have special meaning.
Currently the following special folders exist:

- \_environments\_: This folder contains configuration files for different databases

The general layout looks like the following:

```
.
│   docker-compose.yaml: Optional. Only required for local development
│
└─── migrations: `./migrations` is the default location, but it can be any other folder as well
│   └─── _environments: see remarks above
│       │   development.yaml
│       │   production.yaml
│       │   ...
│   └─── <some_folder>: Any name is possible (typically represents a sub app of a database, just used for better grouping of migrations)
│       │   20171101000001_my_migration.sql
│       │   ...
|       └─── verify
│           │   20171101000001_my_migration.sql
│           │   ...
```

## Config Layout

Configuration files, which are stored in the `_environments` folder (see
[migration layout](#migration-layout)). These files look like this

```yaml
db_type: postgres
host: localhost
port: 35434
db_name: my_db
user: admin
password: admin_pass
```

## Commands

The migration tool includes a `--help` flag, which can be called on the tools itself or on any
subcommand:

```bash
./go_migrations --help
./go_migrations start --help
...
```

## TODOs

- Apply "Prepare" Scripts. Is this used?
