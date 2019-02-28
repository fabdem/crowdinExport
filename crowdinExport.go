package main

//	F.Demurger 2019-02
//  	3 args: key, project name and zip file name
//
//      Option -v version
//      Option -b to build the project
//            Optionnaly build the project and download the zip with all languages.
//      Option -p <proxy url> to use a proxy.
//      Option -t <timeout in second>.
//
//      Overall default timeout set in lib: 40s
//      Returns 1 if there was an error
//      If option built is used, returns "built" or "skipped" if the command is successful and depending if the build was actually done.
//
//
//	cross compilation AMD64:  env GOOS=windows GOARCH=amd64 go build crowdinExport.go


import (
	"flag"
	"fmt"
	"os"
	"github.com/fabdem/go-crowdinProxy"
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
	var timeoutsec int


	const usageVersion   = "Display Version"
	const usageBuild     = "Request a build"
	const usageProxy     = "Use a proxy - followed with url"
	const usageTimeout   = "Set the communication timeout in seconds (default 40s)."

  // Have to create a specific set, the default one is poluted by some test stuff from another lib (?!)
  checkFlags := flag.NewFlagSet("check", flag.ExitOnError)

	checkFlags.BoolVar(&versionFlg, "version", false, usageVersion)
	checkFlags.BoolVar(&versionFlg, "v", false, usageVersion + " (shorthand)")
	checkFlags.IntVar(&timeoutsec, "timeout", 0, usageTimeout)
	checkFlags.IntVar(&timeoutsec, "t", 0, usageTimeout + " (shorthand)")
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
		fmt.Printf("Version %s\n", "2019-02  v1.2.0")
		os.Exit(0)
	}

  // Parse the command parameters
  index := len(os.Args)
  key := os.Args[index - 3]
  project := os.Args[index - 2]
  filename :=  os.Args[index - 1]

  // fmt.Printf("buildFlg=%s\n",buildFlg)
	// fmt.Printf("timeoutsec=%s\n",timeoutsec)
  // fmt.Printf("proxy=%s\n",proxy)
  // fmt.Printf("key=%s\n",key)
  // fmt.Printf("project=%s\n",project)
  // fmt.Printf("filename=%s\n",filename)
  // os.Exit(1)

  // Create a connection
  // var api *crowdin.Crowdin
  crowdinproxy.SetTimeouts(5, timeoutsec) // in sec
  api,err := crowdinproxy.New(key, project, proxy)
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
  }


  // request zip download
  opt := crowdin.DownloadOptions{Package: "all", LocalPath: filename}

  err = api.DownloadTranslations(&opt)
  if err !=nil {
      fmt.Printf("\ncrowdinExport() download error %s\n",err)
      os.Exit(1)
  }

  close(finishChan)  // Stop animation

  // Return "built" or "skipped"
  fmt.Printf("\b%s\r\n",result)
}
