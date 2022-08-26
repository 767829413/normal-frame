package logger

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	logredis "gitlab.oneitfarm.com/bifrost/logrus-redis-hook"

	"github.com/767829413/normal-frame/internal/pkg/options"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const defaultTimestampFormat = "2006-01-02T15:04:05.000-0700"

const (
	fieldAppName       = "appName"       // 微服务appName
	fieldAppID         = "appId"         // 服务appId
	fieldAppVersion    = "appVersion"    // 微服务app版本号
	fieldAppKey        = "appKey"        // appkey
	fieldChannel       = "channel"       // channel
	fieldSubOrgKey     = "subOrgKey"     // 机构唯一码
	fieldTime          = "timestamp"     // 日志时间字符串
	fieldLevel         = "level"         // 日志等级 : DEBUG、INFO 、NOTICE、WARNING、ERR、CRIT、ALERT、 EMERG(系统不可用)
	fieldHostName      = "hostname"      // 主机名
	fieldIP            = "ip"            // ip地址
	fieldPodName       = "podName"       // pod名
	fieldPodIP         = "podIp"         // pod IP
	fieldNodeName      = "nodeName"      // pod内部的node名
	fieldNodeIP        = "nodeIp"        // k8s注入的node节点IP
	fieldContainerName = "containerName" // k8s容器name ，主要进行容器环境区分
	fieldClusterUID    = "clusterUid"    // 集群ID
	fieldImageURL      = "imageUrl"      // 应用镜像URL地址
	fieldUniqueID      = "uniqueId"      // 部署的服务唯一ID
	fieldSiteUID       = "siteUid"       // 可用区唯一标识符
	fieldRunEnvType    = "runEnvType"    // 区分开发环境(development)、测试环境(test)、预发布环境 (pre_release)、生产环境(production) 从环境变量中获取
	fieldMessage       = "message"       // 日志内容
	fieldLogger        = "logger"        // 日志来源函数名
	fieldType          = "type"          // 当前日志的所处动作环境，ACCESS|EVENT|RPC|OTHER
	fieldTitle         = "title"         // 日志标题，不传就是message前100个字符
	fieldPID           = "pid"           // 进程id
	fieldThreadID      = "threadId"      // 线程id
	fieldLanguage      = "language"      // 语言标识
	fieldURL           = "url"           // ⻚面/接口URL
	fieldClientIP      = "clientIp"      // 调用者IP
	fieldErrCode       = "errCode"       // 异常码
	fieldTraceID       = "traceID"       // 全链路TraceId
	fieldSpanID        = "spanID"        // 全链路SpanId :在非span产生的上下文环境中，可以留空
	fieldParentID      = "parentID"      // 全链路 上级SpanId :在非span产生的上下文环境中，可以留空
	fieldCustomLog1    = "customLog1"    // 自定义log1
	fieldCustomLog2    = "customLog2"    // 自定义log2
	fieldCustomLog3    = "customLog3"    // 自定义log3
)

const (
	LogNameDefault = "default"
	LogNameRedis   = "redis"
	LogNameMysql   = "mysql"
	LogNameMongodb = "mongodb"
	LogNameAPI     = "api"
	LogNameAo      = "ao"
	LogNameGRpc    = "grpc"
	LogNameEs      = "es"
	LogNameTmq     = "tmq"
	LogNameAmq     = "amq"
	LogNameLogic   = "logic"
	LogNameFile    = "file"
	LogNameNet     = "net"
)

var (
	logNameList = map[string]string{ // 日志分类
		LogNameRedis:   LogNameRedis,
		LogNameMysql:   LogNameMysql,
		LogNameMongodb: LogNameMongodb,
		LogNameAPI:     LogNameAPI,
		LogNameAo:      LogNameAo,
		LogNameGRpc:    LogNameGRpc,
		LogNameEs:      LogNameEs,
		LogNameTmq:     LogNameTmq,
		LogNameAmq:     LogNameAmq,
		LogNameLogic:   LogNameLogic,
		LogNameFile:    LogNameFile,
		LogNameNet:     LogNameNet,
	}
	fields logrus.Fields
)

type appHook struct {
}

func (hook *appHook) Fire(entry *logrus.Entry) error {
	entry.Data[fieldTime] = time.Now().Format(defaultTimestampFormat)
	entry.Data[fieldThreadID] = getGID()
	return nil
}

