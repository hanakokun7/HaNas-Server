# HaNas

[한국어](README.ko.md) | [日本語](README.ja.md)

HaNas is a self-hosted personal file storage application built with Go and SQLite. It provides account management, file and folder operations, trash, sharing links, and media previews in a web browser.

## Features

- User registration and login
- JWT-based authentication cookies
- Folder and file creation
- File upload and overwrite
- Copy, move, rename, and delete operations
- Trash with permanent deletion after 30 days
- Shared file link creation and removal
- Image and video thumbnail generation
- Image, video, audio, and document previews
- Account and associated data deletion

## Technology Stack

- Go `1.25.1`
- `net/http`
- GORM and glebarez SQLite
- JWT
- bcrypt
- React-based static frontend
- Optional: `ffmpeg`

The static frontend files are embedded in the Go binary, so no separate frontend build step is required.

## Requirements

- Go 1.25 or later
- Permission to use port `80`
- No additional system dependency is required for the basic features
- `ffmpeg` must be available on `PATH` for video thumbnails

On macOS or Linux, using port 80 may require administrator privileges.

```bash
brew install go ffmpeg
```

The command above is a macOS example. Go and ffmpeg can be installed with your operating system's package manager.

## Quick Start

Build from the project root:

```bash
cd /Users/user/Desktop/server
go mod download
go build -o hanas .
```

Place `config.js` next to the executable and start the server:

```bash
./hanas
```

The default frontend configuration is:

```javascript
window.APP_CONFIG = {
  BACKEND_URL: 'http://127.0.0.1:80'
}
```

Open `http://127.0.0.1` or the server host address in a browser.

If port 80 is not available to the current user, run the server with temporary elevated permissions:

```bash
sudo ./hanas
```

For production, place it behind a TLS reverse proxy and set `BACKEND_URL` in `config.js` to the public URL.

## Configuration

`config.js` is read from the directory containing the executable. When deploying the binary elsewhere, copy `config.js` next to the executable.

Example:

```text
/opt/hanas/
├── hanas
├── config.js
├── database.db
├── data/
└── thumbnails/
```

The server always listens on `:80`; no environment variable is provided for changing the port.

## Data Storage

The application uses the following paths relative to the current working directory:

| Path | Purpose |
| --- | --- |
| `database.db` | Accounts, nodes, sharing records, and trash metadata |
| `data/` | Original uploaded files |
| `thumbnails/` | Generated thumbnails |

Back up all three locations when backing up application data.

```bash
mkdir -p backups
cp database.db backups/database.db
cp -R data backups/data
cp -R thumbnails backups/thumbnails
```

Stop the server before a backup or use a database-safe backup procedure; copying a live SQLite database may result in an inconsistent backup.

## Project Structure

```text
.
├── app.go            # Server startup and HTTP routes
├── logic.go          # Authentication, file, and sharing API handlers
├── exec.go           # File operations and data manipulation
├── object.go         # Data models
├── go.mod            # Go dependencies
├── index.html        # Web application entry point
├── config.js         # Frontend backend address
├── assets/
│   ├── index.js      # Built frontend
│   ├── index.css     # Frontend styles
│   ├── icon.ico
│   └── apple-touch-icon.png
├── data/             # Uploaded file storage
├── thumbnails/       # Thumbnail storage
└── database.db       # SQLite database
```

## How It Works

1. On startup, the server opens SQLite and automatically migrates the required tables.
2. The first startup generates a JWT signing secret and stores it in the database.
3. Registration creates a root folder for the new user.
4. A successful login stores the JWT in an `HttpOnly`, `SameSite=Strict` cookie.
5. Uploaded files are stored under `data/` using numeric file IDs, while node metadata is stored in SQLite.
6. Items moved to trash are permanently deleted by a background task after 30 days.

The JWT signing key is stored in the database, so preserve `database.db` during redeployment if existing sessions and token compatibility must be retained.

## Main API

Most user APIs require the `token` cookie.

| Method | Route | Description |
| --- | --- | --- |
| `POST` | `/register` | Create an account |
| `POST` | `/login` | Log in |
| `POST` | `/logout` | Log out |
| `GET` | `/me` | Get the current user |
| `GET` | `/node/` | Get the root folder |
| `GET` | `/node/{id}` | Get a folder |
| `GET` | `/file/{nodeId}` | Download a file |
| `GET` | `/thumbnail/{nodeId}` | Get a thumbnail |
| `POST` | `/upload` | Upload a file or create a folder |
| `POST` | `/copy` | Copy an item |
| `POST` | `/move` | Move an item |
| `POST` | `/rename` | Rename an item |
| `POST` | `/delete` | Permanently delete an item |
| `POST` | `/trash/in` | Move an item to trash |
| `GET` | `/trash` | List trash items |
| `GET` | `/trash/info/{id}` | Get trash item details |
| `POST` | `/trash/restore/{id}` | Restore a trash item |
| `POST` | `/share/create` | Create a shared link |
| `POST` | `/share/delete` | Delete a shared link |
| `GET` | `/s/{token}` | Access a shared file |
| `POST` | `/delete-account` | Delete the account and data |

## Upload Examples

Upload a small file with a JSON request:

```bash
curl -b cookie.txt \
  -H 'Content-Type: application/json' \
  -d '{
    "filename": "hello.txt",
    "is_dir": false,
    "data_base64": "aGVsbG8="
  }' \
  http://127.0.0.1/upload
```

Multipart upload example:

```bash
curl -b cookie.txt \
  -F 'filename=photo.jpg' \
  -F 'oya_id=1' \
  -F 'file=@photo.jpg' \
  http://127.0.0.1/upload
```

Use the cookie returned by the actual login response.

## Development and Verification

Download dependencies and build:

```bash
go mod download
go build ./...
```

The project currently has no built-in Go test suite.
