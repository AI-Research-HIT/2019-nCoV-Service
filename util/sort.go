package util

import "github.com/AI-Research-HIT/2019-nCoV-Service/protodef"

type BaiduCitySlice []protodef.BaiduCityT

func (s BaiduCitySlice) Len() int { return len(s) }

func (s BaiduCitySlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s BaiduCitySlice) Less(i, j int) bool { return s[i].Date < s[j].Date }
