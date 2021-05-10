package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"vladusenko.io/home-torrent/torrentfile"
)

func main() {
	var file []byte
	var torrentFile *torrentfile.Torrent
	var err error

	fmt.Println("Hello from home-torrent")
	if file, err = ioutil.ReadFile("./single-file.torrent"); err != nil {
		panic(err)
	}

	if torrentFile, err = torrentfile.Parse(file); err != nil {
		panic(err)
	}

	jsonTorrent, _ := json.Marshal(torrentFile)

	ioutil.WriteFile("./single-file.json", jsonTorrent, 0644)
}
