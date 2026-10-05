// Copyright 2021-present The Atlas Authors. All rights reserved.
// This source code is licensed under the Apache 2.0 license found
// in the LICENSE file in the root directory of this source tree.

package docker

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDockerConfig(t *testing.T) {
	ctx := context.Background()

	// invalid config
	_, err := (&Config{}).Run(ctx)
	require.Error(t, err)

	// MySQL
	cfg, err := MySQL("latest", Out(io.Discard))
	require.NoError(t, err)
	require.Equal(t, &Config{
		Image: "docker.io/arigaio/mysql:latest",
		User:  url.UserPassword("root", pass),
		Env:   []string{"MYSQL_ROOT_PASSWORD=pass"},
		Port:  "3306",
		Out:   io.Discard,
	}, cfg)

	// MariaDB
	cfg, err = MariaDB("latest", Out(io.Discard))
	require.NoError(t, err)
	require.Equal(t, &Config{
		Image: "docker.io/arigaio/mariadb:latest",
		User:  url.UserPassword("root", pass),
		Env:   []string{"MYSQL_ROOT_PASSWORD=pass"},
		Port:  "3306",
		Out:   io.Discard,
	}, cfg)

	// PostgreSQL
	cfg, err = PostgreSQL("latest", Out(io.Discard))
	require.NoError(t, err)
	require.Equal(t, &Config{
		Image:    "docker.io/library/postgres:latest",
		User:     url.UserPassword("postgres", pass),
		Env:      []string{"POSTGRES_PASSWORD=pass"},
		Database: "postgres",
		Port:     "5432",
		Out:      io.Discard,
	}, cfg)

	// SQL Server
	cfg, err = SQLServer("2022-latest", Out(io.Discard))
	require.NoError(t, err)
	require.Equal(t, &Config{
		Image:    "mcr.microsoft.com/mssql/server:2022-latest",
		User:     url.UserPassword("sa", passSQLServer),
		Port:     "1433",
		Database: "master",
		Out:      io.Discard,
		Env: []string{
			"ACCEPT_EULA=Y",
			"MSSQL_PID=Developer",
			"MSSQL_SA_PASSWORD=" + passSQLServer,
		},
	}, cfg)

	// ClickHouse
	cfg, err = ClickHouse("23.11", Out(io.Discard))
	require.NoError(t, err)
	require.Equal(t, &Config{
		Image: "docker.io/clickhouse/clickhouse-server:23.11",
		User:  url.UserPassword("default", pass),
		Port:  "9000",
		Out:   io.Discard,
		Env: []string{
			"CLICKHOUSE_PASSWORD=pass",
		},
	}, cfg)
}

