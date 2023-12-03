package torrent

import (
	"encoding/binary"
	"errors"
	"net"
)

type Peer struct {
	Ip   net.IP `json:"ip"`
	Port uint16 `json:"port"`
}

const peerByteSize = 6 // 4 for IP + 2 for port

func UnmarshalPeers(peersBinary []byte) ([]Peer, error) {
	if len(peersBinary)%peerByteSize != 0 {
		return nil, errors.New("malformed peers binary representation received")
	}
	peersCount := len(peersBinary) / peerByteSize

	peers := make([]Peer, peersCount)

	for i := 0; i < peersCount; i++ {
		offset := i * peerByteSize
		peers[i].Ip = net.IP(peersBinary[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16([]byte(peersBinary[offset+4 : offset+6]))
	}

	return peers, nil
}
