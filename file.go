package storage_default

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/infrago/storage"
	"github.com/infrago/util"
)

//-------------------- defaultBase begin -------------------------

var (
	errBrowseNotSupported = errors.New("Store browse not supported.")
)

type (
	defaultDriver  struct{}
	defaultConnect struct {
		mutex  sync.RWMutex
		health storage.Health

		instance *storage.Instance

		setting defaultSetting
	}
	defaultSetting struct {
		Storage string
	}
)

// 连接
func (driver *defaultDriver) Connect(instance *storage.Instance) (storage.Connect, error) {
	setting := defaultSetting{
		Storage: "store/storage",
	}

	if vv, ok := instance.Setting["storage"].(string); ok {
		setting.Storage = vv
	}

	return &defaultConnect{
		instance: instance, setting: setting,
	}, nil

}

// 打开连接
func (this *defaultConnect) Open() error {
	return nil
}

func (this *defaultConnect) Health() storage.Health {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.health
}

// 关闭连接
func (this *defaultConnect) Close() error {
	return nil
}

func (this *defaultConnect) Upload(orginal string, opt storage.UploadOption) (string, error) {
	stat, err := os.Stat(orginal)
	if err != nil {
		return "", err
	}

	//250327不再支持目录上传
	if stat.IsDir() {
		return "", errors.New("directory upload not supported")
	}

	ext := util.Extension(orginal)

	if opt.Key == "" {
		//如果没有指定key，使用文件的hash
		//使用hash的前4位，生成2级目录
		hash, hex := this.filehash(orginal)
		if opt.Prefix == "" {
			opt.Prefix = path.Join(hex[0:2], hex[2:4])
		} else {
			opt.Prefix = path.Join(opt.Prefix, hex[0:2], hex[2:4])
		}
		opt.Key = hash
	}

	file := this.instance.File(opt.Prefix, opt.Key, ext, stat.Size())
	if file == nil {
		return "", errors.New("create file error")
	}

	//
	_, sFile, err := this.filepath(file)
	if err != nil {
		return "", err
	}

	//如果文件已经存在，直接返回
	//250327 更新，文件存在也覆盖
	// if _, err := os.Stat(sFile); err == nil {
	// 	return nil
	// }

	//打开原始文件
	fff, err := os.Open(orginal)
	if err != nil {
		return "", err
	}
	defer fff.Close()

	//创建文件
	save, err := os.OpenFile(sFile, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}
	defer save.Close()

	//复制文件
	_, err = io.Copy(save, fff)
	if err != nil {
		return "", err
	}

	return file.Code(), nil
}

func (this *defaultConnect) Fetch(file storage.File, opt storage.FetchOption) (storage.Stream, error) {
	///直接返回本地文件存储
	_, sFile, err := this.filepath(file)
	if err != nil {
		return nil, err
	}

	//打开文件
	fff, err := os.Open(sFile)
	if err != nil {
		return nil, err
	}

	return fff, nil
}

func (this *defaultConnect) Download(file storage.File, opt storage.DownloadOption) (string, error) {
	///直接返回本地文件存储
	_, sFile, err := this.filepath(file)
	if err != nil {
		return "", err
	}
	return sFile, nil
}

func (this *defaultConnect) Remove(file storage.File, opt storage.RemoveOption) error {
	_, sFile, err := this.filepath(file)
	if err != nil {
		return err
	}

	return os.Remove(sFile)
}

func (this *defaultConnect) Browse(file storage.File, opt storage.BrowseOption) (string, error) {
	return "", errBrowseNotSupported
}

//-------------------- defaultBase end -------------------------

// filepath 生成存储路径
func (this *defaultConnect) filepath(file storage.File) (string, string, error) {
	//使用hash的hex hash 的前4位，生成2级目录
	//共256*256个目录

	name := file.Key()
	if file.Type() != "" {
		name = fmt.Sprintf("%s.%s", file.Key(), file.Type())
	}

	sfile := path.Join(this.setting.Storage, file.Prefix(), name)
	spath := path.Dir(sfile)

	// //创建目录
	err := os.MkdirAll(spath, 0777)
	if err != nil {
		return "", "", errors.New("生成目录失败")
	}

	return spath, sfile, nil
}

// 算文件的hash
func (this *defaultConnect) filehash(file string) (string, string) {
	if f, e := os.Open(file); e == nil {
		defer f.Close()
		h := sha1.New()
		if _, e := io.Copy(h, f); e == nil {
			hex := fmt.Sprintf("%x", h.Sum(nil))
			hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
			return hash, hex
		}
	}
	return "", ""
}
