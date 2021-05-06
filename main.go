package main

import (
	"fmt"
	"io/ioutil"

	"vladusenko.io/home-torrent/torrentfile"
)

func main() {
	var file []byte
	var torrentFile *torrentfile.TorrentFile
	var err error

	fmt.Println("Hello from home-torrent")
	if file, err = ioutil.ReadFile("./sample.torrent"); err != nil {
		panic(err)
	}

	if torrentFile, err = torrentfile.Parse(file); err != nil {
		panic(err)
	}

	fmt.Println(torrentFile)
}
