package parser

import (
	"go-crawler/engine"
	"regexp"
)
const cityRe = `<a href="(http://album.zhenai.com/u/[0-9]*)"[^>]*>([^<]+)</a>`
func ParserCity(contents []byte) engine.ParseResult {
	re := regexp.MustCompile(cityRe)
	matches := re.FindAllSubmatch(contents, -1)
	result := engine.ParseResult{}
	for _, m := range matches {
		result.Items = append(result.Items, "User " + string(m[2]))
		result.Requests = append(result.Requests, engine.Request{
			Url: string(m[1]),
			ParserFunc: ParserProfile,
		})
	}
	return result
}