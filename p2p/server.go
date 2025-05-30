
package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"blockchain-node/config"
	"blockchain-node/logger"
)

// MessageType represents the type of P2P message
type MessageType string

const (
	MessageTypeVersion     MessageType = "version"
	MessageTypeVerAck      MessageType = "verack"
	MessageTypeGetBlocks   MessageType = "getblocks"
	MessageTypeInv         MessageType = "inv"
	MessageTypeGetData     MessageType = "getdata"
	MessageTypeBlock       MessageType = "block"
	MessageTypeTx          MessageType = "tx"
	MessageTypePing        MessageType = "ping"
	MessageTypePong        MessageType = "pong"
	MessageTypeAddr        MessageType = "addr"
	MessageTypeGetAddr     MessageType = "getaddr"
)

// Message represents a P2P network message
type Message struct {
	Type      MessageType `json:"type"`
	Payload   []byte      `json:"payload"`
	Timestamp int64       `json:"timestamp"`
	Version   uint32      `json:"version"`
}

// Peer represents a connected peer
type Peer struct {
	ID         string
	Address    string
	Connection net.Conn
	Version    uint32
	Connected  time.Time
	LastSeen   time.Time
	Inbound    bool
	mu         sync.RWMutex
}

// Server represents the P2P server
type Server struct {
	config    *config.NetworkConfig
	peers     map[string]*Peer
	listener  net.Listener
	logger    *logger.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	mu        sync.RWMutex
	
	// Message handlers
	messageHandlers map[MessageType]func(*Peer, *Message) error
	
	// Callbacks
	onNewPeer    func(*Peer)
	onPeerLost   func(*Peer)
	onMessage    func(*Peer, *Message)
}

// NewServer creates a new P2P server
func NewServer(config *config.NetworkConfig) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	
	server := &Server{
		config:          config,
		peers:           make(map[string]*Peer),
		logger:          logger.NewLogger("p2p"),
		ctx:             ctx,
		cancel:          cancel,
		messageHandlers: make(map[MessageType]func(*Peer, *Message) error),
	}

	// Register default message handlers
	server.registerDefaultHandlers()

	return server
}

// Start starts the P2P server
func (s *Server) Start() error {
	s.logger.Info("Starting P2P server", "port", s.config.Port, "maxPeers", s.config.MaxPeers)

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.ListenAddr, s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to start P2P listener: %v", err)
	}
	s.listener = listener

	// Start accepting connections
	s.wg.Add(1)
	go s.acceptConnections()

	// Connect to seed nodes
	s.wg.Add(1)
	go s.connectToSeedNodes()

	// Start peer management
	s.wg.Add(1)
	go s.managePeers()

	s.logger.Info("P2P server started successfully")
	return nil
}

// Stop stops the P2P server
func (s *Server) Stop() error {
	s.logger.Info("Stopping P2P server...")

	s.cancel()

	// Close listener
	if s.listener != nil {
		s.listener.Close()
	}

	// Close all peer connections
	s.mu.Lock()
	for _, peer := range s.peers {
		peer.Connection.Close()
	}
	s.mu.Unlock()

	// Wait for all goroutines to finish
	s.wg.Wait()

	s.logger.Info("P2P server stopped")
	return nil
}

// acceptConnections accepts incoming peer connections
func (s *Server) acceptConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.ctx.Err() == nil {
					s.logger.Error("Failed to accept connection", "error", err)
				}
				continue
			}

			// Check peer limit
			if s.GetPeerCount() >= s.config.MaxPeers {
				s.logger.Warning("Rejecting connection, peer limit reached")
				conn.Close()
				continue
			}

			// Handle new peer in goroutine
			go s.handleNewPeer(conn, true)
		}
	}
}

// connectToSeedNodes connects to configured seed nodes
func (s *Server) connectToSeedNodes() {
	defer s.wg.Done()

	for _, seedNode := range s.config.SeedNodes {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.logger.Info("Connecting to seed node", "address", seedNode)
			
			conn, err := net.DialTimeout("tcp", seedNode, time.Duration(s.config.Timeout)*time.Second)
			if err != nil {
				s.logger.Warning("Failed to connect to seed node", "address", seedNode, "error", err)
				continue
			}

			go s.handleNewPeer(conn, false)
		}
	}
}

