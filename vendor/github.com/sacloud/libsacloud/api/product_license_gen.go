package api

/************************************************
  generated by IDE. for [ProductLicenseAPI]
************************************************/

import (
	"github.com/sacloud/libsacloud/sacloud"
)

/************************************************
   To support fluent interface for Find()
************************************************/

// Reset 検索条件のリセット
func (api *ProductLicenseAPI) Reset() *ProductLicenseAPI {
	api.reset()
	return api
}

// Offset オフセット
func (api *ProductLicenseAPI) Offset(offset int) *ProductLicenseAPI {
	api.offset(offset)
	return api
}

// Limit リミット
func (api *ProductLicenseAPI) Limit(limit int) *ProductLicenseAPI {
	api.limit(limit)
	return api
}

// Include 取得する項目
func (api *ProductLicenseAPI) Include(key string) *ProductLicenseAPI {
	api.include(key)
	return api
}

// Exclude 除外する項目
func (api *ProductLicenseAPI) Exclude(key string) *ProductLicenseAPI {
	api.exclude(key)
	return api
}

// FilterBy 指定キーでのフィルター
func (api *ProductLicenseAPI) FilterBy(key string, value interface{}) *ProductLicenseAPI {
	api.filterBy(key, value, false)
	return api
}

// func (api *ProductLicenseAPI) FilterMultiBy(key string, value interface{}) *ProductLicenseAPI {
// 	api.filterBy(key, value, true)
// 	return api
// }

// WithNameLike 名称条件
func (api *ProductLicenseAPI) WithNameLike(name string) *ProductLicenseAPI {
	return api.FilterBy("Name", name)
}

// WithTag タグ条件
func (api *ProductLicenseAPI) WithTag(tag string) *ProductLicenseAPI {
	return api.FilterBy("Tags.Name", tag)
}

// WithTags タグ(複数)条件
func (api *ProductLicenseAPI) WithTags(tags []string) *ProductLicenseAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ProductLicenseAPI) WithSizeGib(size int) *ProductLicenseAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *ProductLicenseAPI) WithSharedScope() *ProductLicenseAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *ProductLicenseAPI) WithUserScope() *ProductLicenseAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

// SortBy 指定キーでのソート
func (api *ProductLicenseAPI) SortBy(key string, reverse bool) *ProductLicenseAPI {
	api.sortBy(key, reverse)
	return api
}

// SortByName 名称でのソート
func (api *ProductLicenseAPI) SortByName(reverse bool) *ProductLicenseAPI {
	api.sortByName(reverse)
	return api
}

// func (api *ProductLicenseAPI) SortBySize(reverse bool) *ProductLicenseAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

// func (api *ProductLicenseAPI) Create(value *sacloud.ProductLicense) (*sacloud.ProductLicense, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.create(api.createRequest(value), res)
// 	})
// }

// Read 読み取り
func (api *ProductLicenseAPI) Read(id int64) (*sacloud.ProductLicense, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

// func (api *ProductLicenseAPI) Update(id int64, value *sacloud.ProductLicense) (*sacloud.ProductLicense, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.update(id, api.createRequest(value), res)
// 	})
// }

// func (api *ProductLicenseAPI) Delete(id int64) (*sacloud.ProductLicense, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.delete(id, nil, res)
// 	})
// }

/************************************************
  Inner functions
************************************************/

func (api *ProductLicenseAPI) setStateValue(setFunc func(*sacloud.Request)) *ProductLicenseAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *ProductLicenseAPI) request(f func(*sacloud.Response) error) (*sacloud.ProductLicense, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.LicenseInfo, nil
}

func (api *ProductLicenseAPI) createRequest(value *sacloud.ProductLicense) *sacloud.Request {
	req := &sacloud.Request{}
	req.LicenseInfo = value
	return req
}