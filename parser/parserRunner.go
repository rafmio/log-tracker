package parser

func ParserRunner() error {
	fileConfig, err := ReadFileConfig("/home/raf/log-tracker/config/fileConfig.json")
	if err != nil {
		return err
	}

	file, err := SelectAndOpen(fileConfig)
	if err != nil {
		return err
	}
	defer file.Close()

	fp := new(FilePosition)
	err = fp.ReadFPFromFile()
}
