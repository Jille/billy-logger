// Package logger logs all Billy calls made.
package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/helper/polyfill"
)

func Wrap(underlying billy.Basic, log func(msg ...interface{})) billy.Filesystem {
	return logger{
		underlying: polyfill.New(underlying),
		log:        log,
	}
}

var _ billy.Filesystem = logger{}
var _ billy.Change = logger{}

type logger struct {
	underlying billy.Filesystem
	log        func(msg ...interface{})
}

func (l logger) logf(f string, args ...interface{}) {
	l.log(fmt.Sprintf(f, args...))
}

func (l logger) Capabilities() billy.Capability {
	ret := billy.Capabilities(l.underlying)
	l.logf("Capabilities(): %d", ret)
	return ret
}

func (l logger) Join(elem ...string) string {
	ret := l.underlying.Join(elem...)
	l.logf("Join(%v): %q", elem, ret)
	return ret
}

func (l logger) Create(p string) (billy.File, error) {
	ret, err := l.underlying.Create(p)
	l.logf("Create(%q): %v", p, err)
	return ret, err
}

func (l logger) Open(p string) (billy.File, error) {
	ret, err := l.underlying.Open(p)
	l.logf("Open(%q): %v", p, err)
	return ret, err
}

func (l logger) OpenFile(p string, flag int, mode os.FileMode) (billy.File, error) {
	ret, err := l.underlying.OpenFile(p, flag, mode)
	l.logf("OpenFile(%q, %d, %#o): %v", p, flag, mode, err)
	return ret, err
}

func (l logger) Stat(p string) (os.FileInfo, error) {
	ret, err := l.underlying.Stat(p)
	if err != nil {
		l.logf("Stat(%q): -, %v", p, err)
		return ret, err
	}
	l.logf("Stat(%q): %s, %v", p, formatFileInfo(ret), err)
	return ret, err
}

func formatFileInfo(fi os.FileInfo) string {
	return fmt.Sprintf("{name: %q, mode: %#o, dir: %v}", fi.Name(), fi.Mode(), fi.IsDir())
}

func (l logger) Rename(from, to string) error {
	err := l.underlying.Rename(from, to)
	l.logf("Rename(%q, %q): %v", from, to, err)
	return err
}

func (l logger) Remove(p string) error {
	err := l.underlying.Remove(p)
	l.logf("Remove(%q): %v", p, err)
	return err
}

func (l logger) ReadDir(p string) ([]os.FileInfo, error) {
	ret, err := l.underlying.ReadDir(p)
	if err != nil {
		l.logf("ReadDir(%q): -, %v", p, err)
		return nil, err
	}
	descrs := make([]string, len(ret))
	for i, fi := range ret {
		descrs[i] = formatFileInfo(fi)
	}
	l.logf("ReadDir(%q): [%s], %v", p, strings.Join(descrs, ", "), err)
	return ret, nil
}

func (l logger) MkdirAll(p string, perm os.FileMode) error {
	err := l.underlying.MkdirAll(p, perm)
	l.logf("MkdirAll(%q, %#o): %v", p, perm, err)
	return err
}

func (l logger) Symlink(target, link string) error {
	err := billy.ErrNotSupported
	if sfs, ok := l.underlying.(billy.Symlink); ok {
		err = sfs.Symlink(target, link)
	}
	l.logf("Symlink(%q, %q): %v", target, link, err)
	return err
}

func (l logger) Readlink(p string) (string, error) {
	var ret string
	var err error
	if sfs, ok := l.underlying.(billy.Symlink); ok {
		ret, err = sfs.Readlink(p)
	} else {
		err = billy.ErrNotSupported
	}
	l.logf("Readlink(%q): %q, %v", p, ret, err)
	return ret, err
}

func (l logger) Lstat(p string) (os.FileInfo, error) {
	var ret os.FileInfo
	var err error
	if sfs, ok := l.underlying.(billy.Symlink); ok {
		ret, err = sfs.Lstat(p)
	} else {
		err = billy.ErrNotSupported
	}
	if err != nil {
		l.logf("Lstat(%q): -, %v", p, err)
		return ret, err
	}
	l.logf("Lstat(%q): %s, %v", p, formatFileInfo(ret), err)
	return ret, err
}

func (l logger) Chmod(p string, mode os.FileMode) error {
	err := billy.ErrNotSupported
	if cfs, ok := l.underlying.(billy.Change); ok {
		err = cfs.Chmod(p, mode)
	}
	l.logf("Chmod(%q, %#o): %v", p, mode, err)
	return err
}

func (l logger) Chown(p string, uid, gid int) error {
	err := billy.ErrNotSupported
	if cfs, ok := l.underlying.(billy.Change); ok {
		err = cfs.Chown(p, uid, gid)
	}
	l.logf("Chown(%q, %d, %d): %v", p, uid, gid, err)
	return err
}

func (l logger) Lchown(p string, uid, gid int) error {
	err := billy.ErrNotSupported
	if cfs, ok := l.underlying.(billy.Change); ok {
		err = cfs.Lchown(p, uid, gid)
	}
	l.logf("Lchown(%q, %d, %d): %v", p, uid, gid, err)
	return err
}

func (l logger) Chtimes(p string, atime, mtime time.Time) error {
	err := billy.ErrNotSupported
	if cfs, ok := l.underlying.(billy.Change); ok {
		err = cfs.Chtimes(p, atime, mtime)
	}
	l.logf("Chtimes(%q, %s, %s): %v", p, atime, mtime, err)
	return err
}

func (l logger) Chroot(p string) (billy.Filesystem, error) {
	ret, err := l.underlying.Chroot(p)
	l.logf("Chroot(%q): %v", p, err)
	return ret, err
}

func (l logger) TempFile(dir, prefix string) (billy.File, error) {
	ret, err := l.underlying.TempFile(dir, prefix)
	l.logf("TempFile(%q, %q): %v", dir, prefix, err)
	return ret, err
}

func (l logger) Root() string {
	ret := l.underlying.Root()
	l.logf("Root(): %q", ret)
	return ret
}