func TestFromURL(t *testing.T) {
	u, err := url.Parse("docker://mysql")
	require.NoError(t, err)
	cfg, err := FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver: "mysql",
		Image:  "docker.io/arigaio/mysql",
		User:   url.UserPassword("root", pass),
		Env:    []string{"MYSQL_ROOT_PASSWORD=pass"},
		Port:   "3306",
		Out:    io.Discard,
	}, cfg)

	u, err = url.Parse("docker://mysql/8")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver: "mysql",
		Image:  "docker.io/arigaio/mysql:8",
		User:   url.UserPassword("root", pass),
		Env:    []string{"MYSQL_ROOT_PASSWORD=pass"},
		Port:   "3306",
		Out:    io.Discard,
	}, cfg)

	u, err = url.Parse("docker://mysql/latest/test")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "mysql",
		Image:    "docker.io/arigaio/mysql:latest",
		Database: "test",
		Env:      []string{"MYSQL_ROOT_PASSWORD=pass", "MYSQL_DATABASE=test"},
		User:     url.UserPassword("root", pass),
		Port:     "3306",
		Out:      io.Discard,
		setup:    []string{"CREATE DATABASE IF NOT EXISTS `test`"},
	}, cfg)

	u, err = url.Parse("docker://postgres/13")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "postgres",
		Image:    "docker.io/library/postgres:13",
		Database: "postgres",
		Env:      []string{"POSTGRES_PASSWORD=pass"},
		User:     url.UserPassword("postgres", pass),
		Port:     "5432",
		Out:      io.Discard,
	}, cfg)

	// PostGIS.
	u, err = url.Parse("docker://postgis/14-3.4")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "postgres",
		Image:    "docker.io/postgis/postgis:14-3.4",
		Database: "postgres",
		Env:      []string{"POSTGRES_PASSWORD=pass"},
		User:     url.UserPassword("postgres", pass),
		Port:     "5432",
		Out:      io.Discard,
	}, cfg)

	u, err = url.Parse("docker://postgis/14-3.4/dev")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "postgres",
		Image:    "docker.io/postgis/postgis:14-3.4",
		Database: "dev",
		Env:      []string{"POSTGRES_PASSWORD=pass"},
		User:     url.UserPassword("postgres", pass),
		Port:     "5432",
		Out:      io.Discard,
		setup:    []string{`CREATE DATABASE "dev"`},
	}, cfg)

	// PGVector.
	u, err = url.Parse("docker://pgvector/pg17")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "postgres",
		Image:    "docker.io/pgvector/pgvector:pg17",
		Database: "postgres",
		Env:      []string{"POSTGRES_PASSWORD=pass"},
		User:     url.UserPassword("postgres", pass),
		Port:     "5432",
		Out:      io.Discard,
	}, cfg)

	u, err = url.Parse("docker://pgvector/pg17/dev")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "postgres",
		Image:    "docker.io/pgvector/pgvector:pg17",
		Database: "dev",
		Env:      []string{"POSTGRES_PASSWORD=pass", "POSTGRES_DB=dev"},
		User:     url.UserPassword("postgres", pass),
		Port:     "5432",
		Out:      io.Discard,
	}, cfg)

	// SQL Server
	u, err = url.Parse("docker://sqlserver")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "sqlserver",
		Image:    "mcr.microsoft.com/mssql/server",
		Database: "master",
		User:     url.UserPassword("sa", passSQLServer),
		Port:     "1433",
		Out:      io.Discard,
		Env: []string{
			"ACCEPT_EULA=Y",
			"MSSQL_PID=Developer",
			"MSSQL_SA_PASSWORD=" + passSQLServer,
		},
	}, cfg)

	u, err = url.Parse("docker://sqlserver/2022-latest")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "sqlserver",
		Image:    "mcr.microsoft.com/mssql/server:2022-latest",
		Database: "master",
		User:     url.UserPassword("sa", passSQLServer),
		Port:     "1433",
		Out:      io.Discard,
		Env: []string{
			"ACCEPT_EULA=Y",
			"MSSQL_PID=Developer",
			"MSSQL_SA_PASSWORD=" + passSQLServer,
		},
	}, cfg)

	u, err = url.Parse("docker://sqlserver/2019-latest/foo")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "sqlserver",
		setup:    []string{"CREATE DATABASE [foo]"},
		Image:    "mcr.microsoft.com/mssql/server:2019-latest",
		Database: "foo",
		User:     url.UserPassword("sa", passSQLServer),
		Port:     "1433",
		Out:      io.Discard,
		Env: []string{
			"ACCEPT_EULA=Y",
			"MSSQL_PID=Developer",
			"MSSQL_SA_PASSWORD=" + passSQLServer,
		},
	}, cfg)

	// Azure SQL Edge
	u, err = url.Parse("docker+sqlserver://mcr.microsoft.com/azure-sql-edge:1.0.7/foo")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "sqlserver",
		setup:    []string{"CREATE DATABASE [foo]"},
		Image:    "mcr.microsoft.com/azure-sql-edge:1.0.7",
		Database: "foo",
		User:     url.UserPassword("sa", passSQLServer),
		Port:     "1433",
		Out:      io.Discard,
		Env: []string{
			"ACCEPT_EULA=Y",
			"MSSQL_PID=Developer",
			"MSSQL_SA_PASSWORD=" + passSQLServer,
		},
	}, cfg)

	// ClickHouse
	u, err = url.Parse("docker://clickhouse")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver: "clickhouse",
		Image:  "docker.io/clickhouse/clickhouse-server",
		Env:    []string{"CLICKHOUSE_PASSWORD=pass"},
		User:   url.UserPassword("default", pass),
		Port:   "9000",
		Out:    io.Discard,
	}, cfg)

	// ClickHouse with tag
	u, err = url.Parse("docker://clickhouse/23.11")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver: "clickhouse",
		Image:  "docker.io/clickhouse/clickhouse-server:23.11",
		User:   url.UserPassword("default", pass),
		Env:    []string{"CLICKHOUSE_PASSWORD=pass"},
		Port:   "9000",
		Out:    io.Discard,
	}, cfg)
}

