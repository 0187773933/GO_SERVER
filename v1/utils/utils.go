package utils

import (
	"fmt"
	"os"
	"net"
	hex "encoding/hex"
	"encoding/base64"
	filepath "path/filepath"
	yaml "gopkg.in/yaml.v3"
	ioutil "io/ioutil"
	runtime "runtime"
	types "github.com/0187773933/GO_SERVER/v1/types"
	// fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/encryption/v1/encryption"
)

var CONFIG_PATH string = ""

func SetupStackTraceReport() {
	if r := recover(); r != nil {
		stacktrace := make( []byte , 1024 )
		runtime.Stack( stacktrace , true )
		fmt.Printf( "%s\n" , stacktrace )
	}
}

func GenerateNewKeys() {
	admin_login_url := encryption.GenerateRandomString( 16 )
	admin_prefix := encryption.GenerateRandomString( 6 )
	login_url := encryption.GenerateRandomString( 16 )
	prefix := encryption.GenerateRandomString( 6 )
	// https://github.com/gofiber/fiber/blob/main/middleware/encryptcookie/utils.go#L91
	// https://github.com/0187773933/encryption/blob/master/v1/encryption/encryption.go#L46
	// cookie_secret := fiber_cookie.GenerateKey()
	cookie_secret_bytes := encryption.GenerateRandomBytes( 32 )
	cookie_secret := base64.StdEncoding.EncodeToString( cookie_secret_bytes )
	cookie_secret_message := encryption.GenerateRandomString( 16 )
	admin_cookie_secret_message := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	api_key := encryption.GenerateRandomString( 16 )
	encryption_key := encryption.GenerateRandomString( 32 )
	bolt_name := encryption.GenerateRandomString( 6 ) + ".db"
	bolt_prefix := encryption.GenerateRandomString( 6 )
	redis_prefix := encryption.GenerateRandomString( 6 )
	log_name := encryption.GenerateRandomString( 6 ) + ".db"
	log_key := encryption.GenerateRandomString( 6 )
	log_encryption_key := encryption.GenerateRandomString( 32 )
	kyber_private , kyber_public := encryption.KyberGenerateKeyPair()
	kyber_private_string := hex.EncodeToString( kyber_private[ : ] )
	kyber_public_string := hex.EncodeToString( kyber_public[ : ] )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tURL - Admin Login === %s\n" , admin_login_url )
	fmt.Printf( "\tURL - Admin Prefix === %s\n" , admin_prefix )
	fmt.Printf( "\tURL - Login === %s\n" , login_url )
	fmt.Printf( "\tURL - Prefix === %s\n" , prefix )
	fmt.Printf( "\tCOOKIE - Secret === %s\n" , cookie_secret )
	fmt.Printf( "\tCOOKIE - USER - Message === %s\n" , cookie_secret_message )
	fmt.Printf( "\tCOOKIE - ADMIN - Message === %s\n" , admin_cookie_secret_message )
	fmt.Printf( "\tCREDS - Admin Username === %s\n" , admin_username )
	fmt.Printf( "\tCREDS - Admin Password === %s\n" , admin_password )
	fmt.Printf( "\tCREDS - API Key === %s\n" , api_key )
	fmt.Printf( "\tCREDS - Encryption Key === %s\n" , encryption_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n" , admin_password )
	fmt.Printf( "\tLOG - Log Name === %s\n" , log_name )
	fmt.Printf( "\tLOG - Log Key === %s\n" , log_key )
	fmt.Printf( "\tLOG - Encryption Key === %s\n" , log_encryption_key )
	fmt.Printf( "\tBOLT - Name === %s\n" , bolt_name )
	fmt.Printf( "\tBOLT - Prefix === %s\n" , bolt_prefix )
	fmt.Printf( "\tREDIS - Prefix === %s\n" , redis_prefix )
	fmt.Printf( "\tKYBER - Private Key === %s\n" , kyber_private_string )
	fmt.Printf( "\tKYBER - Public Key === %s\n" , kyber_public_string )
	panic( "Exiting" )
}

func WriteConfig( config *types.Config ) {
	fmt.Println( "Writing Config" , CONFIG_PATH )
	config_file , _ := yaml.Marshal( &config )
	ioutil.WriteFile( CONFIG_PATH , config_file , 0644 )
}

func GenerateNewKeysWrite( config *types.Config ) {
	x := config
	x.URLS.AdminLogin = encryption.GenerateRandomString( 16 )
	x.URLS.AdminPrefix = encryption.GenerateRandomString( 6 )
	x.URLS.Login = encryption.GenerateRandomString( 16 )
	x.URLS.Prefix = encryption.GenerateRandomString( 6 )
	cookie_secret_bytes := encryption.GenerateRandomBytes( 32 )
	x.Cookie.Secret = base64.StdEncoding.EncodeToString( cookie_secret_bytes )
	x.Cookie.User.Message = encryption.GenerateRandomString( 16 )
	x.Cookie.Admin.Message = encryption.GenerateRandomString( 16 )
	x.Creds.AdminUsername = encryption.GenerateRandomString( 16 )
	x.Creds.AdminPassword = encryption.GenerateRandomString( 16 )
	x.Creds.APIKey = encryption.GenerateRandomString( 16 )
	x.Creds.EncryptionKey = encryption.GenerateRandomString( 32 )
	kyber_private , kyber_public := encryption.KyberGenerateKeyPair()
	kyber_private_string := hex.EncodeToString( kyber_private[ : ] )
	kyber_public_string := hex.EncodeToString( kyber_public[ : ] )
	x.Creds.Kyber.Private = kyber_private_string
	x.Creds.Kyber.Public = kyber_public_string
	x.Bolt.Prefix = encryption.GenerateRandomString( 6 )
	x.Bolt.Path = encryption.GenerateRandomString( 6 ) + ".db"
	x.Redis.Prefix = encryption.GenerateRandomString( 6 )
	x.Log.LogKey = encryption.GenerateRandomString( 6 )
	x.Log.BoltDBPath = encryption.GenerateRandomString( 6 ) + ".db"
	x.Log.EncryptionKey = encryption.GenerateRandomString( 32 )
	fmt.Println( x )
	fmt.Println( CONFIG_PATH )
	WriteConfig( x )
	panic( "Exiting" )
}

func GetLocalIPAddresses() ( ip_addresses []string ) {
	host , _ := os.Hostname()
	addrs , _ := net.LookupIP( host )
	encountered := make( map[ string ]bool )
	for _ , addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			ip := ipv4.String()
			if !encountered[ ip ] {
				encountered[ ip ] = true
				ip_addresses = append( ip_addresses , ip )
			}
		}
	}
	return
}

func ParseConfig( file_path string ) ( result types.Config ) {
	config_file , _ := ioutil.ReadFile( file_path )
	error := yaml.Unmarshal( config_file , &result )
	if error != nil { panic( error ) }
	return
}

func ParseConfigGeneric() ( result types.ConfigGeneric ) {
	config_file , _ := ioutil.ReadFile( CONFIG_PATH )
    yaml.Unmarshal( []byte( config_file ) , &result )
	return
}

func GetConfig() ( result types.Config ) {
	if len( os.Args ) > 1 {
		CONFIG_PATH , _ = filepath.Abs( os.Args[ 1 ] )
	} else {
		CONFIG_PATH , _ = filepath.Abs( "./SAVE_FILES/config.yaml" )
		if _ , err := os.Stat( CONFIG_PATH ); os.IsNotExist( err ) {
			panic( "Config File Not Found" )
		}
	}
	CONFIG_PATH , _ = filepath.Abs( CONFIG_PATH )
	result = ParseConfig( CONFIG_PATH )
	result.Bolt.Path = filepath.Join( result.SaveFilesPath , result.Bolt.Path )
	result.Log.BoltDBPath = filepath.Join( result.SaveFilesPath , result.Log.BoltDBPath )
	return
}