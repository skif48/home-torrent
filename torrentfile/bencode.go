package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/jackpal/bencode-go"
)

type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (bencodeTorrent *bencodeTorrent) infoHash() ([sha1.Size]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, bencodeTorrent.Info); err != nil {
		return [sha1.Size]byte{}, err
	}
	hash := sha1.Sum(buf.Bytes())
	return hash, nil
}

func (bencodeTorrent *bencodeTorrent) decodePieceHashes() ([][sha1.Size]byte, error) {
	buf := []byte(bencodeTorrent.Info.Pieces)

	if len(buf)%sha1.Size != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}

	hashesAmount := len(buf) / sha1.Size
	hashes := make([][20]byte, hashesAmount)

	for i := 0; i < hashesAmount; i++ {
		copy(hashes[i][:], buf[i*sha1.Size:(i+1)*sha1.Size])
	}
	return hashes, nil
}

func (bencodeTorrent *bencodeTorrent) toTorrentFile() (*TorrentFile, error) {
	var err error
	var infoHash [sha1.Size]byte
	var pieceHashes [][sha1.Size]byte

	if infoHash, err = bencodeTorrent.infoHash(); err != nil {
		return nil, err
	}

	if pieceHashes, err = bencodeTorrent.decodePieceHashes(); err != nil {
		return nil, err
	}

	torrentFile := TorrentFile{
		Announce:    bencodeTorrent.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bencodeTorrent.Info.PieceLength,
		Length:      bencodeTorrent.Info.Length,
		Name:        bencodeTorrent.Info.Name,
	}

	return &torrentFile, nil
}
