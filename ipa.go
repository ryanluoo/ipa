package ipa

// Ipa 用于保存 Ipa 包基本信息
type Ipa struct {
	ipaInfo Info
	path    string
}

// NewIpa 返回一个新 ipa 包，失败返回 nil
func NewIpa(path string) (*Ipa, error) {
	ipaInfo, err := newInfo(path)
	if err != nil {
		return nil, err
	}
	ipa := Ipa{ipaInfo: *ipaInfo, path: path}
	return &ipa, nil
}

// Info 返回 ipa 包的 info.plist 文件
func (ipa Ipa) Info() Info {
	return ipa.ipaInfo
}

// Path 返回 ipa 包的文件路径
func (ipa Ipa) Path() string {
	return ipa.path
}
