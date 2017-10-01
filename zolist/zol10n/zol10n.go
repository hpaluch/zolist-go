package zol10n;

import (
	"net/http"
	"regexp"

	"appengine"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var serverLangs = []language.Tag{
    language.BritishEnglish, // en-GB fallback
    language.Czech,          // de
}

var reAcceptCzLang = regexp.MustCompile(`(cs|cz|sk)`)

func ZoL10n(ctx appengine.Context, r *http.Request ) *message.Printer {
	var acceptLang = r.Header.Get("Accept-Language")
	if acceptLang == "" {
		acceptLang = "en-GB"
	}
	ctx.Infof("Accept-Language: %T %s",acceptLang,acceptLang)
	var tag = language.BritishEnglish

	if reAcceptCzLang.MatchString(acceptLang) {
		tag = language.Czech
	}

	ctx.Infof("Messages: %v",message.DefaultCatalog.Languages())
	p := message.NewPrinter(tag)
	var lText = p.Sprintf("There are %d items",12345)
	ctx.Infof(lText)
	// WARNING: p.Sprint() does not work(?!)
	var text2 = p.Sprintf("Render Time")
	ctx.Infof(text2)
	return p
}

func init() {
	message.SetString(language.BritishEnglish,"There are %d items","There are %d items (GB)")
	message.SetString(language.Czech,"There are %d items","Je tam %d položek (CZ)")

	message.SetString(language.BritishEnglish,"Render Time","Render Time (GB)")
	message.SetString(language.Czech,"Render Time","Doba generování stránky (CZ)")
}