func (hook *appHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Init 初始化logger
func Init(opt *options.LogsOptions) {
	appName := os.Getenv("IDG_SERVICE_NAME")
	appID := os.Getenv("IDG_APPID")
	appVer := os.Getenv("IDG_VERSION")

	// 设置日志格式为json格式
	formatter := &logrus.JSONFormatter{
		DisableTimestamp: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: fieldMessage,
		},
	}
	logrus.SetFormatter(formatter)
	// 设置日志等级
	logrus.SetLevel(logrus.TraceLevel)
	logrus.AddHook(&appHook{})
	if strings.ToLower(opt.OutPut) == "stdout" {
		// 如果指定设置stdout则输出到终端,否则输出到msp redis
		logrus.SetOutput(os.Stdout)
	} else {
		flag, redisHost, redisPort := getMspLogRedis()
		if flag {
			logrus.SetOutput(ioutil.Discard)
			hookConfig := logredis.HookConfig{
				Host:   redisHost,
				Key:    "service_" + appID,
				Format: "origin",
				Port:   redisPort,
			}
			hook, err := logredis.NewHook(hookConfig)
			if err == nil {
				logrus.AddHook(hook)
			} else {
				log.Printf("logredis error: %q", err)
			}
		} else {
			// 获取redis host失败则输出到终端
			logrus.SetOutput(os.Stdout)
		}
	}

	fields = logrus.Fields{
		fieldAppName:       appName,
		fieldAppID:         appID,
		fieldAppVersion:    appVer,
		fieldAppKey:        "appkey",
		fieldChannel:       "1",
		fieldSubOrgKey:     "sub_org_key",
		fieldTime:          "",
		fieldLevel:         "",
		fieldHostName:      getHostname(),
		fieldIP:            getInternetIP(),
		fieldPodName:       os.Getenv("PODNAME"),
		fieldPodIP:         os.Getenv("PODIP"),
		fieldNodeName:      os.Getenv("NODENAME"),
		fieldNodeIP:        os.Getenv("NODEIP"),
		fieldContainerName: os.Getenv("CONTAINERNAME"),
		fieldClusterUID:    os.Getenv("IDG_CLUSTERUID"),
		fieldImageURL:      os.Getenv("IDG_IMAGEURL"),
		fieldUniqueID:      os.Getenv("IDG_UNIQUEID"),
		fieldSiteUID:       os.Getenv("IDG_SITEUID"),
		fieldRunEnvType:    os.Getenv("IDG_RUNTIME"),
		fieldMessage:       "",
		fieldLogger:        "",
		fieldType:          "ACCESS",
		fieldTitle:         "",
		fieldPID:           os.Getpid(),
		fieldLanguage:      "ch",
		fieldURL:           "",
		fieldClientIP:      "",
		fieldErrCode:       "",
		fieldTraceID:       "",
		fieldSpanID:        "",
		fieldParentID:      "",
		fieldCustomLog1:    "",
		fieldCustomLog2:    "",
		fieldCustomLog3:    "",
	}
}

func getLogName(logName string) string {
	if v, ok := logNameList[logName]; ok {
		return v
	} else {
		return LogNameDefault
	}
}

func LogDebugw(c *gin.Context, logName string, msg string) {
	logrus.WithFields(fields).WithFields(getFields(c)).Debug(fmt.Sprintf("%s : %s", getLogName(logName), msg))
}

func LogDebugf(c *gin.Context, logName string, template string, args ...interface{}) {
	logrus.WithFields(fields).WithFields(getFields(c)).Debugf(getLogName(logName)+":"+template, args...)
}

func LogInfow(c *gin.Context, logName string, msg string) {
	logrus.WithFields(fields).WithFields(getFields(c)).WithFields(getFields(c)).Info(fmt.Sprintf("%s : %s", getLogName(logName), msg))
}

func LogInfof(c *gin.Context, logName string, template string, args ...interface{}) {
	logrus.WithFields(fields).WithFields(getFields(c)).Infof(getLogName(logName)+":"+template, args...)
}

func LogWarnw(c *gin.Context, logName string, msg string) {
	logrus.WithFields(fields).WithFields(getFields(c)).Warn(fmt.Sprintf("%s : %s", getLogName(logName), msg))
}

func LogWarnf(c *gin.Context, logName string, template string, args ...interface{}) {
	logrus.WithFields(fields).WithFields(getFields(c)).Warnf(getLogName(logName)+":"+template, args...)
}

func LogError(c *gin.Context, logName string, msg string) {
	logrus.WithFields(fields).WithFields(getFields(c)).WithField(fieldLogger, traceFunc()).Error(fmt.Sprintf("%s : %s", getLogName(logName), msg))
}

func LogErrorw(c *gin.Context, logName string, msg string, err error) {
	logrus.WithFields(fields).WithFields(getFields(c)).WithField(fieldLogger, traceFunc()).Error(fmt.Sprintf("%s : %s, %s", getLogName(logName), msg, err.Error()))
}

func LogErrorf(c *gin.Context, logName string, template string, args ...interface{}) {
	logrus.WithFields(fields).WithFields(getFields(c)).WithField(fieldLogger, traceFunc()).Errorf(getLogName(logName)+":"+template, args...)
}

func LogInfoCustom(c *gin.Context, logName string, fields logrus.Fields, msg string) {
	logrus.WithFields(fields).WithFields(getFields(c)).WithFields(getFields(c)).WithFields(fields).Info(fmt.Sprintf("%s : %s", getLogName(logName), msg))
}

// getInternetIP 用于自动查找本机IP地址
func getInternetIP() (iP string) {
	// 查找本机IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				if ip4[0] == 10 {
					//  赋值新的IP
					iP = ip4.String()
				}
			}
		}
	}
	return
}

// getHostname 用于自动获取本机Hostname信息
func getHostname() (hostname string) {
	// 查找本机hostname
	hostname, err := os.Hostname()
	if err != nil {
		return ""
	}
	return
}

// 获取协程ID
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func traceFunc() string {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return fmt.Sprintf("%s: %d %s", frame.File, frame.Line, frame.Function)
}

func getFields(c *gin.Context) logrus.Fields {
	fields := logrus.Fields{}
	if c != nil {
		fields[fieldTraceID] = c.GetString("sw8")
		fields[fieldURL] = c.Request.Method + "： " + c.Request.URL.Path
	}
	return fields
}

func getMspLogRedis() (flag bool, host string, port int) {
	str := os.Getenv("MSP_LOG_REDIS_HOST")
	strArr := strings.Split(str, ":")
	if len(strArr) != 2 {
		return
	}
	port, err := strconv.Atoi(strArr[1])
	if err != nil {
		return
	}
	flag, host = true, strArr[0]
	return
}