func TestFromURL_CustomImage(t *testing.T) {
	for _, tt := range []struct {
		url, image, db, dialect string
	}{
		// PostgreSQL (local and official images).
		{
			url:     "docker+postgres://local",
			image:   "local",
			db:      "postgres",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres://_/local/dev",
			image:   "local",
			db:      "dev",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres:///local:tag/dev",
			image:   "local:tag",
			db:      "dev",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres://postgres",
			image:   "postgres",
			db:      "postgres",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres://_/postgres/dev",
			image:   "postgres",
			db:      "dev",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres:///postgres:16/dev",
			image:   "postgres:16",
			db:      "dev",
			dialect: "postgres",
		},
		// User images.
		{
			url:     "docker+postgres://postgis/postgis:16",
			image:   "postgis/postgis:16",
			db:      "postgres",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres://postgis/postgis:16/dev",
			image:   "postgis/postgis:16",
			db:      "dev",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres://ghcr.io/namespace/image:tag",
			image:   "ghcr.io/namespace/image:tag",
			db:      "postgres",
			dialect: "postgres",
		},
		{
			url:     "docker+postgres://ghcr.io/namespace/image:tag/dev",
			image:   "ghcr.io/namespace/image:tag",
			db:      "dev",
			dialect: "postgres",
		},
		// MySQL.
		{
			url:     "docker+mysql://local",
			image:   "local",
			dialect: "mysql",
		},
		{
			url:     "docker+mysql:///local/dev",
			image:   "local",
			db:      "dev",
			dialect: "mysql",
		},
		{
			url:     "docker+mysql://user/image",
			image:   "user/image",
			dialect: "mysql",
		},
		{
			url:     "docker+mysql://user/image:tag/dev",
			image:   "user/image:tag",
			db:      "dev",
			dialect: "mysql",
		},
		{
			url:     "docker+mysql://_/mariadb:latest/dev",
			image:   "mariadb:latest",
			db:      "dev",
			dialect: "mysql",
		},
		// SQL Server.
		{
			url:     "docker+sqlserver://mcr.microsoft.com/mssql/server:2022-latest",
			image:   "mcr.microsoft.com/mssql/server:2022-latest",
			db:      "master",
			dialect: "sqlserver",
		},
		{
			url:     "docker+sqlserver://mcr.microsoft.com/mssql/server:2022-latest/dev",
			image:   "mcr.microsoft.com/mssql/server:2022-latest",
			db:      "dev",
			dialect: "sqlserver",
		},
		{
			url:     "docker+sqlserver://mcr.microsoft.com/mssql/server:latest",
			image:   "mcr.microsoft.com/mssql/server:latest",
			db:      "master",
			dialect: "sqlserver",
		},
		// ClickHouse.
		{
			url:     "docker+clickhouse://clickhouse/clickhouse-server:23.11",
			image:   "clickhouse/clickhouse-server:23.11",
			dialect: "clickhouse",
		},
		{
			url:     "docker+clickhouse://clickhouse/clickhouse-server:23.11/dev",
			image:   "clickhouse/clickhouse-server:23.11",
			db:      "dev",
			dialect: "clickhouse",
		},
	} {
		u, err := url.Parse(tt.url)
		require.NoError(t, err)
		cfg, err := FromURL(u)
		require.NoError(t, err)
		require.Equal(t, tt.image, cfg.Image)
		require.Equal(t, tt.db, cfg.Database)
		require.Equal(t, tt.dialect, cfg.driver)
	}
}

func TestFromURL_Podman(t *testing.T) {
	u, err := url.Parse("podman://mysql/8/dev")
	require.NoError(t, err)
	cfg, err := FromURL(u)
	require.NoError(t, err)
	require.Equal(t, &Config{
		driver:   "mysql",
		cli:      "podman",
		Image:    "docker.io/arigaio/mysql:8",
		Database: "dev",
		Env:      []string{"MYSQL_ROOT_PASSWORD=pass", "MYSQL_DATABASE=dev"},
		User:     url.UserPassword("root", pass),
		Port:     "3306",
		Out:      io.Discard,
		setup:    []string{"CREATE DATABASE IF NOT EXISTS `dev`"},
	}, cfg)

	u, err = url.Parse("podman+postgres://docker.io/library/postgres:16/dev")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Equal(t, "podman", cfg.cli)
	require.Equal(t, "postgres", cfg.driver)
	require.Equal(t, "docker.io/library/postgres:16", cfg.Image)
	require.Equal(t, "dev", cfg.Database)

	// Explicit CLI option overrides the scheme.
	cfg, err = FromURL(u, CLI("docker"))
	require.NoError(t, err)
	require.Equal(t, "docker", cfg.cli)

	// Docker URLs leave the CLI to be resolved on Run.
	u, err = url.Parse("docker+postgres://docker.io/library/postgres:16/dev")
	require.NoError(t, err)
	cfg, err = FromURL(u)
	require.NoError(t, err)
	require.Empty(t, cfg.cli)

	u, err = url.Parse("nerdctl://mysql/8")
	require.NoError(t, err)
	_, err = FromURL(u)
	require.EqualError(t, err, `unsupported container runtime "nerdctl"`)
}

