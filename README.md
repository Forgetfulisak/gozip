# GoZIP


Utility to zip/unzip gzip-files
Only unzip is implemented so far.
It's messy and slow, but it works! 


### Installation
```
go install github.com/forgetfulisak/gozip@latest
```

### Usage:
zip:
```
$ gozip <file>
```
unzip:
```
$ gozip -d <file>
```

If no file is provided, it will read from stdin.


Based on:
- https://tools.ietf.org/html/rfc1952
- https://tools.ietf.org/html/rfc1951
- https://www.infinitepartitions.com/art001.html
