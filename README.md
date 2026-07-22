### Gator — RSS Feed Aggregator CLI - Boot.dev project

Gator is a command-line RSS feed aggregator written in Go. It lets you register
users, add RSS feeds, follow/unfollow feeds, periodically scrape them, and
browse aggregated posts — all from the terminal, backed by PostgreSQL.

---

## Prerequisites

| Dependency | Required to run? | Notes |
|---|---|---|
| Go 1.26+ | Yes | Compiles the binary |
| PostgreSQL | Yes | Database |
| goose | Yes (one-time) | Runs schema migrations |
| sqlc | No | Only needed to regenerate DB code from SQL (generated code is committed) |

---

## Installation

### Option 1: Install directly with `go install`

```sh
go install github.com/AbdullahBasir/aggregator@latest
```

This compiles and places the binary in `$GOPATH/bin` (usually `~/go/bin`).
Make sure `~/go/bin` is in your `PATH`.

### Option 2: Clone and build locally

```sh
git clone https://github.com/AbdullahBasir/aggregator.git
cd aggregator
go build -o aggregator
```

---

## Configuration

Two config files are required (both gitignored — create them manually):

### `.env` (project root)

Used by the migration script.

```
DB_URL="postgres://<user>:@localhost:5432/gator"
```

### `~/.gatorconfig.json` (home directory)

Used by the app at runtime for the DB connection and to track the logged-in user.

```json
{
  "db_url": "postgres://<user>:@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

---

## Database Setup

1. Create a PostgreSQL database (e.g., `gator`) and a user with access to it.
2. Install [goose](https://github.com/pressly/goose):
   ```sh
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```
3. Run the migrations:
   ```sh
   ./migrate.sh up
   ```

---

## CLI Commands

### User management

| Command | Description |
|---|---|
| `aggregator register <name>` | Create a user and log in |
| `aggregator login <name>` | Switch to an existing user |
| `aggregator users` | List all users (marks current) |
| `aggregator reset` | Delete all users (cascades to feeds, follows, posts) |

### Feeds

| Command | Description |
|---|---|
| `aggregator addfeed <name> <url>` | Add an RSS feed (requires login) |
| `aggregator feeds` | List all feeds with their creator |
| `aggregator follow <url>` | Follow a feed (requires login) |
| `aggregator following` | List feeds you follow (requires login) |
| `aggregator unfollow <url>` | Unfollow a feed (requires login) |

### Reading

| Command | Description |
|---|---|
| `aggregator agg <duration>` | Continuously scrape feeds (e.g. `30s`, `1m`) |
| `aggregator browse [limit]` | Show recent posts from followed feeds (default: 2) |

---

## Development

If you modify the SQL queries or schema, regenerate the database code with
[sqlc](https://github.com/sqlc-dev/sqlc):

```sh
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc generate
```
