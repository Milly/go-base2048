package base2048

import (
	"strconv"
)

type CorruptInputError int64

func (e CorruptInputError) Error() string {
	return "illegal base2048 data at input " + strconv.FormatInt(int64(e), 10)
}
