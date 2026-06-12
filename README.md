# ytcw - YouTube Channel Watcher

## About

ytcw is a program that calls [yt-dlp](https://github.com/yt-dlp/yt-dlp) to fetch channel information (videos and shorts), store it in a PostgreSQL database, and provide REST API endpoints.

I really like the Holodex project, but it mostly tracks VTuber clippers. I just wanted to see the VTubers’ own content (or some channels that only follow certain YouTubers).  
You can follow channels on Holodex, but I haven’t found a way to filter out live streams.  
This project isn’t limited to the VTuber community - you can use it to fetch any YouTube channel.

## Quick start

Create a database in PostgreSQL, then edit the config file to fit your needs.

First, create a config file at `/etc/ytcw/config.toml`  
For an example, see the default config file: [config/config.toml](./config/config.toml)

```sh
ytcw migrate # Run database migrations (does not create the database itself)
```

```sh
ytcw add-channel -i @MoriCalliope -c Hololive # Add a YouTube channel with a category
```

```sh
ytcw serve # Start the REST API server
```

```sh
ytcw daemon # Start the fetcher daemon
```

The database columns and API fields are mostly aligned with the yt-dlp JSON output.

More information can be found in the [docs](./DOCS.md) file.

Auto-generated Swagger documentation can be found in the [docs/](./docs/) folder or at the `/swagger/` URL (the trailing slash is important).

## systemd service files

Modify these files to suit your own setup, pay special attention to the `User` and `Group` options.  
It’s also expected that the software is installed in `/usr/local/bin/` folder.

```ini
[Unit]
Description=ytcw fetcher daemon
After=network.target
Requires=network-online.target

[Service]
Restart=always
RestartSec=5
Type=exec
ExecStart=/usr/local/bin/ytcw daemon
User=foo
Group=foo

[Install]
WantedBy=multi-user.target
```

```ini
[Unit]
Description=ytcw API server
After=network.target
Requires=network-online.target

[Service]
Restart=always
RestartSec=1
Type=exec
ExecStart=/usr/local/bin/ytcw serve
User=foo
Group=foo

[Install]
WantedBy=multi-user.target
```

Copy the service files into `/etc/systemd/system/` and then run:

```sh
sudo systemctl daemon-reload
```

```sh
sudo systemctl enable --now ytcw.service
```

```sh
sudo systemctl enable --now ytcw-api.service
```

If you make any changes to the config file, please restart the services using systemd.

## Frontend

- [ytcw-frontend](https://github.com/mY9Yd2/ytcw-frontend)

Example caddy file:

```text
example.com {
        # ytcw server

        root * /var/www/example.com

        encode zstd gzip

        file_server

        reverse_proxy /api/* :4040

        log {
                # log file optionally, see the Caddy docs
        }
        
        header {
                # additional headers, if you need them
        }
}
```

Just put the built (`ng build`) files in the `/var/www/example.com` folder.

## Development

You can set up a local configuration inside the `config/` folder named `config.local.toml` and override only the settings you need, such as the database password.
