### Gator — RSS Feed Aggregator CLI

 This project is Gator, a command-line RSS feed aggregator written in Go. It lets users register, add RSS feeds, follow/unfollow feeds, periodically scrape them, and browse aggregated posts — all from the terminal, backed by PostgreSQL.

   ---
  # Required dependencies

  1. Go toolchain (1.26+)

  The project is written in Go. Install the binary directly from the repository:
  go install github.com/AbdullahBasir/aggregator@latest
  This downloads, compiles, and places the aggregator binary in $GOPATH/bin (usually ~/go/bin).
  Make sure ~/go/bin is in your PATH.

  If cloning the repo locally, build with:
  go build -o aggregator
  
  2. PostgreSQL

  A running PostgreSQL instance. The app connects using the connection string in ~/.gatorconfig.json. The .env file (which is gitignored) shows the
  expected format:
  postgres://username:@localhost:5432/gator
  They'll need to create a database (e.g., gator) and a user with access to it.

  3. goose (database migration tool)

  The migrate.sh script uses goose to apply the SQL schema files. Install it with:
  go install github.com/pressly/goose/v3/cmd/goose@latest
  Then run:
  ./migrate.sh up

  4. sqlc (only if modifying queries/schema)

  The internal/database/ package is pre-generated from the SQL in sql/queries/. If they're just running the project, sqlc isn't needed — the generated
  .go files are committed. They'd only need it if they change the queries or schema:
  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

  ---
  # Manual setup steps

  Both of these are gitignored and must be created by hand:

  .env file (project root)

  DB_URL="postgres://<user>:@localhost:5432/gator"
  Used by migrate.sh to run migrations.

  ~/.gatorconfig.json (home directory)

  {
    "db_url": "postgres://<user>:@localhost:5432/gator?sslmode=disable",
    "current_user_name": ""
  }
  Used by the app at runtime to connect to Postgres and track the logged-in user.

  ---
  # Summary checklist

  ┌─────────────────────┬────────────────────────────────────┬──────────────────────────────────┐
  │     Dependency      │                Why                 │         Required to run?         │
  ├─────────────────────┼────────────────────────────────────┼──────────────────────────────────┤
  │ Go 1.26+            │ Compile the binary                 │ Yes                              │
  ├─────────────────────┼────────────────────────────────────┼──────────────────────────────────┤
  │ PostgreSQL          │ Database                           │ Yes                              │
  ├─────────────────────┼────────────────────────────────────┼──────────────────────────────────┤
  │ goose               │ Run schema migrations              │ Yes (one-time setup)             │
  ├─────────────────────┼────────────────────────────────────┼──────────────────────────────────┤
  │ sqlc                │ Regenerate DB code from SQL        │ No (generated code is committed) │
  ├─────────────────────┼────────────────────────────────────┼──────────────────────────────────┤
  │ .env file           │ DB_URL for migrations              │ Yes                              │
  ├─────────────────────┼────────────────────────────────────┼──────────────────────────────────┤
  │ ~/.gatorconfig.json │ Runtime DB connection + user state │ Yes                              │
  

# CLI commands

  aggregator register <name>       Create a user + log in
  aggregator login <name>          Switch to an existing user
  aggregator users                 List all users (marks current)
  aggregator addfeed <name> <url>  Add an RSS feed (requires login)
  aggregator feeds                 List all feeds with creator
  aggregator follow <url>          Follow a feed (requires login)
  aggregator following             List feeds you follow (requires login)
  aggregator unfollow <url>        Unfollow a feed (requires login)
  aggregator agg <duration>        Continuously scrape feeds (e.g. "30s", "1m")
  aggregator browse [limit]        Show recent posts from followed feeds (default 2)
  aggregator reset                 Delete all users (cascades to feeds, follows, posts)

