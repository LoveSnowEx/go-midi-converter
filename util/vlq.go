package util

const (
	vlqContinue = 0x80
	vlqMask     = 0x7f
)

func VlqEncode(n uint32) (out []byte) {
	var quo, rem uint32
	quo = n / vlqContinue
	rem = n % vlqContinue

	out = append(out, byte(rem))

	for quo > 0 {
		out = append(out, byte(quo)|vlqContinue)
		quo = quo / vlqContinue
		// rem = quo % vlqContinue
	}

	reverse(out)
	return
}

func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func VlqDecode(source []byte) (num uint32) {
	for i := 0; i < len(source); i++ {
		var n = uint32(source[i] & vlqMask)
		for (source[i] & vlqContinue) != 0 {
			i++
			n *= 128
			n += uint32(source[i] & vlqMask)
		}
		num += n
	}

	return
}

func Uint32ToBytes(val uint32) (out []byte) {
	out = []byte{
		byte(val >> 24),
		byte(val >> 16),
		byte(val >> 8),
		byte(val),
	}
	return
}
