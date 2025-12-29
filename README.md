# filebrowser-quantum-sync

Extracts users and hashed passwords from [FileBrowser Quantum](https://github.com/gtsteffaniak/filebrowser) and generates a `users` file compatible with htpasswd for multi-container setups.

## Usage (Rootles Mode)

```
podman run -d --name filebrowser-quantum-sync \
  -v FILEBROWSER_DB:/db:ro \
  -v FILEBROWSER_SYNC:/config \
  ghcr.io/cryinkfly/filebrowser-quantum-sync:latest
```
or
```
docker run -d --name filebrowser-quantum-sync \
  -v FILEBROWSER_DB:/db:ro \
  -v FILEBROWSER_SYNC:/config \
  ghcr.io/cryinkfly/filebrowser-quantum-sync:latest
```

### Notes

- `FILEBROWSER_DB` → Path to the FileBrowser Quantum BoltDB, which must be the same volume used by FileBrowser Quantum. It is mounted read-only to prevent accidental writes.
- `FILEBROWSER_SYNC` → Path where the users file will be written. This volume can later be mounted into another container (e.g., `Radicale`) to provide the `htpasswd` file for `authentication`.
- The `FILEBROWSER_DB` volume must be mounted with `:ro` to prevent accidental writes.
- The sync container `permanently overwrites the users file` on each run to keep it in `sync with the database`.

---

#### Example: [Radicale](https://radicale.org/v3.html) container integration (Rootles Mode)

```
podman run -d --name filebrowser-quantum \
  -v FILEBROWSER_FILES:/srv \
  -v FILEBROWSER_DB:/home/filebrowser/sync \
  -p 80:80 \
  docker.io/gtstef/filebrowser:beta
```
or
```
  docker run -d --name filebrowser-quantum \
  -v FILEBROWSER_FILES:/srv \
  -v FILEBROWSER_DB:/home/filebrowser/sync \
  -p 80:80 \
  docker.io/gtstef/filebrowser:beta
```

Change the database file path (config.yaml):

```
server:
  database: "sync/database.db"
```
Create filebrowser-quantum-sync volumes

```
podman volume create RADICALE_CONFIG

wget -O $HOME/.local/share/containers/storage/volumes/RADICALE_CONFIG \
https://raw.githubusercontent.com/cryinkfly/filebrowser-quantum-sync/refs/heads/main/radicale/config
```
or
```
docker volume create RADICALE_CONFIG

wget -O $HOME/.local/share/docker/volumes/RADICALE_CONFIG \
https://raw.githubusercontent.com/cryinkfly/filebrowser-quantum-sync/refs/heads/main/radicale/config
```

Start the filebrowser-quantum-sync container
```
podman run -d --name radicale \
  -p 5232:5232 \
  -v RADICALE_CONFIG:/etc/radicale:ro \
  -v RADICALE_DATA:/var/lib/radicale \
  -v FILEBROWSER_SYNC:/etc/radicale/sync:ro \
  ghcr.io/kozea/radicale:latest
```
or
```
docker run -d --name radicale \
  -p 5232:5232 \
  -v RADICALE_CONFIG:/etc/radicale:ro \
  -v RADICALE_DATA:/var/lib/radicale \
  -v FILEBROWSER_SYNC:/etc/radicale/sync:ro \
  ghcr.io/kozea/radicale:latest
```

The sync container currently runs permanently in the background and updates the users file every 10 seconds.

This interval is temporary and is planned to be increased to 30 seconds in a future revision.

When a user’s password is changed in FileBrowser Quantum, the updated bcrypt hash is automatically written to the users file during the next sync cycle.

At the moment, when a user is removed in FileBrowser Quantum, the corresponding user is not automatically removed from dependent services such as Radicale. This is intentional to prevent accidental user or data removal.

A separate command or dedicated mode for synchronizing user deletions is planned for a future release, giving administrators explicit control over when removed users are propagated to other services.
