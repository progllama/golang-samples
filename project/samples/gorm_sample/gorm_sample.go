package gorm_sample

import (
	"bufio"
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name string `json:"name"`
}

var (
	db  *gorm.DB
	err error
)

func GormMain() {
	// countryNames := getCountryNames()

	db, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	// db.AutoMigrate(&Country{})

	// countries := make([]Country, len(countryNames))

	// for i, name := range countryNames {
	// 	countries[i] = Country{Name: name}
	// }

	// db.CreateInBatches(countries, 100)

	router := gin.Default()
	router.GET("/countries", func(c *gin.Context) {
		var countries []struct {
			Name string
		}
		db.Model(&Country{}).Select("name").Find(&countries)

		c.JSON(http.StatusOK, gin.H{"countries": countries})
	})
	router.GET("/countries/:name", func(c *gin.Context) {

	})
	router.Run()
}

func getCountryNames() []string {
	visitor := NewVisitor()

	filepath.Walk("./", visitor.Visit)
	// filepath.Globは**に対応していない。

	if !visitor.Contains("Country") {
		fmt.Println("国名データなし。")
		return []string{}
	}

	path := visitor.Find("Country")

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	buf := bytes.NewBuffer(data)
	bs := bufio.NewScanner(buf)
	countryNames := make([]string, 0)

	for bs.Scan() {
		countryNames = append(countryNames, bs.Text())
	}

	return countryNames
}

type Visitor struct {
	paths *list.List
	cache []string
}

func (v *Visitor) Contains(pat string) bool {
	for _, path := range v.Paths() {
		if strings.Contains(path, pat) {
			return true
		}
	}
	return false
}

// こういった処理で検索値がもとまらない時の値はどうすればよい?
// 1 戻り値の型のゼロ値を返す。
// 2 nilを返す。
// 3 戻り値に見つかったかのパラメータをboolとして入れる
// 4 戻り値に見つかったかのパラメータをerrorとして入れる
// 5 panicを起こす
// containsなどを用意して事前にチェックしていることを前提にしているかによっても変わる？

func (v *Visitor) Find(str string) string {
	for _, path := range v.Paths() {
		if strings.Contains(path, str) {
			return path
		}
	}
	return ""
}

func (v *Visitor) Paths() []string {
	if len(v.cache) != 0 {
		return v.cache
	}

	paths := make([]string, 0)
	for path := v.paths.Front(); path != nil; path = path.Next() {
		paths = append(paths, path.Value.(string))
	}
	v.cache = paths
	return paths
}

func NewVisitor() Visitor {
	return Visitor{
		paths: list.New(),
		cache: nil,
	}
}

func (v *Visitor) Visit(path string, info os.FileInfo, err error) error {
	v.paths.PushBack(path)
	return nil
}

func (v *Visitor) Print() {
	for path := v.paths.Front(); path != nil; path = path.Next() {
		fmt.Printf("%s\n", path.Value.(string)) // print out the elements
	}
}
