package config

import "embed"

type Fs struct {
	fgroup map[string][]byte
	egroup map[string]*embed.FS
}

func NewFs() *Fs {
	return &Fs{
		fgroup: make(map[string][]byte),
		egroup: make(map[string]*embed.FS),
	}
}

func (fs *Fs) MountFile(filePath string, content []byte) {
	fs.fgroup[filePath] = content
}

func (fs *Fs) GetFile(filePath string) []byte {
	if v, found := fs.fgroup[filePath]; found {
		return v
	}
	return []byte{}
}

func (fs *Fs) MountEmbedFS(fsName string, efs *embed.FS) {
	fs.egroup[fsName] = efs
}

func (fs *Fs) GetEmbedFS(fsName string) *embed.FS {
	if v, found := fs.egroup[fsName]; found {
		return v
	}
	return nil
}
