# UReplace
Path/Query Replacer

## Installation
```
go get -u github.com/theblackturtle/ureplace
```

## Usage
```
Usage of ureplace:
  -a    Append the value
  -b string
        Additional blacklist extensions (js,css)
  -f string
        Payload list
  -i string
        Where to inject
          all: replace all
          one: replace one by one
          2: replace the second path/param
          -2: replace the second path/param from the end (default "all")
  -m    Ignore media extensions
  -p    Replace in Paths
  -q    Replace in Queries
```

## Basic Usage
- Path
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -p newvalue

https://example.com/newvalue/newvalue/newvalue?param1=1&param2=2
```

- Query
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -q newvalue

https://example.com/path1/path2/path3?param1=newvalue&param2=newvalue
```

## Advanced Usage
#### Replace value one by one
- Path
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -i one -p newvalue

https://example.com/newvalue/path2/path3?param1=1&param2=2
https://example.com/path1/newvalue/path3?param1=1&param2=2
https://example.com/path1/path2/newvalue?param1=1&param2=2
```

- Query
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -i one -q newvalue

https://example.com/path1/path2/path3?param1=newvalue&param2=2
https://example.com/path1/path2/path3?param1=1&param2=newvalue
```

#### Replace the last value
- Path
```
❯  echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -i -1 -p newvalue

https://example.com/path1/path2/newvalue?param1=1&param2=2
```

- Query
```
echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -i -1 -q newvalue

https://example.com/path1/path2/path3?param1=1&param2=newvalue
```


#### Append value
- Path
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -a -p newvalue

https://example.com/path1newvalue/path2newvalue/path3newvalue?param1=1&param2=2
```

- Query
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -a -q newvalue

https://example.com/path1/path2/path3?param1=1newvalue&param2=2newvalue
```

#### Replace with payload list
- file.txt
```
newvalue1
newvalue2
```

- Path
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -p -f file.txt

https://example.com/newvalue1/newvalue1/newvalue1?param1=1&param2=2
https://example.com/newvalue2/newvalue2/newvalue2?param1=1&param2=2
```

- Query
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -q -f file.txt

https://example.com/path1/path2/path3?param1=newvalue1&param2=newvalue1
https://example.com/path1/path2/path3?param1=newvalue2&param2=newvalue2
```

#### Ignore media extensions
If you want to ignore media extensions, use `-m` flag.
- Images Extensions: `png, apng, bmp, gif, ico, cur, jpg, jpeg, jfif, pjp, pjpeg, svg, tif, tiff, webp, xbm`
- Audio Extensions: `3gp, aac, flac, mpg, mpeg, mp3, mp4, m4a, m4v, m4p, oga, ogg, ogv, mov, wav, webm`
- Font Extensions: `eot, woff, woff2, ttf, otf`

If you want to ignore additional extensions, use `-b` flag.
```
❯ echo 'https://example.com/path1/path2/path3?param1=1&param2=2' | ureplace -b js,css,php -q test
```