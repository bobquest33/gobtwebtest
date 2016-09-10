package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/static/.DS_Store": {
		local:   "static/.DS_Store",
		size:    8196,
		modtime: 1473528545,
		compressed: `
H4sIAAAJbogA/+yZT2/TTBDGZxO/L05ycQVIkbj4CFJUtVFDm1saUiQOlaIGtUIUBTs2jZGxo3jTUEWR
cubKpWf+fw8EJz4AfBROMPZOguO4RUhIBbI/a/Vs7Rlbnce71m4AgNUH1jpAEbsqCNUKkIpKbYEM6Up4
P8wHcOAIVqEP5h03/V4SiUQikUguGCZEPeO7L5FIlphwftBJa6QToYyuZ0iVWI5GqpPWSCdCGcVlSBVS
lVQj1UlrpBOhNGkxWnwwejKjFQrTSHXS2i/+0xLJEhGu3XPQAR88eBRfv/sd0/VNEMP5Nrbr3yKirCxm
BBAkoj6FUeHV+ShzGPQo6qvZc52Ar619YZms8t//l9RcvhAe2kr+QavrD1vc4IOgbvTvR385lm1Sv2nw
7rR/1/fdWd8w9x172Nau3PI9bjie3Y+lHh44nuUP6/7As4JDOnngWLzb1q42+/YxpjYNz953Asd0XIef
qLnwaGvF0ahc3Szp5crWuKSPNrdulvSNjfXxWHuuqtduVLZ37/XGz168fPX6zdt37z+Iak5nJ7icKPPH
H+U4bvU91/eOQMyieSy6BxbY8BRL3wUOT8CdLyzbo8IyKqw6l/M41YjTmF2ZxRj2me6pxGOkVcKqqBgJ
owrgYtk7WHQPbbSx8BzLz1NGSrz0ediDHdiGBuyipuck7crhCLSjZ3DMSBtpyfGoYExKXPTinKbF/bVG
V6pobrla/U1GzwqSMDuH5zgY2ELLV3Fa9H9igcg4wZfEPse0uBkSyYWSFVIM1/87Z+//SySSfximNFqN
Osw2BBcI19o6tofTBDh/I4DFfjD8czcC5Pdfsrx8DwAA//9S01gLBCAAAA==
`,
	},

	"/static/css/.DS_Store": {
		local:   "static/css/.DS_Store",
		size:    6148,
		modtime: 1473516871,
		compressed: `
H4sIAAAJbogA/+yYsU7DMBRF7wsZUrF4ZPTIhMQfWFVBYuYHUEu3SJVoGdiy8XP8Ezh5V8ii7sAURO+R
rFOpvi/OYvsFgC1fn2+BkH92cOMdVTqOI5rCNtU44AYb7LF/6Ou1qjUWeME2p7ZlfrdZ97s1l3afx/Xn
xJRa5DkHvKHPmXri40dCCCGEEI65ust5lyGE+IOM+0OkEz24jf83dFtkAh3pRA9u47yGbumODnSkEz24
uWkZmw/jk40digU60umXLy3EmXDhCuP5f4eT/b8Q4h9j7epxtcR3Q3DEeNbGPJ7K0IlLQOMfC6+KeZFO
9ODWRUCIufgKAAD//7Yl8UAEGAAA
`,
	},

	"/static/css/reset.css": {
		local:   "static/css/reset.css",
		size:    774,
		modtime: 1473497016,
		compressed: `
H4sIAAAJbogA/3xSy67bOgzc3/+4OJsJcE5fC+drKIu22egViTKQGvn30j4JkHbRhcGxHsOZEReNAS77
G7ysaIUSqJTAiux+8qiQqVJkLB9YvmD5iuUblu9YfqDAhTxerj0ro1QGgZyroLHmdIsg7yu3BiczRrEz
Y/YMzwF+SuAIiTMkNVycxxUNjWJBixQCmla58F5ymtG6s69AFStVOHQIRk7KFd7oFN4jB/SAIJiEg29m
YMo1IpCzjoFnTh5KLpgOKio5QQ/bOuWs0IXJ9qsBqAdVldGOUhO/X0grNZOuJKGZdMfe2szdTFt58u1E
pmhn2stcs2mOnDoSrchdS1fU7m5oFux+o/UYqZoGsYQNXUDdS8ZqXfNmC7Ok4f1cLEpJsyGXq3EbmHLS
U5NfPHy8v/9//A6SFq6i55V39RROFGROg6PGQRLf/zT1dPNq4GHpnz4e2jcvrQS6DccQ3Pcot73LaWGZ
Fx0+7seD2FozoXoLPKRsIl5m5rodtf29MTi2p+OX8Rpo2hVdnzvXz4VtNNs2BsPb2/kJD67jnbfPsE5j
DoFK4+EJHimebNjHI9b7f78DAAD//30+jT4GAwAA
`,
	},

	"/static/css/style.css": {
		local:   "static/css/style.css",
		size:    2205,
		modtime: 1473517404,
		compressed: `
H4sIAAAJbogA/7RVYW+bPBD+nl9hNYqU9IWI0JeodT5N0fYb9m0y2AlWXIxsk7Sb9t93ZyAE2rTp0kHB
9vXuuHuex85tQG4p2zhhcJKKjTaC/BoRkuqn0MqfsthSmBsuTAim1ej3aJRq/lz7sGy3NboqOCXhQaQ7
6ULDuGQq3OIoCjcdR/4ii2QSELNN2TQKSPM3I4vlZEYiXFwRf18+4XM+hw+LkyQg3SuaL2ZnksaRL2rx
2SmxzocSMOwj97eIXYPUJyL0z5AJM620oWQc3+M9+CdoU1DoCiLx5YWpWCqUVyaXtlTsGaSrdLbD0FJb
6aQuKGGp1apyAq1KbBwlSTTBhdPlcd7y7gwrLGyKR0r8VDEnpiE4gTSWUDr6ttclvse8rAxzuc0VPK7t
1EeVzACM2I9vh250VtmA1AuWObmvN6iunJIFYFDoQhzdiSzKyvUxGDicAeggucsB0SiqicgFlnZiONXs
eL1ee1t9MiD3lT2C12FtBIAAFa/aEyVnXB8oChjuGMhD1QZQtRWusd6BdSCTeTw79YmhIhgOuaxZ7JHV
fHkeW8/p0HSEoTnsGmRh7U9BD841+XyCTBcOSKTk5mb1ATU2FCQDBtp1DTasoXsIktyDd56Gq8SdfEDc
tW+HRIfku0C0Lb8isaa316FqsblrsHnZa5y82+v36UVtnri91DAeQPVu6O1AmuUi2wlO/ms23KXqjy9V
/9Apmb1RQv+X3YulPV2/Lv//snh4K7ajs8fPIJDmet+yXhmLuUstgX6DLmMYH+0PLvfeA5WBnx9/q0km
ZANKAa48q6lWHIP+BAAA//+65SA1nQgAAA==
`,
	},

	"/static/index.html": {
		local:   "static/index.html",
		size:    518,
		modtime: 1473529224,
		compressed: `
H4sIAAAJbogA/5RRsU7DMBDd+xWHWWCgRmJBwsnSgsRUpLYDE3Ltq3yt61T2URqJj8dxSiMBC4vzXvzu
Xd6LupjOJovXl0dwvPP1SHUPqEcAyqG2Hchwh6zBOB0TciWWi6ebe3G6YmKP9VSTb2H5DJe3d/AJsyBn
6zWkD2LjMp/M59AE3yrZy8vocChPYQsRfSUStx6TQ2QBLuK6e6OZjDQpyYh5/zijvHyY/o9Dufzl8H0o
2Wfu4Kqx7UmjvF6hL51Q2L8zkK0EY9yliUOzXTVHAdzusRLmzJtgPJltJQ7ak9WMi05/df0gZG9UK0m5
bnnyVpYOg+9bZqXg3GL5rszrH4mTibRnSNGcE26SpGDxON7kfEr2ir+D9vFy4PLXvwIAAP//PuDsPQYC
AAA=
`,
	},
	"/index.html": {
		local:   "static/index.html",
		size:    518,
		modtime: 1473529224,
		compressed: `
H4sIAAAJbogA/5RRsU7DMBDd+xWHWWCgRmJBwsnSgsRUpLYDE3Ltq3yt61T2URqJj8dxSiMBC4vzXvzu
Xd6LupjOJovXl0dwvPP1SHUPqEcAyqG2Hchwh6zBOB0TciWWi6ebe3G6YmKP9VSTb2H5DJe3d/AJsyBn
6zWkD2LjMp/M59AE3yrZy8vocChPYQsRfSUStx6TQ2QBLuK6e6OZjDQpyYh5/zijvHyY/o9Dufzl8H0o
2Wfu4Kqx7UmjvF6hL51Q2L8zkK0EY9yliUOzXTVHAdzusRLmzJtgPJltJQ7ak9WMi05/df0gZG9UK0m5
bnnyVpYOg+9bZqXg3GL5rszrH4mTibRnSNGcE26SpGDxON7kfEr2ir+D9vFy4PLXvwIAAP//PuDsPQYC
AAA=
`,
	},

	"/static/js/.DS_Store": {
		local:   "static/js/.DS_Store",
		size:    6148,
		modtime: 1473528538,
		compressed: `
H4sIAAAJbogA/+yYsUoDQRCG/1mvWLDZ0nJLK8E3WEIUrH0BibFQgilEsNzK5/LN4q7zq4HLgVYR838w
fIHMzO01uzsHwGbPy3MgtZ8RbrxiJ5ExItDm0Xq84AwPeLparW939xphH83v8Ygl7rbrF6v1ggu7bHHa
czebH3YVQgghxBTmisf7XYYQ4g/S94dMF7q6jf8HetiqSXSmC13dxrxAD3SkE53pQlc3Ny3j8GF8snFC
sURnuvzypYU4EI5cqZ//F5ic/4UQ/xgb5tfzGb4GghH9rM0tbpjz9lk4cREI/sHwBN95mS50desyIMQ+
eA8AAP//z+aoGgQYAAA=
`,
	},

	"/static/js/index.js": {
		local:   "static/js/index.js",
		size:    801,
		modtime: 1473529266,
		compressed: `
H4sIAAAJbogA/4yST2vcMBDFz/WnGHSJTLNO6NX40qWwgYQtqQu9FVUad0VtyZXG8Yaw37364y1OE0gv
xkjv9zTzZrrJSNLWwIPotRKELbrB8/KpgHDkQDbKymlAQ9VPpE89xt+PjzeKX1BUbg8of/2wx4uyXgj1
BvFd6Yes1h1wWcnogKqE+CSAqrQx6Hbt3S00bG9YXbzzaNQ9/p7QE08kgEOanAFyE8aDE2DvEZ7gpQew
fde94lKsfToR8GRUnMINHbSvVvrgcg4qR5M6nVzfsCsxalafj/6/+SwXJJp1uX9vZK8DGd41OMO3u9sd
0fg8gqyo7IiGs8/7Ly27jBVdpkyeSTzSgu5QKHScba2hcLNpH0cMGBPj2GspYn9Xx808z5vOumET7NBI
q1Cxtw0Npnyineytx38RozjzsyZ5aBi8h2z89f5ma4fRmqDhMY2yrAOUh5P2Y8FJ0OShaeDD9XWewIs5
L1KHPvj5sMhHShWc4ieuxxl7nctPrKi0DcWfAAAA//+9f2MeIQMAAA==
`,
	},

	"/static/scss/.DS_Store": {
		local:   "static/scss/.DS_Store",
		size:    6148,
		modtime: 1473522565,
		compressed: `
H4sIAAAJbogA/+yYsUoEMRCG/4lbBG1SWqa0EnyDcJyCtS8gd9otXHFa2KXyuXwzTZxfPdg90OpE/w+G
7+BmZrNNklkAtni8uwBS+xnhxjNmiYwJgTaP1uMe59hi3WI732tCrz1u2Q94wrhbfz1u1qtxs+Lirlqc
vb7zzc5CCCGEmMNc8eSwyxBC/EL6/pDpQle38f9ADzs1ic50oavbmBfogY50ojNd6OrmpmUcPoxPNk4o
luhMlx++tBD/hCNX6uf/JfbO/0KIP4wNy5vlAp8DwYR+1uYWt8x5+SjccxEI/sHwFF95mS50desyIMQh
eAsAAP//iJoEagQYAAA=
`,
	},

	"/static/scss/style.scss": {
		local:   "static/scss/style.scss",
		size:    1870,
		modtime: 1473497016,
		compressed: `
H4sIAAAJbogA/6RUa2+bPBT+zq84atQq6RvCpS9R63x63zTdKlXJtFba+tGAE6xQjIzpZVP++46BUJPL
1GlyDM65GD/PeY7Ph3BO6FIxqRchWwrJ4KcFEIpXu+A/eLYiuJYxkzaaJtbGshwHpovr2ZfZHO4fHu9u
559gMb97BCsU8VudTKP1SooyiwlIGnOa2iv9Zpnqhyk6wQtOh6AkzYqcSjSDNz4dgAvu8I8SLvNXPfeT
5CqkfT8Ihts58gb7e/hu9VHvL3fQp7jKkRwTuR2JVEjS8y/12PEhs4x4Y8zTD+3EH/L6fxsCmVAQSUYV
iyF8gydWR3xlhShlxAgkSuXEcVJGR88Mc0ZPzImK4iKnCsuZFU4vojIUmb3koWRN4Wbz673iWZVnMXcW
Nzdw/+32YfrZSmnI0qqWMS/ylL6hClIRrSf6oLkouOIiI0DDQqSlYhpBypaKQOCe6j9K5O9rzRcK64nU
yxQx9W10DsEeI386xn5h4ZorW9HcTvgqSXGqhkKoSuEOq1FFn5GliMpiiAsaKf5cSxZAlCrlGXKTiaw6
0wYnz/JSNQEtGDNAx+z6a7CV7YXHKiHguW5dKYCE6dN1TKbie9PptLHWjaMlVRYtHyaBkiEbCGAbj02X
0Fi8EN0LOHyUSM91sS14VjDVWC/QapIy8gdmgI/HwldHwVA9vEH9oaoOzRFGflEbz5r+r2jVV0JDypHw
TSepjY1EprAzCJycTBrTIVaPqchgPGjZfafctNXsog3R4gY8rpjqeve53xFqBc8Qawv3iGCDRrBb/A0J
Xb4+zMEWqQHKFJIB5zhbW2ouDGoOQ/SD30L83jfQHRKjvq/exb0x2ke3GIkSFq3xtvqnbaePK9r/uKK3
EcFgcliBTeWby6M3G//7n3fV1Wy3XB3Od8MbiGckEc9tTlTKQm+eC46VlvVNsrF+BQAA//+I1DtJTgcA
AA==
`,
	},

	"/static/style.css": {
		local:   "static/style.css",
		size:    2205,
		modtime: 1473517404,
		compressed: `
H4sIAAAJbogA/7RVYW+bPBD+nl9hNYqU9IWI0JeodT5N0fYb9m0y2AlWXIxsk7Sb9t93ZyAE2rTp0kHB
9vXuuHuex85tQG4p2zhhcJKKjTaC/BoRkuqn0MqfsthSmBsuTAim1ej3aJRq/lz7sGy3NboqOCXhQaQ7
6ULDuGQq3OIoCjcdR/4ii2QSELNN2TQKSPM3I4vlZEYiXFwRf18+4XM+hw+LkyQg3SuaL2ZnksaRL2rx
2SmxzocSMOwj97eIXYPUJyL0z5AJM620oWQc3+M9+CdoU1DoCiLx5YWpWCqUVyaXtlTsGaSrdLbD0FJb
6aQuKGGp1apyAq1KbBwlSTTBhdPlcd7y7gwrLGyKR0r8VDEnpiE4gTSWUDr6ttclvse8rAxzuc0VPK7t
1EeVzACM2I9vh250VtmA1AuWObmvN6iunJIFYFDoQhzdiSzKyvUxGDicAeggucsB0SiqicgFlnZiONXs
eL1ee1t9MiD3lT2C12FtBIAAFa/aEyVnXB8oChjuGMhD1QZQtRWusd6BdSCTeTw79YmhIhgOuaxZ7JHV
fHkeW8/p0HSEoTnsGmRh7U9BD841+XyCTBcOSKTk5mb1ATU2FCQDBtp1DTasoXsIktyDd56Gq8SdfEDc
tW+HRIfku0C0Lb8isaa316FqsblrsHnZa5y82+v36UVtnri91DAeQPVu6O1AmuUi2wlO/ms23KXqjy9V
/9Apmb1RQv+X3YulPV2/Lv//snh4K7ajs8fPIJDmet+yXhmLuUstgX6DLmMYH+0PLvfeA5WBnx9/q0km
ZANKAa48q6lWHIP+BAAA//+65SA1nQgAAA==
`,
	},

	"/": {
		isDir: true,
		local: "/",
	},

	"/static": {
		isDir: true,
		local: "/static",
	},

	"/static/css": {
		isDir: true,
		local: "/static/css",
	},

	"/static/js": {
		isDir: true,
		local: "/static/js",
	},

	"/static/scss": {
		isDir: true,
		local: "/static/scss",
	},
}
