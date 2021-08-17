A very simple key/value container with fixed expire time.

Use linked list to avoid full scan while cleaning.

Use fixed expire time to ensure O(1) insertion.