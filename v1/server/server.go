package server

import (
	"embed"
	"fmt"
	fs "io/fs"
	"time"
	"strings"
	// "strconv"
	bolt "github.com/boltdb/bolt"
	redis "github.com/redis/go-redis/v9"
	// logrus "github.com/sirupsen/logrus"
	logger "github.com/0187773933/Logger/v1/logger"
	types "github.com/0187773933/GO_SERVER/v1/types"
	utils "github.com/0187773933/GO_SERVER/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	fiber_favicon "github.com/gofiber/fiber/v2/middleware/favicon"
	// bolt "github.com/boltdb/bolt"
)

//go:embed cdn/*
var CDNFiles embed.FS

//go:embed html/*
var HTMLFiles embed.FS

var CDNFilesFS fs.FS
var HTMLFilesFS fs.FS

var ADMIN_HTML_FILE fs.File
var HOME_HTML_FILE fs.File
var LOGIN_HTML_FILE fs.File

var ADMIN_HTML_FILE_SIZE int
var HOME_HTML_FILE_SIZE int
var LOGIN_HTML_FILE_SIZE int

type Server struct {
	FiberApp *fiber.App `yaml:"fiber_app"`
	Config *types.Config `yaml:"config"`
	Location *time.Location `yaml:"-"`
	DB *bolt.DB `yaml:"-"`
	REDIS *redis.Client `yaml:"-"`
	LOG *logger.Wrapper `yaml:"-"`
}

var log *logger.Wrapper

func ( s *Server ) LogRequest( context *fiber.Ctx ) ( error ) {
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	c_method := context.Method()
	c_path := context.Path()
	if s.Config.URLS.Prefix == "" {
		if strings.HasPrefix( c_path , "/favicon" ) {
			return context.Next()
		}
	} else {
		if strings.HasPrefix( c_path , fmt.Sprintf( "/%s/favicon" , s.Config.URLS.Prefix ) ) {
			return context.Next()
		}
	}
	log_message := fmt.Sprintf( "%s === %s === %s" , ip_address , c_method , c_path )
	log.Info( log_message )
	return context.Next()
}

func ( s *Server ) Start() {
	if s.Config.URLS.AdminPrefix == "" {
		if s.Config.URLS.AdminLogin == "" {
			fmt.Printf( "Admin Login @ http://localhost:%s/login\n" , s.Config.Port )
		} else {
			fmt.Printf( "Admin Login @ http://localhost:%s/%s\n" , s.Config.Port , s.Config.URLS.AdminLogin )
		}
	} else {
		if s.Config.URLS.AdminLogin == "" {
			fmt.Printf( "Admin Login @ http://localhost:%s/%s/login\n" , s.Config.Port , s.Config.URLS.AdminPrefix )
		} else {
			fmt.Printf( "Admin Login @ http://localhost:%s/%s/%s\n" , s.Config.Port , s.Config.URLS.AdminPrefix , s.Config.URLS.AdminLogin )
		}
	}
	fmt.Printf( "Admin Username === %s\n" , s.Config.Creds.AdminUsername )
	fmt.Printf( "Admin Password === %s\n" , s.Config.Creds.AdminPassword )
	fmt.Printf( "Admin API Key === %s\n" , s.Config.Creds.APIKey )
	local_ip_addresses := utils.GetLocalIPAddresses()
	for _ , ip_address := range local_ip_addresses {
		fmt.Printf( "Listening @ http://%s:%s\n" , ip_address , s.Config.Port )
	}
	listen_address := fmt.Sprintf( ":%s" , s.Config.Port )
	log.Info( fmt.Sprintf( "%s Ready" , s.Config.Name ) )
	s.FiberApp.Listen( listen_address )
}

func New( config *types.Config , w_log *logger.Wrapper , db *bolt.DB ) ( server Server ) {
	server.Location , _ = time.LoadLocation( config.TimeZone )
	server.FiberApp = fiber.New()
	server.Config = config
	log = w_log
	server.LOG = w_log
	server.DB = db
	if config.Redis.Enabled == true {
		server.REDIS = redis.NewClient( &redis.Options{
			Addr: fmt.Sprintf( "%s:%s" , config.Redis.Host , config.Redis.Port ) ,
			Password: config.Redis.Password ,
			DB: config.Redis.Number ,
		})
		log.Info( "Redis Connected" )
	}
	server.FiberApp.Use( server.LogRequest )
	server.FiberApp.Use( fiber_favicon.New() )
	server.FiberApp.Use( fiber_cookie.New( fiber_cookie.Config{
		Key: server.Config.Cookie.Secret ,
	}))
	allow_origins_string := strings.Join( config.AllowOrigins , "," )
	server.FiberApp.Use( fiber_cors.New( fiber_cors.Config{
		AllowOrigins: allow_origins_string ,
		AllowHeaders:  "Origin, Content-Type, Accept, key, k" ,
	}))
	CDNFilesFS , _ = fs.Sub( CDNFiles , "cdn" )
	HTMLFilesFS , _ = fs.Sub( HTMLFiles , "html" )
	ADMIN_HTML_FILE , _ = HTMLFilesFS.Open( "admin.html" )
	HOME_HTML_FILE , _ = HTMLFilesFS.Open( "home.html" )
	LOGIN_HTML_FILE , _ = HTMLFilesFS.Open( "login.html" )
	defer ADMIN_HTML_FILE.Close()
	defer HOME_HTML_FILE.Close()
	defer LOGIN_HTML_FILE.Close()
	ADMIN_HTML_FILE_INFO , _ := ADMIN_HTML_FILE.Stat()
	HOME_HTML_FILE_INFO , _ := HOME_HTML_FILE.Stat()
	LOGIN_HTML_FILE_INFO , _ := LOGIN_HTML_FILE.Stat()
	ADMIN_HTML_FILE_SIZE = int( ADMIN_HTML_FILE_INFO.Size() )
	HOME_HTML_FILE_SIZE = int( HOME_HTML_FILE_INFO.Size() )
	LOGIN_HTML_FILE_SIZE = int( LOGIN_HTML_FILE_INFO.Size() )

	server.SetupPublicRoutes()
	server.SetupAdminRoutes()
	return
}