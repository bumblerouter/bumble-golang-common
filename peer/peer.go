package peer

import (
	"bumbleserver.org/common/key"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var publicKeys map[string]*rsa.PublicKey = make(map[string]*rsa.PublicKey)

var peerStringMatch *regexp.Regexp = regexp.MustCompile("^([A-Za-z0-9\\-]+)://([A-Za-z0-9\\-\\.]+)/([A-Za-z0-9\\-\\.]+)/([A-Za-z0-9\\-\\.]+)(/([A-Za-z0-9\\-\\.]+)(/(.*))?)?")
var peerUnquoteMatch *regexp.Regexp = regexp.MustCompile(`^"(.*)"$`)

type Peer struct {
	Scheme   string
	Domain   string
	Class    string
	User     string
	Instance string
	Extra    string
}

func (p Peer) Valid() bool {
	return p.Scheme == "bumble" && p.Domain != "" && p.Class != "" && p.User != ""
}

func (p Peer) String() string {
	if p.Scheme == "bumble" && p.Domain != "" && p.Class != "" && p.User != "" && p.Instance != "" && p.Extra != "" {
		return strings.ToLower(fmt.Sprintf("%s://%s", p.Scheme, strings.Join([]string{p.Domain, p.Class, p.User, p.Instance, p.Extra}, "/")))
	}
	return p.StringInstance()
}

func (p Peer) StringInstance() string {
	if p.Scheme == "bumble" && p.Domain != "" && p.Class != "" && p.User != "" && p.Instance != "" {
		return strings.ToLower(fmt.Sprintf("%s://%s", p.Scheme, strings.Join([]string{p.Domain, p.Class, p.User, p.Instance}, "/")))
	}
	return p.StringUser()
}

func (p Peer) StringUser() string {
	if p.Scheme == "bumble" && p.Domain != "" && p.Class != "" && p.User != "" {
		return strings.ToLower(fmt.Sprintf("%s://%s", p.Scheme, strings.Join([]string{p.Domain, p.Class, p.User}, "/")))
	}
	return ""
}

func (p *Peer) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, p.String())), nil
}

func (p *Peer) UnmarshalJSON(v []byte) error {
	matchedPeer := peerUnquoteMatch.FindStringSubmatch(string(v))
	if len(matchedPeer) < 2 {
		return errors.New("unable to match")
	}
	return p.SetFromString(matchedPeer[1])
}

func (p *Peer) SHA1() string {
	hash := sha1.New()
	hash.Write([]byte(p.StringUser()))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func NewFromString(s string) *Peer {
	peer := new(Peer)
	peer.SetFromString(s)
	return peer
}

func (p *Peer) SetFromString(s string) error {
	parts := peerStringMatch.FindStringSubmatch(s)
	if parts == nil {
		return errors.New("malformed or missing peer name")
	}
	p.Scheme = parts[1]
	p.Domain = parts[2]
	p.Class = parts[3]
	p.User = parts[4]
	p.Instance = parts[6]
	p.Extra = parts[8]
	return nil
}

func (p Peer) PublicKey() *rsa.PublicKey {
	if pubkey, ok := publicKeys[p.StringUser()]; ok && pubkey != nil {
		return pubkey
	}
	pubkey := key.GetPublicKey(p.Domain, p.SHA1())
	if pubkey != nil {
		publicKeys[p.StringUser()] = pubkey
	}
	return pubkey
}

func (p Peer) PublicKeyURL() (*url.URL, error) {
	return key.GetPublicKeyURL(p.Domain, p.SHA1())
}
