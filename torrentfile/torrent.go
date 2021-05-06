package torrentfile

import (
	"bytes"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func Parse(rawFile []byte) (*TorrentFile, error) {
	var err error
	bto := bencodeTorrent{}
	if err = bencode.Unmarshal(bytes.NewReader(rawFile), &bto); err != nil {
		return nil, err
	}

	return bto.toTorrentFile()
}
