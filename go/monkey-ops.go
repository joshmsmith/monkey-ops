package main

import (
	"io/ioutil"
	"log"
	"net/http"

	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {

	//TODO: API_SERVER env variable never actually used, can be removed
	flag.String("API_SERVER", "", "API Server URL")
	flag.String("PROJECT_NAME", "", "Project to get crazy")
	//TODO: TOKEN env variable never actually used, can be removed
	flag.String("TOKEN", "", "Bearer token with edit grants to access to the Openshift project")
	flag.Float64("INTERVAL", 30, "interval time in seconds")
	flag.String("MODE", "background", "Execution mode: background or rest")

	//Binding flags and env vars
	viper.BindPFlag("API_SERVER", flag.Lookup("API_SERVER"))
	viper.BindPFlag("PROJECT_NAME", flag.Lookup("PROJECT_NAME"))
	viper.BindPFlag("TOKEN", flag.Lookup("TOKEN"))
	viper.BindPFlag("INTERVAL", flag.Lookup("INTERVAL"))
	viper.BindPFlag("MODE", flag.Lookup("MODE"))

	viper.BindEnv("KUBERNETES_SERVICE_HOST")
	viper.BindEnv("KUBERNETES_SERVICE_PORT")
	viper.BindEnv("API_SERVER")
	viper.BindEnv("PROJECT_NAME")
	viper.BindEnv("TOKEN")
	viper.BindEnv("INTERVAL")
	viper.BindEnv("MODE")

	flag.Parse()

	//set configuration
	var apiServer string
	if viper.GetString("KUBERNETES_SERVICE_HOST") != "" && viper.GetString("KUBERNETES_SERVICE_PORT") != "" {
		apiServer = "https://" + viper.GetString("KUBERNETES_SERVICE_HOST") + ":" + viper.GetString("KUBERNETES_SERVICE_PORT")
	} else {
		apiServer = viper.GetString("API_SERVER")
	}
	project := viper.GetString("PROJECT_NAME")
	token := viper.GetString("TOKEN")
	interval := viper.GetFloat64("INTERVAL")
	mode := viper.GetString("MODE")

	if mode == "background" {
		// read the service account secret token file at once
		tokenBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
		if err != nil {
			log.Println("Not Service Account Token available")
		} else {
			token = string(tokenBytes[:])
		}

		//validating some required parameters
		if (apiServer == "" || project == "" || token == "") {
			log.Fatal("Required Input Parameters not valid")
		}

		chaosInput := ChaosInput{
			Url:       apiServer,
			Project:   project,
			Token:     token,
			Interval:  interval,
			TotalTime: 0,
		}

		//Launch the chaos
		go ExecuteChaos(&chaosInput, mode)
	}

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
