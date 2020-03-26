package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	path2 "path"
	"sort"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
)

var (
	ImagesExt = []string{"png", "apng", "bmp", "gif", "ico", "cur", "jpg", "jpeg", "jfif", "pjp", "pjpeg", "svg", "tif", "tiff", "webp", "xbm"}
	AudioExt  = []string{"3gp", "aac", "flac", "mpg", "mpeg", "mp3", "mp4", "m4a", "m4v", "m4p", "oga", "ogg", "ogv", "mov", "wav", "webm"}
	FontExt   = []string{"eot", "woff", "woff2", "ttf", "otf"}
	OtherExt  = []string{}
)

var (
	appendMode     bool
	query          bool
	path           bool
	removeMediaExt bool
	place          string
	blacklistExt   string
	payloadFile    string
	payloadList    []string
)

func main() {
	flag.BoolVar(&appendMode, "a", false, "Append the value")
	flag.BoolVar(&removeMediaExt, "m", false, "Ignore media extensions")
	flag.BoolVar(&query, "q", false, "Replace in Queries")
	flag.BoolVar(&path, "p", false, "Replace in Paths")
	flag.StringVar(&blacklistExt, "b", "", "Additional blacklist extensions (js,css)")
	flag.StringVar(&place, "i", "all", "Where to inject\n  all: replace all\n  one: replace one by one\n  2: replace the second path/param\n  -2: replace the second path/param from the end")
	flag.StringVar(&payloadFile, "f", "", "Payload list")
	flag.Parse()

	if !query && !path {
		fmt.Fprintln(os.Stderr, "Choose Query or Path")
		os.Exit(1)
	}

	if payloadFile != "" {
		pf, err := os.Open(payloadFile)
		if err != nil {
			panic(err)
		}
		defer pf.Close()
		scPayload := bufio.NewScanner(pf)
		for scPayload.Scan() {
			line := strings.TrimSpace(scPayload.Text())
			if line != "" {
				payloadList = append(payloadList, line)
			}
		}
	} else {
		payloadList = append(payloadList, flag.Arg(0))
	}

	if blacklistExt != "" {
		bl := strings.Split(blacklistExt, ",")
		for _, e := range bl {
			e = strings.TrimSpace(e)
			if e == "" {
				continue
			}
			OtherExt = append(OtherExt, e)
		}
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		u, err := url.Parse(strings.TrimSpace(sc.Text()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse url %s [%s]\n", sc.Text(), err)
			continue
		}

		if removeMediaExt {
			if BlacklistExt(u) {
				continue
			}
		}

		for _, payload := range payloadList {
			if query {
				replaceQuery(u, payload)
			}
			if path {
				replacePath(u, payload)
			}
		}

	}
}

func replaceQuery(u *url.URL, payload string) {
	pp := make([]string, 0)

	if len(u.Query()) == 0 {
		return
	}

	for p := range u.Query() {
		pp = append(pp, p)
	}
	sort.Strings(pp)

	if place == "all" {
		qs := url.Values{}
		for param, vv := range u.Query() {
			if appendMode {
				qs.Set(param, vv[0]+payload)
			} else {
				qs.Set(param, payload)
			}
		}
		u.RawQuery = qs.Encode()
		fmt.Printf("%s\n", u)
	} else if place == "one" {
		for i := 0; i < len(pp); i++ {
			cloneURL := &url.URL{}
			err := copier.Copy(cloneURL, u)
			if err != nil {
				panic(err)
			}
			qs := cloneURL.Query()

			if appendMode {
				qs.Set(pp[i], qs.Get(pp[i])+payload)
			} else {
				qs.Set(pp[i], payload)
			}
			cloneURL.RawQuery = qs.Encode()
			fmt.Printf("%s\n", cloneURL)
		}
	} else {
		qs := u.Query()
		var toReplacePlace int

		if strings.HasPrefix(place, "-") {
			p, err := strconv.Atoi(place[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert \"place\" string to int\n")
				p = 0
			}
			toReplacePlace = len(pp[:len(pp)-p])
		} else {
			p, err := strconv.Atoi(place)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert \"place\" string to int\n")
				p = 0
			}
			toReplacePlace = p
		}

		if toReplacePlace >= len(pp) {
			toReplacePlace = len(pp) - 1
		}

		if appendMode {
			qs.Set(pp[toReplacePlace], qs.Get(pp[toReplacePlace])+payload)
		} else {
			qs.Set(pp[toReplacePlace], payload)
		}
		u.RawQuery = qs.Encode()
		fmt.Printf("%s\n", u)
	}
}

func replacePath(u *url.URL, payload string) {
	path := strings.TrimPrefix(u.EscapedPath(), "/")
	paths := strings.Split(path, "/")

	if len(paths) == 0 {
		return
	}

	if place == "all" {
		for i := range paths {
			if appendMode {
				paths[i] = paths[i] + payload
			} else {
				paths[i] = payload
			}
		}
		u.Path = strings.Join(paths, "/")
		fmt.Printf("%s\n", u)
	} else if place == "one" {
		for i := 0; i < len(paths); i++ {
			cloneURL := &url.URL{}
			err := copier.Copy(cloneURL, u)
			if err != nil {
				panic(err)
			}
			pathClone := append(paths[:0:0], paths...)
			if appendMode {
				pathClone[i] = pathClone[i] + payload
			} else {
				pathClone[i] = payload
			}
			cloneURL.Path = strings.Join(pathClone, "/")
			fmt.Printf("%s\n", cloneURL)
		}
	} else {
		var toReplacePlace int
		if strings.HasPrefix(place, "-") {
			p, err := strconv.Atoi(place[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert \"place\" string to int\n")
				p = 0
			}
			toReplacePlace = len(paths[:len(paths)-p])
		} else {
			p, err := strconv.Atoi(place)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to convert \"place\" string to int\n")
				p = 0
			}
			toReplacePlace = p
		}
		if toReplacePlace >= len(paths) {
			toReplacePlace = len(paths) - 1
		}
		if appendMode {
			paths[toReplacePlace] = paths[toReplacePlace] + payload
		} else {
			paths[toReplacePlace] = payload
		}
		u.Path = strings.Join(paths, "/")
		fmt.Printf("%s\n", u)
	}
}

// Return true if in blacklist
func BlacklistExt(u *url.URL) bool {
	e := strings.TrimPrefix(path2.Ext(u.Path), ".")
	if inSlice(e, ImagesExt) || inSlice(e, AudioExt) || inSlice(e, FontExt) || inSlice(e, OtherExt) {
		return true
	}
	return false
}

func inSlice(s string, slice []string) bool {
	for _, e := range slice {
		if s == e {
			return true
		}
	}
	return false
}
