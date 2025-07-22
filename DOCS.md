# Docs

## Add channel

Add a YouTube channel using either the @handle or the channel ID (starts with 'UC') for `-i`.  
The category `-c` is optional; it will be created automatically if it doesn't exist.

```sh
ytcw add-channel -i "@IRyS" -c "hololive"
```

## Modify channel

Change the specified channel's category to a different one.

```sh
ytcw modify-channel -i "@IRyS" -c "Something else"
```

To delete or unset a category for the channel, use:

```sh
ytcw modify-channel -i "@IRyS" --unset-category
```

## Delete channel

Soft delete a channel and its videos.

```sh
ytcw delete-channel -i "@IRyS"
```

## Disable channel

Disable a channel for a specific duration. The fetcher daemon won't fetch new information for the specified channel.

```sh
ytcw disable-channel -i "@IRyS" -d "180h"
```

Here, `-d` represents duration. Using hours is typically what you want.
