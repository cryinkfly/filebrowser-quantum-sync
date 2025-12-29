# filebrowser-quantum-sync

> Author: Steve Zabka <br/>
> Author-URL: https://cryinkfly.com <br/>
> License: Apache 2.0

Extracts users and hashed passwords from FileBrowser Quantum and generates a `users` file compatible with htpasswd for multi-container setups.

## Usage

```
podman run --rm \
  -v FILEBROWSER_DB:/db:ro \
  -v FILEBROWSER_SYNC:/sync \
  docker.io/filebrowser-quantum-sync:latest
```

### Notes

- `FILEBROWSER_DB` → Path to the FileBrowser Quantum BoltDB (read-only).
- `FILEBROWSER_SYNC` → Path where the users file will be written. This volume can later be mounted into another container (e.g., `Radicale`) to provide the `htpasswd` file for `authentication`.
- The `FILEBROWSER_DB` volume must be mounted with `:ro` to prevent accidental writes.
- The sync container `permanently overwrites the users file` on each run to keep it in `sync with the database`.

---

#### Example: Radicale container integration

```
mkdir -p RADICALE_CONFIG
```
```
wget -O RADICALE_CONFIG/config \
https://raw.githubusercontent.com/cryinkfly/filebrowser-quantum-sync/refs/heads/main/radicale/config
```
```
podman run -d --name radicale \
  -p 5232:5232 \
  -v RADICALE_CONFIG:/etc/radicale:ro \
  -v RADICALE_DATA:/var/lib/radicale \
  -v FILEBROWSER_SYNC:/etc/radicale/sync:ro \
  ghcr.io/kozea/radicale:latest
```
