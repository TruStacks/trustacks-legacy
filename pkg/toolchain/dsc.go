package toolchain

type DesiredStateConfiguration struct {
}

func (dsc *DesiredStateConfiguration) Store(config interface{}) error {
	return nil
}

func (dsc *DesiredStateConfiguration) Load(config interface{}) (interface{}, error) {
	return nil, nil
}
