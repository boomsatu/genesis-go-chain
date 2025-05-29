
package p2p

import (
	"fmt"
	"net"
	"sync"
	"time"

	"blockchain-node/config"
)

// Server represents the P2P server
type Server struct {
	config   *config.NetworkConfig
	listener net.Listener
	peers    map[string]*Peer
	mu       sync.RWMutex
	quit     chan struct{}
}

// Peer represents a connected peer
type Peer struct {
	conn     net.Conn
	address  string
	lastSeen time.Time
}

// NewServer creates a new P2P server
func NewServer(cfg *config.NetworkConfig) *Server {
	return &Server{
		config: cfg,
		peers:  make(map[string]*Peer),
		quit:   make(chan struct{}),
	}
}

// Start starts the P2P server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.ListenAddr, s.config.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to start P2P server: %v", err)
	}

	s.listener = listener
	fmt.Printf("P2P server listening on %s\n", addr)

	// Start accepting connections
	go s.acceptConnections()

	// Connect to seed nodes
	go s.connectToSeeds()

	return nil
}

// Stop stops the P2P server
func (s *Server) Stop() error {
	close(s.quit)
	
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all peer connections
	s.mu.Lock()
	for _, peer := range s.peers {
		peer.conn.Close()
	}
	s.mu.Unlock()

	return nil
}

// acceptConnections accepts incoming peer connections
func (s *Server) acceptConnections() {
	for {
		select {
		case <-s.quit:
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection: %v\n", err)
				continue
			}

			go s.handlePeer(conn)
		}
	}
}

// handlePeer handles a peer connection
func (s *Server) handlePeer(conn net.Conn) {
	defer conn.Close()

	peerAddr := conn.RemoteAddr().String()
	fmt.Printf("New peer connected: %s\n", peerAddr)

	peer := &Peer{
		conn:     conn,
		address:  peerAddr,
		lastSeen: time.Now(),
	}

	s.mu.Lock()
	s.peers[peerAddr] = peer
	s.mu.Unlock()

	// Handle peer messages
	s.handlePeerMessages(peer)

	// Clean up when done
	s.mu.Lock()
	delete(s.peers, peerAddr)
	s.mu.Unlock()

	fmt.Printf("Peer disconnected: %s\n", peerAddr)
}

// handlePeerMessages handles messages from a peer
func (s *Server) handlePeerMessages(peer *Peer) {
	buffer := make([]byte, 1024)
	
	for {
		n, err := peer.conn.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from peer %s: %v\n", peer.address, err)
			return
		}

		message := buffer[:n]
		s.processMessage(peer, message)
		peer.lastSeen = time.Now()
	}
}

// processMessage processes a message from a peer
func (s *Server) processMessage(peer *Peer, message []byte) {
	// Simple message handling - in production, implement proper protocol
	fmt.Printf("Received message from %s: %s\n", peer.address, string(message))
	
	// Echo message back (placeholder)
	peer.conn.Write([]byte("ACK"))
}

// connectToSeeds connects to seed nodes
func (s *Server) connectToSeeds() {
	for _, seedAddr := range s.config.SeedNodes {
		go s.connectToPeer(seedAddr)
	}
}

// connectToPeer connects to a specific peer
func (s *Server) connectToPeer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Failed to connect to peer %s: %v\n", address, err)
		return
	}

	fmt.Printf("Connected to peer: %s\n", address)
	s.handlePeer(conn)
}

// BroadcastMessage broadcasts a message to all connected peers
func (s *Server) BroadcastMessage(message []byte) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for addr, peer := range s.peers {
		if _, err := peer.conn.Write(message); err != nil {
			fmt.Printf("Failed to send message to peer %s: %v\n", addr, err)
		}
	}
}

// GetPeerCount returns the number of connected peers
func (s *Server) GetPeerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.peers)
}

// GetPeers returns a list of connected peer addresses
func (s *Server) GetPeers() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var peers []string
	for addr := range s.peers {
		peers = append(peers, addr)
	}
	return peers
}
