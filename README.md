# busybee
A discord bot for determining who can get ramen with me. Ingests `.ics` files and adds roles to users depending on which class they are currently attending.

## Usage

To develop on busybee, install Go, clone the project, and run `go run ./src` from the root directory. Be sure to provide an `env.yaml` configuration file with the following structure:
```yaml
BOT:
  APP_ID: string
  TOKEN: string
  ACTIVE_SERVER: string   # One of <SERVER1> or <SERVER2>, etc.
  GUILD_IDS:
    <SERVER1>: string
    <SERVER2: string
    ...
  CHANNEL_IDS:
    <SERVER1>: string
    <SERVER2: string
    ...
```
