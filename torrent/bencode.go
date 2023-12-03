package torrent

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
	Length      int           `bencode:"length,omitempty"`
	Files       []bencodeFile `bencode:"files,omitempty"`
	Name        string        `bencode:"name"`
}

type bencodeTorrent struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func (b *bencodeTorrent) infoHash() ([sha1.Size]byte, error) {
	var buf bytes.Buffer
	var hash [sha1.Size]byte

	if err := bencode.Marshal(&buf, b.Info); err != nil {
		return [sha1.Size]byte{}, err
	}
	var hasher = sha1.New()
	hasher.Write(buf.Bytes())
	copy(hash[:], hasher.Sum(nil))

	return hash, nil
}

func (b *bencodeTorrent) decodePieceHashes() ([][sha1.Size]byte, error) {
	buf := []byte(b.Info.Pieces)

	if len(buf)%sha1.Size != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}

	hashesAmount := len(buf) / sha1.Size
	hashes := make([][sha1.Size]byte, hashesAmount)

	for i := 0; i < hashesAmount; i++ {
		copy(hashes[i][:], buf[i*sha1.Size:(i+1)*sha1.Size])
	}
	return hashes, nil
}

func (b *bencodeTorrent) toTorrentFile() (*Torrent, error) {
	var err error
	var infoHash [sha1.Size]byte
	var pieceHashes [][sha1.Size]byte

	if infoHash, err = b.infoHash(); err != nil {
		return nil, err
	}

	if pieceHashes, err = b.decodePieceHashes(); err != nil {
		return nil, err
	}

	if b.Info.Files != nil && b.Info.Length != 0 {
		return nil, errors.New("Found both `files` and `length` keys")
	}

	torrent := &Torrent{
		Announce:    b.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: b.Info.PieceLength,
	}

	if b.Info.Length != 0 {
		file := &TorrentFile{
			Path:   []string{b.Info.Name},
			Length: b.Info.Length,
		}

		torrent.Files = []*TorrentFile{file}
	} else {
		files := make([]*TorrentFile, len(b.Info.Files))
		for i, file := range b.Info.Files {
			files[i] = &TorrentFile{
				Length: file.Length,
				Path:   file.Path,
			}
		}
		torrent.Files = files
	}

	return torrent, nil
}
