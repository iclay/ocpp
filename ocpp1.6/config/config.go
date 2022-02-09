package config

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

//GConf The label name must be consistent with the configuration file, otherwise it cannot be resolved
type GConf struct {
	Include           string   `label:"include" parse_func:"parse_file"`
	ServiceAddr       string   `label:"service_addr"`
	ServiceURI        string   `label:"service_uri"`
	WsEnable          bool     `label:"ws_enable" parse_func:"parse_bool"`
	WsPort            int      `label:"ws_port"`
	WssEnable         bool     `label:"wss_enable" parse_func:"parse_bool"`
	WssPort           int      `label:"wss_port"`
	TLSCertificate    string   `label:"tls_cert"`
	TLSCertificateKey string   `label:"tls_key"`
	HeartbeatTimeout  int      `label:"heartbeat_timeout"`
	ETCDList          []string `label:"etcd_list" parse_func:"parse_string_list"`
	ETCDBasePath      string   `label:"etcd_base_path"`
	RPCAddress        string   `label:"rpc_addr"`
	UsePool           bool     `label:"use_pool" parse_func:"parse_bool"`
	UseEpoll          bool     `label:"use_epoll" parse_func:"parse_bool"`
	LogPath           string   `label:"log_path"`
	LogLevel          string   `label:"log_level"` // trace, debug, info, warn[ing], error, fatal, panic
	LogMaxDiskUsage   int64    `label:"log_max_disk_usage" parse_func:"parse_bytes"`
	LogMaxFileNum     int64    `label:"log_max_file_num" parse_func:"parse_bytes"`
}

var (
	GCONF     GConf
	GConfItem = &GCONF
)

func parseBool(value string) int {
	if value == "yes" || value == "on" || value == "1" {
		return 1
	}
	return 0
}

//parse string like 2B, 1M, 1G to bytes
func parseAsBytes(value string) int64 {
	if len(value) == 0 {
		return 0
	}

	last := value[len(value)-1]
	if last >= '0' && last <= '9' {
		i, _ := strconv.ParseInt(value, 10, 64)
		return i
	}
	first := value[:len(value)-1]
	i, _ := strconv.ParseInt(first, 10, 64)
	switch last {
	case 'b':
		return i / 8
	case 'B':
		return i
	case 'k', 'K':
		return i * 1024
	case 'M', 'm':
		return i * 1024 * 1024
	case 'G', 'g':
		return i * 1024 * 1024 * 1024
	}
	return i
}

func parseStringList(value string) []string {
	return strings.Split(value, ",")
}

//ParseFile Parse File
func ParseFile(filePath string) {
	confile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer confile.Close()

	lineNum := 0

	scanner := bufio.NewScanner(confile)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if len(line) < 3 || line[0] == '#' || line[0] == '[' {
			continue
		}

		var key string
		var value []string
		v := strings.Fields(line)
		if len(v) < 2 {
			fmt.Println("error parseing file ", filePath, "line", lineNum, line)
			panic(err)
		}
		key = v[0]
		value = v[1:]

		object := reflect.ValueOf(GConfItem)
		myref := object.Elem()
		typeOfType := myref.Type()
		for i := 0; i < myref.NumField(); i++ {
			fieldInfo := myref.Type().Field(i)
			tag := fieldInfo.Tag
			tagName := tag.Get("label")
			parseFunc := tag.Get("parse_func")
			variableName := typeOfType.Field(i).Name

			if strings.Compare(string(key), tagName) == 0 {
				switch parseFunc {
				case "":
					if myref.FieldByName(variableName).Kind() == reflect.Int {
						*(*int)(unsafe.Pointer(myref.FieldByName(variableName).Addr().Pointer())), _ = strconv.Atoi(value[0])
					} else if myref.FieldByName(variableName).Kind() == reflect.Int64 {
						*(*int64)(unsafe.Pointer(myref.FieldByName(variableName).Addr().Pointer())), _ = strconv.ParseInt(value[0], 10, 64)
					} else {
						*(*string)(unsafe.Pointer(myref.FieldByName(variableName).Addr().Pointer())) = value[0]
					}
				case "parse_file":
					ParseFile(value[0])
				case "parse_bool":
					*(*int)(unsafe.Pointer(myref.FieldByName(variableName).Addr().Pointer())) = parseBool(value[0])
				case "parse_string_list":
					*(*[]string)(unsafe.Pointer(myref.FieldByName(variableName).Addr().Pointer())) = parseStringList(value[0])
				case "parse_bytes":
					*(*int64)(unsafe.Pointer(myref.FieldByName(variableName).Addr().Pointer())) = parseAsBytes(value[0])

				}

				break
			}
		}
	}
}

func printInterfaceDepth(depth int, inter interface{}) {
	v := reflect.ValueOf(inter)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		for j := 0; j < depth; j++ {
			fmt.Printf("\t")
		}
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			fmt.Printf("[%s]\n", t.Field(i).Name)
			printInterfaceDepth(depth+1, f.Interface())
		} else {
			fmt.Printf("%s %s = %v\n", t.Field(i).Name, f.Type(), f.Interface())
		}
	}
}

func printInterface(inter interface{}) {
	printInterfaceDepth(0, inter)
}

//Print print the structure content for debugging
func Print() {
	printInterface(*GConfItem)
}
