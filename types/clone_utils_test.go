package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const modifiedText = "modified"

func TestDeepCopyAny(t *testing.T) {
	tests := []struct {
		name       string
		input      any
		want       any
		validateFn func(t *testing.T, original, clone any)
	}{
		{
			name:  "nil value",
			input: nil,
			want:  nil,
		},
		{
			name:  "string value",
			input: "test string",
			want:  "test string",
		},
		{
			name:  "int value",
			input: 42,
			want:  42,
		},
		{
			name:  "float value",
			input: 3.14,
			want:  3.14,
		},
		{
			name:  "bool value",
			input: true,
			want:  true,
		},
		{
			name:  "slice of strings",
			input: []string{"a", "b", "c"},
			want:  []string{"a", "b", "c"},
			validateFn: func(t *testing.T, original, clone any) {
				origSlice := original.([]string)
				cloneSlice := clone.([]string)

				// Modify original
				origSlice[0] = modifiedText

				// Ensure clone is unaffected
				assert.Equal(t, "a", cloneSlice[0])
				assert.NotSame(t, &origSlice, &cloneSlice)
			},
		},
		{
			name:  "slice of any",
			input: []any{"string", 42, true, nil},
			want:  []any{"string", 42, true, nil},
			validateFn: func(t *testing.T, original, clone any) {
				origSlice := original.([]any)
				cloneSlice := clone.([]any)

				// Modify original
				origSlice[0] = modifiedText

				// Ensure clone is unaffected
				assert.Equal(t, "string", cloneSlice[0])
				assert.NotSame(t, &origSlice, &cloneSlice)
			},
		},
		{
			name: "map[string]any",
			input: map[string]any{
				"string": "value",
				"number": 42,
				"bool":   true,
				"null":   nil,
			},
			want: map[string]any{
				"string": "value",
				"number": 42,
				"bool":   true,
				"null":   nil,
			},
			validateFn: func(t *testing.T, original, clone any) {
				origMap := original.(map[string]any)
				cloneMap := clone.(map[string]any)

				// Modify original
				origMap["string"] = modifiedText

				// Ensure clone is unaffected
				assert.Equal(t, "value", cloneMap["string"])
				assert.NotSame(t, &origMap, &cloneMap)
			},
		},
		{
			name: "nested structures",
			input: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": []any{"deep", "nested", "data"},
					},
					"siblings": []string{"a", "b"},
				},
				"array": []any{
					map[string]any{"key": "value1"},
					map[string]any{"key": "value2"},
				},
			},
			want: map[string]any{
				"level1": map[string]any{
					"level2": map[string]any{
						"level3": []any{"deep", "nested", "data"},
					},
					"siblings": []string{"a", "b"},
				},
				"array": []any{
					map[string]any{"key": "value1"},
					map[string]any{"key": "value2"},
				},
			},
			validateFn: func(t *testing.T, original, clone any) {
				origMap := original.(map[string]any)
				cloneMap := clone.(map[string]any)

				// Modify deeply nested values
				level1 := origMap["level1"].(map[string]any)
				level2 := level1["level2"].(map[string]any)
				level3 := level2["level3"].([]any)
				level3[0] = modifiedText

				array := origMap["array"].([]any)
				firstItem := array[0].(map[string]any)
				firstItem["key"] = modifiedText

				// Ensure clone is unaffected
				clonedLevel1 := cloneMap["level1"].(map[string]any)
				clonedLevel2 := clonedLevel1["level2"].(map[string]any)
				clonedLevel3 := clonedLevel2["level3"].([]any)
				assert.Equal(t, "deep", clonedLevel3[0])

				clonedArray := cloneMap["array"].([]any)
				clonedFirstItem := clonedArray[0].(map[string]any)
				assert.Equal(t, "value1", clonedFirstItem["key"])
			},
		},
		{
			name: "struct type",
			input: struct {
				Name  string
				Value int
			}{
				Name:  "test",
				Value: 123,
			},
			want: struct {
				Name  string
				Value int
			}{
				Name:  "test",
				Value: 123,
			},
		},
		{
			name:  "pointer to int",
			input: func() *int { i := 42; return &i }(),
			want:  func() *int { i := 42; return &i }(),
			validateFn: func(t *testing.T, original, clone any) {
				origPtr := original.(*int)
				clonePtr := clone.(*int)

				// Modify original
				*origPtr = 999

				// Ensure clone is unaffected
				assert.Equal(t, 42, *clonePtr)
				assert.NotSame(t, origPtr, clonePtr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deepCopyAny(tt.input)
			assert.Equal(t, tt.want, got)

			if tt.validateFn != nil {
				tt.validateFn(t, tt.input, got)
			}
		})
	}
}

func TestDeepCopyAny_EdgeCases(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		input := []string{}
		got := deepCopyAny(input)

		result, ok := got.([]string)
		assert.True(t, ok)
		assert.Empty(t, result)
		assert.NotSame(t, &input, &result)
	})

	t.Run("empty map", func(t *testing.T) {
		input := map[string]any{}
		got := deepCopyAny(input)

		result, ok := got.(map[string]any)
		assert.True(t, ok)
		assert.Empty(t, result)
		assert.NotSame(t, &input, &result)
	})

	t.Run("nil slice", func(t *testing.T) {
		var input []string
		got := deepCopyAny(input)

		result, ok := got.([]string)
		assert.True(t, ok)
		assert.Nil(t, result)
	})

	t.Run("nil map", func(t *testing.T) {
		var input map[string]any
		got := deepCopyAny(input)

		result, ok := got.(map[string]any)
		assert.True(t, ok)
		assert.Nil(t, result)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var input *int
		got := deepCopyAny(input)

		result, ok := got.(*int)
		assert.True(t, ok)
		assert.Nil(t, result)
	})
}
