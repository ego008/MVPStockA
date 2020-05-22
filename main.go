package main

import (
	"encoding/json"
	"flag"
	"github.com/jinfeijie/MVPStockA/file"
	"github.com/jinfeijie/localcache"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"time"
)

var (
	detail string
	limit  float64
	cache  *localcache.Cache
	f      = "cache.json"
)

func init() {
	cache = localcache.NewCache()
	if !file.IsExist(f) {
		file.Create(f)
	}

	for key, val := range gjson.Parse(file.Read(f)).Map() {
		x := &file.Node{}
		_ = json.Unmarshal([]byte(val.String()), x)
		cache.Set(key, x)
	}

	go func() {
		for {
			if cache.Size() > 0 {
				data, err := json.Marshal(cache.All())
				if err != nil {
					panic(err.Error())
				}
				file.Write(f, string(data))
			}
			time.Sleep(time.Second)
		}
	}()
}

func main() {
	flag.StringVar(&detail, "d", "no", "yes/no")
	flag.Float64Var(&limit, "l", 100000000, "用于过滤的最低交易值")
	flag.Parse()
	go func() {
		for {
			fetch()
			time.Sleep(time.Second)
		}
	}()

	for {
		Show()
		time.Sleep(time.Second)
	}
}

func Show() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"名称", "代码", "开盘价", "当前价", "变动金额", "变动比例", "限制金额"})
	p := message.NewPrinter(language.English)

	var code []string
	for _, val := range cache.All() {
		_node := val.(*file.Node)
		code = append(code, _node.Symbol)
	}
	sort.Strings(code)

	for _, c := range code {
		n := cache.Get(c)
		_node := n.(*file.Node)
		if _node.Percent > 9.9 {
			table.Append([]string{_node.Name, _node.Symbol, _node.Open, _node.Price, p.Sprintf("%.2f", _node.Change), "🔥🔥🔥" + p.Sprintf("%.2f", _node.Percent) + "%", p.Sprintf("%.2f", file.StrToFloat64(_node.TotalAmountMF))})
		} else {
			table.Append([]string{_node.Name, _node.Symbol, _node.Open, _node.Price, p.Sprintf("%.2f", _node.Change), "👍👍👍" + p.Sprintf("%.2f", _node.Percent) + "%", p.Sprintf("%.2f", file.StrToFloat64(_node.TotalAmountMF))})
		}
	}
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	_ = cmd.Run()

	table.SetCaption(true, time.Now().Format("2006-01-02 15:04:05"))
	table.Render()

	data, _ := json.Marshal(cache.All())
	go file.Report(string(data))
}

func fetch() {
	for i := 1; i < 100; i++ {
		requestUrl := "https://quotes.sina.cn/hq/api/openapi.php/StockV2Service.getNodeList"
		query := url.Values{}
		query.Set("num", "20")
		query.Set("sort", "ddjl")
		query.Set("asc", "0")
		query.Set("mtype", "lv2")
		query.Set("page", strconv.Itoa(i))
		data := file.Get(requestUrl, query, nil, 5000)

		if len(gjson.Get(data, "result.data.data").Array()) == 0 {
			goto RETURN
		}
		for _, val := range gjson.Get(data, "result.data.data").Array() {
			_node := &file.Node{}
			err := json.Unmarshal([]byte(val.Get("ext").String()), &_node)
			if err != nil {
				panic(err.Error())
			}
			_node.UpdateTime = time.Now().Unix()
			if file.StrToFloat64(_node.TotalAmountMF) > limit &&
				_node.RpNetRate > 0.1 &&
				_node.RxlNetRate > 0.1 &&
				_node.RlNet > 0 &&
				_node.RmNet <= 0 &&
				_node.RsNet <= 0 &&
				(_node.RxlNet+_node.RlNet)-(_node.RmNet+_node.RsNet) > 0 {
				cache.Set(_node.Symbol, _node)
			} else {
				cache.Delete(_node.Symbol)
			}
		}
	}
RETURN:
	return
}
