package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/jacobsa/go-serial/serial"
)

// GamepadBlock function returns the CLI commands for pbupdater gamepadblock
func GamepadBlock() cli.Command {
	return cli.Command{
		Name:  "gamepadblock",
		Usage: "Install avrdude, and update firmware of your GamepadBlock",
		Action: func(c *cli.Context) {
			valid := false
			for _, s := range []string{"readversion", "prepare", "update"} {
				if s == c.Args().First() {
					valid = true
				}
			}

			usage := func() {
				fmt.Println("Invalid/no subcommand supplied.")
				fmt.Println()
				fmt.Println("Usage:")
				fmt.Println("  petrockutil scan serial")
				fmt.Println("  # List all serial device ports")
				fmt.Println()
				fmt.Println("  petrockutil gamepadblock readversion <port>")
				fmt.Println("  # Reads the version number of the current firmware on the GamepadBlock at given port")
				fmt.Println()
				fmt.Println("  petrockutil gamepadblock prepare")
				fmt.Println("  # Installs avrdude to allow uploading of firmware .hex files to GamepadBlock")
				fmt.Println()
				fmt.Println("  petrockutil gamepadblock update <port>")
				fmt.Println("  # Updates firmware of GamepadBlock at provided port to most recent version")
				fmt.Println()
			}

			if valid == false {
				usage()
				return
			}

			switch c.Args().First() {
			case "readversion":
				if len(c.Args()) < 2 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

				portPath := c.Args()[1]
				// Set up options.
				options := serial.OpenOptions{
					PortName:              portPath,
					BaudRate:              115200,
					DataBits:              8,
					StopBits:              1,
					MinimumReadSize:       0,
					InterCharacterTimeout: 1000,
				}

				// Open the port.
				port, err := serial.Open(options)
				if err != nil {
					log.Fatalf("serial.Open: %v", err)
				}

				defer port.Close() // Make sure to close it later.

				b := []byte{0x76} // "v"
				_, err = port.Write(b)
				if err != nil {
					log.Fatalf("port.Write: %v", err)
				}
				buf := make([]byte, 16)
				numBytes, err := port.Read(buf)
				if err != nil {
					fmt.Println("Could not read firmware version on port " + portPath)
					return
				}
				versionString := buf[:numBytes]

				fmt.Println("Found GamepadBlock at", portPath, ". The firmware version of it is", string(versionString))

			case "prepare":
				switch runtime.GOOS {
				case "linux":
					fmt.Println("Attempting to install avrdude with apt-get...")
					cmd := exec.Command("sudo", "apt-get", "-y", "install", "avrdude")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "darwin":
					fmt.Println("Attempting to install avrdude with Homebrew...")
					cmd := exec.Command("brew", "install", "avrdude")
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				case "windows":
					_, err := exec.Command("NET", "SESSION").Output()
					if err != nil {
						fmt.Println("Please run cmd.exe as administrator and try again")
						os.Exit(1)
					}

					fmt.Println("Installing winavr...")
					fmt.Printf("Initial value for $PATH: %s\n", os.Getenv("PATH"))
					dirName, _ := createPetrockblockDirectory()
					exeFile := "https://s3.amazonaws.com/gort-io/support/WinAVR-20100110-install.exe"
					fileName := downloadFromUrl(dirName, exeFile)
					cmd := exec.Command(petrockblockDirName() + "\\" + fileName)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}

				default:
					fmt.Println("OS not yet supported.")
				}
			case "update":
				if len(c.Args()) < 1 {
					fmt.Println("Invalid number of arguments.")
					usage()
					return
				}

				port := c.Args()[1]
				file, _ := ioutil.TempFile(os.TempDir(), "")
				defer file.Close()
				defer os.Remove(file.Name())

				fmt.Println("Downloading the latest firmware")
				filename := downloadLatestFirmware()
				fmt.Println("Downloaded firmware file", filename)

				fmt.Println("Press the RESET button on the GameapadBlock to activate its update mode.")
				fmt.Print("Press 'Enter' to continue...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')

				switch runtime.GOOS {
				case "darwin", "linux", "windows":
					cmd := exec.Command("avrdude", "-pm32u2", "-cavr109", fmt.Sprintf("-P%v", port), "-D", fmt.Sprintf("-Uflash:w:%v:a", filename))
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
						fmt.Println("An error occurred during the firmware update process")
					}
					fmt.Println("Finished firmware update on port", port)

				default:
					fmt.Println("OS not yet supported.")
				}
			}
		},
	}
}

func downloadLatestFirmware() string {
	req, err := http.NewRequest("GET", "https://github.com/petrockblog/GamepadBlockUpdater/releases/latest", nil)
	if err != nil {
		panic(err)
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("Redirect")
	}

	response, err := client.Do(req)
	currentVersion := ""
	if err != nil {
		if response.StatusCode == http.StatusFound { //status code 302
			tempurl, _ := response.Location()
			s := strings.Split(tempurl.Path, "tag/")
			_, currentVersion = s[0], s[1]
			if currentVersion != "" {
				fmt.Println("Most current firmware version is", currentVersion)
				downloadURL := tempurl.Scheme + "://" + tempurl.Host + tempurl.Path
				fmt.Println(downloadURL)
				downloadURL = strings.Replace(downloadURL, "releases/tag", "releases/download", 1)
				downloadURL += "/firmware.hex"
				downloadedFile := downloadFromUrl(".", downloadURL)
				if downloadedFile == "" {
					panic("Could not download latest firmware file.")
				} else {
					return downloadedFile
				}
			} else {
				fmt.Println("Could not determine most recent firmware version")
			}
		} else {
			panic(err)
		}
	}
	return ""
}
