package torrent

import (
	"bytes"

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

func Parse(rawFile []byte) (*Torrent, error) {
	var err error
	bto := bencodeTorrent{}
	if err = bencode.Unmarshal(bytes.NewReader(rawFile), &bto); err != nil {
		return nil, err
	}

	return bto.toTorrentFile()
}
