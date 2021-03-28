# GoZIP


Utility to zip/unzip gzip-files
Only unzip is implemented so far.
It's messy, but it works. 


### Usage:
zip:
```
$ gobase64 <file>
```
unzip:
```
$ gobase64 -d <file>
```

If no file is provided, it will read from stdin.


Based on:
- https://tools.ietf.org/html/rfc1952
- https://tools.ietf.org/html/rfc1951
- https://www.infinitepartitions.com/art001.html
