package jsonmask

import (
	"errors"
)

// Mask selects the specific parts of an object, according to the "mask".
func Mask(obj interface{}, mask string) (res interface{}, err error) {
	defer func() {
		if r := recover(); err != nil {
			res = obj

			switch r.(type) {
			case error:
				err = r.(error)
			case string:
				err = errors.New(r.(string))
			default:
				err = errors.New("json mask panic")
			}
		}
	}()

	compiledMask, err := compile(mask)
	if err != nil {
		return nil, err
	}

	res, err = filter(obj, compiledMask)
	if err != nil {
		return nil, err
	}

	return res, nil
}