// handleNewPeer handles a new peer connection
func (s *Server) handleNewPeer(conn net.Conn, inbound bool) {
	peerAddr := conn.RemoteAddr().String()
	peerID := fmt.Sprintf("%s-%d", peerAddr, time.Now().UnixNano())

	peer := &Peer{
		ID:         peerID,
		Address:    peerAddr,
		Connection: conn,
		Version:    1,
		Connected:  time.Now(),
		LastSeen:   time.Now(),
		Inbound:    inbound,
	}

	s.logger.Info("New peer connection", "peerID", peerID, "address", peerAddr, "inbound", inbound)

	// Add peer to the list
	s.mu.Lock()
	s.peers[peerID] = peer
	s.mu.Unlock()

	// Notify new peer callback
	if s.onNewPeer != nil {
		s.onNewPeer(peer)
	}

	// Start peer communication
	if !inbound {
		// Send version message for outbound connections
		s.sendVersionMessage(peer)
	}

	// Handle peer messages
	s.handlePeerMessages(peer)
}

// handlePeerMessages handles messages from a peer
func (s *Server) handlePeerMessages(peer *Peer) {
	defer func() {
		// Clean up when peer disconnects
		s.mu.Lock()
		delete(s.peers, peer.ID)
		s.mu.Unlock()

		peer.Connection.Close()
		
		s.logger.Info("Peer disconnected", "peerID", peer.ID, "address", peer.Address)

		// Notify peer lost callback
		if s.onPeerLost != nil {
			s.onPeerLost(peer)
		}
	}()

	// Set connection timeout
	peer.Connection.SetReadDeadline(time.Now().Add(time.Duration(s.config.Timeout) * time.Second))

	decoder := json.NewDecoder(peer.Connection)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			var message Message
			if err := decoder.Decode(&message); err != nil {
				s.logger.Debug("Failed to decode message from peer", "peerID", peer.ID, "error", err)
				return
			}

			// Update last seen
			peer.mu.Lock()
			peer.LastSeen = time.Now()
			peer.mu.Unlock()

			// Reset read deadline
			peer.Connection.SetReadDeadline(time.Now().Add(time.Duration(s.config.Timeout) * time.Second))

			// Handle message
			if err := s.handleMessage(peer, &message); err != nil {
				s.logger.Warning("Failed to handle message", "peerID", peer.ID, "type", message.Type, "error", err)
			}

			// Notify message callback
			if s.onMessage != nil {
				s.onMessage(peer, &message)
			}
		}
	}
}

// handleMessage handles a specific message type
func (s *Server) handleMessage(peer *Peer, message *Message) error {
	handler, exists := s.messageHandlers[message.Type]
	if !exists {
		s.logger.Debug("No handler for message type", "type", message.Type, "peerID", peer.ID)
		return nil
	}

	return handler(peer, message)
}

// registerDefaultHandlers registers default message handlers
func (s *Server) registerDefaultHandlers() {
	s.messageHandlers[MessageTypeVersion] = s.handleVersionMessage
	s.messageHandlers[MessageTypeVerAck] = s.handleVerAckMessage
	s.messageHandlers[MessageTypePing] = s.handlePingMessage
	s.messageHandlers[MessageTypePong] = s.handlePongMessage
	s.messageHandlers[MessageTypeGetAddr] = s.handleGetAddrMessage
	s.messageHandlers[MessageTypeAddr] = s.handleAddrMessage
}

// Message handlers
func (s *Server) handleVersionMessage(peer *Peer, message *Message) error {
	s.logger.Debug("Received version message", "peerID", peer.ID)
	
	// Send verack response
	verackMsg := &Message{
		Type:      MessageTypeVerAck,
		Payload:   []byte{},
		Timestamp: time.Now().Unix(),
		Version:   1,
	}
	
	return s.sendMessage(peer, verackMsg)
}

func (s *Server) handleVerAckMessage(peer *Peer, message *Message) error {
	s.logger.Debug("Received verack message", "peerID", peer.ID)
	// Version handshake completed
	return nil
}

func (s *Server) handlePingMessage(peer *Peer, message *Message) error {
	s.logger.Debug("Received ping message", "peerID", peer.ID)
	
	// Send pong response
	pongMsg := &Message{
		Type:      MessageTypePong,
		Payload:   message.Payload, // Echo the payload
		Timestamp: time.Now().Unix(),
		Version:   1,
	}
	
	return s.sendMessage(peer, pongMsg)
}

func (s *Server) handlePongMessage(peer *Peer, message *Message) error {
	s.logger.Debug("Received pong message", "peerID", peer.ID)
	// Pong received, peer is alive
	return nil
}

func (s *Server) handleGetAddrMessage(peer *Peer, message *Message) error {
	s.logger.Debug("Received getaddr message", "peerID", peer.ID)
	
	// Send known peer addresses
	addresses := s.getKnownAddresses()
	addrPayload, _ := json.Marshal(addresses)
	
	addrMsg := &Message{
		Type:      MessageTypeAddr,
		Payload:   addrPayload,
		Timestamp: time.Now().Unix(),
		Version:   1,
	}
	
	return s.sendMessage(peer, addrMsg)
}

