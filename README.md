# Packaging application

THis repo contains a sample application for managing packaging. It is able to define packs size, and determine how many packs are needed for
optimal delivery of a package.

Contains two parts: UI and Backend

## Running locally (native)

### Backend
The backend is present under `server` and is build with Golang.

It is dependent on a database, so that needs to be setup for the app. There are two options to get this working:

#### Container database
Under the server folder, there is a docker compose setup for starting a PostgreSQL container.
Docker is required on the local machine for this to work. Can be started with
```bash
docker-compose -f db-docker-compose.yml logs -f
```

This exposes a database on `localhost:6432`, with name `server`, user `server` and pass `server`.
After this the golang app can be run with (this will set the default OS env vars for the app)
```bash
make run
```

#### Local database
If another database needs to be used, for example some local one on the system, just the make run command needs
to be adjusted to that database props

```bash
make run DB_HOST=val DB_PORT=val DB_USERNAME=val DB_PASSWORD=val DB_NAME=val
```

The application can be accessed on http://localhost:8080

### UI
The UI application is build with Angular, and node is needed to run it.

```bash
npm start
```

The application can be accessed on http://localhost:4200

## Running locally (Docker Compose)
There is also an option to run the UI and backend as a bundle with docker-compose.
The root folder has this file exposed, and can be executed with
```bash
docker-compose build
docker-compose up -d
```

The application can be accessed on http://localhost