func TestResolveCLI(t *testing.T) {
	fakeBin := func(t *testing.T, names ...string) string {
		dir := t.TempDir()
		for _, n := range names {
			require.NoError(t, os.WriteFile(filepath.Join(dir, n), []byte("#!/bin/sh\n"), 0755))
		}
		return dir
	}
	resolve := func(c *Config) string {
		require.NoError(t, c.resolveCLI())
		return c.cli
	}
	t.Setenv(CLIEnv, "")

	// Explicitly set.
	t.Setenv("PATH", fakeBin(t, "docker"))
	require.Equal(t, "podman", resolve(&Config{cli: "podman"}))

	// Docker is preferred.
	t.Setenv("PATH", fakeBin(t, "docker", "podman"))
	require.Equal(t, "docker", resolve(&Config{}))

	// Fallback to podman.
	t.Setenv("PATH", fakeBin(t, "podman"))
	require.Equal(t, "podman", resolve(&Config{}))

	// Environment override.
	t.Setenv(CLIEnv, "/usr/local/bin/podman")
	c := &Config{}
	require.Equal(t, "/usr/local/bin/podman", resolve(c))
	require.True(t, c.isPodman())
	t.Setenv(CLIEnv, "")

	// Nothing found.
	t.Setenv("PATH", fakeBin(t))
	require.EqualError(t, (&Config{}).resolveCLI(), "no container runtime found: install docker or podman")
}

func TestIsPodman(t *testing.T) {
	for cli, want := range map[string]bool{
		"":                      false,
		"docker":                false,
		"/usr/bin/docker":       false,
		"podman":                true,
		"/usr/local/bin/podman": true,
		"podman.exe":            true,
		"Podman.EXE":            true,
	} {
		require.Equal(t, want, (&Config{cli: cli}).isPodman(), cli)
	}
}

func TestImageURL(t *testing.T) {
	for img, u := range map[string]string{
		"postgres:15":                    "docker+postgres://_/postgres:15",
		"postgres":                       "docker+postgres://_/postgres",
		"postgis/postgis:14-3.4":         "docker+postgres://postgis/postgis:14-3.4",
		"ghcr.io/namespace/postgres:tag": "docker+postgres://ghcr.io/namespace/postgres:tag",
	} {
		got, err := ImageURL(DriverPostgres, img)
		require.NoError(t, err)
		require.Equal(t, u, got.String())
	}
	for img, u := range map[string]string{
		"mcr.microsoft.com/azure-sql-edge:1.0.7":     "docker+sqlserver://mcr.microsoft.com/azure-sql-edge:1.0.7",
		"mcr.microsoft.com/mssql/server:2022-latest": "docker+sqlserver://mcr.microsoft.com/mssql/server:2022-latest",
	} {
		got, err := ImageURL(DriverSQLServer, img)
		require.NoError(t, err)
		require.Equal(t, u, got.String())
	}
}

func TestContainerURL(t *testing.T) {
	c := &Container{
		Config: Config{
			driver: "postgres",
			User:   url.UserPassword("postgres", "pass"),
		},
		Port: "5432",
	}
	u, err := c.URL()
	require.NoError(t, err)
	require.Equal(t, "postgres://postgres:pass@localhost:5432/?sslmode=disable", u.String())

	// With DOCKER_HOST
	t.Setenv("DOCKER_HOST", "tcp://host.docker.internal:2375")
	u, err = c.URL()
	require.NoError(t, err)
	require.Equal(t, "postgres://postgres:pass@host.docker.internal:5432/?sslmode=disable", u.String())

	// Podman uses CONTAINER_HOST instead of DOCKER_HOST.
	c.cli = "podman"
	t.Setenv("CONTAINER_HOST", "")
	u, err = c.URL()
	require.NoError(t, err)
	require.Equal(t, "postgres://postgres:pass@localhost:5432/?sslmode=disable", u.String())
	t.Setenv("CONTAINER_HOST", "ssh://user@podman.internal:22/run/podman/podman.sock")
	u, err = c.URL()
	require.NoError(t, err)
	require.Equal(t, "postgres://postgres:pass@podman.internal:5432/?sslmode=disable", u.String())
}
