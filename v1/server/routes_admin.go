package server

import (
	"fmt"
	"strconv"
	// "io/fs"
	// "net/http"
	// net_url "net/url"
	// bolt_api "github.com/boltdb/bolt"
	// encryption "github.com/0187773933/encryption/v1/encryption"
	fiber "github.com/gofiber/fiber/v2"
	// logger "github.com/0187773933/Logger/v1/logger"
)

func ( s *Server ) GetLogMessages( c *fiber.Ctx ) ( error ) {
	count := c.Query( "count" )
	if count == "" {
		count = c.Query( "c" )
	}
	count_int , _ := strconv.Atoi( count )
	if count_int == 0 {
		count_int = -1
	}
	log.Debug( fmt.Sprintf( "Count === %d" , count_int ) )
	messages := log.GetMessages( count_int )
	return c.JSON( fiber.Map{
		"result": true ,
		"url": "/log/:count" ,
		"count": count ,
		"messages": messages ,
	})
}

// func ( s *Server ) ExampleGet( c *fiber.Ctx ) ( error ) {
// 	id := c.Params( "id" )
// 	extra := c.Query( "extra" )
// 	return c.JSON( fiber.Map{
// 		"result": true ,
// 		"url": "/example/:id" ,
// 		"id": id ,
// 		"extra": extra ,
// 	})
// }

// func ( s *Server ) ExamplePost( c *fiber.Ctx ) ( error ) {
// 	var json_body map[ string ]interface{}
// 	c.BodyParser( &json_body )
// 	return c.JSON( fiber.Map{
// 		"result": true ,
// 		"url": "/example/:id" ,
// 		"data": json_body ,
// 	})
// }


func ( s *Server ) SetupAdminRoutes() {
	cdn_group := s.FiberApp.Group( "/cdn" )
	cdn_group.Use( CDNLimter )
	cdn_group.Use( s.ValidateAdminMW )
	// s.FiberApp.Static( "/cdn" , cdn_fs )
	cdn_group.Use( "/" , s.StaticHandler( "/cdn" , CDNFilesFS ) )
	var admin fiber.Router
	if s.Config.URLS.AdminPrefix == "" {
		admin = s.FiberApp.Group( "/admin" )
	} else {
		admin = s.FiberApp.Group( fmt.Sprintf( "/%s" , s.Config.URLS.AdminPrefix ) )
	}
	admin.Use( s.ValidateAdminMW )
	admin.Get( "/log/view" , s.GetLogMessages )
	// admin.Get( "/example/:id" , s.ExampleGet )
	// admin.Post( "/example" , s.ExamplePost )
}