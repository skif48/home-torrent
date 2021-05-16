package torrentfile

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"

	"github.com/jackpal/bencode-go"
)

type bencodeFile struct {
	Path   []string `bencode:"path"`
	Length int      `bencode:"length"`
}

type bencodeInfo struct {
	Pieces      string        `bencode:"pieces"`
	PieceLength int           `bencode:"piece length"`
	Length      int           `bencode:"length"`
	Files       []bencodeFile `bencode:"files"`
	Name        string        `bencode:"name"`
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

func (bencodeTorrent *bencodeTorrent) toTorrentFile() (*Torrent, error) {
	var err error
	var infoHash [sha1.Size]byte
	var pieceHashes [][sha1.Size]byte

	if infoHash, err = bencodeTorrent.infoHash(); err != nil {
		return nil, err
	}

	if pieceHashes, err = bencodeTorrent.decodePieceHashes(); err != nil {
		return nil, err
	}

	if bencodeTorrent.Info.Files != nil && bencodeTorrent.Info.Length != 0 {
		return nil, errors.New("Found both `files` and `length` keys")
	}

	torrent := &Torrent{
		Announce:    bencodeTorrent.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bencodeTorrent.Info.PieceLength,
	}

	if bencodeTorrent.Info.Length != 0 {
		file := &TorrentFile{
			Path:   []string{bencodeTorrent.Info.Name},
			Length: bencodeTorrent.Info.Length,
		}

		torrent.Files = []*TorrentFile{file}
	} else {
		files := make([]*TorrentFile, len(bencodeTorrent.Info.Files))
		for i, file := range bencodeTorrent.Info.Files {
			files[i] = &TorrentFile{
				Length: file.Length,
				Path:   file.Path,
			}
		}
		torrent.Files = files
	}

	return torrent, nil
}
