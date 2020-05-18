package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/rocinax/rigis/pkg/rigis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/net/netutil"
)

func init() {

	// *************** Server Default Setting ***************
	viper.SetDefault("ServerPort", 6443)
	viper.SetDefault("ServerName", "Rocinax Rigis Server")
	viper.SetDefault("ServerMaxConnections", 256)

	// pflag config
	pflag.String("config", "/opt/rocinax/rigis/config", "config: rocinax rigis config directory.")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// name of config file (without extension)
	viper.SetConfigName("rigis")
	viper.AddConfigPath(viper.GetString("config"))
	viper.SetConfigType("yaml")

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {

		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	switch viper.GetString("LogType") {
	case ("Stdout"):
		logrus.SetOutput(os.Stdout)
	case ("File"):

		var logFile *os.File
		var err error

		fileName := path.Join(
			viper.GetString("LogDir"),
			"rigis.log",
		)

		// file check
		_, statErr := os.Stat(fileName)
		if !os.IsNotExist(statErr) {
			logFile, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
		} else {
			logFile, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0644)
		}
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(logFile)
	default:
		logrus.SetOutput(os.Stdout)
	}

	switch viper.GetString("LogLevel") {
	case ("Trace"):
		logrus.SetLevel(logrus.TraceLevel)
	case ("Debug"):
		logrus.SetLevel(logrus.DebugLevel)
	case ("Info"):
		logrus.SetLevel(logrus.InfoLevel)
	case ("Warn"):
		logrus.SetLevel(logrus.WarnLevel)
	case ("Error"):
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
}

func main() {

	// define and get rigis configuration
	var config rigis.Config
	err := viper.UnmarshalKey("Rigis", &config)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}

	// port normalize
	var addr string
	if viper.GetInt("ServerPort") <= 1024 || viper.GetInt("ServerPort") > 65535 {
		logrus.WithFields(logrus.Fields{
			"type": "server",
			"app":  "rigis",
		}).Errorf("server port is out of range. :%d", viper.GetInt("ServerPort"))
		panic(fmt.Errorf("ServerPort is out of range. :%d", viper.GetInt("ServerPort")))
	}
	addr = ":" + strconv.Itoa(viper.GetInt("ServerPort"))

	// create rigis
	rgs := rigis.NewRigis(config)

	logrus.WithFields(logrus.Fields{
		"type": "server",
		"app":  "rigis",
	}).Infof("rocinax rigis is started. port: %s", addr)

	if viper.GetBool("TLSEnable") {

		mux := http.NewServeMux()
		mux.HandleFunc("/", rgs.ServeHTTP)

		tlsCertFiles := viper.GetStringSlice("TLSCertFiles")
		tlsPrivateKeyFiles := viper.GetStringSlice("TLSPrivateKeyFiles")

		if len(tlsCertFiles) <= 0 || len(tlsCertFiles) != len(tlsPrivateKeyFiles) {
			logrus.WithFields(logrus.Fields{
				"type": "server",
				"app":  "rigis",
			}).Errorf("TLSFiles is Not Defined. %d %d", len(tlsCertFiles), len(tlsPrivateKeyFiles))
			panic("TLSFiles is Not Defined.")
		}

		tlsConfig := &tls.Config{}

		for i := 0; i < len(tlsCertFiles); i++ {
			certificate, err := tls.LoadX509KeyPair(tlsCertFiles[i], tlsPrivateKeyFiles[i])
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"type": "server",
					"app":  "rigis",
				}).Fatal(err)
				panic(err)
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
		}

		tlsConfig.BuildNameToCertificate()

		ServerSSL := &http.Server{
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		listener, err := tls.Listen("tcp", addr, tlsConfig)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"type": "server",
				"app":  "rigis",
			}).Fatal(err)
		}

		// Set ServerMaxConnections
		limitListener := netutil.LimitListener(listener, viper.GetInt("ServerMaxConnections"))

		err = ServerSSL.Serve(limitListener)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"type": "server",
				"app":  "rigis",
			}).Fatal(err)
		}

	} else {

		// define handle func
		http.HandleFunc("/", rgs.ServeHTTP)

		// run rigis server
		http.ListenAndServe(addr, nil)
	}
}
