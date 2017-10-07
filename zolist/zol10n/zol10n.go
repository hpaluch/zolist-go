package zol10n

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
	return fmt.Sprintf("/%s", UrlLangs[langIndex])
}

var serverLangs = []language.Tag{
	language.BritishEnglish, // en-GB fallback
	language.Czech,          // de
}

var reUrlBaseLang = regexp.MustCompile(`^/(cz|en)/`)

// Here are all Accept-Language values that should match Czech
var reAcceptCzLang = regexp.MustCompile(`(cs|cz|sk)`)

func LangFromHeader(ctx appengine.Context, r *http.Request) int {
	var acceptLang = r.Header.Get("Accept-Language")
	if acceptLang == "" {
		acceptLang = "en-GB"
	}
	ctx.Infof("Accept-Language: %T %s", acceptLang, acceptLang)
	var langIndex = LangEnIndex

	if reAcceptCzLang.MatchString(acceptLang) {
		langIndex = LangCzIndex
	}
	return langIndex
}

func LangFromUrlBase(ctx appengine.Context, r *http.Request) (int, error) {
	var groups = reUrlBaseLang.FindStringSubmatch(r.URL.Path)
	if len(groups) != 2 {
		ctx.Errorf("Got unexpected number of groups %d <> 2: %v for path '%s'", len(groups), groups, r.URL.Path)
		return -1, errors.New("Invalid path prefix")
	}
	var lang = groups[1]
	for i := range UrlLangs {
		if UrlLangs[i] == lang {
			return i, nil
		}
	}
	ctx.Errorf("Language '%s' not found", lang)
	return -1, errors.New(fmt.Sprintf("Language '%s' not found", lang))
}

func LocFromIndex(langIndex int) *message.Printer {
	var tag = serverLangs[langIndex]
	return message.NewPrinter(tag)
}

func LocFromUrlBase(ctx appengine.Context, r *http.Request) (int, *message.Printer, error) {

	var langIndex, err = LangFromUrlBase(ctx, r)
	if err != nil {
		return -1, nil, err
	}

	var tag = serverLangs[langIndex]

	//ctx.Infof("Messages: %v", message.DefaultCatalog.Languages())
	p := message.NewPrinter(tag)
	var lText = p.Sprintf("There are %d items", 12345)
	ctx.Infof(lText)
	// WARNING: p.Sprint() does not work(?!)
	return langIndex, p, nil
}

func init() {
	message.SetString(language.BritishEnglish, "There are %d items", "There are %d items")
	message.SetString(language.Czech, "There are %d items", "Je tam %d položek")

	message.SetString(language.BritishEnglish, "Render Time", "Render Time")
	message.SetString(language.Czech, "Render Time", "Doba generování stránky")

	message.SetString(language.BritishEnglish, "ZoList", "ZoList")
	message.SetString(language.Czech, "ZoList", "ZoSeznam")

	message.SetString(language.BritishEnglish, "Now is", "Now is")
	message.SetString(language.Czech, "Now is", "Aktualizováno")

	message.SetString(language.BritishEnglish, "See", "See")
	message.SetString(language.Czech, "See", "Podívejte se na")

	const SRC_GIT_HUB = "Source on GitHub"
	message.SetString(language.BritishEnglish, SRC_GIT_HUB, SRC_GIT_HUB)
	message.SetString(language.Czech, SRC_GIT_HUB, "Zdrojáky na GitHubu")

	const MEAL = "Meal"
	message.SetString(language.BritishEnglish, MEAL, MEAL)
	message.SetString(language.Czech, MEAL, "Jídlo")

	const PRICE = "Price"
	message.SetString(language.BritishEnglish, PRICE, PRICE)
	message.SetString(language.Czech, PRICE, "Cena")

	const UPDATED = "Updated on"
	message.SetString(language.BritishEnglish, UPDATED, UPDATED)
	message.SetString(language.Czech, UPDATED, "Aktualizováno")

	const REST = "Restaurant %s"
	message.SetString(language.BritishEnglish, REST, "Restaurant ‘%s’")
	message.SetString(language.Czech, REST, "Restaurace „%s“")

	const MENU_DETAIL = "Menu Detail"
	message.SetString(language.BritishEnglish, MENU_DETAIL, MENU_DETAIL)
	message.SetString(language.Czech, MENU_DETAIL, "Detail nabídky")

	const MENU_PREV = "Prev. Menu"
	message.SetString(language.BritishEnglish, MENU_PREV, MENU_PREV)
	message.SetString(language.Czech, MENU_PREV, "Předchozí nabídka")

	const MENU_NEXT = "Next Menu"
	message.SetString(language.BritishEnglish, MENU_NEXT, MENU_NEXT)
	message.SetString(language.Czech, MENU_NEXT, "Další nabídka")

	const BACK_TO_L = "Back to List"
	message.SetString(language.BritishEnglish, BACK_TO_L, BACK_TO_L)
	message.SetString(language.Czech, BACK_TO_L, "Zpět na seznam")

	const DETAIL_OF = "Detail of %s"
	message.SetString(language.BritishEnglish, DETAIL_OF, DETAIL_OF)
	message.SetString(language.Czech, DETAIL_OF, "Detail nabídky z %s")

	const LIST_TITLE = "Favorite Restaurants menu"
	message.SetString(language.BritishEnglish, LIST_TITLE, LIST_TITLE)
	message.SetString(language.Czech, LIST_TITLE, "Nabídka oblíbených restaurací")
}
