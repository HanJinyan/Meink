package app

//非常简单的Netfarmwork框架
// 非常简单的路由处理器
import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

/*数据结构*/

type MatchMap map[string]string //匹配 ,相配字典

type Handle func(ctx *Context) //处理，控制器

type Context struct {
	http.ResponseWriter //写入回应
	Res                 http.ResponseWriter
	Req                 *http.Request          //请求
	Param               MatchMap               //参数
	Ware                map[string]interface{} //作…用的器皿 ,物品
	Stop                func()                 //匿名函数
}

type Core struct {
	router      map[string][]Handle
	patternAry  [][]string    //模式，方法
	patternRegx regexp.Regexp //patternAry的正则表达式
}

/*内部方法的实现*/
func (core *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchMap, pattern := core.match(r.URL.Path)
	ctx := &Context{w, w, r, matchMap, make(map[string]interface{}), nil}
	allHandle := make([]Handle, 0)
	handleAry1, ok1 := core.router[r.Method+":"+pattern]
	if ok1 {
		allHandle = append(allHandle, handleAry1...)
	}
	handleAry2, ok2 := core.router[r.Method+":*"]
	if ok2 {
		allHandle = append(allHandle, handleAry2...)
	}
	if len(allHandle) != 0 {
		for _, handle := range allHandle {
			keep := true
			ctx.Stop = func() {
				keep = false
			}
			handle(ctx)
			if !keep {
				return
			}
		}
	} else {
		http.NotFound(w, r)
	}
}

func (core *Core) match(path string) (matchMap MatchMap, pattern string) {

	pathAry := core.patternRegx.FindAllString(path, -1)

	matchMap = make(map[string]string)

	for _, patternItem := range core.patternAry {

		if len(pathAry) != len(patternItem) {
			continue
		}
		for j, patternKey := range patternItem {
			if j > len(pathAry)-1 {
				break
			}
			pathKey := pathAry[j]
			if strings.HasPrefix(patternKey, ":") {
				name := strings.TrimPrefix(patternKey, ":")
				matchMap[name] = pathKey
			} else {
				if pathKey != patternKey {
					break
				}
			}
			//匹配成功
			if j == len(patternItem)-1 {
				pattern = strings.Join(patternItem, "/")
				return
			}
		}
	}
	return
}
func (core *Core) addHandle(method string, pattern string, handle Handle) {
	path := method + ":" + pattern
	if _, ok := core.router[path]; !ok {
		core.router[path] = make([]Handle, 0)
	}
	core.router[path] = append(core.router[path], handle)
	core.patternAry = append(core.patternAry, strings.Split(pattern, "/"))
}

/* 对外接口*/

func (core *Core) Use(handle Handle) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	for _, method := range methods {
		core.addHandle(method, "*", handle)
	}
}

func (core *Core) Get(pattern string, handle Handle) {
	core.addHandle("GET", pattern, handle)
}

func (core *Core) Post(pattern string, handle Handle) {
	core.addHandle("POST", pattern, handle)
}

func (core *Core) Put(pattern string, handle Handle) {
	core.addHandle("PUT", pattern, handle)
}

func (core *Core) Delete(pattern string, handle Handle) {
	core.addHandle("DELETE", pattern, handle)
}

func (core *Core) Options(pattern string, handle Handle) {
	core.addHandle("OPTIONS", pattern, handle)
}

func (core *Core) Head(pattern string, handle Handle) {
	core.addHandle("HEAD", pattern, handle)
}
func (core *Core) Listen(addr string) {
	err := http.ListenAndServe(addr, core)
	if err != nil {
		fmt.Println(err)
	}
}
func Cors(ctx *Context) {
	ctx.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	if ctx.Req.Method == "OPTIONS" {
		ctx.Write([]byte{})
		ctx.Stop()
	}
}
func NewHttpServer() (core Core) {
	core = Core{}
	core.router = make(map[string][]Handle)
	core.patternAry = make([][]string, 0)
	patternRegx, _ := regexp.Compile("([^/])*")
	core.patternRegx = *patternRegx
	return
}
func Static(root string) func(ctx *Context) {
	return func(ctx *Context) {
		reqURL := ctx.Req.URL.Path
		if reqURL == "" || reqURL == "/" {
			indexPath := filepath.Join(root, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(ctx.Res, ctx.Req, indexPath)
				ctx.Stop()
			}
		} else {
			fileName := root + reqURL
			filePath, _ := filepath.Abs(fileName)
			rootParh, _ := filepath.Abs(root)
			fileDir := filepath.Dir(filePath)
			f, err := os.Stat(filePath)
			if err == nil && filepath.HasPrefix(fileDir, rootParh) {
				if f.IsDir() {
					http.ServeFile(ctx.Res, ctx.Req, filepath.Join(filePath, "index.html"))
				} else {
					http.ServeFile(ctx.Res, ctx.Req, filePath)
				}
				ctx.Stop()
			}
		}
	}
}
