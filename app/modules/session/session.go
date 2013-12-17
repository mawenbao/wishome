package session

type Session interface {
    Encrypt() bool
    Decrypt() bool
    IsEncrypted() bool
    IsExpired() bool
    // load session from cookie and decrypt
    Load(map[string]string) bool
    // encrypt and save session in cookie
    Save(map[string]string) bool
}

