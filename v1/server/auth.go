package server

import (
	"fmt"
	"time"
	"strings"
	fiber "github.com/gofiber/fiber/v2"
	// rate_limiter "github.com/gofiber/fiber/v2/middleware/limiter"
	bcrypt "golang.org/x/crypto/bcrypt"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

// https://github.com/gofiber/fiber/blob/main/middleware/encryptcookie/utils.go#L16

func ( s *Server ) ValidateAdmin( context *fiber.Ctx ) ( result bool ) {
	result = false
	admin_cookie := context.Cookies( s.Config.Cookie.Admin.Name )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( s.Config.Creds.EncryptionKey , admin_cookie )
		if admin_cookie_value == s.Config.Cookie.Admin.Message {
			result = true
			return
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == s.Config.Creds.APIKey {
			result = true
			return
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == s.Config.Creds.APIKey {
			result = true
			return
		}
	}
	return
}

func ( s *Server ) ValidateAdminMW( context *fiber.Ctx ) ( error ) {
	admin_cookie := context.Cookies( s.Config.Cookie.Admin.Name )
	if admin_cookie != "" {
		admin_cookie_value := encryption.SecretBoxDecrypt( s.Config.Creds.EncryptionKey , admin_cookie )
		if admin_cookie_value == s.Config.Cookie.Admin.Message {
			return context.Next()
		}
	}
	admin_api_key_header := context.Get( "key" )
	if admin_api_key_header != "" {
		if admin_api_key_header == s.Config.Creds.APIKey {
			return context.Next()
		}
	}
	admin_api_key_query := context.Query( "k" )
	if admin_api_key_query != "" {
		if admin_api_key_query == s.Config.Creds.APIKey {
			return context.Next()
		}
	}
	ip_address := context.Get( "x-forwarded-for" )
	if ip_address == "" { ip_address = context.IP() }
	log.Debug( fmt.Sprintf( "%s === %s === %s === %s" , ip_address , context.Method() , context.Path() , "NL" ) )
	return context.Status( fiber.StatusUnauthorized ).SendString( "why" )
}

func ( s *Server ) ValidateAdminLoginCredentials( context *fiber.Ctx ) ( result bool ) {
	result = false
	uploaded_username := context.FormValue( "username" )
	if uploaded_username == "" { fmt.Println( "username empty" ); return }
	if uploaded_username != s.Config.Creds.AdminUsername { fmt.Println( "username not correct" ); return }
	uploaded_password := context.FormValue( "password" )
	if uploaded_password == "" { fmt.Println( "password empty" ); return }
	fmt.Println( "uploaded_username ===" , uploaded_username )
	fmt.Println( "uploaded_password ===" , uploaded_password )
	password_matches := bcrypt.CompareHashAndPassword( []byte( uploaded_password ) , []byte( s.Config.Creds.AdminPassword ) )
	if password_matches != nil { fmt.Println( "bcrypted password doesn't match" ); return }
	fmt.Println( "password matched" )
	result = true
	return
}

// POST http://localhost:5950/admin/login
func ( s *Server ) AdminLogin( context *fiber.Ctx ) ( error ) {
	valid_login := s.ValidateAdminLoginCredentials( context )
	if valid_login == false { return s.RenderFailedLogin( context ) }
	host := context.Hostname()
	domain := strings.Split( host , ":" )[ 0 ] // setting this leaks url-prefix and locks to specific domain
	context.Cookie(
		&fiber.Cookie{
			Name: s.Config.Cookie.Admin.Name ,
			Value: encryption.SecretBoxEncrypt( s.Config.Creds.EncryptionKey , s.Config.Cookie.Admin.Message ) ,
			Secure: true ,
			Path: "/" ,
			Domain: domain ,
			HTTPOnly: true ,
			SameSite: "Lax" ,
			Expires: time.Now().AddDate( 10 , 0 , 0 ) , // aka 10 years from now
		} ,
	)
	return context.Redirect( "/" )
}

func ( s *Server ) AdminLogout( context *fiber.Ctx ) ( error ) {
	context.Cookie( &fiber.Cookie{
		Name: s.Config.Cookie.Admin.Name ,
		Value: "" ,
		Expires: time.Now().Add( -time.Hour ) , // set the expiration to the past
		HTTPOnly: true ,
		Secure: true ,
	})
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>Logged Out</h1>" )
}

func ( s *Server ) RenderFailedLogin( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "<h1>no</h1>" )
}

func ( s *Server ) LoginPage( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendFile( "./v1/server/html/login.html" )
}