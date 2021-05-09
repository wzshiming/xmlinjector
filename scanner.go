package xmlinjector

import (
	"bytes"
	"fmt"
)

const (
	prefix  = `<!--`
	suffix  = `-->`
	endChar = `/`
)

var (
	errNotMatch = fmt.Errorf("not match")
	errEmptyKey = fmt.Errorf("key cannot be empty")
)

func scanPairAnnotationElement(key []byte, data []byte) (args []byte, begin int, end int, single bool, err error) {
	if len(key) == 0 {
		return nil, 0, 0, false, errEmptyKey
	}
	keyIndex := 0
	prefixContent, _, prefixEnd, err := scanTryAnnotationElement(data, 0,
		func(content []byte) bool {
			index := strIndexIgnoreSpace(content, key)
			if index == -1 {
				return false
			}
			if index+len(key) != len(content) &&
				content[index+len(key)] != ' ' {
				return false
			}
			keyIndex = index
			return true
		})
	if err != nil {
		return nil, 0, 0, false, fmt.Errorf("%w: prefix element %q", err, key)
	}

	if bytes.HasSuffix(prefixContent, []byte(endChar)) {
		args = prefixContent[keyIndex+len(key) : len(prefixContent)-1]
		args = bytes.TrimSpace(args)
		end = prefixEnd - len(suffix)
		begin = end - len(endChar)
		return args, begin, end, true, nil
	}

	_, suffixBegin, _, err := scanTryAnnotationElement(data[prefixEnd:], 0,
		func(content []byte) bool {
			index := strIndexIgnoreSpace(content, []byte(endChar))
			if index == -1 {
				return false
			}
			if !bytes.HasPrefix(content[index+1:], key) {
				return false
			}
			keyEnd := index + 1 + len(key)
			if keyEnd != len(content) &&
				content[keyEnd] != ' ' {
				return false
			}
			return true
		})
	if err != nil {
		return nil, 0, 0, false, fmt.Errorf("%w: suffix element %q", err, key)
	}

	args = prefixContent[keyIndex+len(key):]
	args = bytes.TrimSpace(args)
	begin = prefixEnd
	end = prefixEnd + suffixBegin
	return args, begin, end, false, nil
}

func scanTryAnnotationElement(data []byte, off int, check func(content []byte) bool) (content []byte, begin int, end int, err error) {
	content, begin, end, err = scanAnnotationElement(data)
	if err != nil {
		return nil, 0, 0, err
	}
	if !check(content) {
		return scanTryAnnotationElement(data[end:], end+off, check)
	}
	return content, begin + off, end + off, nil
}

func scanAnnotationElement(data []byte) (content []byte, begin int, end int, err error) {
	begin = bytes.Index(data, []byte(prefix))
	if begin == -1 {
		return nil, 0, 0, fmt.Errorf("%w: prefix %q", errNotMatch, prefix)
	}
	contentBegin := begin + len(prefix)
	offEnd := bytes.Index(data[contentBegin:], []byte(suffix))
	if offEnd == -1 {
		return nil, 0, 0, fmt.Errorf("%w: suffix %q", errNotMatch, suffix)
	}
	contentEnd := contentBegin + offEnd
	end = contentEnd + len(suffix)
	return data[contentBegin:contentEnd], begin, end, nil
}

func strIndexIgnoreSpace(s, sep []byte) int {
	off := 0
	for ; off != len(s); off++ {
		if !isSpace(s[off]) {
			break
		}
	}
	if !bytes.HasPrefix(s[off:], sep) {
		return -1
	}
	return off
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t'
}
