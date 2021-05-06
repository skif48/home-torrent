package torrentfile

import (
	"bytes"
	"net/url"
	"strconv"

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

func (torrentFile *TorrentFile) buildTrackerURL(peerID string, port uint16) (string, error) {
	var base *url.URL
	var err error

	if base, err = url.Parse(torrentFile.Announce); err != nil {
		return "", err
	}

	queryParams := url.Values{
		"info_hash":  []string{string(torrentFile.InfoHash[:])},
		"peer_id":    []string{peerID},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(torrentFile.Length)},
	}

	base.RawQuery = queryParams.Encode()
	return base.String(), nil
}
