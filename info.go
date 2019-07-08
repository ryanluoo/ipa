package ipa

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

//Info 储存一个 ipa 包中的 info.plist 信息
type Info struct {
	originalInfo map[string]interface{}
}

// OriginalInfo 返回 info.plist 的原始信息
func (info Info) OriginalInfo() map[string]interface{} {
	return info.originalInfo
}

//DisplayName 一个 iOS App 的显示名称
func (info Info) DisplayName() string {
	return fmt.Sprintf("%v", info.originalInfo["CFBundleDisplayName"])
}

//BundleID iOS App 的 bundle id
func (info Info) BundleID() string {
	return fmt.Sprintf("%v", info.originalInfo["CFBundleIdentifier"])
}

//BuildVersion iOS App 的BuildVersion
func (info Info) BuildVersion() string {
	return fmt.Sprintf("%v", info.originalInfo["CFBundleVersion"])
}

//Version iOS App 的版本号
func (info Info) Version() string {
	return fmt.Sprintf("%v", info.originalInfo["CFBundleShortVersionString"])
}

// URLScheme 查找相关 url scheme
func (info Info) URLScheme(schemeName string) (urlScheme string) {
	urlSchemes, succeed := info.originalInfo["CFBundleURLTypes"].([]interface{})
	if !succeed {
		return
	}

	for _, item := range urlSchemes {
		urlSchemeMap, succeed := item.(map[string]interface{})
		if !succeed {
			continue
		}

		name, succeed := urlSchemeMap["CFBundleURLName"].(string)
		if !succeed {
			continue
		}
		fmt.Println("urlScheme name " + name)

		if name == schemeName {
			urls, succeed := urlSchemeMap["CFBundleURLSchemes"].([]interface{})
			if !succeed {
				return
			}

			if len(urls) == 0 {
				return
			}

			urlScheme, succeed = urls[0].(string)
			if !succeed {
				return
			}

			break
		}
	}

	return
}

func newInfo(ipaPath string) (*Info, error) {
	info, err := findPlist(ipaPath)

	if err != nil {
		return nil, err
	}

	return &Info{originalInfo: info}, nil
}

func findPlist(ipaPath string) (map[string]interface{}, error) {
	reader, err := zip.OpenReader(ipaPath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	for _, file := range reader.File {
		if !strings.HasSuffix(file.Name, ".app/Info.plist") {
			continue
		}

		if len(strings.Split(file.Name, string(filepath.Separator))) != 3 {
			continue
		}

		reader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		plistBytes, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Println("plist read error:", err)
			return nil, err
		}

		var result map[string]interface{}
		if _, err := plist.Unmarshal(plistBytes, &result); err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, errors.New("未找到 Info.plist 文件")
}
