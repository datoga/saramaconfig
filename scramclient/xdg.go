package scramclient

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"

	"github.com/xdg/scram"
)

var SHA256 scram.HashGeneratorFcn = func() hash.Hash { return sha256.New() }
var SHA512 scram.HashGeneratorFcn = func() hash.Hash { return sha512.New() }

type XDG struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

//Begin creates the client.
func (x *XDG) Begin(userName, password, authzID string) (err error) {
	x.Client, err = x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.ClientConversation = x.Client.NewConversation()
	return nil
}

//Step is the adapter with the step for the client.
func (x *XDG) Step(challenge string) (response string, err error) {
	response, err = x.ClientConversation.Step(challenge)
	return
}

//Done ends the connection.
func (x *XDG) Done() bool {
	return x.ClientConversation.Done()
}
