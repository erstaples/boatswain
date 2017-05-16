package lib

// type ServiceMapConfig struct {
// 	ServiceMaps []ServiceMap      `yaml:"ServiceMaps"`
// 	Ingress     ServiceMapIngress `yaml:"Ingress"`
// }

// type ServiceMap struct {
// 	Name                   string   `yaml:"Name"`
// 	Test                   []string `yaml:"Test"`
// 	Staging                []string `yaml:"Staging"`
// 	CloudFormationTemplate string   `yaml:"CloudFormationTemplate"`
// }

// type ServiceMapIngress struct {
// 	Template string `yaml:"Template"`
// 	Service  string `yaml:"Service"`
// 	Port     string `yaml:"Port"`
// }

// func loadServiceMap() {
// 	//get the service serviceMap file
// 	path := viper.GetString("release") //TODO: capitalize this
// 	fullPath := path + "/.servicemap/staging.yaml"
// 	valuesBytes, err := ioutil.ReadFile(fullPath)

// 	if err != nil {
// 		panic(err)
// 	}

// 	err = yaml.Unmarshal(valuesBytes, &serviceMapConfig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	for _, sMap := range serviceMapConfig.ServiceMaps {
// 		if serviceMapName == sMap.Name {
// 			serviceMap = sMap
// 		}
// 	}
// }

// func convertMapToEnvVars(serviceMap ServiceMap) map[string]string {
// 	env := make(map[string]string)

// 	for _, testSvc := range serviceMap.Test {
// 		env[testSvc] = branchName + "-" + testSvc
// 	}

// 	for _, stagingSvc := range serviceMap.Staging {
// 		env[stagingSvc] = "staging-" + stagingSvc
// 	}

// 	return env
// }
