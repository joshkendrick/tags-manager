## Tags Manager

This is a small utility to manipulate the tags I've put on my photos and videos and back them up to a boltDB database so they don't get lost (like when Windows changes something and I lose them, or it turns out there isn't a way to read them anymore... I don't know... it makes me feel better?)

Makes use of the [exiftool](https://www.sno.phy.queensu.ca/~phil/exiftool/) executable v11.10 (saved in this repo, and already outdated)

I built the Windows exporter executable from reading [this](https://github.com/golang/go/wiki/WindowsCrossCompiling)

#### If you run this on your target machine, it will default to the correct $GOOS and $GOARCH:
```
go build -o tags-manager_v0.2.0.exe main.go
```

### Usage
Run the program: `go run main.go`

#### Commands:
```
index <path> -> adds file tags in that path to the database
list -> displays all data we have
list <tag_or_absolute_filepath> -> displays data about that key
clear -> clears out the database, starts fresh
exit -> PEACE
```