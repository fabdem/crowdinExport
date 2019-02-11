package main

//	F.Demurger 2019-02
//  	3 args: key, project name and zip file name
//
//      Option -v version
//      Option -b to build the project
//            Optionnaly build the project and download the zip with all languages. 
//      Option -p <proxy url> to use a proxy. 
//
//      To be noted: comm timeout is 40s when not using a proxy and 300s with proxy. The no-proxy timeout is hard coded in the lib :(
//      Returns 1 if there was an error
//      If option built is used, returns "built" or "skipped" if the command is successful and depending if the build was actually done.
//
//       
//	cross compilation AMD64:  env GOOS=windows GOARCH=amd64 go build crowdinExport.go


import (
	"flag"
	"fmt"
	"os"
	"github.com/fabdem/go-crowdinproxy"
	//"go-crowdinproxy"
	"github.com/medisafe/go-crowdin"
    "time"
)

var idx int = 0
var finishChan chan struct{}

func animation() {    
    sequence := [...]string {"\b|","\b/","\b-","\b\\"}
    // sequence := [...]string {" 1"," 2"," 3"," 4"}
                            
    for {
        select {
            default: 
                fmt.Printf("%s",sequence[idx])
                idx = (idx + 1) % len(sequence) 
                amt := time.Duration(100)
                time.Sleep(time.Millisecond * amt) 
                
            case <-finishChan:
                return
        }
    }
}

func main() {

	var versionFlg bool
	var buildFlg bool
	var proxy string
	

	const usageVersion   = "Display Version"
	const usageBuild   = "Request a build"
	const usageProxy   = "Use a proxy followed with url"
    
    // Have to create a spbyecific set, the default one is poluted by some test stuff from another lib (?!) 
    checkFlags := flag.NewFlagSet("check", flag.ExitOnError)
    
	checkFlags.BoolVar(&versionFlg, "version", false, usageVersion)
	checkFlags.BoolVar(&versionFlg, "v", false, usageVersion + " (shorthand)")
	checkFlags.BoolVar(&buildFlg, "build", false, usageBuild)
	checkFlags.BoolVar(&buildFlg, "b", false, usageBuild + " (shorthand)")
	checkFlags.StringVar(&proxy, "proxy", "", usageProxy)
	checkFlags.StringVar(&proxy, "p", "", usageProxy + " (shorthand)")
	checkFlags.Usage = func() {
        fmt.Printf("Usage: %s <opt> <key> <project name> <path and name of zip>\n",os.Args[0])
        checkFlags.PrintDefaults()
    }

    // Check parameters
	checkFlags.Parse(os.Args[1:])
	
	if versionFlg {
		fmt.Printf("Version %s\n", "2019-02  v1.1.0")
		os.Exit(0)
	}
	
	proxyFlg := ( len(proxy) > 0 )

	// Check syntax
	// 4 crowdinExport <key> <proj_name> <path>
	// 5 crowdinExport -b <key> <proj_name> <path>
	// 6 crowdinExport -p <proxy> <key> <proj_name> <path>
	// 7 crowdinExport -b -p <proxy> <key> <proj_name> <path>
	switch nbArgs := len(os.Args); {
        case nbArgs <= 3: 
            checkFlags.Usage()  // Display usage
            fmt.Printf("Missing parameters\n")
            os.Exit(1)
        case nbArgs == 4:
            if buildFlg || proxyFlg {
                checkFlags.Usage()  // Display usage
                fmt.Printf("Missing parameters\n")
                os.Exit(1)
            }
        case nbArgs == 5:
            if proxyFlg {
                checkFlags.Usage()  // Display usage
                fmt.Printf("Missing parameters\n")
                os.Exit(1)
            }
        case nbArgs == 6:
            if buildFlg || !proxyFlg {
                checkFlags.Usage()  // Display usage
                fmt.Printf("Invalid or too many parameters: %d\n",nbArgs)
                os.Exit(1)
            }			
        case nbArgs == 7:
            if !buildFlg || !proxyFlg {
                checkFlags.Usage()  // Display usage
                fmt.Printf("Invalid or too many parameters: %d\n",nbArgs)
                os.Exit(1)
            }
    }
        
	
    // Parse the command parameters
    index := 0
	if buildFlg {
        index++
    }
	if proxyFlg {
        index += 2
    }
    key := os.Args[1 + index]
    project := os.Args[2 + index]
    filename :=  os.Args[3 + index]

    // fmt.Printf("proxyFlg=%s\n",proxyFlg)
    // fmt.Printf("buildFlg=%s\n",buildFlg)
    // fmt.Printf("proxy=%s\n",proxy)
    // fmt.Printf("key=%s\n",key)
    // fmt.Printf("project=%s\n",project)
    // fmt.Printf("filename=%s\n",filename)
    // os.Exit(1)
    
    // Create a connection with or without proxy
    var err error
    var api *crowdin.Crowdin
    if proxyFlg {
        crowdinproxy.SetTimeouts(5, 300) // insec
        api,err = crowdinproxy.New(key, project, proxy)
    } else {
        api = crowdin.New(key, project)       
    }
    if err !=nil {
        fmt.Printf("\ncrowdinExport() - connection problem %s\n",err)
        os.Exit(1)
    }
    
    //api.SetDebug(true, nil)
    finishChan = make(chan struct{})
    go animation()

    //time.Sleep(time.Millisecond * 5000)                   

    var result string
    
    if buildFlg {
                
        // Request a build
        response,err := api.ExportTranslations()
        
        if err !=nil {
            fmt.Printf("\ncrowdinExport() build request error %s\n",err)
            os.Exit(1)
        }
        
        result = response.Success.Status
        
        // If there is no build necessary let's do a download anyway
        if response.Success.Status == "skipped" { 
            buildFlg = false
        }
    } 

    if !buildFlg {
        opt := crowdin.DownloadOptions{Package: "all", LocalPath: filename}

        // request zip download
        err := api.DownloadTranslations(&opt)
        if err !=nil {
            fmt.Printf("\ncrowdinExport() download error %s\n",err)
            os.Exit(1)
        }
    }

    close(finishChan)  // Stop animation

    // Return "built" or "skipped"
    fmt.Printf("\b%s\r\n",result)
}
