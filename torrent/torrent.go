package torrent

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Path   []string `json:"path"`
	Length int      `json:"length"`
}

type Torrent struct {
	Announce    string         `json:"announce"`
	InfoHash    [20]byte       `json:"info_hash"`
	PieceHashes [][20]byte     `json:"piece_hashes"`
	PieceLength int            `json:"piece_length"`
	Files       []*TorrentFile `json:"files"`
}

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func Parse(rawFile []byte) (*Torrent, error) {
	var err error
	bto := bencodeTorrent{}
	if err = bencode.Unmarshal(bytes.NewReader(rawFile), &bto); err != nil {
		return nil, err
	}

	return bto.toTorrentFile()
}

func (torrent *Torrent) TotalLength() int {
	var length int

	for _, file := range torrent.Files {
		length += file.Length
	}

	return length
}

func (torrent *Torrent) buildHttpTrackerUrl(peerId [20]byte, port uint16) (string, error) {
	base, err := url.Parse(torrent.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(torrent.InfoHash[:])},
		"peer_id":    []string{string(peerId[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(torrent.TotalLength())},
	}

	base.RawQuery = base.RawQuery + params.Encode()

	return base.String(), nil
}

func (torrent *Torrent) RequestPeers(peerId [20]byte, port uint16) ([]Peer, error) {
	url, err := torrent.buildHttpTrackerUrl(peerId, port)
	if err != nil {
		return nil, err
	}

	fmt.Println(url)

	httpClient := &http.Client{Timeout: 15 * time.Second}
	response, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	trackerResponse := TrackerResponse{}
	if err = bencode.Unmarshal(response.Body, &trackerResponse); err != nil {
		return nil, err
	}

	return UnmarshalPeer([]byte(trackerResponse.Peers))
}
