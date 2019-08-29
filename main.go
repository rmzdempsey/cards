package main

func main() {

	config := CardServerConfig{
		port : 8000,
		gameConfigs : []CardGameConfig{
				CardGameConfig{ Name:"poker" }, 
			},	
	}

	startCardServer(config)
}
