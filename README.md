# SS-freeURI
A Shadowsocks free URI generator.
## Usage
#### Run directly
```shell
go run main.go
```
#### Build from source
```shell
git clone https://github.com/goFindAlex/SS-freeURI.git && cd SS-freeURI
go build
```
#### Or download compiled binary
[Release]https://github.com/goFindAlex/SS-freeURI/releases
## Notice
- For v1.0, you will need to use proxy listen on "http://127.0.0.1/1087" (which is default setting in ShadowsocksX-NG).
- Create "/Users/username/Desktop/nodes.txt" after running 
## Feature
- The URIs in nodes.txt can be imported to ShadowsocksX-NG directly through "Import Server URIs..."
- The URIs are sorted by average round-trip time (from low to high). So theoretically the URI in the 1st line would be the fastest one for you.