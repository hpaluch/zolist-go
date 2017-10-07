package zol10n;

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"appengine"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const LangEnIndex = 0
const LangCzIndex = 1

var HtmlLangs = []string{
	"en",
	"cs",
}
var LangNames = []string{
	"English",
	"Česky",
}
var UrlLangs = []string{
	"en",
	"cz",
}

func LangUrlBase(langIndex int) string {
	return fmt.Sprintf("/%s",UrlLangs[ langIndex ])
}

var serverLangs = []language.Tag{
    language.BritishEnglish, // en-GB fallback
    language.Czech,          // de
}

var reUrlBaseLang = regexp.MustCompile(`^/(cz|en)/`)

// Here are all Accept-Language values that should match Czech
var reAcceptCzLang = regexp.MustCompile(`(cs|cz|sk)`)

func LangFromHeader(ctx appengine.Context, r *http.Request ) int {
	var acceptLang = r.Header.Get("Accept-Language")
	if acceptLang == "" {
		acceptLang = "en-GB"
	}
	ctx.Infof("Accept-Language: %T %s",acceptLang,acceptLang)
	var langIndex = LangEnIndex 

	if reAcceptCzLang.MatchString(acceptLang) {
		langIndex = LangCzIndex
	}
	return langIndex
}

func LangFromUrlBase(ctx appengine.Context, r *http.Request) (int,error) {
	var groups = reUrlBaseLang.FindStringSubmatch(r.URL.Path)
	if len(groups) != 2 {
		ctx.Errorf("Got unexpected number of groups %d <> 2: %v for path '%s'", len(groups), groups,r.URL.Path)
		return -1,errors.New("Invalid path prefix")
	}
	var lang = groups[1]
	for i := range UrlLangs {
		if UrlLangs[i] == lang {
			return i,nil
		}
	}
	ctx.Errorf("Language '%s' not found",lang)
	return -1,errors.New(fmt.Sprintf("Language '%s' not found",lang))
}


func LocFromUrlBase(ctx appengine.Context, r *http.Request ) (int,*message.Printer,error) {

	var langIndex,err = LangFromUrlBase(ctx,r)
	if err != nil {
		return -1,nil,err
	}

	var tag = serverLangs[ langIndex ]

	ctx.Infof("Messages: %v",message.DefaultCatalog.Languages())
	p := message.NewPrinter(tag)
	var lText = p.Sprintf("There are %d items",12345)
	ctx.Infof(lText)
	// WARNING: p.Sprint() does not work(?!)
	var text2 = p.Sprintf("Render Time")
	ctx.Infof(text2)
	return langIndex,p,nil
}

func init() {
	message.SetString(language.BritishEnglish,"There are %d items","There are %d items (GB)")
	message.SetString(language.Czech,"There are %d items","Je tam %d položek (CZ)")

	message.SetString(language.BritishEnglish,"Render Time","Render Time (GB)")
	message.SetString(language.Czech,"Render Time","Doba generování stránky (CZ)")
}
