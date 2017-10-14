package services

import (
	"errors"
	"fmt"
	"myMovies/models"
	"myMovies/utils"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var prePath = models.HttpPath

//创建movie对象  一个对象表示一个小电影的所有查询信息
func UpdateMovies(catelog string) error {
	fmt.Println("更新代号为" + catelog + "的模块...")

	firstUrl := prePath + "/Html/" + catelog + "/"
	total, totalPage, err := getTotalNums(firstUrl)
	if err != nil {
		return errors.New("查询影片数量和页数失败：" + err.Error())
	}

	currentNums, err := models.Engine.Where("catelog = ? ", catelog).Count(&models.MyMovie{})
	if err != nil {
		return errors.New("查询当前已存在总数失败")
	}

	if int64(total) == currentNums {
		fmt.Println(catelog + "模块下的影片地址已是最新")
		return nil
	}

	//fmt.Println(total, currentNums)
	count := 0
	for i := 1; i <= totalPage && count <= (total-int(currentNums)); i++ {

		fmt.Println("\t 获取第", i, "页信息...")

		suffix := fmt.Sprint("index-", i, ".html")
		if i == 1 {
			suffix = ""
		}
		urls, err := getFirstUrlByPage(firstUrl + suffix)
		if err != nil {
			return errors.New("获取第" + strconv.Itoa(i) + "页的一级url失败")
		}
		for _, url := range urls {
			count++
			err = createMyMovieByTagUrl(url, catelog)
			if err != nil {
				return errors.New("创建movie对象失败: " + err.Error())
			}
		}
	}

	fmt.Println("++ 更新" + catelog + "模块成功!!!")
	return nil
}

//获取当前项的最大影片数和最大页数
func getTotalNums(firstUrl string) (int, int, error) {
	bs, err := utils.GetForContents(firstUrl)
	if err != nil {
		return 0, 0, errors.New("获取" + firstUrl + "的内容失败: " + err.Error())
	}
	reg := regexp.MustCompile("共.*?</strong>")
	//获取第一步  "共318部&nbsp;1/27</strong>"
	nums := reg.Find(bs)
	//获取所有数字
	reg = regexp.MustCompile("\\d+")
	result := reg.FindAll(nums, len(nums)-1)
	if len(result) < 3 {
		return 0, 0, errors.New("获取有效数字失败; ")
	}
	//总数
	total, err := strconv.Atoi(string(result[0]))
	if err != nil {
		return 0, 0, errors.New("获取影片总数失败: " + err.Error())
	}
	maxPage, err := strconv.Atoi(string(result[2]))
	if err != nil {
		return 0, 0, errors.New("获取最大页数失败: " + err.Error())
	}
	return total, maxPage, nil
}

//通过带页数的url 获取当前页所有的一级url  返回url数组（一般是6个）
func getFirstUrlByPage(cacheUrl string) ([]string, error) {
	datas, err := utils.GetForContents(cacheUrl)
	if err != nil {
		return nil, errors.New("获取第一个页面失败: " + err.Error())
	}
	result := getMessageByReg(datas)
	urls := []string{}
	for _, s := range result {
		urls = append(urls, prePath+s)
	}
	return urls, nil
}

//利用正则获取需要的信息(第一级的网址)
func getMessageByReg(body []byte) []string {
	//fmt.Println(string(body))
	reg := regexp.MustCompile("/Html/.*?html")
	bses := reg.FindAll(body, -1)
	result := []string{}
	for _, bs := range bses {
		if strings.Contains(string(bs), "index") {
			continue
		}
		result = append(result, string(bs))
	}
	return result
}

//利用正则获取下载链接，返回值为下载连接和names的map 这里是二级的网址了
func getSecondUrlAndNames(body []byte) (string, string) {
	//获取link
	reg := regexp.MustCompile("https:.*?.mp4")
	bs := reg.Find(body)
	//获取name=link

	title := getNameWithLink(body)
	return string(bs), title
}

//对name进行处理
func getNameWithLink(body []byte) string {
	reg := regexp.MustCompile(`<dd class="film_title"><h1>(.*?)</h1></dd>`)
	bs := reg.Find(body)
	//文件名不能有空格，所以把空格去掉
	str := strings.Trim(string(bs), `<dd class="film_title"><h1>`)
	str = strings.Trim(str, `</h1></dd>`)
	str = strings.Replace(str, " ", "", len(bs)-1)
	return str
}

//创建myMovie对象
func createMyMovieByTagUrl(tagUrl, catelog string) error {
	movie, err := getMoviesByTag(tagUrl, catelog)
	if err != nil {
		return errors.New("从" + tagUrl + "获取movie对象失败: " + err.Error())
	}
	//fmt.Println(movie.Name, movie.FirstUrl, movie.SecondUrl)
	_, err = models.Engine.Insert(&movie)
	if err != nil {
		return errors.New("创建movie对象失败： " + err.Error())
	}
	return nil
}

//获取到当前url对应的movie对象
func getMoviesByTag(tagUrl, catelog string) (models.MyMovie, error) {

	datas, err := utils.GetForContents(tagUrl)
	if err != nil {
		return models.MyMovie{}, errors.New("获取" + tagUrl + "内容失败：" + err.Error())
	}
	secondUrl, name := getSecondUrlAndNames(datas)
	movie := models.MyMovie{
		Name:      name,
		FirstUrl:  tagUrl,
		SecondUrl: secondUrl,
		Catelog:   catelog,
		UpdateAt:  time.Now(),
	}
	return movie, nil
}
