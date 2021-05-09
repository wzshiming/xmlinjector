package xmlinjector

import (
	"errors"
)

func Inject(key []byte, data []byte, inject func(args, origin []byte) []byte) ([]byte, error) {
	var list []injectItem
	size := len(data)
	off := 0
	for {
		args, begin, end, single, err := scanPairAnnotationElement(key, data[off:])
		if err != nil {
			if errors.Is(err, errNotMatch) {
				break
			}
			return nil, err
		}
		begin += off
		end += off
		content := inject(args, data[begin:end])

		size += len(content) - (end - begin)
		if single {
			size += len(suffix) + len(prefix) + len(endChar) + len(key) + 2
		}
		item := injectItem{
			begin:   begin,
			end:     end,
			single:  single,
			content: content,
		}
		list = append(list, item)
		off = end + len(suffix)
	}

	out := make([]byte, 0, size)
	end := 0
	for _, item := range list {
		out = append(out, data[end:item.begin]...)
		if item.single {
			out = append(out, suffix...)
			out = append(out, item.content...)
			out = append(out, prefix...)
			out = append(out, ' ')
			out = append(out, endChar...)
			out = append(out, key...)
			out = append(out, ' ')
		} else {
			out = append(out, item.content...)
		}
		end = item.end
	}
	out = append(out, data[end:]...)
	return out, nil
}

type injectItem struct {
	begin   int
	end     int
	single  bool
	content []byte
}
