package main

import (
	"errors"
	"io/ioutil"
	"path"
)

// ErrNoAvatar 는 Avatar 인스턴스가 아바타 URL을 제공할 수 없을 때 리턴되는 에러다
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

// Avatar 는 사용자 프로필 사진을 표현할 수 있는 타입을 나타낸다
type Avatar interface {
	// GetAvatarURL 은 지정된 클라이언트에 대한 아바타 URL을 가져오고, 문제가 발생하면 에러를 리턴한다.
	// 객체가 지정된 클라이언트의 URL을 가져올 수 없는 경우 ErrNoAvatarURL이 리턴된다
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
var UseAuthAvatar AuthAvatar
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}
var UseGravatar GravatarAvatar
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

type FileSystemAvatar struct{}
var UseFileSystemAvatar FileSystemAvatar
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := path.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}