package jsonmask

// Mask selects the specific parts of an object, according to the "mask".
func Mask(obj interface{}, mask string) (interface{}, error) {
	compiledMask, err := compile(mask)
	if err != nil {
		return nil, err
	}

	filteredObj, err := filter(obj, compiledMask)
	if err != nil {
		return nil, err
	}

	return filteredObj, nil
}
