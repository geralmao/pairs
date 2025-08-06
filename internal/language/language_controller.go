package language

import (
	"encoding/json"
	"fmt"

	"github.com/programatta/pairs/internal/assets/lang"
)

func Init() {
	langCtrl = newLanguageController()
}

func LoadById(langId string) {
	langCtrl.load(langId)
	Value = langCtrl.Data()
}

var Value *lang.LanguageData

func newLanguageController() *languageController {
	return &languageController{
		langData: nil,
	}
}

func (lc *languageController) load(langId string) {
	langfile := fmt.Sprintf("%s.json", langId)
	data, err := lang.LanguagesFS.ReadFile(langfile)
	if err != nil {
		fmt.Printf("Error loading language with id [%s]! \n\b- Error:%v \n\b- Trying load default language...", langId, err)
		langfile = fmt.Sprintf("%s.json", lang.LangDefault)
		data, err = lang.LanguagesFS.ReadFile(langfile)
		if err != nil {
			panic(fmt.Sprintf("Error loading default language [%s] \n\b- Error: %v!", lang.LangDefault, err))
		}
		fmt.Println("\nDefault language loaded!")
	}

	lc.langData = &lang.LanguageData{}
	err = json.Unmarshal(data, lc.langData)
	if err != nil {
		panic(fmt.Sprintf("Error unmarsal data from [%s] \n\b- Error:%v", langfile, err))
	}
}

func (lc *languageController) Data() *lang.LanguageData {
	return lc.langData
}

type languageController struct {
	langData *lang.LanguageData
}

var langCtrl *languageController
