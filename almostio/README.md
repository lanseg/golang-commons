# Almostio

Almost usable utility classes that help with io and byte processing.

## Marshal

Reader and writer for a given format.

## FixedSizeWriter

Forwards at most maxCapacity byte to the underlying writer

## MultiWriteCloser

MultiWriteCloser is similar to the io.MultiWriter, but with the Close() method.

## Overlay

Overlay is an extra layer between the filesystem (io) and the user code. For example, when you
want to keep files with extra long filenames with uncommon characters. Allows to workaround several
filesystem limitations while keeping files stored in the system as is and accessible by other 
programs.

### Example 1: Http downloader cache

```go
package main

import (
        "fmt"
        "io"
        "os"

        ol "github.com/lanseg/golang-commons/almostio"
)

func main() {
        lo, err := ol.NewLocalOverlay("cache", ol.NewJsonMarshal[ol.OverlayMetadata]())
        if err != nil {
                fmt.Printf("Cannot create overlay: %s\n", err)
                os.Exit(-1)
        }

        afile := "http://someurl.domain/привет こんにちは"
        ow, _ := lo.OpenWrite(afile)
        ow.Write([]byte("Hello world"))
        ow.Close()

        r, _ := lo.OpenRead(afile)
        data, _ := io.ReadAll(r)
        fmt.Printf("Data: %s\n", string(data))
}
```

And in the result you will get a new folder with the following file structure:
```
./cache
./cache/8bfb0bf7_http_someurl.domain_
./cache/.overlay
./cache/.overlay/metadata.json
```

And the metadata.json with the system information:
```json
{
    "fileMetadata": {
        "http://someurl.domain/привет こんにちは": {
            "name": "http://someurl.domain/привет こんにちは",
            "localName": "8bfb0bf7_http_someurl.domain_",
            "sha256": "64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c",
            "mime": "text/plain; charset=utf-8"
        }
    }
}
```
