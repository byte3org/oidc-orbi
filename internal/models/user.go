package models

import (
	"fmt"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ModelBase
	Username      string
	Password      string
	FirstName     string
	LastName      string
	Email         string
	EmailVerified bool
	Phone         string
	PhoneVerified bool
	IsAdmin       bool
	credentials   []webauthn.Credential
}

func (u *User) ValidatePassword(input string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input)); err != nil {
		return false
	}

	return true
}

func (u User) WebAuthnID() []byte {
	return []byte(u.ID.String())
}

// WebAuthnName returns the user's username
func (u User) WebAuthnName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// WebAuthnDisplayName returns the user's display name
func (u User) WebAuthnDisplayName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

// WebAuthnIcon is not (yet) implemented
func (u User) WebAuthnIcon() string {
	return ""
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) {
	u.credentials = append(u.credentials, cred)
}

// WebAuthnCredentials returns credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

// CredentialExcludeList returns a CredentialDescriptor array filled
// with all the user's credentials
func (u User) CredentialExcludeList() []protocol.CredentialDescriptor {

	credentialExcludeList := []protocol.CredentialDescriptor{}
	for _, cred := range u.credentials {
		descriptor := protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
		credentialExcludeList = append(credentialExcludeList, descriptor)
	}

	return credentialExcludeList
}
