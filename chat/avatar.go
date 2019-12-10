package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

//ErrNoAvatarURL is the error that is returned
//when the Avatar instance is unable to provide an Avatar UTL

var ErrNoAvatarURL = errors.New("chat: Unable to get an Avatar URL.")

//GetAvatarURL gets the avatar URL for specified client,
// or returns an error if something goes wrong.
//ErrNoAvatarURL is returned
//if the object is unable to get a URL for specified client

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

// can assign the UseAuthAvatar to any field looking for an Avatar interface
var UseAuthAvatar AuthAvatar

// tobe refactor for line of sight
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

////////////////////////////////////////////////////////////////////////
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
