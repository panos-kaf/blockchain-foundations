package peer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"marabu/internal/logs"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	BOOTSTRAP_PEERS = Peers{
		"95.179.158.137:18018",
		"95.179.132.22:18018",
		"45.32.235.245:18018",
	}
	PEERS_FILE      = filepath.Join(".", "db", "peers.csv")
	knownPeers      = make(map[BuPeer]string)
	knownPeersMutex sync.Mutex
)

func init() {
	loadPeers()
	if _, err := os.Stat(PEERS_FILE); errors.Is(err, os.ErrNotExist) {
		savePeers()
	}
}

// Load peers from file and bootstrap list
func loadPeers() {
	knownPeersMutex.Lock()
	for _, peer := range BOOTSTRAP_PEERS {
		knownPeers[peer] = "bootstrap"
	}
	knownPeersMutex.Unlock()
	file, err := os.Open(PEERS_FILE)
	if err != nil {
		return
	}
	defer file.Close()
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return
	}
	for _, rec := range records {
		if len(rec) < 2 || rec[0] == "Address" {
			continue
		}
		knownPeers[BuPeer(rec[0])] = rec[1]
	}
}

// Save peers to file
func savePeers() {
	knownPeersMutex.Lock()
	defer knownPeersMutex.Unlock()
	file, err := os.Create(PEERS_FILE)
	if err != nil {
		logs.GlobalLog(fmt.Sprintf("Failed to save peers file: %v", err))
		return
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()
	w.Write([]string{"Address", "Source"})
	for peer, source := range knownPeers {
		w.Write([]string{string(peer), string(source)})
	}
}

// Get all known peers
func GetKnownPeers() Peers {
	knownPeersMutex.Lock()
	defer knownPeersMutex.Unlock()
	keys := make(Peers, 0, len(knownPeers))
	for k := range knownPeers {
		keys = append(keys, k)
	}
	return keys
}

// Validate and sanitize peer address
func sanitizePeer(peer BuPeer) (BuPeer, bool) {

	peerStr := strings.TrimSpace(string(peer))
	ipv4 := regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?):([0-9]{1,5})$`)
	ipv6 := regexp.MustCompile(`^\[([a-fA-F0-9:]+)\]:([0-9]{1,5})$`)
	domain := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}:([0-9]{1,5})$`)

	isIPv4 := ipv4.MatchString(peerStr)
	isIPv6 := ipv6.MatchString(peerStr)
	isDomain := domain.MatchString(peerStr)

	if !isIPv4 && !isIPv6 && !isDomain {
		return "", false
	}

	lastColon := strings.LastIndex(peerStr, ":")
	if lastColon == -1 {
		return "", false
	}
	portStr := peerStr[lastColon+1:]
	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 || port > 65535 {
		return "", false
	}
	host := peerStr[:lastColon]

	if isIPv6 {
		if host == "::1" || host == "[::1]" ||
			strings.HasPrefix(host, "[fe80:") || strings.HasPrefix(host, "[fc00:") {
			return "", false
		}
	}
	if host == "localhost" {
		return "", false
	}
	if isIPv4 {
		if strings.HasPrefix(host, "127.") || strings.HasPrefix(host, "0.") ||
			strings.HasPrefix(host, "192.168.") || strings.HasPrefix(host, "10.") {
			return "", false
		}
		octets := strings.Split(host, ".")
		if len(octets) != 4 {
			return "", false
		}
		for _, octet := range octets {
			num, err := strconv.Atoi(octet)
			if err != nil || num < 0 || num > 255 {
				return "", false
			}
		}
		if strings.HasPrefix(host, "172.") {
			second, err := strconv.Atoi(octets[1])
			if err != nil || second < 16 || second > 31 {
				return "", false
			}
		}
	}
	return BuPeer(peerStr), true
}

// Add new peers
func AppendPeers(peers Peers, server string) {
	knownPeersMutex.Lock()
	defer knownPeersMutex.Unlock()
	changed := false
	for _, peer := range peers {
		if sanitized, ok := sanitizePeer(peer); ok {
			if _, exists := knownPeers[sanitized]; !exists {
				knownPeers[sanitized] = server
				logs.GlobalLog(fmt.Sprintf("Added new peer: %s from server %s", sanitized, server))
				changed = true
			}
		}
	}
	if changed {
		logs.GlobalLog(fmt.Sprintf("Saving %d peers to disk...", len(knownPeers)))
		savePeers()
	}
}

// Select random peers per source
func SelectRandomPeersPerSource(count int) []string {
	knownPeersMutex.Lock()
	defer knownPeersMutex.Unlock()

	peersBySource := make(map[string][]string)
	for peer, source := range knownPeers {
		peersBySource[source] = append(peersBySource[source], string(peer))
	}
	selected := []string{}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for _, peers := range peersBySource {
		if len(peers) <= count {
			selected = append(selected, peers...)
		} else {
			perm := rng.Perm(len(peers))
			for i := range count {
				selected = append(selected, peers[perm[i]])
			}
		}
	}
	return selected
}
