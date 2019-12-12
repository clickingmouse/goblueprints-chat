package main

import (
	"errors"
	"io/ioutil"
	"path"
)

//ErrNoAvatarURL is the error that is returned
//when the Avatar instance is unable to provide an Avatar UTL

var ErrNoAvatarURL = errors.New("chat: Unable to get an Avatar URL.")

//GetAvatarURL gets the avatar URL for specified client,
// or returns an error if something goes wrong.
//ErrNoAvatarURL is returned
//if the object is unable to get a URL for specified client

type Avatar interface {
	//GetAvatarURL(c *client) (string, error)
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

// can assign the UseAuthAvatar to any field looking for an Avatar interface
var UseAuthAvatar AuthAvatar

// tobe refactor for line of sight
//func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {

	// if url, ok := c.userData["avatar_url"]; ok {
	// 	if urlStr, ok := url.(string); ok {
	// 		return urlStr, nil
	// 	}
	// }
	//	return "", ErrNoAvatarURL

	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil

}

////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////
type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

//func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {

	// if email, ok := c.userData["email"]; ok {
	// 	if emailStr, ok := email.(string); ok {
	// 		m := md5.New()
	// 		io.WriteString(m, strings.ToLower(emailStr))
	// 		return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
	// 	}
	// }
	//+++
	//+ if userid, ok := c.userData["userid"]; ok {
	//+ 	if useridStr, ok := userid.(string); ok {
	//+ 		return "//www.gravatar.com/avatar/" + useridStr, nil
	//+ 	}
	//+ }
	//+ return "", ErrNoAvatarURL

	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil

}

///////////////////////////////////////////////////////////
//
//
//
///////////////////////////////////////////////////////////

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

//+func (FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {

	//+ if userid, ok := c.userData["userid"]; ok {
	//+ 	if useridStr, ok := userid.(string); ok {
	///////////////////////////////////////////
	//+ files, err := ioutil.ReadDir("avatars")
	//+ if err != nil {
	//+ 	return "", ErrNoAvatarURL
	//+ }

	if files, err := ioutil.ReadDir("Avatars"); err == nil {

		for _, file := range files {
			if file.IsDir() {
				continue
			}
			//+			if match, _ := path.Match(useridStr+"*", file.Name()); match
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
			// {
			// 	/////////////////////////////////////////////////////
			// 	return "/avatars/" + useridStr + ".jpg", nil
			// } //
		}
	}
	//+ 	}
	//+ }
	return "", ErrNoAvatarURL
}
