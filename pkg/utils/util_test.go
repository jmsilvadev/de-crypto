package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHexUint64_ValidHexString(t *testing.T) {
	t.Run("parse_hex_uint64_with_0x_prefix", func(t *testing.T) {
		result, err := ParseHexUint64("0x1a2b")
		assert.NoError(t, err)
		assert.Equal(t, uint64(6699), result)
	})

	t.Run("parse_hex_uint64_without_0x_prefix", func(t *testing.T) {
		result, err := ParseHexUint64("1a2b")
		assert.NoError(t, err)
		assert.Equal(t, uint64(6699), result)
	})

	t.Run("parse_hex_uint64_uppercase", func(t *testing.T) {
		result, err := ParseHexUint64("0x1A2B")
		assert.NoError(t, err)
		assert.Equal(t, uint64(6699), result)
	})

	t.Run("parse_hex_uint64_zero", func(t *testing.T) {
		result, err := ParseHexUint64("0x0")
		assert.NoError(t, err)
		assert.Equal(t, uint64(0), result)
	})

	t.Run("parse_hex_uint64_max_uint64", func(t *testing.T) {
		result, err := ParseHexUint64("0xffffffffffffffff")
		assert.NoError(t, err)
		assert.Equal(t, uint64(18446744073709551615), result)
	})
}

func TestParseHexUint64_InvalidHexString(t *testing.T) {
	t.Run("parse_hex_uint64_invalid_hex", func(t *testing.T) {
		_, err := ParseHexUint64("0xgg")
		assert.Error(t, err)
	})

	t.Run("parse_hex_uint64_empty_string", func(t *testing.T) {
		_, err := ParseHexUint64("")
		assert.Error(t, err)
	})

	t.Run("parse_hex_uint64_invalid_characters", func(t *testing.T) {
		_, err := ParseHexUint64("0x123g")
		assert.Error(t, err)
	})

	t.Run("parse_hex_uint64_negative_sign", func(t *testing.T) {
		_, err := ParseHexUint64("-0x123")
		assert.Error(t, err)
	})
}

func TestParseHexUint64_EdgeCases(t *testing.T) {
	t.Run("parse_hex_uint64_single_digit", func(t *testing.T) {
		result, err := ParseHexUint64("0xa")
		assert.NoError(t, err)
		assert.Equal(t, uint64(10), result)
	})

	t.Run("parse_hex_uint64_with_whitespace", func(t *testing.T) {

		_, err := ParseHexUint64(" 0x1a2b ")
		assert.Error(t, err)
	})

	t.Run("parse_hex_uint64_very_long_hex", func(t *testing.T) {
		longHex := "0x" + strings.Repeat("f", 20)
		_, err := ParseHexUint64(longHex)
		assert.Error(t, err) // Should overflow uint64
	})
}
