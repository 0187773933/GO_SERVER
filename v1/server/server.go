package server

import (
	"fmt"
	"time"
	"strings"
	bolt "github.com/boltdb/bolt"
	// logrus "github.com/sirupsen/logrus"
	logger "github.com/0187773933/Logger/v1/logger"
	types "github.com/0187773933/BLANK_SERVER/v1/types"
	utils "github.com/0187773933/BLANK_SERVER/v1/utils"
	fiber "github.com/gofiber/fiber/v2"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	fiber_cors "github.com/gofiber/fiber/v2/middleware/cors"
	fiber_favicon "github.com/gofiber/fiber/v2/middleware/favicon"
	// bolt "github.com/boltdb/bolt"
)

type Server struct {
	FiberApp *fiber.App `yaml:"fiber_app"`
	Config *types.Config `yaml:"config"`
	Location *time.Location `yaml:"-"`
	DB *bolt.DB `yaml:"-"`
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
	server.DB = db
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
	server.SetupPublicRoutes()
	server.SetupAdminRoutes()
	return
}