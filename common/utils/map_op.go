package utils

import "encoding/json"

func MapGetStringOr(m map[string]interface{}, key string, value string) string {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		v, typeOk := r.(string)
		if typeOk {
			return v
		}
	}
	return value
}

func MapGetStringOr2(m map[string]string, key string, value string) string {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		return r
	}
	return value
}

func MapStringGetOr(m map[string]string, key string, value string) string {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		return r
	}

	return value
}

func MapStringGet(m map[string]string, key string) string {
	return MapStringGetOr(m, key, "")
}

func MapGetRaw(m map[string]interface{}, key string) interface{} {
	return MapGetRawOr(m, key, nil)
}

func MapGetRawOr(m map[string]interface{}, key string, value interface{}) interface{} {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		return r
	} else {
		return value
	}
}

func MapGetString(m map[string]interface{}, key string) string {
	return MapGetStringOr(m, key, "")
}

func MapGetString2(m map[string]string, key string) string {
	return MapGetStringOr2(m, key, "")
}

func MapGetMapRaw(m map[string]interface{}, key string) map[string]interface{} {
	return MapGetMapRawOr(m, key, make(map[string]interface{}))
}

func MapGetMapRawOr(m map[string]interface{}, key string, value map[string]interface{}) map[string]interface{} {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		data, typeOk := r.(map[string]interface{})
		if typeOk {
			return data
		}
	}
	return value
}

func MapGetIntOr(m map[string]interface{}, key string, value int) int {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		v, typeOk := r.(int)
		if typeOk {
			return v
		}
	}
	return value
}

func MapGetInt(m map[string]interface{}, key string) int {
	return MapGetIntOr(m, key, 0)
}

func MapGetFloat64Or(m map[string]interface{}, key string, value float64) float64 {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		v, typeOk := r.(float64)
		if typeOk {
			return v
		}
	}
	return value
}

func MapGetFloat64(m map[string]interface{}, key string) float64 {
	return MapGetFloat64Or(m, key, 0)
}

func MapGetFloat32Or(m map[string]interface{}, key string, value float32) float32 {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		v, typeOk := r.(float32)
		if typeOk {
			return v
		}
	}
	return value
}

func MapGetFloat32(m map[string]interface{}, key string) float32 {
	return MapGetFloat32Or(m, key, 0)
}

func MapGetBoolOr(m map[string]interface{}, key string, value bool) bool {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		v, typeOk := r.(bool)
		if typeOk {
			return v
		}
	}
	return value
}

func MapGetBool(m map[string]interface{}, key string) bool {
	return MapGetBoolOr(m, key, false)
}

func MapGetInt64Or(m map[string]interface{}, key string, value int64) int64 {
	if m == nil {
		return value
	}

	r, ok := m[key]
	if ok {
		v, typeOk := r.(int64)
		if typeOk {
			return v
		}
	}
	return value
}

func MapGetInt64(m map[string]interface{}, key string) int64 {
	return MapGetInt64Or(m, key, 0)
}

func ToMapParams(params interface{}) (map[string]interface{}, error) {
	raw, err := json.Marshal(params)
	if err != nil {
		return nil, Errorf("marshal params failed: %s", err)
	}

	var p = map[string]interface{}{}
	err = json.Unmarshal(raw, &p)
	if err != nil {
		return nil, Errorf("unmarshal map params failed: %s", err)
	}

	return p, nil
}

func MergeStringMap(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
