package main

//	F.Demurger 2019-01
//  	3 args: key, project name and zip file name
//
//      Option -v version
//      Option -b to build the project
//            Optionnaly build the project and download the zip with all languages. 
//      Returns 1 if there was an error
//      If option built is used, returns "built" or "skipped" if the command is successful and depending if the build was actually done.
//
//
//	cross compilation AMD64:  env GOOS=windows GOARCH=amd64 go build crowdinExport.go


import (
	"flag"
	"fmt"
	"os"
	"github.com/medisafe/go-crowdin"
)

func main() {

	var versionFlg bool
	var buildFlg bool

	const usageVersion   = "Display Version"
	const usageBuild   = "Request a build"
    
    // Have to create a specific set, the default one is poluted by some test stuff from another lib (?!) 
    checkFlags := flag.NewFlagSet("check", flag.ExitOnError)
    
	checkFlags.BoolVar(&versionFlg, "version", false, usageVersion)
	checkFlags.BoolVar(&versionFlg, "v", false, usageVersion + " (shorthand)")
	checkFlags.BoolVar(&buildFlg, "build", false, usageBuild)
	checkFlags.BoolVar(&buildFlg, "b", false, usageBuild + " (shorthand)")
	checkFlags.Usage = func() {
        fmt.Printf("Usage: %s <key> <project name> <path and name of zip>\n",os.Args[0])
        checkFlags.PrintDefaults()
    }

    // Check parameters
	checkFlags.Parse(os.Args[1:])
	
	if versionFlg {
		fmt.Printf("Version %s\n", "2019-01  v0.1.0")
		os.Exit(0)
	}

	if (buildFlg && len(os.Args) < 4) ||
       (!buildFlg && len(os.Args) < 3)   {
		checkFlags.Usage()  // Display usage
		fmt.Printf("Missing parameters\n")
        os.Exit(1)
	}
	
    // Parse the command parameters
	index := 0
	if buildFlg {
        index = 1
    }
    key := os.Args[1 + index]
    project := os.Args[2 + index]
    filename :=  os.Args[3 + index]

    //fmt.Printf("key=%s\n",key)
    //fmt.Printf("project=%s\n",key)
    //fmt.Printf("filename=%s\n",key)
    
    // Create a connection
    api := crowdin.New(key, project)
    
    //api.SetDebug(true, nil)

    if buildFlg {
        
        // Request a build
        response,err := api.ExportTranslations()
        
        if err !=nil {
            fmt.Printf("crowdinExport() build request error %s\n",err)
            os.Exit(1)
        }
        
        // Return "built" or "skipped"
        fmt.Printf("%s\r\n",response.Success.Status)
        
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
            fmt.Printf("crowdinExport() download error %s\n",err)
            os.Exit(1)
        }
    }
}
