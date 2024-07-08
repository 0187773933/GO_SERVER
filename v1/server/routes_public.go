package server

import (
	"fmt"
	"time"
	// "strings"
	fiber "github.com/gofiber/fiber/v2"
	rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	// bcrypt "golang.org/x/crypto/bcrypt"
	// encryption "github.com/0187773933/encryption/v1/encryption"
	// try "github.com/manucorporat/try"
)

func CDNMaxedOut( c *fiber.Ctx ) error {
	ip_address := c.IP()
	log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
	log.Info( log_message )
	c.Set( "Content-Type" , "text/html" )
	return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6000);</script></html>" )
}

var CDNLimter = rate_limiter.New( rate_limiter.Config{
	Max: 6 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: CDNMaxedOut ,
	LimiterMiddleware: rate_limiter.SlidingWindow{} ,
})

func PublicMaxedOut( c *fiber.Ctx ) error {
	ip_address := c.IP()
	log_message := fmt.Sprintf( "%s === %s === %s === PUBLIC RATE LIMIT REACHED !!!" , ip_address , c.Method() , c.Path() );
	log.Info( log_message )
	c.Set( "Content-Type" , "text/html" )
	return c.SendString( "<html><h1>loading ...</h1><script>setTimeout(function(){ window.location.reload(1); }, 6000);</script></html>" )
}

var PublicLimter = rate_limiter.New( rate_limiter.Config{
	Max: 3 ,
	Expiration: 1 * time.Second ,
	KeyGenerator: func( c *fiber.Ctx ) string {
		return c.Get( "x-forwarded-for" )
	} ,
	LimitReached: PublicMaxedOut ,
	LimiterMiddleware: rate_limiter.SlidingWindow{} ,
})

func ( s *Server ) RenderHomePage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	admin_logged_in := s.ValidateAdmin( context )
	if admin_logged_in == true {
		// fmt.Println( "RenderHomePage() --> Admin" )
		return context.SendFile( "./v1/server/html/admin.html" )
	}
	return context.SendFile( "./v1/server/html/home.html" )
}

func ( s *Server ) SetupPublicRoutes() {
	home_url := "/"
	// login_url := "/login"
	// logout_url := "/logout"
	admin_login_url := "/admin/login"
	admin_logout_url := "/admin/logout"
	if s.Config.URLS.Prefix != "" {
		home_url = fmt.Sprintf( "/%s" , s.Config.URLS.Prefix )
		// login_url = fmt.Sprintf( "/%s/login" , s.Config.URLS.Prefix )
		// logout_url = fmt.Sprintf( "/%s/logout" , s.Config.URLS.Prefix )
	}
	if s.Config.URLS.AdminPrefix != "" {
		if s.Config.URLS.AdminLogin != "" {
			admin_login_url = fmt.Sprintf( "/%s/%s" , s.Config.URLS.AdminPrefix , s.Config.URLS.AdminLogin )
		} else {
			admin_login_url = fmt.Sprintf( "/%s/login" , s.Config.URLS.AdminPrefix )
		}
		admin_logout_url = fmt.Sprintf( "/%s/logout" , s.Config.URLS.AdminPrefix )
	}
	s.FiberApp.Get( home_url , PublicLimter , s.RenderHomePage )
	// s.FiberApp.Get( login_url , PublicLimter , s.LoginPage )
	// s.FiberApp.Post( login_url , PublicLimter , s.Login )
	// s.FiberApp.Get( logout_url , PublicLimter , s.Logout )
	s.FiberApp.Get( admin_login_url , PublicLimter , s.LoginPage )
	s.FiberApp.Post( admin_login_url , PublicLimter , s.AdminLogin )
	s.FiberApp.Get( admin_logout_url , PublicLimter , s.AdminLogout )
}