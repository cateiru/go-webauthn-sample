package src

import "github.com/go-webauthn/webauthn/webauthn"

type User struct {
	ID          []byte
	Name        string
	DisplayName string

	Credentials []webauthn.Credential
}

func (u *User) WebAuthnID() []byte {
	return u.ID
}

func (u *User) WebAuthnName() string {
	return u.Name
}

func (u *User) WebAuthnDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Name
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *User) WebAuthnIcon() string {
	return ""
}
