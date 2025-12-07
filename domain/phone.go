package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"

	"github.com/nyaruka/phonenumbers"
)

var ErrInvalidPhoneNumber error = errors.New("invalid phone number")

type Phone struct {
	number *phonenumbers.PhoneNumber
}

func NewPhone(number string) (Phone, error) {
	pnum, err := phonenumbers.Parse(strings.TrimSpace(number), "PL")
	if err != nil || !phonenumbers.IsValidNumber(pnum) {
		return Phone{}, ErrInvalidPhoneNumber
	}
	return Phone{
		number: pnum,
	}, nil
}

func (p Phone) Hash() string {
	sum := sha256.Sum256([]byte(p.String()))
	return hex.EncodeToString(sum[:])
}

func (p Phone) Masked() string {
	return maskExceptLastNDigitsWithPrefix2(p.String(), 3, '*')
}

func (p Phone) String() string {
	return phonenumbers.Format(
		p.number,
		phonenumbers.E164,
	)
}

func (p Phone) Code() string {
	return phonenumbers.GetRegionCodeForNumber(p.number)
}

func maskExceptLastNDigitsWithPrefix2(s string, n int, mask rune) string {
	s = strings.TrimSpace(s)
	if n < 0 {
		n = 0
	}
	rs := []rune(s)
	L := len(rs)
	if L == 0 {
		return s
	}

	// 1) Policz wszystkie cyfry
	totalDigits := 0
	for _, r := range rs {
		if unicode.IsDigit(r) {
			totalDigits++
		}
	}

	// 2) Policz, ile cyfr jest w "prefiksie", którego nie maskujemy:
	// pozycje 1 i 2 (rs[0] to '+', zgodnie z założeniem)
	prefixDigits := 0
	for i := 1; i <= 2 && i < L; i++ {
		if unicode.IsDigit(rs[i]) {
			prefixDigits++
		}
	}

	// 3) Liczba cyfr, które POTENCJALNIE mogą być maskowane
	maskableDigitsTotal := totalDigits - prefixDigits

	// Jeżeli poza prefiksem jest <= n cyfr, nic nie będziemy maskować
	if maskableDigitsTotal <= n {
		return s
	}

	// 4) Iteracja i budowa wyniku
	out := make([]rune, 0, L)
	maskableDigitsLeft := maskableDigitsTotal
	keep := n

	for i, r := range rs {
		// Zawsze przepuszczamy '+' i DWA KOLEJNE ZNAKI
		if i <= 2 {
			out = append(out, r)
			continue
		}

		if unicode.IsDigit(r) {
			// Maskujemy tylko CYFRY poza prefiksem
			if maskableDigitsLeft > keep {
				out = append(out, mask)
			} else {
				out = append(out, r) // wśród ostatnich n cyfr
			}
			maskableDigitsLeft--
		} else {
			// Inne znaki przepuszczamy
			out = append(out, r)
		}
	}
	return string(out)
}
