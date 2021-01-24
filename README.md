# Ariadne
Ariadne is a file indexing and finding tool

## To compile the tool on Linux:
* install golang
* install sqlite3
* go get github.com/mattn/go-sqlite3
* go get github.com/rjeczalik/notify
* git clone https://github.com/Matnsz/Ariadne (or download the zip and extract)
* sqlite3 Ariadne/ariadne-daemon/files.db < Ariadne/ariadne-daemon/files.db.sql
* sqlite3 Ariadne/ariadne-daemon/watched_dirs.db < Ariadne/ariadne-daemon/watched_dirs.db.sql
* cd Ariadne/ariadne-daemon/ && go build
* cd ../ariadne-cli/ && go build
