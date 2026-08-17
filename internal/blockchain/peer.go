package blockchain

import (
	"errors"
	"net/url"
	"strings"
)

// PeerSet stores known peer node URLs.
type PeerSet struct {
	peers map[string]struct{}
}

// NewPeerSet creates an empty peer registry.
func NewPeerSet() *PeerSet {
	return &PeerSet{
		peers: make(map[string]struct{}),
	}
}

// Add validates and stores a peer URL.
func (p *PeerSet) Add(peer string) error {
	peer = strings.TrimSpace(peer)

	if peer == "" {
		return errors.New("peer cannot be empty")
	}

	parsed, err := url.Parse(peer)
	if err != nil {
		return err
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("peer must use http or https")
	}

	if parsed.Host == "" {
		return errors.New("peer host is required")
	}

	p.peers[peer] = struct{}{}

	return nil
}

// List returns all registered peers.
func (p *PeerSet) List() []string {
	result := make([]string, 0, len(p.peers))

	for peer := range p.peers {
		result = append(
			result,
			peer,
		)
	}

	return result
}