func (s *Server) handleAddrMessage(peer *Peer, message *Message) error {
	s.logger.Debug("Received addr message", "peerID", peer.ID)
	
	var addresses []string
	if err := json.Unmarshal(message.Payload, &addresses); err != nil {
		return fmt.Errorf("failed to unmarshal addresses: %v", err)
	}
	
	// Process received addresses (could connect to new peers)
	s.logger.Info("Received peer addresses", "count", len(addresses), "from", peer.ID)
	
	return nil
}

// sendVersionMessage sends a version message to a peer
func (s *Server) sendVersionMessage(peer *Peer) error {
	versionMsg := &Message{
		Type:      MessageTypeVersion,
		Payload:   []byte("lumina-node-v1.0"),
		Timestamp: time.Now().Unix(),
		Version:   1,
	}
	
	return s.sendMessage(peer, versionMsg)
}

// sendMessage sends a message to a peer
func (s *Server) sendMessage(peer *Peer, message *Message) error {
	peer.mu.Lock()
	defer peer.mu.Unlock()

	encoder := json.NewEncoder(peer.Connection)
	if err := encoder.Encode(message); err != nil {
		return fmt.Errorf("failed to send message to peer %s: %v", peer.ID, err)
	}

	s.logger.Debug("Sent message to peer", "type", message.Type, "peerID", peer.ID)
	return nil
}

// BroadcastMessage broadcasts a message to all connected peers
func (s *Server) BroadcastMessage(data []byte) {
	s.mu.RLock()
	peers := make([]*Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	s.mu.RUnlock()

	message := &Message{
		Type:      MessageTypeBlock, // Assuming it's a block broadcast
		Payload:   data,
		Timestamp: time.Now().Unix(),
		Version:   1,
	}

	for _, peer := range peers {
		if err := s.sendMessage(peer, message); err != nil {
			s.logger.Warning("Failed to broadcast message to peer", "peerID", peer.ID, "error", err)
		}
	}

	s.logger.Debug("Broadcasted message to peers", "peerCount", len(peers))
}

// managePeers manages peer connections and performs periodic maintenance
func (s *Server) managePeers() {
	defer s.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.performPeerMaintenance()
		}
	}
}

// performPeerMaintenance performs periodic peer maintenance
func (s *Server) performPeerMaintenance() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for peerID, peer := range s.peers {
		peer.mu.RLock()
		timeSinceLastSeen := now.Sub(peer.LastSeen)
		peer.mu.RUnlock()

		// Remove peers that haven't been seen for too long
		if timeSinceLastSeen > time.Duration(s.config.Timeout*2)*time.Second {
			s.logger.Info("Removing inactive peer", "peerID", peerID, "lastSeen", timeSinceLastSeen)
			peer.Connection.Close()
			delete(s.peers, peerID)
		}
	}

	s.logger.Debug("Peer maintenance completed", "activePeers", len(s.peers))
}

// GetPeerCount returns the number of connected peers
func (s *Server) GetPeerCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.peers)
}

// GetPeers returns a list of connected peers
func (s *Server) GetPeers() []*Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peers := make([]*Peer, 0, len(s.peers))
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	return peers
}

// getKnownAddresses returns known peer addresses
func (s *Server) getKnownAddresses() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	addresses := make([]string, 0, len(s.peers))
	for _, peer := range s.peers {
		addresses = append(addresses, peer.Address)
	}

	return addresses
}

// RegisterMessageHandler registers a handler for a specific message type
func (s *Server) RegisterMessageHandler(messageType MessageType, handler func(*Peer, *Message) error) {
	s.messageHandlers[messageType] = handler
}

// SetCallbacks sets callback functions for peer events
func (s *Server) SetCallbacks(onNewPeer, onPeerLost func(*Peer), onMessage func(*Peer, *Message)) {
	s.onNewPeer = onNewPeer
	s.onPeerLost = onPeerLost
	s.onMessage = onMessage
}

// SendToPeer sends a message to a specific peer
func (s *Server) SendToPeer(peerID string, messageType MessageType, payload []byte) error {
	s.mu.RLock()
	peer, exists := s.peers[peerID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("peer not found: %s", peerID)
	}

	message := &Message{
		Type:      messageType,
		Payload:   payload,
		Timestamp: time.Now().Unix(),
		Version:   1,
	}

	return s.sendMessage(peer, message)
}

// GetPeerInfo returns information about a specific peer
func (s *Server) GetPeerInfo(peerID string) *Peer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.peers[peerID]
}
