package client

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/ProtonMail/gluon/async"
	"github.com/ProtonMail/go-proton-api"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
)

type ProtonAddresses struct {
	User    proton.User
	Addrs   []proton.Address
	AddrKRs map[string]*crypto.KeyRing
}

// GetProtonAddresses Retrieves user information and decrypted addresses
func GetProtonAddresses(ctx context.Context, client *proton.Client, password string) (*ProtonAddresses, error) {
	user, err := client.GetUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the user: %w", err)
	}

	addrs, err := client.GetAddresses(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the addresses.: %w", err)
	}

	salts, err := client.GetSalts(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the salts: %w", err)
	}

	keyPass, err := salts.SaltForKey([]byte(password), user.Keys.Primary().ID)
	if err != nil {
		return nil, fmt.Errorf("unable to generate keypass: %w", err)
	}

	_, addrKRs, err := proton.Unlock(user, addrs, keyPass, async.NoopPanicHandler{})
	if err != nil {
		return nil, fmt.Errorf("unable to unlock address: %w", err)
	}

	return &ProtonAddresses{
		User:    user,
		Addrs:   addrs,
		AddrKRs: addrKRs,
	}, nil
}

// GetDecryptedMessage retrieves and deciphers a specific message
func GetDecryptedMessage(ctx context.Context, client *proton.Client, messageID string, addrKR *crypto.KeyRing) (string, error) {
	message, err := client.GetFullMessage(
		ctx,
		messageID,
		proton.NewParallelScheduler(runtime.NumCPU()/2, async.NoopPanicHandler{}),
		proton.NewDefaultAttachmentAllocator(),
	)
	if err != nil {
		return "", fmt.Errorf("impossible de récupérer le message: %w", err)
	}

	// Decipher
	plainBytes, err := message.Decrypt(addrKR)
	if err != nil {
		return "", fmt.Errorf("impossible de déchiffrer le message: %w", err)
	}

	// Construct a complete message
	fullMessage := fmt.Sprintf("Headers:\n%v\n\nContenu:\n%s", message.Header, string(plainBytes))
	return fullMessage, nil
}

// StreamMessages creates an event stream for new messages
func StreamMessages(ctx context.Context, client *proton.Client, period, jitter time.Duration) (<-chan string, error) {
	messageIDs := make(chan string)
	stream := client.NewEventStream(ctx, period, jitter, "latest")

	go func() {
		defer close(messageIDs)
		for e := range stream {
			for _, m := range e.Messages {
				fmt.Println(m.Action, m.ID, m.EventItem)
				if m.Action == proton.EventCreate {
					messageIDs <- m.ID
				}
			}
		}
	}()

	return messageIDs, nil
}
